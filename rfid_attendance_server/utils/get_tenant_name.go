package utils

import (
	"regexp"
	"strings"
)

// GetTenantName converts an organization name into a valid tenant schema name like "tenant_x"
func GetTenantName(orgName string) string {
	// Convert to lowercase
	lower := strings.ToLower(orgName)

	// Replace all non-alphanumeric characters with underscores
	reg, _ := regexp.Compile(`[^a-z0-9]+`)
	cleaned := reg.ReplaceAllString(lower, "_")

	// Trim underscores from start and end
	cleaned = strings.Trim(cleaned, "_")

	// Prefix with "tenant_"
	return "tenant_" + cleaned
}
