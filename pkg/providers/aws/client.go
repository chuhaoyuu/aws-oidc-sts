package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// TODO:
// 1. Create s3 bucket (Done in S3Service.Create())
// 2. Upload JWKSFileName to S3 bucket
// 3. Upload openid-configuration to S3 bucket
// 4. Create IAM role with trust policy to allow sts:AssumeRole with OIDC provider
// 5. Create IAM policy to allow sts:AssumeRole with OIDC provider

// AwsClient defines an interface for interacting with AWS services.
// It provides methods for obtaining caller identity information via STS
// and creating S3 buckets.
//
// Methods:
//   - AwsClientIdentity: Retrieves the AWS account and user identity information
//     using the STS GetCallerIdentity API.
//   - CreateS3Bucket: Creates a new S3 bucket using the S3 CreateBucket API.
type AwsClient interface {
	// STS
	AwsClientIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error)

	// s3
	CreateS3Bucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
}

// AwsServiceClient represents a client for interacting with AWS services.
// It encapsulates clients for specific AWS services, such as STS and S3,
// allowing for streamlined access and operations.
//
// Fields:
//   - stsClient: A client for interacting with AWS Security Token Service (STS),
//     used for managing temporary credentials and identity federation.
//   - s3Client: A client for interacting with Amazon Simple Storage Service (S3),
//     used for object storage and related operations.
type AwsServiceClient struct {
	stsClient *sts.Client
	s3Client  *s3.Client
}

// AwsClientIdentity retrieves the AWS caller identity information using the provided STS client.
// It accepts a GetCallerIdentityInput object and returns a GetCallerIdentityOutput object along with an error, if any.
// This function uses the context.TODO() for the request context.
//
// Parameters:
//   - input: A pointer to sts.GetCallerIdentityInput containing the parameters for the GetCallerIdentity API call.
//
// Returns:
//   - *sts.GetCallerIdentityOutput: The output containing the caller identity information.
//   - error: An error object if the operation fails.
func (c *AwsServiceClient) AwsClientIdentity(input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return c.stsClient.GetCallerIdentity(context.TODO(), input)

}

// CreateS3Bucket creates a new S3 bucket using the provided input parameters.
// It utilizes the AWS SDK's S3 client to perform the operation.
// 
// Parameters:
//   - input: A pointer to s3.CreateBucketInput containing the configuration for the bucket creation.
//
// Returns:
//   - A pointer to s3.CreateBucketOutput containing the details of the created bucket.
//   - An error if the bucket creation fails.
//
// Note:
//   - Ensure that the input parameters meet the requirements for bucket creation in AWS.
//   - The operation is performed using a context with no timeout (context.TODO).
func (c *AwsServiceClient) CreateS3Bucket(input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return c.s3Client.CreateBucket(context.TODO(), input)

}

// newClient initializes a new AWS SDK client configuration for the specified region.
// It loads the default configuration using the AWS SDK for Go v2 and applies the provided region.
// 
// Parameters:
//   - region: A string specifying the AWS region to configure the client for.
//
// Returns:
//   - aws.Config: The loaded AWS SDK configuration.
//   - error: An error if the configuration could not be loaded.
//
// Example usage:
//   cfg, err := newClient("us-west-2")
//   if err != nil {
//       log.Fatalf("Failed to create AWS client: %v", err)
//   }
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

// NewAwsFromConfig creates a new AWS client configured for the specified region.
// It initializes the AWS STS and S3 clients using the provided region configuration.
//
// Parameters:
//   - region: The AWS region to configure the client for.
//
// Returns:
//   - AwsClient: An interface representing the AWS client with STS and S3 capabilities.
//   - error: An error if the client creation fails, wrapped with additional context.
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
