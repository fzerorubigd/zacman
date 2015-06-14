package main

import (
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vada-ir/semaphore"
)

func updateEntry(md *cobra.Command, args []string) {
	makeRoot()
	p, err := loadSnapShot("master")
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}

	if concurrent < 1 {
		concurrent = 1
	}
	s := semaphore.NewSemaphore(int(concurrent))
	wg := sync.WaitGroup{}
	wg.Add(len(p.Plugins))
	lock := sync.RWMutex{}
	for i := range p.Plugins {
		go func(i string) {
			s.Acquire(1)
			defer func() {
				s.Release(1)
				wg.Done()
			}()
			lock.RLock()
			pl := p.Plugins[i]
			lock.RUnlock()
			plUp, err := doBundle(pl.Repo, pl.SubPath, true)
			if err != nil {
				logrus.Warn(err)
			}

			lock.Lock()
			defer lock.Unlock()

			p.Plugins[i] = *plUp
		}(i)
	}

	wg.Wait()
}

func initUpdateCommand() *cobra.Command {
	var update = &cobra.Command{
		Use:   "update",
		Short: "update all plugins",
		Long:  `update all plugins at once`,
		Run:   updateEntry,
	}

	update.Flags().UintVarP(
		&concurrent,
		"concurrent",
		"c",
		5,
		"how many repo to update in parallel",
	)

	return update
}
