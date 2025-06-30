package cmd

import (
	"github.com/chuhaoyuu/aws-oidc-sts/cmd"
	"github.com/spf13/cobra"
)

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Create resources such as key pairs or identity provider configurations",
	Long: `The create command allows you to generate various resources, such as RSA key pairs 
or JSON Web Key Sets (JWKS), and save them to a specified target directory. This command 
is useful for initializing cryptographic assets or identity provider configurations.

Example usage:
  aws-oidc-sts aws rsa-key-pair --target-dir /path/to/directory
  aws-oidc-sts aws identity-provider --target-dir /path/to/directory`,
}

func init() {
	cmd.RootCmd.AddCommand(awsCmd)
	awsCmd.AddCommand(rsaKeyPairCmd)
	awsCmd.AddCommand(identityProviderCmd)

}
