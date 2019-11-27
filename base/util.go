package base

import (
	"bufio"
	"golang.org/x/sys/unix"
	"os"
	"path"
	"syscall"
)

func mkdir(file_name string) error {
	return os.MkdirAll(path.Dir(file_name), 0775)
}

func Mkfifo(named_pipe string) error {
	fn := func(named_pipe string) error {
		err := syscall.Mkfifo(named_pipe, 0600)
		if err == syscall.EEXIST {
			return nil
		} else if err != nil {
			return err
		}
		file, err := os.OpenFile(named_pipe, os.O_RDWR, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = unix.FcntlInt(file.Fd(), unix.F_SETPIPE_SZ, 262144)
		return err
	}
	err := fn(named_pipe)
	if err == syscall.ENOENT {
		if nil == mkdir(named_pipe) {
			return fn(named_pipe)
		}
	}

	return err
}

func Symlink(old_path, new_path string) error {
	fn := func(old_path, new_path string) error {
		err := syscall.Symlink(old_path, new_path)
		if err == syscall.EEXIST {
			return nil
		}
		return err
	}

	err := fn(old_path, new_path)
	if err == syscall.ENOENT {
		if nil == mkdir(new_path) {
			return fn(old_path, new_path)
		}
	}

	return err
}

func Exp2(n uint64) uint64 {
	result := uint64(1)
	for i := 0; i < int(n); i++ {
		result *= 2
	}

	return result
}

func ReadLineFile(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	values := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		values = append(values, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return values, nil
}
