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
