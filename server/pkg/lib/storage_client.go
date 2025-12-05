package libfx

import (
	"fmt"

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

func NewStorageClient(params StorageClientParams) (*minio.Client, error) {
	useSecure := params.FlagConfig.Environment == "production"

	minioClient, err := minio.New(
		params.AppConfig.StorageEndpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(params.AppConfig.StorageAccessKeyID, params.AppConfig.StorageSecretAccessKey, ""),
			Secure: useSecure,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to storage: %w", err)
	}

	return minioClient, nil
}
