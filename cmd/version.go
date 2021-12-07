package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version    = "Unknown"
	BuildStamp = ""
	GitHash    = ""
	GoVersion  = ""
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "the version number of replay uploader",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:             %s\n", Version)
		fmt.Printf("Git Commit Hash:     %s\n", GitHash)
		fmt.Printf("UTC Build Time :     %s\n", BuildStamp)
		fmt.Printf("Go Version:          %s\n", GoVersion)
	},
}
