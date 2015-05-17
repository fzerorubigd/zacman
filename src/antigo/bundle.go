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

var ignore = regexp.MustCompile(`$(http://|https://|git://|git@|/).*`)

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

func cloneRepo(repo, target string) error {
	is, err := exists(target)
	panicOnErr(err)

	if is {
		if update {
			logrus.Print("Try to update repository...")
			session := sh.NewSession()
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

	// If there is an plugin.zsh file then simply source it
	res = append(addMatch(files, `.*\.plugin\.zsh$`, target))

	if len(res) == 0 {
		res = append(addMatch(files, `.*\.zsh`, target))
	}

	if len(res) == 0 {
		res = append(addMatch(files, `.*\.sh`, target))
	}

	if len(res) == 0 && recursive {
		for i := range files {
			if files[i].IsDir() && files[i].Name() == "src" {
				return createIndexCache(target+"/src", false)
			}
		}
	}
	cmd := sh.NewSession()
	hash, err := cmd.SetDir(target).Command("git", "rev-parse", "--verify", "HEAD").Output()
	if err != nil {
		return "", "", nil, err
	}
	return strings.Trim(string(hash), " \n\t"), target, res, nil
}

func bundleEntry(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}

	makeRoot()
	p, err := loadSnapShot("master")
	if err != nil {
		logrus.Warn(err)
		p = newSnapShot()
	}

	for i := range args {
		git, target := findTargetRepo(args[i])
		if err := cloneRepo(git, target); err != nil {
			logrus.Warn(err)
			continue
		}
		hash, fpath, res, err := createIndexCache(target, true)
		if err != nil {
			logrus.Warn(err)
			continue
		}

		pl := Plugin{
			Repo:  git,
			Path:  target,
			Hash:  hash,
			Files: res,
			FPath: fpath,
			Order: order,
		}

		p.Plugins[git] = pl
	}

	err = saveSnapShot("master", p)
	if err != nil {
		logrus.Warn(err)
	}
}

func initBundleCommand() *cobra.Command {
	var bundle = &cobra.Command{
		Use:   "bundle",
		Short: "try to bundle a plugin",
		Long: `try to bundle a plugin, the full git repo,
or "user/repo" for github user is accepted usage :
antigo bundle target
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

	return bundle
}
