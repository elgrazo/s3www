package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	minio "github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/credentials"
	"github.com/minio/minio-go/v6/pkg/s3utils"
)

// S3 - A S3 implements FileSystem using the minio client
// allowing access to your S3 buckets and objects.
//
// Note that S3 will allow all access to files in your private
// buckets, If you have any sensitive information please make
// sure to not sure this project.
type S3 struct {
	*minio.Client
	bucket string
}

// Open - implements http.Filesystem implementation.
func (s3 *S3) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, pathSeparator) {
		return &httpMinioObject{
			client: s3.Client,
			object: nil,
			isDir:  true,
			bucket: bucket,
			prefix: strings.TrimSuffix(name, pathSeparator),
		}, nil
	}

	name = strings.TrimPrefix(name, pathSeparator)
	obj, err := getObject(s3, name)
	if err != nil {
		if notfound == "" {
			return nil, os.ErrNotExist
		} else {
			obj, err = getObject(s3, notfound)
			return &httpMinioObject{
				client: s3.Client,
				object: obj,
				isDir:  false,
				bucket: bucket,
				prefix: name,
			}, os.ErrNotExist
		}

	}

	return &httpMinioObject{
		client: s3.Client,
		object: obj,
		isDir:  false,
		bucket: bucket,
		prefix: name,
	}, nil
}

func getObject(s3 *S3, name string) (*minio.Object, error) {
	names := [3]string{name, name + "/index.html", name + "/index.htm"}
	for _, n := range names {
		obj, err := s3.Client.GetObject(s3.bucket, n, minio.GetObjectOptions{})
		if err == nil {
			if _, err = obj.Stat(); err == nil {
				return obj, nil
			}
		}
	}
	return nil, os.ErrNotExist
}

var (
	endpoint  string
	accessKey string
	secretKey string
	address   string
	bucket    string
	notfound  string
)

func init() {

	endpoint = getEnv("ENDPOINT", "")
	accessKey = getEnv("ACCESSKEY", "")
	secretKey = getEnv("SECRETKEY", "")
	bucket = getEnv("BUCKET", "")
	address = getEnv("ADDRESS", "127.0.0.1:8080")
	notfound = getEnv("404PAGE", "")
}

func getEnv(key, dflt string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		val = dflt
	}
	return val
}

// NewCustomHTTPTransport returns a new http configuration
// used while communicating with the cloud backends.
// This sets the value for MaxIdleConnsPerHost from 2 (go default)
// to 100.
func NewCustomHTTPTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          1024,
		MaxIdleConnsPerHost:   1024,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
	}
}

func FileServerWithCustom404(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			if notfound != "" {
				w.WriteHeader(http.StatusNotFound)
				r.URL.Path = path.Clean(notfound)
				fsh.ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}

			return
		}
		fsh.ServeHTTP(w, r)
	})
}

func main() {

	if strings.TrimSpace(bucket) == "" {
		log.Fatalln(`Bucket name cannot be empty, please provide BUCKET environmental variable`)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	// Chains all credential types, in the following order:
	//  - AWS env vars (i.e. AWS_ACCESS_KEY_ID)
	//  - AWS creds file (i.e. AWS_SHARED_CREDENTIALS_FILE or ~/.aws/credentials)
	//  - IAM profile based credentials. (performs an HTTP
	//    call to a pre-defined endpoint, only valid inside
	//    configured ec2 instances)
	var defaultAWSCredProviders = []credentials.Provider{
		&credentials.EnvAWS{},
		&credentials.FileAWSCredentials{},
		&credentials.IAM{
			Client: &http.Client{
				Transport: NewCustomHTTPTransport(),
			},
		},
		&credentials.EnvMinio{},
	}
	if accessKey != "" && secretKey != "" {
		defaultAWSCredProviders = []credentials.Provider{
			&credentials.Static{
				Value: credentials.Value{
					AccessKeyID:     accessKey,
					SecretAccessKey: secretKey,
				},
			},
		}
	}

	// If we see an Amazon S3 endpoint, then we use more ways to fetch backend credentials.
	// Specifically IAM style rotating credentials are only supported with AWS S3 endpoint.
	creds := credentials.NewChainCredentials(defaultAWSCredProviders)

	client, err := minio.NewWithOptions(u.Host, &minio.Options{
		Creds:        creds,
		Secure:       u.Scheme == "https",
		Region:       s3utils.GetRegionFromURL(*u),
		BucketLookup: minio.BucketLookupAuto,
	})
	if err != nil {
		log.Fatalln(err)
	}

	//mux := http.FileServer(&S3{client, bucket})
	mux := FileServerWithCustom404(&S3{client, bucket})
	log.Printf("Started listening on http://%s\n", address)
	log.Fatalln(http.ListenAndServe(address, mux))
}
