package storage

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryStorage struct {
	cld    *cloudinary.Cloudinary
	folder string
}

func NewCloudinaryStorage(
	cloudName string,
	apiKey string,
	apiSecret string,
	folder string,
) (Storage, error) {

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)

	if err != nil {
		return nil, err
	}

	return &cloudinaryStorage{
		cld:    cld,
		folder: folder,
	}, nil
}

func (s *cloudinaryStorage) Upload(
	file io.Reader,
	filename string,
) (*UploadResult, error) {

	result, err := s.cld.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			Folder:         s.folder,
			ResourceType:   "image",
			UniqueFilename: boolPtr(true),
		},
	)

	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
	}, nil
}

func (s *cloudinaryStorage) Delete(publicID string) error {
	_, err := s.cld.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID:     publicID,
			ResourceType: "image",
		},
	)

	return err
}

func boolPtr(v bool) *bool {
	return &v
}
