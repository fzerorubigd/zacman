package main

import (
	"assert"
	"common"
	"flagx"
	"os"

	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/expand"
)

func main() {
	err := flagx.Parse()
	if err != nil {
		logrus.Fatal(err)
	}
}

func getRootDir() string {
	root := os.Getenv("ZACMAN_ROOT")
	if root == "" {
		var err error
		root, err = expand.HomeDir()
		assert.Nil(err)
		root += "/.zacman"
	}

	return root
}

func makeRoot() string {
	root := getRootDir()
	if b, err := common.Exists(root); !b || err != nil {
		assert.Nil(os.Mkdir(root, 0750))
	}

	snapshots := filepath.Join(root, "snapshots", "default")
	if b, err := common.Exists(snapshots); !b || err != nil {
		assert.Nil(os.MkdirAll(snapshots, 0750))
		active := filepath.Join(root, "snapshots", "active")
		assert.Nil(os.Symlink(snapshots, active))
	}

	return root
}
