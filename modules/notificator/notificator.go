package notificator

import (
	"fmt"
	"user_management/common"
	"user_management/components/appctx"

	socketio "github.com/googollee/go-socket.io"
)

func Handler(appctx appctx.AppContext) {
	defer common.Recovery()
	socketService := appctx.GetSocketService()
	socketService.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("Connected:", s.ID())
		return nil
	})
	socketService.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		fmt.Println("Receive Message : " + msg)
		s.Emit("reply", "OK")
	})
	socketService.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		fmt.Println("Receive Message : " + msg)
		s.Emit("reply", "OK")
	})
	socketService.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("Somebody just close the connection ")
	})
	go socketService.Serve()
	defer socketService.Close()
}
