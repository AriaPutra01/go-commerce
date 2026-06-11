package product

import "mime/multipart"

type UploadProductImageRequest struct {
	IsPrimary bool                  `form:"is_primary"`
	File      *multipart.FileHeader `form:"file" binding:"required"`
}

type UploadProductImageResponse struct {
	ID            string `json:"id"`
	SavedFileName string `json:"saved_file_name"`
}
