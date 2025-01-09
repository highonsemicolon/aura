package utils

import (
	"os"

	"gopkg.in/yaml.v2"
)

type RolePrivileges struct {
	Roles map[string]Role `yaml:"roles"`
}

type Role struct {
	Inherits   []string `yaml:"inherits"`
	Privileges []Action `yaml:"privileges"`
}

type Action struct {
	Action string `yaml:"action"`
}

func LoadPrivileges(filename string) (*RolePrivileges, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var privileges RolePrivileges
	if err := yaml.Unmarshal(data, &privileges); err != nil {
		return nil, err
	}

	return &privileges, nil
}

func ComputeRolePrivilegesDFS(role string, privileges *RolePrivileges, visited map[string]bool) map[string]struct{} {
	if visited[role] {
		return nil
	}
	visited[role] = true

	roleData, exists := privileges.Roles[role]
	if !exists {
		return nil
	}

	effective := make(map[string]struct{})
	for _, privilege := range roleData.Privileges {
		effective[privilege.Action] = struct{}{}
	}

	for _, inheritedRole := range roleData.Inherits {
		inheritedPrivileges := ComputeRolePrivilegesDFS(inheritedRole, privileges, visited)
		for action := range inheritedPrivileges {
			effective[action] = struct{}{}
		}
	}

	return effective
}
