package cmd

import (
	"log/slog"

	rootCmd "github.com/chuhaoyuu/aws-oidc-sts/cmd"
	"github.com/chuhaoyuu/aws-oidc-sts/pkg/providers"
	"github.com/spf13/cobra"
)

var (
	bucketName string
	region     string
)

var identityProviderCmd = &cobra.Command{
	Use:   "identity-provider",
	Short: "Generate a JSON Web Key Set (JWKS) for an identity provider",
	Long: `The identity-provider command generates a JSON Web Key Set (JWKS) 
and saves it to the specified target directory. This is useful for setting up 
or configuring an identity provider that requires a JWKS for token signing 
and verification.

Example usage:
  aws-oidc-sts create identity-provider --target-dir /path/to/directory --bucket-name my-s3-bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := providers.CreateIdentityProvider(rootCmd.TargetDir, bucketName, region); err != nil {
			cmd.PrintErrln("Failed to create identity provider:", err)
			cmd.SilenceUsage = true
		} else {
			slog.Info("Identity provider created successfully.")
		}

	},
}

func init() {
	identityProviderCmd.Flags().StringVarP(&bucketName, "bucket-name", "b", "", "S3 bucket name to store the JWKS and openid-configuration (required)")
	identityProviderCmd.Flags().StringVarP(&region, "region", "r", "", "AWS region (required)")
	identityProviderCmd.MarkFlagRequired("bucket-name")
	identityProviderCmd.MarkFlagRequired("region")
}
