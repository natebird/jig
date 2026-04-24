package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Project struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
}

type Target struct {
	Scheme      string `toml:"scheme"`
	Destination string `toml:"destination"`
}

type Lint struct {
	Executable string   `toml:"executable"`
	Args       []string `toml:"args"`
}

type Config struct {
	Dir     string
	Project Project            `toml:"project"`
	Targets map[string]Target  `toml:"targets"`
	Lint    Lint               `toml:"lint"`
}

func Load() (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	configPath, err := findConfig(dir)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, err
	}

	cfg.Dir = filepath.Dir(configPath)

	if cfg.Lint.Executable == "" {
		cfg.Lint.Executable = "swiftlint"
	}

	return &cfg, nil
}

func findConfig(dir string) (string, error) {
	for {
		candidate := filepath.Join(dir, "jig.toml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("jig.toml not found in current directory or any parent directory")
}
