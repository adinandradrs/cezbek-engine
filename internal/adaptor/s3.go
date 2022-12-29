package adaptor

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
	"io"
)

type (
	S3Bucket struct {
		Bucket   string
		Uploader *s3manager.Uploader
		Logger   *zap.Logger
	}

	S3UploadRequest struct {
		ContentType string
		Source      io.Reader
		Destination string
	}

	S3UploadResponse struct {
		Id      string
		Version string
	}
)

type S3Watcher interface {
	Upload(req *S3UploadRequest) (*S3UploadResponse, *model.TechnicalError)
}

func NewS3(b *S3Bucket) S3Watcher {
	return b
}

func (b S3Bucket) Upload(req *S3UploadRequest) (*S3UploadResponse, *model.TechnicalError) {
	v, err := b.Uploader.Upload(&s3manager.UploadInput{
		Bucket:      &b.Bucket,
		Key:         &req.Destination,
		Body:        req.Source,
		ContentType: &req.ContentType,
	})

	if err != nil {
		return nil, apps.Exception("failed to upload file", err, zap.Any("data", req), b.Logger)
	}

	return &S3UploadResponse{
		Id:      v.UploadID,
		Version: *v.VersionID,
	}, nil
}
