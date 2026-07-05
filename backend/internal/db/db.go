package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/vectorsighttechnologies/serverless-orchestrator/backend/internal/types"
	_ "modernc.org/sqlite"
)

// DB wrap the database client and encryption keys.
type DB struct {
	SQLDB         *sql.DB
	encryptionKey []byte
}

// NewDB initializes a connection pool and migrates database schemas.
func NewDB(driver, connectionURL, encKeyHex string) (*DB, error) {
	// Parse the encryption key
	key, err := hex.DecodeString(encKeyHex)
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("invalid encryption key: must be a 32-byte hex string (64 characters)")
	}

	// Normalise driver name
	sqlDriver := driver
	if sqlDriver == "sqlite" {
		sqlDriver = "sqlite"
	} else if sqlDriver == "postgres" {
		sqlDriver = "postgres"
	} else {
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := sql.Open(sqlDriver, connectionURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	instance := &DB{
		SQLDB:         db,
		encryptionKey: key,
	}

	if err := instance.migrate(); err != nil {
		return nil, fmt.Errorf("failed database migrations: %w", err)
	}

	return instance, nil
}

// Close terminates connection pool operations.
func (d *DB) Close() error {
	return d.SQLDB.Close()
}

// migrate creates tables dynamically on startup.
func (d *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS user_preferences (
			user_id VARCHAR(36) PRIMARY KEY,
			selected_provider VARCHAR(50) DEFAULT 'newrelic',
			nr_account_id VARCHAR(255),
			nr_api_key_encrypted TEXT,
			nr_license_key_encrypted TEXT,
			nr_region VARCHAR(10),
			lambda_api_url_encrypted TEXT,
			lambda_api_key_encrypted TEXT,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS user_connections (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			name VARCHAR(255) NOT NULL,
			aws_region VARCHAR(50) NOT NULL,
			lambda_api_url_encrypted TEXT NOT NULL,
			lambda_api_key_encrypted TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			action VARCHAR(100) NOT NULL,
			target TEXT,
			status VARCHAR(50) NOT NULL,
			details TEXT,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, q := range queries {
		if _, err := d.SQLDB.Exec(q); err != nil {
			return err
		}
	}

	_ = d.addColumnIfNotExists("user_preferences", "selected_provider", "VARCHAR(50) DEFAULT 'newrelic'")
	_ = d.addColumnIfNotExists("user_preferences", "dd_api_key_encrypted", "TEXT")
	_ = d.addColumnIfNotExists("user_preferences", "dd_site", "VARCHAR(50)")

	return nil
}

func (d *DB) addColumnIfNotExists(table, column, colType string) error {
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, colType)
	_, err := d.SQLDB.Exec(query)
	if err != nil {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "duplicate column") || strings.Contains(errStr, "already exists") {
			return nil
		}
		return err
	}
	return nil
}

// ─────────────────────────────────────────────────────────────
// Cryptographic Utility (AES-256-GCM)
// ─────────────────────────────────────────────────────────────

// Encrypt encrypts plain text into a Base64-encoded GCM block.
func (d *DB) Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	block, err := aes.NewCipher(d.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decodes Base64 data and decrypts the GCM block.
func (d *DB) Decrypt(cipherTextB64 string) (string, error) {
	if cipherTextB64 == "" {
		return "", nil
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(d.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, encryptedPayload := cipherText[:nonceSize], cipherText[nonceSize:]
	plainTextBytes, err := gcm.Open(nil, nonce, encryptedPayload, nil)
	if err != nil {
		return "", err
	}

	return string(plainTextBytes), nil
}

// ─────────────────────────────────────────────────────────────
// User Registration & Authentication Queries
// ─────────────────────────────────────────────────────────────

// CreateUser inserts a new user record.
func (d *DB) CreateUser(email, passwordHash string) (string, error) {
	userID := generateUUID()
	query := `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)`
	_, err := d.SQLDB.Exec(query, userID, strings.ToLower(email), passwordHash)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// GetUserByEmail retrieves credentials by user email.
func (d *DB) GetUserByEmail(email string) (string, string, error) {
	var id, hash string
	query := `SELECT id, password_hash FROM users WHERE email = $1`
	err := d.SQLDB.QueryRow(query, strings.ToLower(email)).Scan(&id, &hash)
	if err != nil {
		return "", "", err
	}
	return id, hash, nil
}

// ─────────────────────────────────────────────────────────────
// User Preference Queries
// ─────────────────────────────────────────────────────────────

// GetUserPreferences fetches user configs and decrypts parameters.
func (d *DB) GetUserPreferences(userID string) (*types.UserPreferences, error) {
	var (
		selectedProvider         sql.NullString
		nrAccountID              sql.NullString
		nrApiKeyEncrypted        sql.NullString
		nrLicenseKeyEncrypted    sql.NullString
		nrRegion                 sql.NullString
		ddApiKeyEncrypted        sql.NullString
		ddSite                   sql.NullString
		lambdaApiUrlEncrypted    sql.NullString
		lambdaApiKeyEncrypted    sql.NullString
	)

	query := `SELECT selected_provider, nr_account_id, nr_api_key_encrypted, nr_license_key_encrypted, 
	                 nr_region, dd_api_key_encrypted, dd_site, lambda_api_url_encrypted, lambda_api_key_encrypted 
	          FROM user_preferences WHERE user_id = $1`

	err := d.SQLDB.QueryRow(query, userID).Scan(
		&selectedProvider, &nrAccountID, &nrApiKeyEncrypted, &nrLicenseKeyEncrypted,
		&nrRegion, &ddApiKeyEncrypted, &ddSite, &lambdaApiUrlEncrypted, &lambdaApiKeyEncrypted,
	)

	if err == sql.ErrNoRows {
		// Return empty struct if preferences aren't set yet
		return &types.UserPreferences{}, nil
	} else if err != nil {
		return nil, err
	}

	// Decrypt sensitive credentials
	nrApiKey, _ := d.Decrypt(nrApiKeyEncrypted.String)
	nrLicenseKey, _ := d.Decrypt(nrLicenseKeyEncrypted.String)
	ddApiKey, _ := d.Decrypt(ddApiKeyEncrypted.String)
	lambdaApiUrl, _ := d.Decrypt(lambdaApiUrlEncrypted.String)
	lambdaApiKey, _ := d.Decrypt(lambdaApiKeyEncrypted.String)

	provider := selectedProvider.String
	if provider == "" {
		provider = "newrelic"
	}

	return &types.UserPreferences{
		SelectedProvider: provider,
		NRAccountID:      nrAccountID.String,
		NRApiKey:         nrApiKey,
		NRLicenseKey:     nrLicenseKey,
		NRRegion:         nrRegion.String,
		DDApiKey:         ddApiKey,
		DDSite:           ddSite.String,
		LambdaAPIURL:     lambdaApiUrl,
		LambdaAPIKey:     lambdaApiKey,
	}, nil
}

// SaveUserPreferences saves and encrypts user configs.
func (d *DB) SaveUserPreferences(userID string, prefs *types.UserPreferences) error {
	nrApiKeyEncrypted, err := d.Encrypt(prefs.NRApiKey)
	if err != nil {
		return err
	}
	nrLicenseKeyEncrypted, err := d.Encrypt(prefs.NRLicenseKey)
	if err != nil {
		return err
	}
	ddApiKeyEncrypted, err := d.Encrypt(prefs.DDApiKey)
	if err != nil {
		return err
	}
	lambdaApiUrlEncrypted, err := d.Encrypt(prefs.LambdaAPIURL)
	if err != nil {
		return err
	}
	lambdaApiKeyEncrypted, err := d.Encrypt(prefs.LambdaAPIKey)
	if err != nil {
		return err
	}

	provider := prefs.SelectedProvider
	if provider == "" {
		provider = "newrelic"
	}

	// Check if preferences already exist for user to decide INSERT or UPDATE (upsert)
	var count int
	err = d.SQLDB.QueryRow(`SELECT COUNT(*) FROM user_preferences WHERE user_id = $1`, userID).Scan(&count)
	if err != nil {
		return err
	}

	var query string
	if count > 0 {
		query = `UPDATE user_preferences 
		         SET selected_provider = $1, nr_account_id = $2, nr_api_key_encrypted = $3, 
		             nr_license_key_encrypted = $4, nr_region = $5, dd_api_key_encrypted = $6,
		             dd_site = $7, lambda_api_url_encrypted = $8, lambda_api_key_encrypted = $9, 
		             updated_at = $10 
		         WHERE user_id = $11`
		_, err = d.SQLDB.Exec(query,
			provider, prefs.NRAccountID, nrApiKeyEncrypted,
			nrLicenseKeyEncrypted, prefs.NRRegion, ddApiKeyEncrypted,
			prefs.DDSite, lambdaApiUrlEncrypted, lambdaApiKeyEncrypted,
			time.Now(), userID,
		)
	} else {
		query = `INSERT INTO user_preferences 
		         (user_id, selected_provider, nr_account_id, nr_api_key_encrypted, 
		          nr_license_key_encrypted, nr_region, dd_api_key_encrypted, dd_site, 
		          lambda_api_url_encrypted, lambda_api_key_encrypted) 
		         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		_, err = d.SQLDB.Exec(query,
			userID, provider, prefs.NRAccountID, nrApiKeyEncrypted,
			nrLicenseKeyEncrypted, prefs.NRRegion, ddApiKeyEncrypted, prefs.DDSite,
			lambdaApiUrlEncrypted, lambdaApiKeyEncrypted,
		)
	}

	return err
}

// ─────────────────────────────────────────────────────────────
// Audit Log Queries
// ─────────────────────────────────────────────────────────────

// CreateAuditLog creates a record in the audit logs table.
func (d *DB) CreateAuditLog(userID, action, target, status, details string) error {
	logID := generateUUID()
	query := `INSERT INTO audit_logs (id, user_id, action, target, status, details) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := d.SQLDB.Exec(query, logID, userID, action, target, status, details)
	return err
}

// generateUUID returns a standard formatted UUIDv4-like string using crypto/rand.
func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// GetUserConnections fetches all connections for a user and decrypts endpoints.
func (d *DB) GetUserConnections(userID string) ([]types.UserConnection, error) {
	query := `SELECT id, name, aws_region, lambda_api_url_encrypted, lambda_api_key_encrypted 
	          FROM user_connections WHERE user_id = $1 ORDER BY created_at ASC`
	rows, err := d.SQLDB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conns []types.UserConnection
	for rows.Next() {
		var c types.UserConnection
		var urlEnc, keyEnc string
		if err := rows.Scan(&c.ID, &c.Name, &c.AWSRegion, &urlEnc, &keyEnc); err != nil {
			return nil, err
		}
		c.LambdaAPIURL, _ = d.Decrypt(urlEnc)
		c.LambdaAPIKey, _ = d.Decrypt(keyEnc)
		c.HasAPIKey = c.LambdaAPIKey != ""
		if c.HasAPIKey {
			c.LambdaAPIKey = "••••••••"
		}
		conns = append(conns, c)
	}
	return conns, nil
}

// GetUserConnection fetches a specific connection by ID for verification.
func (d *DB) GetUserConnection(userID string, connID string) (*types.UserConnection, error) {
	query := `SELECT id, name, aws_region, lambda_api_url_encrypted, lambda_api_key_encrypted 
	          FROM user_connections WHERE user_id = $1 AND id = $2`
	var c types.UserConnection
	var urlEnc, keyEnc string
	err := d.SQLDB.QueryRow(query, userID, connID).Scan(&c.ID, &c.Name, &c.AWSRegion, &urlEnc, &keyEnc)
	if err != nil {
		return nil, err
	}
	c.LambdaAPIURL, _ = d.Decrypt(urlEnc)
	c.LambdaAPIKey, _ = d.Decrypt(keyEnc)
	c.HasAPIKey = c.LambdaAPIKey != ""
	return &c, nil
}

// SaveUserConnection encrypts credentials and updates or inserts the connection.
func (d *DB) SaveUserConnection(userID string, c *types.UserConnection) (string, error) {
	urlEnc, err := d.Encrypt(c.LambdaAPIURL)
	if err != nil {
		return "", err
	}

	var keyEnc string
	if c.LambdaAPIKey != "••••••••" && c.LambdaAPIKey != "" {
		keyEnc, err = d.Encrypt(c.LambdaAPIKey)
		if err != nil {
			return "", err
		}
	}

	if c.ID == "" {
		c.ID = generateUUID()
		query := `INSERT INTO user_connections (id, user_id, name, aws_region, lambda_api_url_encrypted, lambda_api_key_encrypted) 
		          VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = d.SQLDB.Exec(query, c.ID, userID, c.Name, c.AWSRegion, urlEnc, keyEnc)
	} else {
		// If key is obscured, preserve existing key
		if c.LambdaAPIKey == "••••••••" || c.LambdaAPIKey == "" {
			existing, err := d.GetUserConnection(userID, c.ID)
			if err == nil {
				keyEnc, _ = d.Encrypt(existing.LambdaAPIKey)
			}
		}
		query := `UPDATE user_connections 
		          SET name = $1, aws_region = $2, lambda_api_url_encrypted = $3, lambda_api_key_encrypted = $4 
		          WHERE user_id = $5 AND id = $6`
		_, err = d.SQLDB.Exec(query, c.Name, c.AWSRegion, urlEnc, keyEnc, userID, c.ID)
	}
	return c.ID, err
}

// DeleteUserConnection deletes a user connection by ID.
func (d *DB) DeleteUserConnection(userID string, connID string) error {
	query := `DELETE FROM user_connections WHERE user_id = $1 AND id = $2`
	_, err := d.SQLDB.Exec(query, userID, connID)
	return err
}

