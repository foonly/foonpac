package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/foonly/foonpac/internal/config"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group [groupname] [packagename]",
	Short: "View and manage group files",
	Long: `View and manage group files.
- foonpac group: List active groups and package counts.
- foonpac group <groupname>: List packages in a group.
- foonpac group clean: Sort all group files alphabetically.
- foonpac group <groupname> <packagename>: Add a package to a group.`,
	Run: func(cmd *cobra.Command, args []string) {
		configDir := config.GetConfigDir()
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		hostname, err := os.Hostname()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(args) == 0 {
			listGroups(cfg, hostname, configDir)
		} else if len(args) == 1 {
			if args[0] == "clean" {
				cleanGroups(configDir)
			} else {
				listGroupPackages(args[0], configDir)
			}
		} else if len(args) == 2 {
			addPackageToGroup(args[0], args[1], configDir)
		} else {
			cmd.Help()
		}
	},
}

func listGroups(cfg *config.Config, hostname string, configDir string) {
	groups, ok := cfg.Hosts[hostname]
	if !ok {
		fmt.Printf("No groups defined for host: %s\n", hostname)
		return
	}

	for _, group := range groups {
		groupPath := filepath.Join(configDir, "groups", group+".txt")
		packages, err := config.ReadGroupFile(groupPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s: 0 packages (file not found)\n", group)
			} else {
				fmt.Printf("%s: error reading group: %v\n", group, err)
			}
			continue
		}
		fmt.Printf("%s: %d packages\n", group, len(packages))
	}
}

func listGroupPackages(groupName string, configDir string) {
	groupPath := filepath.Join(configDir, "groups", groupName+".txt")
	packages, err := config.ReadGroupFile(groupPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: group file not found: %s\n", groupPath)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error reading group file: %v\n", err)
		os.Exit(1)
	}

	for _, pkg := range packages {
		fmt.Println(pkg)
	}
}

func cleanGroups(configDir string) {
	groupsDir := filepath.Join(configDir, "groups")
	files, err := os.ReadDir(groupsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading groups directory: %v\n", err)
		os.Exit(1)
	}

	skipped := 0
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".txt" {
			continue
		}

		path := filepath.Join(groupsDir, file.Name())
		packages, err := config.ReadGroupFile(path)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file.Name(), err)
			continue
		}

		// Check if the file would actually change
		originalContent, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file.Name(), err)
			continue
		}

		sort.Strings(packages)
		var newContent strings.Builder
		for _, pkg := range packages {
			newContent.WriteString(pkg)
			newContent.WriteString("\n")
		}

		if string(originalContent) == newContent.String() {
			skipped++
			continue
		}

		err = config.WriteGroupFile(path, packages)
		if err != nil {
			fmt.Printf("Error writing %s: %v\n", file.Name(), err)
			continue
		}
		fmt.Printf("Cleaned %s\n", file.Name())
	}

	if skipped > 0 {
		fmt.Printf("%d files skipped\n", skipped)
	}
}

func addPackageToGroup(groupName string, packageName string, configDir string) {
	groupPath := filepath.Join(configDir, "groups", groupName+".txt")

	// Create groups directory if it doesn't exist
	groupsDir := filepath.Dir(groupPath)
	if _, err := os.Stat(groupsDir); os.IsNotExist(err) {
		err = os.MkdirAll(groupsDir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating groups directory: %v\n", err)
			os.Exit(1)
		}
	}

	packages, err := config.ReadGroupFile(groupPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error reading group file: %v\n", err)
		os.Exit(1)
	}

	// Check if package already exists
	if slices.Contains(packages, packageName) {
		fmt.Printf("Package %s already in group %s\n", packageName, groupName)
		return
	}

	packages = append(packages, packageName)
	err = config.WriteGroupFile(groupPath, packages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing group file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Added %s to group %s\n", packageName, groupName)
}

func init() {
	rootCmd.AddCommand(groupCmd)
}
