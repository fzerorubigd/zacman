package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/Sirupsen/logrus"
)

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func panicOnErr(err error) {
	if err != nil {
		logrus.Panic(err)
	}
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func safeRemove(folder, root string) error {
	if strings.HasPrefix(folder, root) {
		return os.RemoveAll(folder)
	}

	return errors.New("folder is not inside the ZACMAN_ROOT")
}

func sha1Sum(i string) string {
	s := sha1.New()
	io.WriteString(s, i)
	return fmt.Sprintf("%x", s.Sum(nil))
}
