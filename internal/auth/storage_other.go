//go:build !darwin && !windows

package auth

func newPlatformStore(profile string) (TokenStore, error) {
	return newFileStore(profile)
}
