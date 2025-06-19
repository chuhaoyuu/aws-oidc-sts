package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// TODO:
// 1. Create s3 bucket (Done in S3Service.Create())
// 2. Upload JWKSFileName to S3 bucket
// 3. Upload openid-configuration to S3 bucket
// 4. Create IAM role with trust policy to allow sts:AssumeRole with OIDC provider
// 5. Create IAM policy to allow sts:AssumeRole with OIDC provider

// AwsService defines an interface for interacting with AWS services.
// It provides a method to create resources or perform operations
// that may result in an error.
type AwsService interface {
	Create() error
}

// AwsClient initializes and returns an AWS SDK configuration object.
// It loads the default configuration with the specified AWS region and
// retrieves the AWS client identity for logging purposes.
//
// The function logs the AWS client identity details, including the account,
// ARN, region, and user ID.
//
// Returns:
//   - aws.Config: The AWS SDK configuration object.
//   - error: An error if the configuration loading or identity retrieval fails.
func AwsClient(region string) (aws.Config, error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load SDK config: %w", err)
	}

	identity, err := clientIdentity(cfg)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to get AWS client identity: %w", err)
	}

	slog.Info("AWS Client Identity",
		"Account", aws.ToString(identity.Account),
		"Arn", aws.ToString(identity.Arn),
		"Region", cfg.Region,
		"UserId", aws.ToString(identity.UserId),
	)

	return cfg, nil
}

// clientIdentity retrieves the AWS caller identity using the provided AWS configuration.
// It utilizes the AWS SDK's STS (Security Token Service) client to fetch the caller identity.
//
// Parameters:
//   - cfg: An aws.Config object containing the AWS configuration.
//
// Returns:
//   - *sts.GetCallerIdentityOutput: The output containing the caller identity details.
//   - error: An error if the operation fails, otherwise nil.
func clientIdentity(cfg aws.Config) (*sts.GetCallerIdentityOutput, error) {

	// Use the AWS SDK to get the identity
	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return &sts.GetCallerIdentityOutput{}, fmt.Errorf("failed to get caller identity: %w", err)
	}

	return identity, nil

}

// Create initializes and creates an AWS resource using the provided AwsService.
// It returns an error if the service is nil or if the creation process fails.
//
// Parameters:
//   - service: An implementation of the AwsService interface that defines the
//     resource creation logic.
//
// Returns:
//   - error: An error indicating why the creation failed, or nil if the operation
//     was successful.
func Create(service AwsService) error {
	if service == nil {
		return fmt.Errorf("resource is nil")
	}
	err := service.Create()
	if err != nil {
		return fmt.Errorf("failed to create AWS resource: %v", err)
	}
	return nil
}
