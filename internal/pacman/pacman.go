package pacman

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetExplicitlyInstalledNative() ([]string, error) {
	return runPacman("-Qenq")
}

func GetForeign() ([]string, error) {
	return runPacman("-Qmq")
}

func GetInstalledDependencies() ([]string, error) {
	return runPacman("-Qdq")
}

func runPacman(args ...string) ([]string, error) {
	cmd := exec.Command("pacman", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var packages []string
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		packages = append(packages, scanner.Text())
	}
	return packages, scanner.Err()
}

func Sync(toInstall []string, toRemove []string) error {
	if len(toInstall) > 0 {
		fmt.Printf("Installing: %s\n", strings.Join(toInstall, " "))
		if err := runSudoPacman("-S", toInstall...); err != nil {
			return err
		}
	}

	if len(toRemove) > 0 {
		fmt.Printf("Removing: %s\n", strings.Join(toRemove, " "))
		if err := runSudoPacman("-Rs", toRemove...); err != nil {
			return err
		}
	}

	return nil
}

func runSudoPacman(op string, packages ...string) error {
	args := append([]string{"pacman", op}, packages...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetTransitiveDependencies returns all dependencies (transitive) for the given packages.
// Using pactree if available, or parsing pacman -Qi if not.
// For now, let's try pacman -Qi as it's more likely to be there.
func GetTransitiveDependencies(packages []string) (map[string]bool, error) {
	if len(packages) == 0 {
		return make(map[string]bool), nil
	}

	// This is a bit complex to do perfectly without pactree.
	// For now, let's use pactree -u if available, otherwise return error or simplified version.
	_, err := exec.LookPath("pactree")
	if err == nil {
		return getDepsWithPactree(packages)
	}

	return nil, fmt.Errorf("pactree is required for dependency listing")
}

func getDepsWithPactree(packages []string) (map[string]bool, error) {
	deps := make(map[string]bool)
	for _, pkg := range packages {
		cmd := exec.Command("pactree", "-u", pkg)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			// Some packages might not be installed, which pactree might complain about.
			// Or they might not exist.
			continue
		}
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			dep := scanner.Text()
			if dep != pkg {
				deps[dep] = true
			}
		}
	}
	return deps, nil
}
