package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindConfig(t *testing.T) {
	// Create a temp dir tree: /tmp/a/b/c
	root := t.TempDir()
	sub := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}

	// Place jig.toml at root
	configPath := filepath.Join(root, "jig.toml")
	content := `
[project]
name = "TestProject"

[targets.ios]
scheme = "iOS"
destination = "platform=iOS Simulator,name=iPhone 17"

[targets.macos]
scheme = "macOS"
destination = "platform=macOS"

[lint]
executable = "swiftlint"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// findConfig should find it from a subdirectory
	found, err := findConfig(sub)
	if err != nil {
		t.Fatalf("expected to find config, got error: %v", err)
	}
	if found != configPath {
		t.Errorf("expected %q, got %q", configPath, found)
	}
}

func TestFindConfigNotFound(t *testing.T) {
	root := t.TempDir()
	_, err := findConfig(root)
	if err == nil {
		t.Fatal("expected error when jig.toml is absent")
	}
}

func TestLoadConfig(t *testing.T) {
	root := t.TempDir()
	content := `
[project]
name = "Quotebook"
path = "Quotebook.xcodeproj"

[targets.ios]
scheme = "iOS"
destination = "platform=iOS Simulator,name=iPhone 17"

[targets.macos]
scheme = "macOS"
destination = "platform=macOS"

[lint]
executable = "swiftlint"
args = []
`
	configPath := filepath.Join(root, "jig.toml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to root so Load() finds the config
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(root)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Project.Name != "Quotebook" {
		t.Errorf("expected project name 'Quotebook', got %q", cfg.Project.Name)
	}
	if cfg.Project.Path != "Quotebook.xcodeproj" {
		t.Errorf("expected project path 'Quotebook.xcodeproj', got %q", cfg.Project.Path)
	}
	if len(cfg.Targets) != 2 {
		t.Errorf("expected 2 targets, got %d", len(cfg.Targets))
	}
	ios, ok := cfg.Targets["ios"]
	if !ok {
		t.Fatal("missing ios target")
	}
	if ios.Scheme != "iOS" {
		t.Errorf("expected scheme 'iOS', got %q", ios.Scheme)
	}
	if cfg.Lint.Executable != "swiftlint" {
		t.Errorf("expected lint executable 'swiftlint', got %q", cfg.Lint.Executable)
	}
	realRoot, _ := filepath.EvalSymlinks(root)
	if cfg.Dir != realRoot {
		t.Errorf("expected Dir %q, got %q", realRoot, cfg.Dir)
	}
}

func TestLintExecutableDefault(t *testing.T) {
	root := t.TempDir()
	content := `
[project]
name = "Test"

[targets.ios]
scheme = "iOS"
destination = "platform=iOS Simulator,name=iPhone 17"
`
	configPath := filepath.Join(root, "jig.toml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(root)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Lint.Executable != "swiftlint" {
		t.Errorf("expected default lint executable 'swiftlint', got %q", cfg.Lint.Executable)
	}
}
