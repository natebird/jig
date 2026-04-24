package cmd

import (
	"fmt"
	"os"

	"github.com/natebird/jig/config"
	"github.com/spf13/cobra"
)

var (
	cfg         *config.Config
	targetFlag  string
	verboseFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "jig",
	Short: "Apple dev CLI for xcodebuild and swiftlint",
	Long:  "jig simplifies running xcodebuild and swiftlint commands using a jig.toml config file.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&targetFlag, "target", "all", "Target to operate on (e.g. ios, macos, all)")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Stream full command output instead of summary")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return nil
	}

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(lintCmd)
}
