package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// TODO:
// 1. Create s3 bucket (Done in S3Service.Create())
// 2. Upload JWKSFileName to S3 bucket
// 3. Upload openid-configuration to S3 bucket
// 4. Create IAM role with trust policy to allow sts:AssumeRole with OIDC provider
// 5. Create IAM policy to allow sts:AssumeRole with OIDC provider

type AwsClient interface {
	// STS
	AwsClientIdentity() (*sts.GetCallerIdentityOutput, error)

	// s3
	CreateS3Bucket(bucketName, region string) error
}

type AwsServiceClient struct {
	stsClient *sts.Client
	s3Client  *s3.Client
}

func (c *AwsServiceClient) AwsClientIdentity() (*sts.GetCallerIdentityOutput, error) {

	// Use the AWS SDK to get the identity
	identity, err := c.stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return &sts.GetCallerIdentityOutput{}, fmt.Errorf("failed to get caller identity: %w", err)
	}

	return identity, nil

}

func (c *AwsServiceClient) CreateS3Bucket(bucketName, region string) error {
	// Create the bucket
	slog.Info("Creating S3 bucket", "BucketName", bucketName)
	_, err := c.s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
	}

	slog.Info("S3 Bucket created successfully", "BucketName", bucketName)

	return nil
}

// newClient initializes and returns an AWS SDK configuration object.
// It loads the default configuration with the specified AWS region and
// retrieves the AWS client identity for logging purposes.
//
// The function logs the AWS client identity details, including the account,
// ARN, region, and user ID.
//
// Returns:
//   - aws.Config: The AWS SDK configuration object.
//   - error: An error if the configuration loading or identity retrieval fails.
func newClient(region string) (aws.Config, error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return cfg, nil

}

// NewAwsFromConfig creates a new AWS service client configured for the specified region.
// It initializes the AWS client using the provided region and returns an instance of AwsServiceClient.
// If there is an error during the client creation, it returns an error.
//
// Parameters:
//   - region: The AWS region to configure the client for.
//
// Returns:
//   - Client: An interface representing the AWS service client.
//   - error: An error if the client creation fails.
func NewAwsFromConfig(region string) (AwsClient, error) {

	config, err := newClient(region)
	if err != nil {
		return &AwsServiceClient{}, fmt.Errorf("failed to create AWS client: %w", err)
	}

	return &AwsServiceClient{
		stsClient: sts.NewFromConfig(config),
		s3Client:  s3.NewFromConfig(config),
	}, nil

}
