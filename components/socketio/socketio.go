package socketprovider

import (
	"context"
	"log"

	socketio "github.com/googollee/go-socket.io"
	engineio "github.com/googollee/go-socket.io/engineio"

	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

type SocketService interface {
	OnConnect(f func(socketio.Conn) error)
	OnDisconnect(f func(socketio.Conn, string))
	OnEvent(event string, f func(socketio.Conn, string))
	OnError(f func(socketio.Conn, error))
	JoinRoom(room string, connection socketio.Conn) bool
	LeaveRoom(room string, connection socketio.Conn) bool
	BroadcastToRoom(room string, event string, args ...interface{}) bool
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

	Namespace string
}

type socketService struct {
	server *socketio.Server
	config *SocketOptions
}

type key string

var SocketServiceKey key = "SocketServiceKey"

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
	provider.config = options
	return provider
}

func WithContext(ctx context.Context, socketService SocketService) context.Context {
	return context.WithValue(ctx, SocketServiceKey, socketService)
}

func FromContext(ctx context.Context) (*socketService, bool) {
	socketServiceCtx := ctx.Value(SocketServiceKey)
	if so, ok := socketServiceCtx.(*socketService); ok {
		return so, true
	}
	return nil, false
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

func WithNamespace(namespace string) func(*SocketOptions) {
	return func(so *SocketOptions) {
		so.Namespace = namespace
	}
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

	if options.Namespace == "" {
		options.Namespace = "/"
	}

	return options
}

func (s *socketService) OnConnect(f func(socketio.Conn) error) {
	s.server.OnConnect(s.config.Namespace, f)
}

func (s *socketService) OnDisconnect(f func(socketio.Conn, string)) {
	s.server.OnDisconnect(s.config.Namespace, f)
}

func (s *socketService) OnEvent(event string, f func(socketio.Conn, string)) {
	s.server.OnEvent(s.config.Namespace, event, f)
}

func (s *socketService) OnError(f func(socketio.Conn, error)) {
	s.server.OnError(s.config.Namespace, f)
}

func (s *socketService) JoinRoom(room string, connection socketio.Conn) bool {
	return s.server.JoinRoom(s.config.Namespace, room, connection)
}

func (s *socketService) LeaveRoom(room string, connection socketio.Conn) bool {
	return s.server.LeaveRoom(s.config.Namespace, room, connection)
}

func (s *socketService) BroadcastToRoom(room string, event string, args ...interface{}) bool {
	return s.server.BroadcastToRoom(s.config.Namespace, room, event, args...)
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
