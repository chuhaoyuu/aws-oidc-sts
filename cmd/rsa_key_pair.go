package cmd

import (
	"github.com/chuhaoyuu/aws-oidc-sts/pkg/providers"
	"github.com/spf13/cobra"
)

var rsaKeyPairCmd = &cobra.Command{
	Use:   "rsa-key-pair",
	Short: "Generate an RSA key pair and save it to the target directory",
	Long: `The rsa-key-pair command generates a new RSA key pair and saves the keys 
to the specified target directory. This command is useful for creating secure 
key pairs for cryptographic operations. 

Example usage:
  aws-oidc-sts create rsa-key-pair --target-dir /path/to/directory`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := providers.CreateRSAKeyPair(TargetDir); err != nil {
			cmd.PrintErrln("Failed to create RSA key pair:", err)
			cmd.SilenceUsage = true
		}
	},
}
