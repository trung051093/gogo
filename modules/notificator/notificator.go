package notificator

import (
	"fmt"
	"gogo/common"
	"gogo/components/appctx"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func Handler(appctx appctx.AppContext) {
	defer common.Recovery()
	socketService := appctx.GetSocketService()
	socketService.OnConnect(func(s socketio.Conn) error {
		log.Println("Connected:", s.ID())
		return nil
	})
	socketService.OnEvent("msg", func(s socketio.Conn, msg interface{}) {
		log.Println("Receive Message : " + msg.(string))
		s.Emit("reply", "OK")
	})
	socketService.OnEvent("msg", func(s socketio.Conn, msg interface{}) {
		log.Println("Receive Message : " + msg.(string))
		s.Emit("reply", "OK")
	})
	socketService.OnEvent("join:notification", func(s socketio.Conn, msg interface{}) {
		log.Println("join notification: " + msg.(string))
		ok := socketService.JoinRoom("notification", s)
		s.Emit("reply", fmt.Sprintf("%s join room 'notification' %v", s.ID(), ok))
	})
	socketService.OnEvent("leave:notification", func(s socketio.Conn, msg interface{}) {
		log.Println("leave notification: " + msg.(string))
		ok := socketService.LeaveRoom("notification", s)
		s.Emit("reply", fmt.Sprintf("%s leave room 'notification' %v", s.ID(), ok))
	})
	socketService.Serve()
	defer socketService.Close()
}
