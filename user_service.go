package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// userService provides user info for the currently "logged in" user.
// We don't really have "users" in this demo though and neither do we have login.
// So it's basically just a thing that provides a set of attributes from a live file
type userService struct {
	mu   sync.Mutex
	user user
}

const (
	userInfoPath = "./data/userinfo.yaml"
)

func newUserService(ctx context.Context) (*userService, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	u := &userService{}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("user reload returned:")

				if event.Has(fsnotify.Write) && filepath.Ext(event.Name) == ".yaml" {
					log.Println("user reload returned:", u.reload())
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watch error:", err)
			}
		}
	}()

	err = watcher.Add(userInfoPath)
	if err != nil {
		return nil, err
	}

	log.Println("watching for events on", userInfoPath)
	return u, u.reload()
}

func (u *userService) reload() error {
	f, err := os.Open(userInfoPath)
	if err != nil {
		return err
	}
	defer f.Close()

	y := yaml.NewDecoder(f)
	y.SetStrict(true)

	u2 := user{}
	err = y.Decode(&u2)
	if err != nil {
		return err
	}

	u.mu.Lock()
	u.user = u2
	u.mu.Unlock()

	return nil
}

func (u *userService) get() user {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.user
}
