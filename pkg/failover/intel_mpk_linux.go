//go:build linux

// SPDX-License-Identifier: Apache-2.0

package failover

import "syscall"

type SiliconVault struct {
	Key int
}

func AllocateSecureVault() (*SiliconVault, error) {
	pkey, _, err := syscall.Syscall(
		329,
		0,
		0,
		0,
	)

	if err != 0 {
		return nil, err
	}

	return &SiliconVault{
		Key: int(pkey),
	}, nil
}
