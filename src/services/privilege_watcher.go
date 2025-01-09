package services

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	filename        string
	privilegeLoader *PrivilegeLoader
}

func NewFileWatcher(filename string) *FileWatcher {
	privilegeLoader := NewPrivilegeLoader(filename)
	return &FileWatcher{
		filename:        filename,
		privilegeLoader: privilegeLoader,
	}
}

func (f *FileWatcher) GetEffectivePrivileges(role string) ([]string, bool) {
	return f.privilegeLoader.GetEffectivePrivileges(role)
}

func (f *FileWatcher) GetEffectivePrivilegess() *sync.Map {
	return f.privilegeLoader.effectivePrivilegesCache
}

func (f *FileWatcher) Start() {
	if err := f.privilegeLoader.LoadAndComputePrivileges(); err != nil {
		log.Fatal("error loading privileges:", err)
	}

	if err := f.watch(); err != nil {
		log.Fatal("error watching file:", err)
	}
}

func (f *FileWatcher) watch() error {
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
				if err := f.privilegeLoader.LoadAndComputePrivileges(); err != nil {
					log.Println("error reloading privileges:", err)
				} else {
					log.Println("privileges reloaded successfully")
				}
			}
		case err := <-watcher.Errors:
			log.Println("watcher error:", err)
		}
	}
}
