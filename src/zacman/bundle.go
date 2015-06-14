package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codeskyblue/go-sh"
	"github.com/spf13/cobra"
)

var ignore = regexp.MustCompile(`^(http://|https://|git://|git@|/).*`)

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
			logrus.Print("Try to update repository...")
			session := sh.NewSession()
			_, err = session.SetDir(target).Command("git", "checkout", "master").Output()
			session = sh.NewSession()
			_, err = session.SetDir(target).Command("git", "pull").Output()
			return err
		}

		logrus.Info("already cloned, use --update to update it")
		return nil
	}

	logrus.Print("Clone new repository...")

	return sh.Command("git", "clone", repo, target).Run()
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

	cmd := sh.NewSession()
	hash, err := cmd.SetDir(target).Command("git", "rev-parse", "--verify", "HEAD").Output()
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

func bundleEntry(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		cmd.Usage()
		os.Exit(1)
	}

	subpath := ""
	if len(args) == 2 {
		subpath = args[1]
	}

	makeRoot()
	p, err := loadSnapShot("master")
	if err != nil {
		logrus.Warn(err)
		p = newSnapShot()
	}
	pl, err := doBundle(args[0], subpath, update)
	if err != nil {
		logrus.Fatal(err)
	}
	name := strings.Trim(pl.Repo+"/"+pl.SubPath, "/")
	p.Plugins[name] = *pl

	err = saveSnapShot("master", p)
	if err != nil {
		logrus.Fatal(err)
	}
}

func initBundleCommand() *cobra.Command {
	var bundle = &cobra.Command{
		Use:   "bundle",
		Short: "try to bundle a plugin",
		Long: `try to bundle a plugin, the full git repo,
or "user/repo" for github user is accepted usage :
zacman bundle target
`,
		Run: bundleEntry,
	}

	bundle.Flags().BoolVarP(
		&update,
		"update",
		"u",
		false,
		"update the repo if already there?",
	)

	bundle.Flags().IntVarP(
		&order,
		"order",
		"o",
		0,
		"order of loading the bundle, greater value means load sooner",
	)

	bundle.Flags().BoolVarP(&theme, "theme", "t", false, "is this a theme?")

	return bundle
}
