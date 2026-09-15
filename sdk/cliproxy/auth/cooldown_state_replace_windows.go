//go:build windows

package auth

import "golang.org/x/sys/windows"

// os.Rename cannot replace an existing destination on Windows. The cooldown
// writer intentionally creates the temporary file beside the destination, so
// MoveFileEx preserves the same-volume atomic replacement contract.
func replaceCooldownStateFile(src, dst string) error {
	srcPtr, errSrc := windows.UTF16PtrFromString(src)
	if errSrc != nil {
		return errSrc
	}
	dstPtr, errDst := windows.UTF16PtrFromString(dst)
	if errDst != nil {
		return errDst
	}
	return windows.MoveFileEx(srcPtr, dstPtr, windows.MOVEFILE_REPLACE_EXISTING)
}
