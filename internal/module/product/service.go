package product

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
)

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) UploadProductImage(ctx context.Context, req *UploadProductImageRequest) (*UploadProductImageResponse, error) {
	newID := uuid.New()
	ext := filepath.Ext(req.File.Filename)
	savedFileName := fmt.Sprintf("%s%s", newID.String(), ext)

	if err := s.repository.CreateProductImage(ctx, &ProductImage{
		ID:        newID,
		UrlPath:   savedFileName,
		IsPrimary: req.IsPrimary,
	}); err != nil {
		return nil, err
	}

	return &UploadProductImageResponse{
		ID:            newID.String(),
		SavedFileName: savedFileName,
	}, nil
}
