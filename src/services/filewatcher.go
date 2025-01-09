package services

import (
	"aura/src/utils"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	filename                 string
	effectivePrivilegesCache *sync.Map
}

func NewFileWatcher(filename string) *FileWatcher {
	return &FileWatcher{
		filename:                 filename,
		effectivePrivilegesCache: &sync.Map{},
	}
}

func (f *FileWatcher) Start() error {
	if err := f.loadPrivileges(); err != nil {
		return err
	}

	go func() {
		if err := f.watch(); err != nil {
			log.Fatal("error watching file:", err)
		}
	}()
	return nil
}

func (f *FileWatcher) loadPrivileges() error {

	privileges, err := utils.LoadPrivileges(f.filename)
	if err != nil {
		return err
	}

	f.computeEffectivePrivileges(privileges)
	return nil
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
				if err := f.loadPrivileges(); err != nil {
					log.Println("error reloading privileges:", err)
				} else {
					log.Println("privileges reloaded successfully")
				}
			}
		case err := <-watcher.Errors:
			log.Println(err)
		}
	}
}

func (f *FileWatcher) computeEffectivePrivileges(privileges *utils.RolePrivileges) {
	var newPrivileges sync.Map
	for role := range privileges.Roles {
		effective := computeRolePrivilegesDFS(role, privileges, make(map[string]bool))
		newPrivileges.Store(role, effective)
	}

	f.effectivePrivilegesCache.Clear()
	f.effectivePrivilegesCache = &newPrivileges
}

func computeRolePrivilegesDFS(role string, privileges *utils.RolePrivileges, visited map[string]bool) []string {
	if visited[role] {
		return nil
	}
	visited[role] = true

	roleData, exists := privileges.Roles[role]
	if !exists {
		return nil
	}

	effective := make([]string, 0)
	for _, privilege := range roleData.Privileges {
		effective = append(effective, privilege.Action)
	}

	for _, inheritedRole := range roleData.Inherits {
		effective = append(effective, computeRolePrivilegesDFS(inheritedRole, privileges, visited)...)
	}

	return effective
}
