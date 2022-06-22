package main

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/saeidee/trek/internal"
	"github.com/thatisuday/commando"
	"github.com/xuri/excelize/v2"
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
			configParser := internal.NewConfigParser()
			config, err := configParser.Parse("config.yml")
			if err != nil {
				log.Fatalf("Unable to read the file with error: %v", err)
			}

			factory := internal.NewFactory(config)

			cld, err := factory.NewCloudinaryInstance()
			if err != nil {
				log.Fatalf("Unable to connect to Cloudinary, error: %v\n", err)
			}

			sess, err := factory.NewAWSSession()
			if err != nil {
				log.Fatalf("Unble to connect to AWS, error: %v\n", err)
			}

			startedAt := time.Now()
			logChannel := make(chan internal.Log)
			migrator := internal.NewMigrator(s3.New(sess), cld)
			logger := internal.NewLogger(excelize.NewFile())

			log.Println("Migration started! 🔥🔥🔥")

			go migrator.Migrate(config, logChannel)

			for l := range logChannel {
				logger.Log(l)
				log.Println(l.Bucket, l.Item, l.Error)
			}

			_ = logger.SaveFile("s3-to-cloudinary-logs.xlsx")

			log.Printf("Migration done! 🎉🎉🎉\n Duration: %v seconds", time.Since(startedAt).Seconds())
		})

	commando.Parse(nil)
}
