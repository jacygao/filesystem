package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type s3Backed struct {
	s3client     *s3.S3
	bucketName   string
	bucketPrefix string
}

func NewS3Store(session *session.Session, bucketName, bucketPrefix string) *s3Backed {
	return &s3Backed{
		s3client:     s3.New(session),
		bucketName:   bucketName,
		bucketPrefix: bucketPrefix,
	}
}

func (s *s3Backed) getS3Key(ResourceKey string) string {
	return s.bucketPrefix + ResourceKey
}

func (s *s3Backed) Get(ctx context.Context, ResourceKey string) (io.ReadCloser, error) {
	log.Printf("Sending Get request to S3. Method: Get(%s)", ResourceKey)
	params := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s.getS3Key(ResourceKey)),
	}
	result, err := s.s3client.GetObject(params)
	if err != nil {
		return nil, err
	}
	log.Println("Finished sending Get request to S3")
	// Consider the issue of relying on the client to Close result.Body. Given that
	// the only client is handlers.go at the moment, and due to its use of image.Decode it can't
	// benefit from a streaming API, should we be robust here and read the entire stream?
	return result.Body, nil
}

func (s *s3Backed) Put(ctx context.Context, ResourceKey string, body io.Reader) error {
	log.Printf("Sending Put request to S3. Method: Put(%s)", ResourceKey)
	// The s3 API currently requires a ReadSeeker for the PutObject() endpoint. Hence we read
	// the entire body and wrap in a bytes.NewReader(), which has the Seek() method implemented.
	buffer, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	params := &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s.getS3Key(ResourceKey)),
		Body:   bytes.NewReader(buffer),
	}

	if _, err := s.s3client.PutObject(params); err != nil {
		return err
	}
	return nil
}

func (s *s3Backed) Delete(ctx context.Context, ResourceKey string) error {
	log.Printf("Sending Get request to S3. Method: Delete(%s)", ResourceKey)

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s.getS3Key(ResourceKey)),
	}
	if _, err := s.s3client.DeleteObject(params); err != nil {
		return err
	}
	log.Println("Finished sending Delete request to S3")
	return nil
}

// List reads the directory named by dir and returns a list of directory entries.
// Deprecated: The S3 List function has not been implemented.  If needed, please update gopexa library.
func (s *s3Backed) List(ctx context.Context, dir string) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("error: S3backed List function not implemented")
}
