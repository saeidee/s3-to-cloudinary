package support

import (
   "context"
   "fmt"
   "github.com/aws/aws-sdk-go/aws/session"
   "github.com/aws/aws-sdk-go/service/s3"
   "github.com/cloudinary/cloudinary-go"
   "github.com/cloudinary/cloudinary-go/api/uploader"
   "path/filepath"
   "strings"
   "sync"
)

type migrator struct {
   session    *session.Session
   cloudinary *cloudinary.Cloudinary
   wg         *sync.WaitGroup
}

func NewMigrator(sess *session.Session, cld *cloudinary.Cloudinary) *migrator {
   var wq sync.WaitGroup

   return &migrator{session: sess, cloudinary: cld, wg: &wq}
}

func (m *migrator) Migrate(config *Config, logChannel chan interface{}) {
   s3Service := s3.New(m.session)
   var token *string = nil

   for _, bucket := range config.Buckets {
      for {
         resp, err := m.listObjects(s3Service, &bucket, token)
         if err != nil {
            logChannel <- err

            continue
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

func (m *migrator) listObjects(
   s3Service *s3.S3,
   bucket *string,
   continuationToken *string,
) (*s3.ListObjectsV2Output, error) {
   resp, err := s3Service.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: bucket, ContinuationToken: continuationToken})
   if err != nil {
      return nil, err
   }

   return resp, nil
}

func (m *migrator) upload(item *s3.Object, bucket string, wg *sync.WaitGroup, logChannel chan<- interface{}) {
   defer wg.Done()

   basename := *item.Key
   file := fmt.Sprintf("s3://%s/%s", bucket, *item.Key)
   fileName := strings.TrimSuffix(basename, filepath.Ext(basename))
   publicID := fmt.Sprintf("%s/%s", bucket, fileName)

   uploadRes, err := m.cloudinary.Upload.Upload(context.Background(), file, uploader.UploadParams{PublicID: publicID})
   if err != nil {
      logChannel <- err

      return
   }

   logChannel <- uploadRes.Error
}