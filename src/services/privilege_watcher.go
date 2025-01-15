package services

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type ReadOnlyMap interface {
	Load(key interface{}) (value interface{}, ok bool)
	Range(f func(key, value interface{}) bool)
}

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

func (f *FileWatcher) GetEffectivePrivilegesCache() ReadOnlyMap {
	return f.privilegeLoader.GetEffectivePrivilegesCache()
}

func (f *FileWatcher) Load() *FileWatcher {
	if err := f.privilegeLoader.LoadAndComputePrivileges(); err != nil {
		log.Fatal("error loading privileges:", err)
	}
	return f
}

func (f *FileWatcher) Watch() {
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
