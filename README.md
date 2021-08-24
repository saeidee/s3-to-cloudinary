## Trek(S3 -> Cloudinary)

__It is an implementation of migrating all of your media from S3 to Cloudinary, using goroutines it has the highest performance__

### Usage

You need to enable lazy loading and map your buckets in your Cloudinary dashboard you can find more information about lazy migration here https://cloudinary.com/documentation/fetch_remote_images#auto_upload_remote_files.

After that update your keys in the config.yml and then since this project is a CLI tool, so you can run the start command to start migration.

```console
foo@bar:~$ go run trek.go start
```

### Contributors

- __Richard Bajuzik__
