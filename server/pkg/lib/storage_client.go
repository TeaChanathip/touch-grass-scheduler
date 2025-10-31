package libfx

import (
	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type StorageClientParams struct {
	fx.In
	FlagConfig *configfx.FlagConfig
	AppConfig  *configfx.AppConfig
	Logger     *zap.Logger
}

func NewStorageClient(params StorageClientParams) *minio.Client {
	var useSecure bool = false
	if params.FlagConfig.Environment == "production" {
		useSecure = true
	}

	minioClient, err := minio.New(
		params.AppConfig.StorageEndpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(params.AppConfig.StorageAccessKeyID, params.AppConfig.StorageSecretAccessKey, ""),
			Secure: useSecure,
		},
	)
	if err != nil {
		params.Logger.Fatal("Error connecting to Storage", zap.Error(err))
	}

	return minioClient
}
