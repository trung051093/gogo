package notificator

import (
	"context"
	"gogo/common"
	"gogo/components/appctx"
	"log"

	"github.com/minio/minio-go/v7/pkg/notification"
)

func FileHandler(appCtx appctx.AppContext) {
	storageService := appCtx.GetStorageService()
	storageService.ListenNotification(
		context.Background(),
		common.ImageBucket,
		"",
		"",
		[]string{
			"s3:ObjectCreated:*",
			"s3:ObjectRemoved:*",
		},
		func(noti *notification.Info) {
			log.Println("minio notification info:", noti)
		},
	)
}
