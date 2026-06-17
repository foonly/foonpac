package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/foonly/foonpac/internal/logic"
	"github.com/foonly/foonpac/internal/pacman"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync installed packages with declarations",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := logic.GetState()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		toInstall := logic.Difference(state.Declared, state.Explicitly)
		toRemove := logic.Difference(state.Explicitly, state.Declared)

		if len(toInstall) == 0 && len(toRemove) == 0 {
			fmt.Println("Everything is in sync.")
			return
		}

		if len(toInstall) > 0 {
			fmt.Println("Packages to install:")
			for _, pkg := range toInstall {
				fmt.Printf("  + %s\n", pkg)
			}
		}

		if len(toRemove) > 0 {
			fmt.Println("Packages to remove:")
			for _, pkg := range toRemove {
				fmt.Printf("  - %s\n", pkg)
			}
		}

		fmt.Print("\nProceed with sync? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" || input == "yes" {
			if err := pacman.Sync(toInstall, toRemove); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Sync complete.")
		} else {
			fmt.Println("Sync cancelled.")
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
