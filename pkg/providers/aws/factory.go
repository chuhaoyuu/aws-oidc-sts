package aws

// Builder constructs and returns an instance of an AWS service based on the provided service type.
// It takes an AwsService interface as input and performs a type switch to determine the specific
// implementation of the service. Depending on the type, it initializes and returns the appropriate
// service instance.
//
// Parameters:
//   - serviceType: An implementation of the AwsService interface representing the desired AWS service.
//
// Returns:
//   - An initialized instance of the specific AWS service (e.g., *S3Service, *AWSCloudFront) if the
//     type matches, or nil if the service type is not recognized.
func Builder(serviceType AwsService) AwsService {
	switch service := serviceType.(type) {
	case *S3Service:
		return &S3Service{
			Client:     service.Client,
			BucketName: service.BucketName,
		}
	case *AWSCloudFront:
		return &AWSCloudFront{
			Name: service.Name,
		}
	default:
		return nil
	}
}
