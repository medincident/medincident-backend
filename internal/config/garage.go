package config

// GarageConfig is the S3-compatible object storage block used by
// command-server. Garage exposes an S3-compatible API; we connect via
// the standard AWS SDK with static credentials and path-style
// addressing.
type GarageConfig struct {
	Endpoint        string `yaml:"endpoint"          validate:"required,url"`
	AccessKeyID     string `yaml:"access_key_id"     validate:"required"`
	SecretAccessKey string `yaml:"secret_access_key" validate:"required"`
	Region          string `yaml:"region"            validate:"required"`
	Bucket          string `yaml:"bucket"            validate:"required"`
}
