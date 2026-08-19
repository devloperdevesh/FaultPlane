//go:build windows

// SPDX-License-Identifier: Apache-2.0

package failover

type SiliconVault struct {
	Key int
}

func AllocateSecureVault() (*SiliconVault, error) {
	return &SiliconVault{
		Key: -1,
	}, nil
}
