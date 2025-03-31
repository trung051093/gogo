package service

import (
	"context"
	"gogo/common"
	"gogo/components/appctx"
	"gogo/components/hasher"
	"gogo/modules/snake_battle/dto"
	"gogo/modules/snake_battle/entity"
	"gogo/modules/snake_battle/repository"
	"log"

	socketio "github.com/googollee/go-socket.io"
	"github.com/mitchellh/mapstructure"
)

type GameService interface {
	CreateRoom(ctx context.Context, req *dto.CreateRoomReq) (*dto.CreateRoomRes, error)
	GameSocketListener(ctx context.Context) error
}

type gameService struct {
	appCtx         appctx.AppContext
	hashService    hasher.HashService
	gameRepo       repository.GameRepository
	roomRepo       repository.RoomRepository
	userRepo       repository.UserRepository
	roomMemberRepo repository.RoomMemberRepository
}

func NewGameService(appCtx appctx.AppContext) GameService {
	hashService := appCtx.GetHashService()
	db := appCtx.GetMainDBConnection()
	gameRepo := repository.NewGameRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	userRepo := repository.NewUserRepository(db)
	roomMemberRepo := repository.NewRoomMemberRepository(db)
	return &gameService{
		appCtx:         appCtx,
		hashService:    hashService,
		gameRepo:       gameRepo,
		roomRepo:       roomRepo,
		userRepo:       userRepo,
		roomMemberRepo: roomMemberRepo,
	}
}

func (s *gameService) CreateRoom(ctx context.Context, req *dto.CreateRoomReq) (*dto.CreateRoomRes, error) {
	roomCode := s.hashService.GenerateRandomString(10)
	room := &entity.Room{
		RoomCode:   roomCode,
		MaxPlayers: *req.MaxPlayers,
	}

	if _, err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}

	return &dto.CreateRoomRes{
		Room: *room,
	}, nil
}

// --- Socket Event Listener Registration ---
func (s *gameService) GameSocketListener(ctx context.Context) error {
	defer common.Recovery()

	socketService := s.appCtx.GetSocketService()

	socketService.OnConnect(func(conn socketio.Conn) error {
		log.Println("Client connected")
		return nil
	})

	socketService.OnDisconnect(func(conn socketio.Conn, reason string) {
		log.Println("Client disconnected: ", reason)
		return
	})

	socketService.OnEvent(entity.GameEventCreateRoom, func(conn socketio.Conn, payload interface{}) {
		defer common.Recovery()
		req := &dto.CreateRoomReq{}
		err := mapstructure.Decode(payload, req)
		if err != nil {
			conn.Emit(entity.GameEventReply, dto.Reply{Ok: false, Ref: req.Ref, Error: err.Error()})
			return
		}

		room, err := s.CreateRoom(ctx, req)
		if err != nil {
			conn.Emit(entity.GameEventReply, dto.Reply{Ok: false, Ref: req.Ref, Error: err.Error()})
			return
		}

		conn.Join(room.RoomCode)
		conn.Emit(entity.GameEventReply, dto.Reply{Ok: true, Ref: req.Ref})
	})

	socketService.OnEvent(entity.GameEventJoinRoom, func(conn socketio.Conn, payload interface{}) {
		defer common.Recovery()
		req := &dto.JoinRoomReq{}
		err := mapstructure.Decode(payload, req)
		if err != nil {
			conn.Emit(entity.GameEventReply, dto.Reply{Ok: false, Ref: req.Ref, Error: err.Error()})
			return
		}

		conn.Join(req.RoomCode)
		conn.Emit(entity.GameEventReply, dto.Reply{Ok: true, Ref: req.Ref})
	})

	socketService.OnEvent(entity.GameEventLeaveRoom, func(conn socketio.Conn, payload interface{}) {
		defer common.Recovery()
		req := &dto.LeaveRoomReq{}
		err := mapstructure.Decode(payload, req)
		if err != nil {
			conn.Emit(entity.GameEventReply, dto.Reply{Ok: false, Ref: req.Ref, Error: err.Error()})
			return
		}

		conn.Leave(req.RoomCode)
		conn.Emit(entity.GameEventReply, dto.Reply{Ok: true, Ref: req.Ref})
	})

	socketService.OnEvent(entity.GameEventMove, func(conn socketio.Conn, payload interface{}) {
		defer common.Recovery()
		req := &dto.MoveReq{}
		err := mapstructure.Decode(payload, req)
		if err != nil {
			conn.Emit(entity.GameEventReply, dto.Reply{Ok: false, Ref: req.Ref, Error: err.Error()})
			return
		}

		conn.Emit(entity.GameEventMove, req)
		conn.Emit(entity.GameEventReply, dto.Reply{Ok: true, Ref: req.Ref})
	})

	return nil
}
