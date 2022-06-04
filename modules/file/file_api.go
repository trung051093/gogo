package file

import (
	"errors"
	"net/http"
	"time"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	ctx *gin.Context
}

func GetUploadPresignedUrl(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		fileName := ginCtx.Query("fileName")
		fileType := ginCtx.Query("fileType")
		objectName := uuid.New().String()

		if fileType == "" {
			panic(common.NewCustomError(errors.New("Cannot found file type"), "Cannot found file type", "FILE"))
		}

		if fileName == "" {
			fileName = uuid.New().String()
		}

		storageService := appCtx.GetStorageService()
		policy := storage.NewPostPolicy()
		policy.SetBucket(common.PhotoBucket)
		policy.SetKey(objectName)
		policy.SetContentLengthRange(0, 10000000)
		policy.SetUserMetadata("fileName", fileName)
		policy.SetContentType(fileType)
		// expire in 1 day
		policy.SetExpires(time.Now().UTC().AddDate(0, 0, 1))

		url, formData, presignErr := storageService.PresignedPostObject(
			ginCtx.Request.Context(),
			policy,
		)

		if presignErr != nil {
			panic(common.NewCustomError(presignErr, "Cannot get presign url", "PRESIGN_URL"))
		}

		res := &map[string]interface{}{
			"url":    url.String(),
			"fields": formData,
		}

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(res))
	}
}
