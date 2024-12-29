package main

import (
	"bufio"
	"io"
	"os"
)

func atomicWrite(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := w.Write(data); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	return nil
}

func atomicCopy(dst, src string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	if _, err := io.Copy(df, sf); err != nil {
		return err
	}

	if err := df.Sync(); err != nil {
		return err
	}

	return nil
}

func checkExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
