package cmd

import (
	"fmt"
	"os"

	"github.com/CodeSeedx/scodecounter/formatter"
	"github.com/CodeSeedx/scodecounter/gitinfo"
	"github.com/CodeSeedx/scodecounter/scanner"
	"github.com/spf13/cobra"
)

var (
	outputJSON bool
	noGit      bool
	version    = "dev"
)

// SetVersion sets the version string (called from main)
func SetVersion(v string) {
	version = v
}

var rootCmd = &cobra.Command{
	Use:   "scodecounter [path]",
	Short: "Count lines of code, files, and git commits in a project",
	Long:  "scodecounter - A fast code statistics tool that counts lines, files, and git commits across all programming languages.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runScan,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "Output as JSON")
	rootCmd.Flags().BoolVar(&noGit, "no-git", false, "Skip git statistics")
	rootCmd.Version = version
}

func runScan(cmd *cobra.Command, args []string) error {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	// Scan
	result, err := scanner.Scan(root)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Git info
	var gitStat *gitinfo.CommitInfo
	if !noGit {
		gitStat, _ = gitinfo.GetCommitInfo(root)
	}

	// Output
	if outputJSON {
		return formatter.FormatJSON(os.Stdout, result, gitStat)
	}
	formatter.FormatTerminal(os.Stdout, result, gitStat)
	return nil
}
