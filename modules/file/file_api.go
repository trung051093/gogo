package file

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"user_management/common"
	"user_management/components/appctx"
	storageprovider "user_management/components/storage"
	filemodel "user_management/modules/file/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUploadPresignedUrl godoc
// @Summary      get a presigned url to upload
// @Description  get a presigned url to upload
// @Tags         file
// @Accept       json
// @Produce      json
// @Param        fileName  query     string                                               true   "fileName"
// @Param        fileType  query     string                                               true   "fileType"
// @Param        proxy     query     string                                               false  "proxy"
// @Success      200       {object}  common.Response{data=filemodel.PresignedPostObject}  "desc"
// @Failure      400       {object}  common.AppError
// @Router       /api/v1/file/presign-url [get]
func GetUploadPresignedUrl(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		config := appCtx.GetConfig()
		configStorage := config.GetStorageConfig()

		fileName := ginCtx.Query("fileName")
		fileType := ginCtx.Query("fileType")
		proxy := ginCtx.Query("proxy")

		objectName := uuid.New().String()

		if fileType == "" {
			panic(common.NewCustomError(errors.New("cannot found file type"), "Cannot found file type", "FILE"))
		}

		if fileName == "" {
			fileName = uuid.New().String()
		}

		storageService := appCtx.GetStorageService()
		policy := storageprovider.NewPostPolicy()

		policy.SetBucket(common.ImageBucket)
		policy.SetKey(objectName)
		policy.SetContentLengthRange(0, 10000000)
		policy.SetUserMetadata("fileName", fileName)
		policy.SetContentType(fileType)
		// expire in 1 day
		policy.SetExpires(time.Now().UTC().AddDate(0, 0, 1))

		_, formData, presignErr := storageService.PresignedPostObject(
			ginCtx.Request.Context(),
			policy,
		)

		if presignErr != nil {
			panic(common.NewCustomError(presignErr, "Cannot get presign url", "PRESIGN_URL"))
		}

		publicUrl := configStorage.PublicUrl
		if proxy != "" {
			publicUrl = proxy
		}
		res := &filemodel.PresignedPostObject{
			Url:    fmt.Sprintf("%s/%s", publicUrl, common.ImageBucket),
			Fields: formData,
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(res))
	}
}
