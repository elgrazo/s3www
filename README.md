# s3www
Serve static files from any S3 compatible object storage endpoints.

Fork from [harshavardhana/s3www](https://github.com/harshavardhana/s3www) version v0.3.0 with following modifications:

- Removed Let's Encrypt support (should be handled by a load balancer like traefik)
- Settings will be set by using environmental variables instead of flags
- Custom 404 page can be loaded from bucket

## Run on Docker
Pull from Dockerhub:
```

```

Set the environment variables:
```
ENDPOINT: Address of your S3 instance (e.g. https://minio.example.com)
ACCESSKEY: Access Key
SECRETKEY: Secret Key
BUCKET: Bucket name
ADDRESS: Address for S3WWW to listen to (Default: 127.0.0.1:8080)
404PAGE: Name of the 404 error page in the bucket (e.g. 404.html)
```


## Binaries
Released binaries are available [here](https://github.com/elgrazo/s3www/releases), or you can compile yourself from source.

```
go get github.com/elgrazo/s3www
```




## License
This project is distributed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0), see [LICENSE](./LICENSE) for more information.

