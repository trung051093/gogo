package socketprovider

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
	engineio "github.com/googollee/go-socket.io/engineio"
)

type SocketService interface {
	OnConnect(namespace string, f func(socketio.Conn) error)
	OnDisconnect(namespace string, f func(socketio.Conn, string))
	OnEvent(namespace string, event string, f func(socketio.Conn, string))
	Serve() error
	Close() error
	GetServer() *socketio.Server
}

type SocketOptions struct {
	RedisConfig        *socketio.RedisAdapterOptions
	RedisAdapterEnable bool
}

type socketService struct {
	server *socketio.Server
}

func NewSocketProvider(options *engineio.Options, optionsFunc ...func(*SocketOptions)) *socketService {
	server := socketio.NewServer(options)
	provider := &socketService{server: server}
	provider.handlerOptions(optionsFunc...)
	return provider
}

func WithRedisAdapter(redisConfig *socketio.RedisAdapterOptions) func(*SocketOptions) {
	return func(so *SocketOptions) {
		so.RedisConfig = redisConfig
		so.RedisAdapterEnable = true
	}
}

func initOptions() *SocketOptions {
	option := &SocketOptions{}
	return option
}

func (s *socketService) handlerOptions(optionsFuncs ...func(*SocketOptions)) {
	options := initOptions()

	for _, optionFunc := range optionsFuncs {
		optionFunc(options)
	}

	if options.RedisAdapterEnable {
		if _, err := s.server.Adapter(options.RedisConfig); err != nil {
			log.Println("Socket connect Redis Adapter failed")
		}
	}
}

func (s *socketService) OnConnect(namespace string, f func(socketio.Conn) error) {
	s.server.OnConnect(namespace, f)
}

func (s *socketService) OnDisconnect(namespace string, f func(socketio.Conn, string)) {
	s.server.OnDisconnect(namespace, f)
}

func (s *socketService) OnEvent(namespace string, event string, f func(socketio.Conn, string)) {
	s.server.OnEvent(namespace, event, f)
}

func (s *socketService) Serve() error {
	return s.server.Serve()
}

func (s *socketService) Close() error {
	return s.server.Close()
}

func (s *socketService) GetServer() *socketio.Server {
	return s.server
}
