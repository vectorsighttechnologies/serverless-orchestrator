package datadog

import (
	"fmt"
	"strings"
)

// GetAccountID returns the publishing account ID based on the region.
func GetAccountID(region string) string {
	if strings.HasPrefix(region, "us-gov-") {
		return GovCloudAccountID
	}
	return LayerAccountID
}

// BuildExtensionLayerARN builds the Datadog Extension Layer ARN for a region and CPU architecture.
func BuildExtensionLayerARN(region, architecture string, version int) string {
	account := GetAccountID(region)
	layerName := "Datadog-Extension"
	if strings.EqualFold(architecture, "arm64") {
		layerName = "Datadog-Extension-ARM"
	}
	return fmt.Sprintf("arn:aws:lambda:%s:%s:layer:%s:%d", region, account, layerName, version)
}

// BuildLibraryLayerARN builds the Datadog Library Layer ARN.
func BuildLibraryLayerARN(region, libraryName string, version int) string {
	account := GetAccountID(region)
	return fmt.Sprintf("arn:aws:lambda:%s:%s:layer:%s:%d", region, account, libraryName, version)
}

// IsDatadogLayer checks if a layer ARN belongs to Datadog.
func IsDatadogLayer(layerARN, region string) bool {
	account := GetAccountID(region)
	prefix := fmt.Sprintf("arn:aws:lambda:%s:%s:layer:", region, account)
	return strings.HasPrefix(layerARN, prefix)
}

// stripVersion removes the ":123" version suffix from a layer ARN.
func stripVersion(arn string) string {
	idx := strings.LastIndex(arn, ":")
	if idx < 0 {
		return arn
	}
	return arn[:idx]
}
