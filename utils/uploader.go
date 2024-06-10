package utils

import (
	"context"
	"mime/multipart"
	"movie-ticket-booking/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type ImageInfo struct {
	ImageURL string
	PublicID string
}

func UploadToCloudinary(file multipart.File, filePath string, folder string) (ImageInfo, error) {
	ctx := context.Background()
	cld, err := config.SetupCloudinary()
	if err != nil {
		return ImageInfo{}, err
	}
	publicID := uuid.New().String()

	uploadParams := uploader.UploadParams{
		PublicID: folder + "/" + publicID,
		Folder:   folder,
	}
	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return ImageInfo{}, err
	}

	imageInfo := ImageInfo{
		ImageURL: result.SecureURL,
		PublicID: result.PublicID,
	}

	return imageInfo, nil
}
func DeleteFromCloudinary(publicID string) error {
	ctx := context.Background()
	cld, err := config.SetupCloudinary()
	if err != nil {
		return err
	}

	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return err
	}

	return nil
}
