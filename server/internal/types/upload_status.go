package types

type UploadStatus BaseStringEnum

const (
	UploadStatusPending   UploadStatus = "pending"
	UploadStatusCompleted UploadStatus = "completed"
)
