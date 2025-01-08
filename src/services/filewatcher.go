package services

import (
	"aura/src/utils"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	filename   string
	privileges *utils.RolePrivileges
	mu         sync.RWMutex
}

func NewFileWatcher(filename string) *FileWatcher {
	return &FileWatcher{
		filename: filename,
	}
}

func (f *FileWatcher) Start() error {
	if err := f.loadPrivileges(); err != nil {
		return err
	}
	if err := f.Watch(); err != nil {
		return err
	}
	return nil
}

func (f *FileWatcher) loadPrivileges() error {

	privileges, err := utils.LoadPrivileges(f.filename)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.privileges = privileges
	log.Println("privileges loaded successfully")
	return nil
}

func (f *FileWatcher) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(f.filename)
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				newPrivileges, err := utils.LoadPrivileges(f.filename)
				if err != nil {
					log.Println("error reloading privileges:", err)
				} else {
					f.privileges = newPrivileges
					log.Println("privileges reloaded successfully")

				}
			}
		case err := <-watcher.Errors:
			log.Println(err)
		}
	}
}

func (f *FileWatcher) GetPrivileges() *utils.RolePrivileges {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.privileges
}
