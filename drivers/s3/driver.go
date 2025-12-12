package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	filesystem "github.com/donnigundala/dg-filesystem"
)

// S3Disk implements the Disk interface for AWS S3 compatible storage.
type S3Disk struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	uploader   *manager.Uploader
	downloader *manager.Downloader
	bucket     string
	root       string // Optional prefix
	endpoint   string // Public URL endpoint
}

// NewS3Disk creates a new S3 disk instance.
func NewS3Disk(cfg map[string]interface{}) (filesystem.Disk, error) {
	bucket, ok := cfg["bucket"].(string)
	if !ok {
		return nil, fmt.Errorf("s3 driver requires 'bucket' config")
	}

	region, _ := cfg["region"].(string)
	// accessKey, _ := cfg["key"].(string)
	// secretKey, _ := cfg["secret"].(string)
	// We prefer environment variables or IAM roles, but could support static creds here if needed
	// For "key" and "secret", we would use credentials.NewStaticCredentialsProvider

	endpoint, _ := cfg["url"].(string)         // Public URL base
	apiEndpoint, _ := cfg["endpoint"].(string) // API Endpoint (e.g. for MinIO)
	usePathStyle, _ := cfg["use_path_style"].(bool)

	// Load default config
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	if apiEndpoint != "" {
		// Custom endpoint resolver for MinIO or other S3 compatibles
		awsCfg.BaseEndpoint = aws.String(apiEndpoint)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = usePathStyle
	})

	presigner := s3.NewPresignClient(client)
	uploader := manager.NewUploader(client)
	downloader := manager.NewDownloader(client)

	return &S3Disk{
		client:     client,
		presigner:  presigner,
		uploader:   uploader,
		downloader: downloader,
		bucket:     bucket,
		endpoint:   endpoint,
		root:       "", // we can add root prefix support later
	}, nil
}

func (d *S3Disk) Put(path string, content []byte) error {
	reader := bytes.NewReader(content)
	_, err := d.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
		Body:   reader,
	})
	return err
}

func (d *S3Disk) PutStream(path string, content io.Reader) error {
	_, err := d.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
		Body:   content,
	})
	return err
}

func (d *S3Disk) Get(path string) ([]byte, error) {
	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := d.downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (d *S3Disk) GetStream(path string) (io.ReadCloser, error) {
	out, err := d.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (d *S3Disk) Exists(path string) (bool, error) {
	_, err := d.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		// Check for 404
		// Error handling for specific status codes in aws-sdk-v2 is a bit verbose
		// Assuming error means not found for simple exists check, but technically could be permission denied
		return false, nil
	}
	return true, nil
}

func (d *S3Disk) Delete(path string) error {
	_, err := d.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	return err
}

func (d *S3Disk) Url(path string) string {
	if d.endpoint != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(d.endpoint, "/"), strings.TrimPrefix(path, "/"))
	}
	// Fallback to standard S3 URL (virtual-hosted style)
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", d.bucket, strings.TrimPrefix(path, "/"))
}

func (d *S3Disk) SignedUrl(path string, expiration time.Duration) (string, error) {
	req, err := d.presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	}, func(o *s3.PresignOptions) {
		o.Expires = expiration
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (d *S3Disk) MakeDirectory(path string) error {
	// S3 is flat, no action needed for directories usually
	// We could upload a 0-byte object with trailing slash if strict folder emulation is needed
	return nil
}

func (d *S3Disk) DeleteDirectory(path string) error {
	// Need to list objects with prefix and delete them
	// This is expensive, implementing minimal version
	// Use ListObjectsV2 and DeleteObjects
	return fmt.Errorf("DeleteDirectory not implemented for S3 yet")
}
