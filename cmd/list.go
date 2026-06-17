package cmd

import (
	"fmt"
	"os"

	"github.com/foonly/foonpac/internal/logic"
	"github.com/foonly/foonpac/internal/pacman"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List packages based on management status",
}

var unmanagedCmd = &cobra.Command{
	Use:   "unmanaged",
	Short: "Lists installed packages not in the declaration files",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := logic.GetState()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		unmanaged := logic.Difference(state.Explicitly, state.Declared)
		for _, pkg := range unmanaged {
			fmt.Println(pkg)
		}
	},
}

var managedCmd = &cobra.Command{
	Use:   "managed",
	Short: "Lists installed packages in the declaration files",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := logic.GetState()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		managed := logic.Intersection(state.Explicitly, state.Declared)
		for _, pkg := range managed {
			fmt.Println(pkg)
		}
	},
}

var dependenciesCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Lists installed dependencies required by declared packages",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := logic.GetState()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		managedInstalled := logic.Intersection(state.Declared, state.Explicitly)
		allDepsOfManaged, err := pacman.GetTransitiveDependencies(managedInstalled)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, dep := range state.Deps {
			if allDepsOfManaged[dep] {
				fmt.Println(dep)
			}
		}
	},
}

var missingCmd = &cobra.Command{
	Use:   "missing",
	Short: "Lists declared packages that are not installed",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := logic.GetState()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		missing := logic.Difference(state.Declared, state.Explicitly)
		for _, pkg := range missing {
			fmt.Println(pkg)
		}
	},
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Lists installed packages not in the official repositories (AUR/manual)",
	Run: func(cmd *cobra.Command, args []string) {
		local, err := pacman.GetForeign()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		for _, pkg := range local {
			fmt.Println(pkg)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(unmanagedCmd)
	listCmd.AddCommand(managedCmd)
	listCmd.AddCommand(dependenciesCmd)
	listCmd.AddCommand(missingCmd)
	listCmd.AddCommand(localCmd)
}
