package middleware

import (
	"context"
	"fmt"
	"golectro-product/internal/constants"
	"golectro-product/internal/utils"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type UploadOptions struct {
	FieldName     string
	MaxFileSizeMB int64
	MaxFiles      int
	BucketName    string
	AllowedTypes  []string
}

func SingleFileUpload(minioClient *minio.Client, opts UploadOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := GetUser(c)
		file, err := c.FormFile(opts.FieldName)
		if err != nil {
			res := utils.FailedResponse(c, http.StatusBadRequest, constants.FileNotFound, nil)
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		if file.Size > opts.MaxFileSizeMB*1024*1024 {
			res := utils.FailedResponse(c, http.StatusBadRequest, constants.FileSizeExceeded, nil)
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		src, err := file.Open()
		if err != nil {
			res := utils.FailedResponse(c, http.StatusInternalServerError, constants.InvalidOpenFile, nil)
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}
		defer src.Close()

		buffer := make([]byte, 512)
		if _, err := src.Read(buffer); err != nil {
			res := utils.FailedResponse(c, http.StatusInternalServerError, constants.InvalidReadFile, nil)
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}
		contentType := http.DetectContentType(buffer)

		if _, err := src.Seek(0, io.SeekStart); err != nil {
			res := utils.FailedResponse(c, http.StatusInternalServerError, constants.InvalidResetPosition, nil)
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		if len(opts.AllowedTypes) > 0 {
			allowed := slices.Contains(opts.AllowedTypes, contentType)
			if !allowed {
				res := utils.FailedResponse(c, http.StatusInternalServerError, constants.InvalidFileType, nil)
				c.AbortWithStatusJSON(res.StatusCode, res)
				return
			}
		}

		now := time.Now()
		ext := filepath.Ext(file.Filename)
		formattedTime := now.Format("02-01-2006-15-04-05-000")
		uniqueName := fmt.Sprintf("%v-%s%s", auth.ID, formattedTime, ext)

		metadata := map[string]string{
			"original_name": file.Filename,
			"author_id":     auth.ID.String(),
			"created_at":    now.Format(time.RFC3339Nano),
			"updated_at":    now.Format(time.RFC3339Nano),
			"file_size":     fmt.Sprintf("%d", file.Size),
		}

		info, err := minioClient.PutObject(
			context.Background(),
			opts.BucketName,
			uniqueName,
			src,
			file.Size,
			minio.PutObjectOptions{
				ContentType:  contentType,
				UserMetadata: metadata,
			},
		)
		if err != nil {
			res := utils.FailedResponse(c, http.StatusInternalServerError, constants.FailedUploadObject, nil)
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		c.Set("uploadedFile", map[string]any{
			"file_name":     uniqueName,
			"original_name": file.Filename,
			"file_size":     file.Size,
			"url":           fmt.Sprintf("/%s/%s", info.Bucket, info.Key),
			"author_id":     auth.ID,
			"created_at":    now,
			"updated_at":    now,
		})

		c.Next()
	}
}

func MultipleFileUpload(minioClient *minio.Client, opts UploadOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := GetUser(c)

		form, err := c.MultipartForm()
		if err != nil {
			res := utils.FailedResponse(c, http.StatusBadRequest, constants.FileNotFound, nil)
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		files := form.File[opts.FieldName]
		if len(files) == 0 {
			res := utils.FailedResponse(c, http.StatusBadRequest, constants.FileNotFound, nil)
			c.AbortWithStatusJSON(res.StatusCode, res)
			return
		}

		var uploadedFiles []map[string]any

		for _, file := range files {
			if file.Size > opts.MaxFileSizeMB*1024*1024 {
				res := utils.FailedResponse(c, http.StatusBadRequest, constants.FileSizeExceeded, nil)
				c.AbortWithStatusJSON(res.StatusCode, res)
				return
			}

			src, err := file.Open()
			if err != nil {
				res := utils.FailedResponse(c, http.StatusInternalServerError, constants.InvalidOpenFile, nil)
				c.AbortWithStatusJSON(res.StatusCode, res)
				return
			}

			buffer := make([]byte, 512)
			if _, err := src.Read(buffer); err != nil {
				src.Close()
				res := utils.FailedResponse(c, http.StatusInternalServerError, constants.InvalidReadFile, nil)
				c.AbortWithStatusJSON(res.StatusCode, res)
				return
			}
			contentType := http.DetectContentType(buffer)

			if _, err := src.Seek(0, io.SeekStart); err != nil {
				src.Close()
				res := utils.FailedResponse(c, http.StatusInternalServerError, constants.InvalidResetPosition, nil)
				c.AbortWithStatusJSON(res.StatusCode, res)
				return
			}

			if len(opts.AllowedTypes) > 0 {
				allowed := slices.Contains(opts.AllowedTypes, contentType)
				if !allowed {
					src.Close()
					res := utils.FailedResponse(c, http.StatusInternalServerError, constants.InvalidFileType, nil)
					c.AbortWithStatusJSON(res.StatusCode, res)
					return
				}
			}

			now := time.Now()
			ext := filepath.Ext(file.Filename)
			formattedTime := now.Format("02-01-2006-15-04-05-000")
			uniqueName := fmt.Sprintf("%v-%s%s", auth.ID, formattedTime, ext)

			metadata := map[string]string{
				"original_name": file.Filename,
				"author_id":     auth.ID.String(),
				"created_at":    now.Format(time.RFC3339Nano),
				"updated_at":    now.Format(time.RFC3339Nano),
				"file_size":     fmt.Sprintf("%d", file.Size),
			}

			info, err := minioClient.PutObject(
				context.Background(),
				opts.BucketName,
				uniqueName,
				src,
				file.Size,
				minio.PutObjectOptions{
					ContentType:  contentType,
					UserMetadata: metadata,
				},
			)
			src.Close()

			if err != nil {
				res := utils.FailedResponse(c, http.StatusInternalServerError, constants.FailedUploadObject, nil)
				c.AbortWithStatusJSON(res.StatusCode, res)
				return
			}

			uploadedFiles = append(uploadedFiles, map[string]any{
				"file_name":     uniqueName,
				"original_name": file.Filename,
				"file_size":     file.Size,
				"url":           fmt.Sprintf("/%s/%s", info.Bucket, info.Key),
				"author_id":     auth.ID,
				"created_at":    now,
				"updated_at":    now,
			})
		}

		c.Set("uploadedFiles", uploadedFiles)
		c.Next()
	}
}
