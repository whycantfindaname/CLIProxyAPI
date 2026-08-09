package management

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRuntimeManagementSecret(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		t.Setenv(managementPasswordEnv, "")
		t.Setenv(managementPasswordFileEnv, "")

		secret, errSecret := LoadRuntimeManagementSecret()
		if errSecret != nil {
			t.Fatalf("LoadRuntimeManagementSecret() error = %v", errSecret)
		}
		if secret != "" {
			t.Fatalf("LoadRuntimeManagementSecret() = %q, want empty", secret)
		}
	})

	t.Run("direct value takes precedence", func(t *testing.T) {
		t.Setenv(managementPasswordEnv, " direct-secret \n")
		t.Setenv(managementPasswordFileEnv, filepath.Join(t.TempDir(), "missing"))

		secret, errSecret := LoadRuntimeManagementSecret()
		if errSecret != nil {
			t.Fatalf("LoadRuntimeManagementSecret() error = %v", errSecret)
		}
		if secret != "direct-secret" {
			t.Fatalf("LoadRuntimeManagementSecret() = %q, want %q", secret, "direct-secret")
		}
	})

	t.Run("file value", func(t *testing.T) {
		t.Setenv(managementPasswordEnv, "")
		path := filepath.Join(t.TempDir(), "management-password")
		if errWrite := os.WriteFile(path, []byte(" file-secret \n"), 0o600); errWrite != nil {
			t.Fatalf("write management password: %v", errWrite)
		}
		t.Setenv(managementPasswordFileEnv, path)

		secret, errSecret := LoadRuntimeManagementSecret()
		if errSecret != nil {
			t.Fatalf("LoadRuntimeManagementSecret() error = %v", errSecret)
		}
		if secret != "file-secret" {
			t.Fatalf("LoadRuntimeManagementSecret() = %q, want %q", secret, "file-secret")
		}
	})

	t.Run("missing file", func(t *testing.T) {
		t.Setenv(managementPasswordEnv, "")
		path := filepath.Join(t.TempDir(), "missing")
		t.Setenv(managementPasswordFileEnv, path)

		_, errSecret := LoadRuntimeManagementSecret()
		if errSecret == nil || !strings.Contains(errSecret.Error(), path) {
			t.Fatalf("LoadRuntimeManagementSecret() error = %v, want path context", errSecret)
		}
	})

	t.Run("empty file", func(t *testing.T) {
		t.Setenv(managementPasswordEnv, "")
		path := filepath.Join(t.TempDir(), "management-password")
		if errWrite := os.WriteFile(path, []byte(" \n"), 0o600); errWrite != nil {
			t.Fatalf("write management password: %v", errWrite)
		}
		t.Setenv(managementPasswordFileEnv, path)

		_, errSecret := LoadRuntimeManagementSecret()
		if errSecret == nil || !strings.Contains(errSecret.Error(), "is empty") {
			t.Fatalf("LoadRuntimeManagementSecret() error = %v, want empty-file error", errSecret)
		}
	})
}
