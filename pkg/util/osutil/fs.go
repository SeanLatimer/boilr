package osutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileExists checks whether the given path exists and belongs to a file.
func FileExists(path string) (bool, error) {
	path = filepath.Clean(path)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	if info.IsDir() {
		return false, fmt.Errorf("%v: is a directory, expected file", path)
	}

	return true, nil
}

// DirExists checks whether the given path exists and belongs to a directory.
func DirExists(path string) (bool, error) {
	path = filepath.Clean(path)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	if !info.IsDir() {
		return false, fmt.Errorf("%v: is a file, expected directory", path)
	}

	return true, nil
}

// CreateDirs creates directories from the given directory path arguments.
func CreateDirs(dirPaths ...string) error {
	for _, path := range dirPaths {
		path = filepath.Clean(path)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

// CopyRecursively copies a given directory to the destination.
// Creates the directory if the destination doesn't exist.
func CopyRecursively(srcPath, dstPath string) error {
	srcPath = filepath.Clean(srcPath)
	dstPath = filepath.Clean(dstPath)
	if err := os.Mkdir(dstPath, 0755); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %q doesn't exist", filepath.Dir(dstPath))
		}

		if !os.IsExist(err) {
			return err
		}
	}

	return filepath.Walk(srcPath, func(fname string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcPath, fname)
		if err != nil {
			return err
		}

		mirrorPath := filepath.Join(dstPath, relPath)

		if info.IsDir() {
			if err := os.Mkdir(mirrorPath, info.Mode()); err != nil {
				if !os.IsExist(err) {
					return err
				}
			}
		} else {
			if err := CopyFile(fname, mirrorPath); err != nil {
				return err
			}
		}

		return nil
	})
}

// CopyFile copies a single file
func CopyFile(srcPath, dstPath string) error {
	fi, err := os.Lstat(srcPath)
	if err != nil {
		return err
	}

	srcf, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcf.Close()

	dstf, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, fi.Mode())
	if err != nil {
		return err
	}
	defer dstf.Close()

	if _, err := io.Copy(dstf, srcf); err != nil {
		return err
	}

	return nil
}