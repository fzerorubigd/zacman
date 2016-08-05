package plugin

import (
	"assert"
	"common"
	"git"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"fmt"

	"errors"

	"strings"

	"time"

	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/expand"
)

const (
	definitionRepo = "https://github.com/fzerorubigd/zacman-repository.git"
)

var (
	ignore = regexp.MustCompile(`^(http://|https://|git://|git@|/).*`)
)

type Script struct {
	Src string
}

type Repository struct {
	Type     string
	Address  string
	Target   string
	Checkout string
}

type Config struct {
	Source      []string
	FPath       []string
	Interactive bool
}

type InstallData struct {
	Repo Repository
	Zsh  Config
}

// Plugin is used for a single plugin
type Definition struct {
	Author      string `toml:"author"`
	Name        string `toml:"-"`
	Version     string `toml:"version"`
	Description string
	Install     InstallData
}

// Load the plugin definition from an io.Reader
func Load(reader io.Reader) (*Definition, error) {
	var p Definition
	_, err := toml.DecodeReader(reader, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// LoadFromFile try to load definition from a file
func LoadFromFile(f string) (*Definition, error) {
	fl, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	return Load(fl)
}

// Save the plugin into a writer
func Save(writer io.Writer, p *Definition) error {
	e := toml.NewEncoder(writer)
	return e.Encode(p)
}

// SaveToFile try to save definition into a file
func SaveToFile(f string, p *Definition) error {
	fl, err := os.Create(f)
	if err != nil {
		return err
	}

	return Save(fl, p)
}

func UpdateIndex(root string) error {
	idx := filepath.Join(root, "index")
	is, err := common.Exists(idx)
	assert.Nil(err)

	if is {
		return git.Pull(idx)
	}

	return git.Clone(definitionRepo, idx)
}

func InstallPlugin(root string, fl string, update bool) (err error) {
	target := filepath.Join(root, "index", fl+".toml")
	d, err := LoadFromFile(target)
	if err != nil {
		return err
	}
	d.Name = fl
	var clone string
	defer func() {
		if err == nil {
			// there is no error, so save them inside the cache
			var hash string
			hash, err = git.CommitHash(clone)
			if err != nil {
				return
			}
			d.Install.Repo.Checkout = hash
			snaps := filepath.Join(root, "snapshots", "active", fl+".toml")
			err = SaveToFile(snaps, d)
		}
	}()

	if d.Install.Repo.Type != "git" {
		return fmt.Errorf("only git repository is accepted, provided : %s", d.Install.Repo.Type)
	}
	clone, _ = expand.Path(d.Install.Repo.Target)
	if clone == "" {
		clone = filepath.Join(root, fl)
	}
	d.Install.Repo.Target = clone
	is, err := common.Exists(clone)
	assert.Nil(err)
	if !is {
		err = git.Clone(d.Install.Repo.Address, clone)
		if err == nil && d.Install.Repo.Checkout != "" {
			return git.Checkout(clone, d.Install.Repo.Checkout)
		}
		return err
	}
	if update {
		err = git.Fetch(clone)
		if err == nil && d.Install.Repo.Checkout != "" {
			return git.Checkout(clone, d.Install.Repo.Checkout)
		}
		return err
	}

	return errors.New("Already cloned. use --update to update")
}

func RemovePlugin(root, fl string) error {
	// TODO : better error
	target := filepath.Join(root, "snapshots", "active", fl+".toml")
	_, err := LoadFromFile(target)
	if err != nil {
		return err
	}

	// ok the file is there. I don't think deleting the source file is a good idea :)
	// so leave them.
	return os.Remove(target)
}

func MkZshRC(root string) (err error) {
	defer func() {
		if err != nil {
			logrus.Warn(err)
		} else {
			logrus.Info("Create zshrc file done")
		}
	}()
	link := filepath.Join(root, "snapshots", "active")
	target, err := os.Readlink(link)
	if err != nil {
		return err
	}
	defs := []*Definition{}
	filepath.Walk(
		target,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if strings.ToLower(filepath.Ext(path)) == ".toml" {
				d, err := LoadFromFile(path)
				// Can not load it?
				if err != nil {
					logrus.Warnf("can not load %s reason %s", path, err.Error())
					return nil
				}
				defs = append(defs, d)
			}
			return nil
		})

	res := buildLoadScript(defs)
	zshrc := filepath.Join(target, "zshrc")
	return ioutil.WriteFile(zshrc, []byte(res), 0644)
}

func buildLoadScript(p []*Definition) string {
	final := fmt.Sprintf(
		`# Auto generated file, DO NOT EDIT!
# build with zacman
# %s
`,
		time.Now().Format(time.RFC3339))
	res := ""
	var fpath []string
	//TODO handle order
	for i := range p {
		for _, f := range p[i].Install.Zsh.FPath {
			fpath = append(fpath, filepath.Join(p[i].Install.Repo.Target, f))
		}
		res += fmt.Sprintf("# %s %s\n", p[i].Name, p[i].Version)
		for _, f := range p[i].Install.Zsh.Source {
			if p[i].Install.Zsh.Interactive {
				res += "[[ $- = *i* ]] && source " + filepath.Join(p[i].Install.Repo.Target, f) + "\n"
			} else {
				res += "source " + filepath.Join(p[i].Install.Repo.Target, f) + "\n"
			}
			fpath = append(fpath)
		}
	}

	final += "fpath=(" + strings.Join(fpath, " ") + " $fpath)\nautoload -U compinit\ncompinit -i\n" + res + "\n"

	return final
}
