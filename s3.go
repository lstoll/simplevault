package main

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

type S3 interface {
	Get(string) ([]byte, error)
	Put(string, []byte) error
}

type s3client struct {
	client *s3.S3
	bucket *s3.Bucket
	prefix string
}

func NewS3(key, secret, bucket, prefix string) S3 {
	auth, _ := aws.GetAuth(key, secret)
	region := aws.USEast // haha for life. TODO - configurable?
	client := s3.New(auth, region)
	return &s3client{
		client: client,
		bucket: client.Bucket(bucket),
		prefix: prefix,
	}
}
func (s *s3client) Get(key string) ([]byte, error) {
	return s.bucket.Get(s.prefix + "/" + key)
}

func (s *s3client) Put(key string, data []byte) error {
	return s.bucket.Put(s.prefix+"/"+key, data, "application/octet-stream", s3.BucketOwnerFull)
}
