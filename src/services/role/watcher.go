package services

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type ReadOnlyMap interface {
	Load(key interface{}) (value interface{}, ok bool)
	Range(f func(key, value interface{}) bool)
}

type FileWatcher interface {
	GetEffectivePrivileges(role string) ([]string, bool)
	GetEffectivePrivilegesCache() ReadOnlyMap
	Load() FileWatcher
	Watch()
}

type fileWatcher struct {
	filename        string
	privilegeLoader *privilegeLoader
}

func NewFileWatcher(filename string) *fileWatcher {
	privilegeLoader := NewPrivilegeLoader(filename)
	return &fileWatcher{
		filename:        filename,
		privilegeLoader: privilegeLoader,
	}
}

func (f *fileWatcher) GetEffectivePrivileges(role string) ([]string, bool) {
	return f.privilegeLoader.GetEffectivePrivileges(role)
}

func (f *fileWatcher) GetEffectivePrivilegesCache() ReadOnlyMap {
	return f.privilegeLoader.GetEffectivePrivilegesCache()
}

func (f *fileWatcher) Load() FileWatcher {
	if err := f.privilegeLoader.LoadAndComputePrivileges(); err != nil {
		log.Fatal("error loading privileges:", err)
	}
	return f
}

func (f *fileWatcher) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("error creating watcher:", err)
	}
	defer watcher.Close()

	err = watcher.Add(f.filename)
	if err != nil {
		log.Fatal("error adding file to watcher:", err)
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
