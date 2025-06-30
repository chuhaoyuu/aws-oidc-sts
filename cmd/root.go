package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	TargetDir string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "aws-oidc-sts",
	Short: "A CLI tool for managing OIDC and STS resources in AWS",
	Long: `aws-oidc-sts is a command-line tool designed to simplify the management 
of OpenID Connect (OIDC) and Security Token Service (STS) resources in AWS. 
It provides commands to generate cryptographic assets, configure identity providers, 
and manage related resources.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	RootCmd.Root().CompletionOptions.DisableDefaultCmd = true
	RootCmd.PersistentFlags().StringVarP(&TargetDir, "output-dir", "o", pwd, "Target directory for the generated files")

}
