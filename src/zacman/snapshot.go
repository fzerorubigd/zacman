package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/vada-ir/semaphore"
)

var ignore = regexp.MustCompile(`^(http://|https://|git://|git@|/).*`)

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
		if err := cloneRepo(pl.Repo, pl.Path, false); err != nil {
			return err
		}
		pull = false
	} else if err != nil {
		return err
	}
	cmd := exec.Command("git", "cat-file", "-t", pl.Hash)
	cmd.Dir = pl.Path

	_, err := cmd.CombinedOutput()
	if err != nil {
		if pull {
			logrus.Warnf("failed with error %s , try to pull", err.Error())
			cmd = exec.Command("git", "pull")
			cmd.Dir = pl.Path
			out, err := cmd.Output()
			if err != nil {
				logrus.Warn(string(out))
				return err
			}

			return checkoutCommit(pl, false)
		}
		return err
	}

	cmd = exec.Command("git", "checkout", pl.Hash)
	cmd.Dir = pl.Path
	out, err := cmd.Output()
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
	final += "fpath=(" + strings.Join(fpath, " ") + " $fpath)\nautoload -U compinit\ncompinit -i\n" + res + "\n"

	return final
}

func findTargetRepo(short string) (string, string) {
	if !ignore.MatchString(short) {
		short = "https://github.com/" + short
	}

	dir := strings.Replace(short, ":", "-COLON-", -1)
	dir = strings.Replace(dir, "/", "-SLASH-", -1)
	dir = strings.Replace(dir, ".", "-DOT-", -1)
	dir = strings.Replace(dir, "@", "-AT-", -1)

	return short, root + "/repo/" + dir
}

func cloneRepo(repo, target string, update bool) error {
	is, err := exists(target)
	panicOnErr(err)

	if is {
		if update {
			logrus.Printf("Try to update repository %s", repo)
			cmd := exec.Command("git", "checkout", "master")
			cmd.Dir = target
			_, err = cmd.CombinedOutput()
			cmd = exec.Command("git", "pull")
			cmd.Dir = target
			_, err = cmd.CombinedOutput()
			return err
		}

		logrus.Info("already cloned, use --update to update it")
		return nil
	}

	logrus.Print("Clone new repository...")

	return exec.Command("git", "clone", repo, target).Wait()
}

func addMatch(files []os.FileInfo, pattern, target string) []string {
	var res []string
	re := regexp.MustCompile(pattern)
	for i := range files {
		if !files[i].IsDir() && re.MatchString(files[i].Name()) {
			res = append(res, target+"/"+files[i].Name())
		}
	}

	return res
}

func createIndexCache(target string, recursive bool) (string, string, []string, error) {
	var res []string
	files, err := ioutil.ReadDir(target)
	if err != nil {
		return "", "", nil, err
	}

	if theme {
		// this is a theme bundle, try to find zsh-theme
		res = append(res, addMatch(files, `.*\.zsh-theme$`, target)...)
	}

	if len(res) == 0 {
		// If there is an plugin.zsh file then simply source it
		res = append(res, addMatch(files, `.*\.plugin\.zsh$`, target)...)
	}

	if len(res) == 0 {
		res = append(res, addMatch(files, `.*\.zsh$`, target)...)
	}

	if len(res) == 0 {
		res = append(res, addMatch(files, `.*\.sh$`, target)...)
	}

	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	cmd.Dir = target
	hash, err := cmd.Output()
	if err != nil {
		return "", "", nil, err
	}
	return strings.Trim(string(hash), " \n\t"), target, res, nil
}

func doBundle(path, subpath string, update bool) (*Plugin, error) {
	git, target := findTargetRepo(path)

	if err := cloneRepo(git, target, update); err != nil {
		return nil, err
	}
	hash, fpath, res, err := createIndexCache(target+"/"+subpath, true)
	if err != nil {
		return nil, err
	}

	pl := Plugin{
		Repo:    git,
		SubPath: subpath,
		Path:    target,
		Hash:    hash,
		Files:   res,
		FPath:   fpath,
		Order:   order,
		Theme:   theme,
	}

	return &pl, nil
}
