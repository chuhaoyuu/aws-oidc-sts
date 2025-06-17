package cmd

import (
	"github.com/chuhaoyuu/aws-oidc-sts/pkg/providers"
	"github.com/spf13/cobra"
)

var identityProviderCmd = &cobra.Command{
	Use:   "identity-provider",
	Short: "Generate a JSON Web Key Set (JWKS) for an identity provider",
	Long: `The identity-provider command generates a JSON Web Key Set (JWKS) 
and saves it to the specified target directory. This is useful for setting up 
or configuring an identity provider that requires a JWKS for token signing 
and verification.

Example usage:
  aws-oidc-sts create identity-provider --target-dir /path/to/directory`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := providers.CreateJSONWebKeySet(TargetDir); err != nil {
			cmd.PrintErrln("Error creating JSON Web Key Set:", err)
		}
	},
	DisableFlagParsing: true,
}
