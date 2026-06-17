package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var AppVersion = "dev"

type Config struct {
	Hosts map[string][]string `mapstructure:"hosts"`
}

func LoadConfig() (*Config, error) {
	if viper.ConfigFileUsed() == "" {
		return nil, fmt.Errorf("no configuration file found")
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}
	return &cfg, nil
}

func GetDeclaredPackages(hostname string, configDir string) ([]string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	groups, ok := cfg.Hosts[hostname]
	if !ok {
		return nil, fmt.Errorf("no groups defined for host: %s", hostname)
	}

	packageMap := make(map[string]bool)
	for _, group := range groups {
		groupPath := filepath.Join(configDir, "groups", group+".txt")
		packages, err := readGroupFile(groupPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Warning: group file not found: %s\n", groupPath)
				continue
			}
			return nil, fmt.Errorf("failed to read group file %s: %w", groupPath, err)
		}
		for _, pkg := range packages {
			packageMap[pkg] = true
		}
	}

	var allPackages []string
	for pkg := range packageMap {
		allPackages = append(allPackages, pkg)
	}

	return allPackages, nil
}

func readGroupFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packages []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		packages = append(packages, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return packages, nil
}

func GetConfigDir() string {
	if cfgFile := viper.ConfigFileUsed(); cfgFile != "" {
		return filepath.Dir(cfgFile)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "foonpac")
}
