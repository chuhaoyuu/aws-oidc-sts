package aws

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Service represents a service for interacting with an S3 bucket.
// It contains an S3 client and the name of the bucket to operate on.
type S3Service struct {
	Client     *s3.Client
	BucketName string
	Region     string
}

// Create creates an S3 bucket using the AWS SDK for Go v2.
// It logs the process of bucket creation and returns an error if the operation fails.
//
// The bucket name is specified by the S3Service's BucketName field, and the AWS region
// is determined by the providers.AWSRegion constant.
//
// Returns:
//   - nil if the bucket is created successfully.
//   - an error if the bucket creation fails, including the bucket name and the underlying error.
func (s *S3Service) Create() error {
	// Create the bucket
	slog.Info("Creating S3 bucket", "BucketName", s.BucketName)
	_, err := s.Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(s.BucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(s.Region),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", s.BucketName, err)
	}

	slog.Info("S3 Bucket created successfully", "BucketName", s.BucketName)

	return nil
}
