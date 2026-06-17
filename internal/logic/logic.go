package logic

import (
	"os"

	"github.com/foonly/foonpac/internal/config"
	"github.com/foonly/foonpac/internal/pacman"
)

type State struct {
	Declared   []string
	Explicitly []string
	Deps       []string
}

func GetState() (*State, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	configDir := config.GetConfigDir()
	declared, err := config.GetDeclaredPackages(hostname, configDir)
	if err != nil {
		return nil, err
	}

	explicitly, err := pacman.GetExplicitlyInstalledNative()
	if err != nil {
		return nil, err
	}

	deps, err := pacman.GetInstalledDependencies()
	if err != nil {
		return nil, err
	}

	return &State{
		Declared:   declared,
		Explicitly: explicitly,
		Deps:       deps,
	}, nil
}

func Difference(a, b []string) []string {
	mb := make(map[string]bool, len(b))
	for _, x := range b {
		mb[x] = true
	}
	var diff []string
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}

func Intersection(a, b []string) []string {
	mb := make(map[string]bool, len(b))
	for _, x := range b {
		mb[x] = true
	}
	var res []string
	for _, x := range a {
		if _, ok := mb[x]; ok {
			res = append(res, x)
		}
	}
	return res
}
