package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	TargetDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
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
	err := rootCmd.Execute()
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

	rootCmd.Root().CompletionOptions.DisableDefaultCmd = false
	rootCmd.AddCommand(createCmd)
	rootCmd.PersistentFlags().StringVarP(&TargetDir, "output-dir", "o", pwd, "Target directory for the generated files")

}
