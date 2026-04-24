package cmd

import (
	"fmt"

	"github.com/natebird/jig/runner"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build one or all targets",
	RunE: func(cmd *cobra.Command, args []string) error {
		targets := resolveTargets(cfg.Targets, targetFlag)
		if len(targets) == 0 {
			return fmt.Errorf("no targets found for --target=%q", targetFlag)
		}

		for name, target := range targets {
			fmt.Printf("==> Building target: %s\n", name)
			xargs := xcodebuildArgs(cfg, target, "build")
			if err := runner.Run(cfg.Dir, xargs...); err != nil {
				return fmt.Errorf("build failed for target %q: %w", name, err)
			}
		}
		return nil
	},
}
