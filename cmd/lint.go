package cmd

import (
	"github.com/natebird/jig/runner"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Run swiftlint",
	RunE: func(cmd *cobra.Command, args []string) error {
		largs := append([]string{cfg.Lint.Executable}, cfg.Lint.Args...)
		return runner.Run(cfg.Dir, verboseFlag, largs...)
	},
}
