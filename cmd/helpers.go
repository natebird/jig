package cmd

import (
	"github.com/natebird/jig/config"
)

// resolveTargets returns the targets to operate on based on the --target flag.
func resolveTargets(targets map[string]config.Target, target string) map[string]config.Target {
	if target == "all" {
		return targets
	}
	if t, ok := targets[target]; ok {
		return map[string]config.Target{target: t}
	}
	return nil
}

// xcodebuildArgs constructs the xcodebuild argument list for a given action and target.
func xcodebuildArgs(cfg *config.Config, target config.Target, action string) []string {
	args := []string{"xcodebuild", action}
	if cfg.Project.Path != "" {
		args = append(args, "-project", cfg.Project.Path)
	}
	args = append(args, "-scheme", target.Scheme)
	args = append(args, "-destination", target.Destination)
	return args
}
