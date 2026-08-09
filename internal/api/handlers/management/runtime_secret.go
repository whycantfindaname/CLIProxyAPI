package management

import (
	"fmt"
	"os"
	"strings"
)

const (
	managementPasswordEnv     = "MANAGEMENT_PASSWORD"
	managementPasswordFileEnv = "MANAGEMENT_PASSWORD_FILE"
)

// LoadRuntimeManagementSecret resolves the runtime management password from the
// direct environment value or an owner-managed file path.
func LoadRuntimeManagementSecret() (string, error) {
	if secret := strings.TrimSpace(os.Getenv(managementPasswordEnv)); secret != "" {
		return secret, nil
	}

	path := strings.TrimSpace(os.Getenv(managementPasswordFileEnv))
	if path == "" {
		return "", nil
	}
	contents, errRead := os.ReadFile(path)
	if errRead != nil {
		return "", fmt.Errorf("read management password file %q: %w", path, errRead)
	}
	secret := strings.TrimSpace(string(contents))
	if secret == "" {
		return "", fmt.Errorf("management password file %q is empty", path)
	}
	return secret, nil
}
