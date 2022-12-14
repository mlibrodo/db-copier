package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mlibrodo/db-copier/aws/config"
	"github.com/mlibrodo/db-copier/log"
	"os"
	"path/filepath"
)

type S3Object struct {
	Bucket string
	Key    string
}

func (in *S3Object) GetFileName() string {
	_, filename := filepath.Split(in.Key)
	return filename
}

func Download(s3Object S3Object, file *os.File) error {
	svc := s3.NewFromConfig(*config.AWSConfig)
	downloader := manager.NewDownloader(svc)

	numBytes, err := downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: aws.String(s3Object.Bucket),
		Key:    aws.String(s3Object.Key),
	})

	if err != nil {
		log.Error(err)
		return err
	}

	log.WithFields(
		log.Fields{
			"file":     file.Name(),
			"numBytes": numBytes,
		},
	).Debug("Finished Downloading")

	return nil
}

func Upload(s3Object S3Object, f *string) error {

	file, err := os.Open(*f)

	if err != nil {
		return err
	}

	defer func(file *os.File) {
		tmpErr := file.Close()

		if err != nil {
			err = tmpErr
		}

	}(file)

	svc := s3.NewFromConfig(*config.AWSConfig)

	uploader := manager.NewUploader(svc)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3Object.Bucket),
		Key:    aws.String(s3Object.Key),
		Body:   file,
	})

	if err != nil {
		log.Error(err)
		return err
	}

	log.WithFields(log.Fields{"file": file.Name(), "key": *result.Key}).Debug("Finished Uploading")

	return err
}

func List(bucket *string, prefix *string) ([]string, error) {
	svc := s3.NewFromConfig(*config.AWSConfig)

	var files []string

	notDone := true

	var out *s3.ListObjectsV2Output
	var err error
	var continuationToken *string

	for notDone {
		in := &s3.ListObjectsV2Input{
			Bucket:            bucket,
			Prefix:            prefix,
			ContinuationToken: continuationToken,
		}

		if out, err = svc.ListObjectsV2(context.TODO(), in); err != nil {
			return nil, err
		}

		for _, obj := range out.Contents {
			files = append(files, *obj.Key)
		}
		if out.ContinuationToken != nil {
			continuationToken = out.ContinuationToken
		} else {
			notDone = false
		}

	}
	return files, nil
}
