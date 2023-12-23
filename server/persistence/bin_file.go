package persistence

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strix-server/system"
)

var MinioClient *minio.Client

func InitBinary() {
	var err error
	binConfig := system.SystemConfig.Binary
	MinioClient, err = minio.New(binConfig.ServerAddress, &minio.Options{
		Creds: credentials.NewStaticV4(binConfig.Username, binConfig.Password, ""),
	})
	if err != nil {
		system.Logger.Fatal(err)
	}
	ctx := context.Background()
	err = MinioClient.MakeBucket(ctx, binConfig.Bucket, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := MinioClient.BucketExists(ctx, binConfig.Bucket)
		if errBucketExists == nil && exists {
			system.Logger.Infof("We already own %s\n", binConfig.Bucket)
		} else {
			system.Logger.Fatal(err)
		}
	} else {
		system.Logger.Infof("Successfully created %s\n", binConfig.Bucket)
	}
}
