package bootstrap

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog"

	"github.com/medincident/medincident-backend/internal/config"
)

// OpenGarage builds an S3 client configured against a Garage endpoint.
// Garage is an S3-compatible object storage backend; the client uses
// static credentials and path-style addressing (required by Garage).
//
// No health check or warm-up is performed — the function returns
// immediately so that a temporarily unreachable Garage instance cannot
// stall the boot sequence (see Hard Rule 13).
func OpenGarage(cfg *config.GarageConfig, logger *zerolog.Logger) (client *s3.Client, cleanup func()) {
	client = s3.New(s3.Options{
		Region:       cfg.Region,
		BaseEndpoint: aws.String(cfg.Endpoint),
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		UsePathStyle: true,
	})

	logger.Info().
		Str("endpoint", cfg.Endpoint).
		Str("bucket", cfg.Bucket).
		Str("region", cfg.Region).
		Msg("garage s3 client configured")

	return client, func() {}
}
