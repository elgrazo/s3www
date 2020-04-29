# s3www
Serve static files from any S3 compatible object storage endpoints.

Fork from [harshavardhana/s3www](https://github.com/harshavardhana/s3www) version 0.3.0 with following modifications:

- Removed Let's Encrypt support (should be handled by a load balancer like traefik)
- Settings will be set by using environmental variables instead of flags
- Custom 404 page can be loaded from bucket



## Binaries
Released binaries are available [here](https://github.com/elgrazo/s3www/releases), or you can compile yourself from source.

```
go get github.com/elgrazo/s3www
```




### Test
Point your web browser to https://example.com ensure your `s3www` is serving your `index.html` successfully.


## Run locally
Make sure you have `index.html` under `website-bucket`
```
s3www -endpoint "https://s3.amazonaws.com" -accessKey "accessKey" \
      -secretKey "secretKey" -bucket "website-bucket"

s3www: Started listening on http://127.0.0.1:8080
```

### Test
Point your web browser to http://127.0.0.1:8080 ensure your `s3www` is serving your `index.html` successfully.

## License
This project is distributed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0), see [LICENSE](./LICENSE) for more information.

