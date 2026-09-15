//go:build !windows

package auth

import "os"

func replaceCooldownStateFile(src, dst string) error {
	return os.Rename(src, dst)
}
