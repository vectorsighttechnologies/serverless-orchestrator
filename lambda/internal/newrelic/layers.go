// Package newrelic — layers.go
//
// Layer discovery from layers.newrelic-external.com + layer selection logic.
// Ported from: newrelic_lambda_cli/layers.py :: index(), layer_selection()
package newrelic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vectorsighttechnologies/serverless-orchestrator/lambda/internal/types"
)

// httpClient is reused across calls (connection pooling, keep-alive).
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// DiscoverLayers fetches available NR Lambda layers for the given region,
// runtime, and architecture from the NR layer registry.
//
// Ported from: layers.py :: index(region, runtime, architecture)
func DiscoverLayers(region, runtime, architecture string) ([]types.LayerInfo, error) {
	url := fmt.Sprintf(LayerRegistryURL, region, runtime)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch NR layers from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NR layer registry returned status %d for %s", resp.StatusCode, url)
	}

	var registry types.LayerRegistryResponse
	if err := json.NewDecoder(resp.Body).Decode(&registry); err != nil {
		return nil, fmt.Errorf("failed to decode NR layer registry response: %w", err)
	}

	// Filter by architecture — mirrors the Python CLI logic
	var compatible []types.LayerInfo
	for _, layer := range registry.Layers {
		archs := layer.LatestMatchingVersion.CompatibleArchitectures
		if len(archs) == 0 {
			archs = []string{"x86_64"} // default same as Python CLI
		}
		for _, a := range archs {
			if a == architecture {
				compatible = append(compatible, layer)
				break
			}
		}
	}

	return compatible, nil
}

// SelectLayer picks the best layer ARN from available layers.
//
// Ported from: layers.py :: layer_selection()
// Since we're an API (not interactive CLI), we always auto-select:
//   - If upgrading, prefer the layer matching the existing base ARN
//   - Otherwise, pick the first available layer
func SelectLayer(layers []types.LayerInfo, existingLayerARN string, upgrade bool) (string, error) {
	if len(layers) == 0 {
		return "", fmt.Errorf("no compatible NR Lambda layers found")
	}

	// On upgrade, try to match the existing layer's base ARN
	if upgrade && existingLayerARN != "" {
		baseARN := stripVersion(existingLayerARN)
		for _, layer := range layers {
			candidateARN := layer.LatestMatchingVersion.LayerVersionArn
			if stripVersion(candidateARN) == baseARN {
				return candidateARN, nil
			}
		}
	}

	// Default: first available layer
	return layers[0].LatestMatchingVersion.LayerVersionArn, nil
}

// stripVersion removes the ":123" version suffix from a layer ARN.
func stripVersion(arn string) string {
	idx := strings.LastIndex(arn, ":")
	if idx < 0 {
		return arn
	}
	return arn[:idx]
}

// GetARNPrefix returns the NR layer ARN prefix for a given region.
// Used to detect existing NR layers on a function.
// Ported from: utils.py :: get_arn_prefix()
func GetARNPrefix(region string) string {
	return fmt.Sprintf(LayerARNPrefix, region)
}

// IsNRLayer returns true if the layer ARN belongs to New Relic.
func IsNRLayer(layerARN, region string) bool {
	return strings.HasPrefix(layerARN, GetARNPrefix(region))
}
