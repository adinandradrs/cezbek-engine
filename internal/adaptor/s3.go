package adaptor

import (
	"github.com/adinandradrs/cezbek-engine/internal/apps"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type S3Bucket struct {
	Bucket   string
	Uploader *s3manager.Uploader
	Logger   *zap.Logger
}

type S3Watcher interface {
	Upload(req *model.S3UploadRequest) (*s3manager.UploadOutput, *model.TechnicalError)
}

func NewS3(b S3Bucket) S3Watcher {
	return &b
}

func (b S3Bucket) Upload(req *model.S3UploadRequest) (*s3manager.UploadOutput, *model.TechnicalError) {
	v, err := b.Uploader.Upload(&s3manager.UploadInput{
		Bucket:      &b.Bucket,
		Key:         &req.Destination,
		Body:        req.Source,
		ContentType: &req.ContentType,
	})

	if err != nil {
		return nil, apps.Exception("failed to upload file", err, zap.Any("data", req), b.Logger)
	}
	b.Logger.Info("s3 response", zap.Any("", v))
	return v, nil
}
