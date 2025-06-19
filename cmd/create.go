package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources such as key pairs or identity provider configurations",
	Long: `The create command allows you to generate various resources, such as RSA key pairs 
or JSON Web Key Sets (JWKS), and save them to a specified target directory. This command 
is useful for initializing cryptographic assets or identity provider configurations.

Example usage:
  create rsa-key-pair --target-dir /path/to/directory
  create identity-provider --target-dir /path/to/directory`,
}

func init() {
	createCmd.AddCommand(rsaKeyPairCmd)
	createCmd.AddCommand(identityProviderCmd)

}
