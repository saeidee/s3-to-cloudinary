package support

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cloudinary/cloudinary-go"
	cloudinaryConfig "github.com/cloudinary/cloudinary-go/config"
)

type Factory struct {
	config *Config
}

func NewFactory(config *Config) *Factory {
	return &Factory{config: config}
}

func (f *Factory) NewCloudinaryInstance() (*cloudinary.Cloudinary, error) {
	if !f.config.HasUploadPrefix() {
		cld, err := cloudinary.NewFromParams(
			f.config.Secrets.Cloudinary.CloudName,
			f.config.Secrets.Cloudinary.ApiKey,
			f.config.Secrets.Cloudinary.ApiSecret,
		)

		return cld, err
	}

	cld, err := cloudinary.NewFromConfiguration(
		cloudinaryConfig.Configuration{
			Cloud: cloudinaryConfig.Cloud{
				CloudName: f.config.Secrets.Cloudinary.CloudName,
				APIKey:    f.config.Secrets.Cloudinary.ApiKey,
				APISecret: f.config.Secrets.Cloudinary.ApiSecret,
			},
			API: cloudinaryConfig.API{
				UploadPrefix: f.config.Secrets.Cloudinary.UploadPrefix,
				Timeout:      60,
				ChunkSize:    20000000,
			},
		},
	)

	return cld, err
}

func (f *Factory) NewAWSSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
		Credentials: credentials.NewStaticCredentials(
			f.config.Secrets.AWS.AccessKeyID,
			f.config.Secrets.AWS.SecretAccessKey,
			"",
		),
	})
}
