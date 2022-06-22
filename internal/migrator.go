package internal

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"path/filepath"
	"strings"
	"sync"
)

type migrator struct {
	svg        *s3.S3
	cloudinary *cloudinary.Cloudinary
	wg         *sync.WaitGroup
}

func NewMigrator(svg *s3.S3, cld *cloudinary.Cloudinary) *migrator {
	var wq sync.WaitGroup

	return &migrator{svg: svg, cloudinary: cld, wg: &wq}
}

func (m *migrator) Migrate(config *Config, logChannel chan Log) {
	var token *string = nil

	for _, bucket := range config.Buckets {
		for {
			resp, err := m.listObjects(&s3.ListObjectsV2Input{
				Bucket:            &bucket,
				ContinuationToken: token,
				MaxKeys:           &config.MaxKeys,
			})
			if err != nil {
				logChannel <- Log{
					Bucket: bucket,
					Item:   nil,
					Error:  err.Error(),
				}

				break
			}

			token = resp.NextContinuationToken

			for _, item := range resp.Contents {
				m.wg.Add(1)
				go m.upload(item, bucket, m.wg, logChannel)
			}

			if token == nil {
				break
			}
		}
	}

	m.wg.Wait()
	close(logChannel)
}

func (m *migrator) listObjects(params *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	resp, err := m.svg.ListObjectsV2(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *migrator) upload(item *s3.Object, bucket string, wg *sync.WaitGroup, logChannel chan<- Log) {
	defer wg.Done()

	basename := *item.Key
	file := fmt.Sprintf("s3://%s/%s", bucket, *item.Key)
	fileName := strings.TrimSuffix(basename, filepath.Ext(basename))
	publicID := fmt.Sprintf("%s/%s", bucket, fileName)

	uploadRes, err := m.cloudinary.Upload.Upload(context.Background(), file, uploader.UploadParams{PublicID: publicID})
	if err != nil {
		logChannel <- Log{
			Bucket: bucket,
			Item:   item,
			Error:  err.Error(),
		}

		return
	}

	logChannel <- Log{
		Bucket: bucket,
		Item:   item,
		Error:  uploadRes.Error.Message,
	}
}
