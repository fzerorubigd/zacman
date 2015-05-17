package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codeskyblue/go-sh"
)

type Plugin struct {
	Repo  string   `json:"repository"`
	Path  string   `json:"path"`
	Hash  string   `json:"commit-hash"`
	Files []string `json:"source-files"`
	FPath string   `json:"fpath"`
	Order int      `json:"order"`
}

type Plugins struct {
	Version int               `json:"version"`
	Date    time.Time         `json:"date"`
	Plugins map[string]Plugin `json:"plugins"`
}

func saveSnapShot(name string, pl *Plugins) error {
	path := root + "/snapshots/" + name

	pl.Version = version
	pl.Date = time.Now()
	j, err := json.Marshal(pl)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, j, 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadSnapShot(name string) (*Plugins, error) {
	path := root + "/snapshots/" + name
	j, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := &Plugins{}
	err = json.Unmarshal(j, p)
	if err != nil {
		return nil, err
	}

	if p.Version > version {
		return nil, fmt.Errorf("the maximum supported version is %d", version)
	}

	return p, nil
}

func newSnapShot() *Plugins {
	return &Plugins{Plugins: make(map[string]Plugin), Version: version, Date: time.Now()}
}

func checkoutCommit(pl Plugin, pull bool) error {
	if b, err := exists(pl.Path); !b && err == nil {
		if err := cloneRepo(pl.Repo, pl.Path); err != nil {
			return err
		}
		pull = false
	} else if err != nil {
		return err
	}
	sess := sh.NewSession()
	_, err := sess.SetDir(pl.Path).Command("git", "cat-file", "-t", pl.Hash).Output()
	if err != nil {
		if pull {
			logrus.Warnf("failed with error %s , try to pull", err.Error())
			sess = sh.NewSession()
			out, err := sess.SetDir(pl.Path).Command("git", "pull").Output()
			if err != nil {
				logrus.Warn(string(out))
				return err
			}

			return checkoutCommit(pl, false)
		}
		return err
	}

	sess = sh.NewSession()
	out, err := sess.SetDir(pl.Path).Command("git", "checkout", pl.Hash).Output()
	if err != nil {
		logrus.Warn(string(out))
		return err
	}

	return nil
}

func compileSnapshot(p *Plugins) {
	for i := range p.Plugins {
		err := checkoutCommit(p.Plugins[i], true)
		if err != nil {
			logrus.Warnf("can not restore state for %s reason is %s", p.Plugins[i].Repo, err.Error())
			continue
		}
	}
}

func buildLoadScipt(p *Plugins) string {
	res := "#build with antigo\n"
	var fpath []string
	//TODO handle order
	for i := range p.Plugins {
		fpath = append(fpath, p.Plugins[i].FPath)
		for _, f := range p.Plugins[i].Files {
			res += "source " + f + "\n"
		}
	}

	res += "fpath=(" + strings.Join(fpath, " ") + " $fpath)"

	return res
}
