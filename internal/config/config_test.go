package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestReadGroupFile(t *testing.T) {
	// Create a temporary group file
	tmpDir, err := os.MkdirTemp("", "foonpac-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	groupPath := filepath.Join(tmpDir, "test-group.txt")
	content := `
# This is a comment
package1
  package2

# Another comment
package3
`
	if err := os.WriteFile(groupPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	want := []string{"package1", "package2", "package3"}
	got, err := readGroupFile(groupPath)
	if err != nil {
		t.Errorf("readGroupFile() error = %v", err)
		return
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("readGroupFile() = %v, want %v", got, want)
	}
}

func TestGetDeclaredPackages_MissingGroup(t *testing.T) {
	// Setup viper for test
	viper.Reset()
	tmpDir, _ := os.MkdirTemp("", "foonpac-config-test")
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte(`
hosts:
  test-host:
    - existing
    - missing
`), 0644)

	groupsDir := filepath.Join(tmpDir, "groups")
	os.Mkdir(groupsDir, 0755)
	os.WriteFile(filepath.Join(groupsDir, "existing.txt"), []byte("pkg1\npkg2"), 0644)

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	packages, err := GetDeclaredPackages("test-host", tmpDir)
	if err != nil {
		t.Fatalf("GetDeclaredPackages() error = %v", err)
	}

	// Should only contain packages from 'existing' group
	if len(packages) != 2 {
		t.Errorf("Expected 2 packages, got %v: %v", len(packages), packages)
	}

	// Check content (order might vary)
	found1, found2 := false, false
	for _, p := range packages {
		if p == "pkg1" {
			found1 = true
		}
		if p == "pkg2" {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Errorf("Expected pkg1 and pkg2, got %v", packages)
	}
}
