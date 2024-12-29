package main

import "os"

const backupExt = ".bak"

type backupSet []string

func (b *backupSet) createBackup(src string) (string, error) {
	dst := src + backupExt
	if err := atomicCopy(dst, src); err != nil {
		return "", err
	}

	*b = append(*b, src)

	return dst, nil
}

func (b *backupSet) restoreBackups() error {
	var restoreErr error

	var bs backupSet

	for _, src := range *b {
		dst := src + backupExt
		if err := atomicCopy(dst, src); err != nil {
			bs = append(bs, src)
			restoreErr = err
		}
	}

	*b = bs

	return restoreErr
}

func (b *backupSet) removeBackups() error {
	var removeErr error

	for _, src := range *b {
		dst := src + backupExt
		if err := os.Remove(dst); err != nil {
			removeErr = err
		}
	}

	return removeErr
}

func checkBackupExists(src string) (bool, error) {
	return checkExists(src + backupExt)
}
