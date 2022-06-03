package file

import (
	"fmt"
	"net/http"
	"time"
	"user_management/common"
	"user_management/components/appctx"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	ctx *gin.Context
}

func GetPresignedUrlToUpload(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		storageService := appCtx.GetStorageService()
		url, presignErr := storageService.PresignedPutObject(
			ginCtx.Request.Context(),
			common.PhotoBucket,
			uuid.New().String(),
			30*time.Minute,
		)

		if presignErr != nil {
			panic(common.NewCustomError(presignErr, "Cannot get presign url", "PRESIGN_URL"))
		}

		presignedUrl := fmt.Sprintf("http://%s%s", url.Host, url.RequestURI())

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(presignedUrl))
	}
}
