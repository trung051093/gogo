package socketprovider

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
	engineio "github.com/googollee/go-socket.io/engineio"

	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

type SocketService interface {
	OnConnect(namespace string, f func(socketio.Conn) error)
	OnDisconnect(namespace string, f func(socketio.Conn, string))
	OnEvent(namespace string, event string, f func(socketio.Conn, string))
	OnError(namespace string, f func(socketio.Conn, error))
	Serve() error
	Close() error
	GetServer() *socketio.Server
}

type SocketRedisAdapterConfig socketio.RedisAdapterOptions

type SocketEngineConfig engineio.Options

type SocketOptions struct {
	// redis adapter config
	RedisAdapterConfig *SocketRedisAdapterConfig

	// transport config
	SocketTransport []transport.Transport
}

type socketService struct {
	server *socketio.Server
}

func NewSocketProvider(optionsFunc ...func(*SocketOptions)) *socketService {
	provider := &socketService{}
	options := provider.handlerOptions(optionsFunc...)

	server := socketio.NewServer(&engineio.Options{
		Transports: options.SocketTransport,
	})

	if options.RedisAdapterConfig != nil {
		redisAdapterConfig := &socketio.RedisAdapterOptions{
			Addr:     options.RedisAdapterConfig.Addr,
			Port:     options.RedisAdapterConfig.Port,
			Prefix:   options.RedisAdapterConfig.Prefix,
			Host:     options.RedisAdapterConfig.Host,
			Password: options.RedisAdapterConfig.Password,
		}
		if _, err := server.Adapter(redisAdapterConfig); err != nil {
			log.Println("Socket connect Redis Adapter failed")
		}
	}

	provider.server = server
	return provider
}

func WithRedisAdapter(redisConfig *SocketRedisAdapterConfig) func(*SocketOptions) {
	return func(so *SocketOptions) {
		so.RedisAdapterConfig = redisConfig
	}
}

func WithWebsocketTransport(so *SocketOptions) {
	so.SocketTransport = append(so.SocketTransport, websocket.Default)
}

func WithPollingTransport(so *SocketOptions) {
	so.SocketTransport = append(so.SocketTransport, polling.Default)
}

func initOptions() *SocketOptions {
	option := &SocketOptions{
		SocketTransport: []transport.Transport{},
	}
	return option
}

func (s *socketService) handlerOptions(optionsFuncs ...func(*SocketOptions)) *SocketOptions {
	options := initOptions()

	for _, optionFunc := range optionsFuncs {
		optionFunc(options)
	}

	return options
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

func (s *socketService) OnError(namespace string, f func(socketio.Conn, error)) {
	s.server.OnError(namespace, f)
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
