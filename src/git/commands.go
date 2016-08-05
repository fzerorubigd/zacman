package git

import (
	"assert"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

var (
	gitPath = "/usr/bin/git"
)

// Clone a repository in target
func Clone(repo, target string) (err error) {
	logrus.Infof("Clone %s into %s", repo, target)
	defer func() {
		if err != nil {
			logrus.Warnf("Clone faile with reason %s", err.Error())
		} else {
			logrus.Info("Clone done.")
		}
	}()
	cmd := exec.Command(gitPath, "clone", repo, target)
	if err = cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

// Pull Changes from repository
func Pull(root string) (err error) {
	logrus.Infof("Pull into %s", root)
	defer func() {
		if err != nil {
			logrus.Warnf("Pull faile with reason %s", err.Error())
		} else {
			logrus.Info("Pull done.")
		}
	}()
	cmd := exec.Command(gitPath, "pull")
	cmd.Dir = root
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func Checkout(root, target string) (err error) {
	logrus.Infof("Checkout into %s", target)
	defer func() {
		if err != nil {
			logrus.Warnf("Checkout faile with reason %s", err.Error())
		} else {
			logrus.Info("Checkout done.")
		}
	}()
	cmd := exec.Command(gitPath, "checkout", target)
	cmd.Dir = root
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func CommitHash(root string) (hash string, err error) {
	logrus.Infof("Get commit hash for %s", root)
	defer func() {
		if err != nil {
			logrus.Warnf("Show failed with reason %s", err.Error())
		} else {
			logrus.Info("Show done.")
		}
	}()
	cmd := exec.Command(gitPath, "show", "--pretty=%H")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err == nil {
		hash = strings.Trim(string(out), " \n\t")
	}

	return hash, err
}

// SetGitPath to change the default path
func SetGitPath(new string) error {
	cmd := exec.Command(new, "version")
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	gitPath = new

	return nil
}

func init() {
	tmp := os.Getenv("ZACMAN_GIT_PATH")
	if tmp != "" {
		assert.Nil(SetGitPath(tmp), tmp)
	}
}
