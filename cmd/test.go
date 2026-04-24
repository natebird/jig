package cmd

import (
	"fmt"

	"github.com/natebird/jig/runner"
	"github.com/spf13/cobra"
)

var (
	onlyTesting string
	testPlan    string
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test one or all targets",
	RunE: func(cmd *cobra.Command, args []string) error {
		targets := resolveTargets(cfg.Targets, targetFlag)
		if len(targets) == 0 {
			return fmt.Errorf("no targets found for --target=%q", targetFlag)
		}

		for name, target := range targets {
			fmt.Printf("==> Testing: %s (scheme=%s)\n", name, target.Scheme)
			xargs := xcodebuildArgs(cfg, target, "test")
			if onlyTesting != "" {
				xargs = append(xargs, "-only-testing", onlyTesting)
			}
			if testPlan != "" {
				xargs = append(xargs, "-testPlan", testPlan)
			}
			if err := runner.Run(cfg.Dir, verboseFlag, xargs...); err != nil {
				return fmt.Errorf("tests failed for target %q: %w", name, err)
			}
		}
		return nil
	},
}

func init() {
	testCmd.Flags().StringVar(&onlyTesting, "only-testing", "", "Limit testing to a specific test bundle/class/method")
	testCmd.Flags().StringVar(&testPlan, "test-plan", "", "Use a specific test plan")
}
