package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cloudinary/cloudinary-go"
	"github.com/saeidee/trek/support"
	"github.com/thatisuday/commando"
	"log"
	"time"
)

func main() {
	commando.
		SetExecutableName("trek").
		SetVersion("0.0.1").
		SetDescription("This tool is for migrating your s3 media to cloudinary.")

	commando.
		Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			for k, v := range args {
				log.Printf("arg -> %v: %v(%T)\n", k, v.Value, v.Value)
			}

			for k, v := range flags {
				log.Printf("flag -> %v: %v(%T)\n", k, v.Value, v.Value)
			}
		})

	commando.
		Register("start").
		SetShortDescription("starting migration for your media").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			configParser := support.NewConfigParser()
			config, err := configParser.Parse("config.yml")
			if err != nil {
				log.Fatalf("Unable to read the file with error: %v", err)
			}

			cld, err := cloudinary.NewFromParams(
				config.Secrets.Cloudinary.CloudName,
				config.Secrets.Cloudinary.ApiKey,
				config.Secrets.Cloudinary.ApiSecret,
			)
			if err != nil {
				log.Fatalf("Unable to connect to Cloudinary, error: %v\n", err)
			}

			sess, err := session.NewSession(&aws.Config{
				Region: aws.String("eu-west-1"),
				Credentials: credentials.NewStaticCredentials(
					config.Secrets.AWS.AccessKeyID,
					config.Secrets.AWS.SecretAccessKey,
					"",
				),
			})
			if err != nil {
				log.Fatalf("Unble to connect to AWS, error: %v\n", err)
			}

			logChannel := make(chan interface{})
			migrator := support.NewMigrator(sess, cld)

			startedAt := time.Now()
			log.Println("Migration started! 🔥🔥🔥")

			go migrator.Migrate(config, logChannel)

			for l := range logChannel {
				log.Println(l)
			}

			log.Printf("Migration done! 🎉🎉🎉\n Duration: %v seconds", time.Since(startedAt).Seconds())
		})

	commando.Parse(nil)
}
