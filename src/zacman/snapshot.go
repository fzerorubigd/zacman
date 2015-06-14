package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codeskyblue/go-sh"
	"github.com/vada-ir/semaphore"
)

// Plugin is used for a single plugin
type Plugin struct {
	Repo    string   `json:"repository"`
	SubPath string   `json:"sub-path"`
	Path    string   `json:"path"`
	Hash    string   `json:"commit-hash"`
	Files   []string `json:"source-files"`
	FPath   string   `json:"fpath"`
	Order   int      `json:"order"`
	Theme   bool     `json:"theme"`
}

// Plugins is used for bunch of plugins
type Plugins struct {
	Version int       `json:"version"`
	Date    time.Time `json:"date"`
	// Maps are not sortable in golang.
	Plugins map[string]Plugin `json:"plugins"`
}

// SortablePlugins is a type for handling sort on plugins in puild
type SortablePlugins []Plugin

func (s SortablePlugins) Len() int           { return len(s) }
func (s SortablePlugins) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortablePlugins) Less(i, j int) bool { return s[i].Order > s[j].Order }

// Sort return a sort list from plugins
func (p Plugins) Sort() SortablePlugins {
	s := make(SortablePlugins, len(p.Plugins))
	i := 0
	for _, plugin := range p.Plugins {
		s[i] = plugin
		i++
	}

	sort.Sort(s)
	return s
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

func listSnapShots() []string {
	var res []string

	path := root + "/snapshots/"

	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		str := []byte(f.Name())
		if len(str) > 41 {
			res = append(res, string(str[:len(str)-41]))
		}
	}

	return res
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

func compileSnapshot(p *Plugins, concurrent uint) {
	wg := sync.WaitGroup{}
	wg.Add(len(p.Plugins))
	// I don't want more than
	s := semaphore.NewSemaphore(int(concurrent))
	for i := range p.Plugins {
		go func(i string) {
			s.Acquire(1)
			defer func() {
				s.Release(1)
				wg.Done()
			}()
			err := checkoutCommit(p.Plugins[i], true)
			if err != nil {
				logrus.Warnf("can not restore state for %s reason is %s", p.Plugins[i].Repo, err.Error())
			}
		}(i)
	}

	wg.Wait()
}

func buildLoadScipt(p *Plugins) string {
	final := fmt.Sprintf(
		`# Auto generated file, DO NOT EDIT!
# build with zacman
# %s

`,
		time.Now().Format(time.RFC3339))
	res := ""
	var fpath []string
	//TODO handle order
	s := p.Sort()
	for i := range s {
		fpath = append(fpath, s[i].FPath)
		for _, f := range s[i].Files {
			res += "source " + f + "\n"
		}
	}
	final += "autoload -U compinit\ncompinit -i\n" + res + "\nfpath=(" + strings.Join(fpath, " ") + " $fpath)\n"

	return final
}
