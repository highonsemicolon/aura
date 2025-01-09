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

func ComputeRolePrivilegesDFS(role string, privileges *RolePrivileges, visited map[string]bool) []string {
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
		effective = append(effective, ComputeRolePrivilegesDFS(inheritedRole, privileges, visited)...)
	}

	return effective
}
