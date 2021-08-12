package handler

import (
	"fmt"
	"sync"

	"github.com/kousuketk/GogRPC2/build"
	"github.com/kousuketk/GogRPC2/game"
	"github.com/kousuketk/GogRPC2/gen/api"
)

type GameHandler struct {
	sync.RWMutex
	games  map[int32]*game.Game
	client map[int32][]api.GameService_PlayServer
}

func NewGameHnadler() *GameHandler {
	return &GameHandler{
		games:  make(map[int32]*game.Game),
		client: make(map[int32][]api.GameService_PlayServer),
	}
}

func (h *GameHandler) Play(stream api.GameService_PlayServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		roomID := req.GetRoomId()
		player := build.Player(req.GetPlayer())

		switch req.GetAction().(type) {
		case *api.PlayRequest_Start:
			err := h.start(stream, roomID, player)
			if err != nil {
				return err
			}
		case *api.PlayRequest_Move:
			action := req.GetMove()
			x := action.GetMove().GetX()
			y := action.GetMove().GetY()
			err := h.move(roomID, x, y, player)
			if err != nil {
				return err
			}
		}
	}
}

func (h *GameHandler) start(stream api.GameService_PlayServer, roomID int32, me *game.Player) error {
	h.Lock()
	defer h.Unlock()

	g := h.games[roomID]
	if g == nil {
		g = game.NewGame(game.None)
		h.games[roomID] = g
		h.client[roomID] = make([]api.GameService_PlayServer, 0, 2)
	}

	h.client[roomID] = append(h.client[roomID], stream)

	if len(h.client[roomID]) == 2 {
		for _, s := range h.client[roomID] {
			err := s.Send(&api.PlayResponse{
				Event: &api.PlayResponse_Ready{
					Ready: &api.PlayResponse_ReadyEvent{},
				},
			})
			if err != nil {
				return err
			}
		}
		fmt.Printf("game has started room_id=%v\n", roomID)
	} else {
		err := stream.Send(&api.PlayResponse{
			Event: &api.PlayResponse_Waiting{
				Waiting: &api.PlayResponse_WaitingEvent{},
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *GameHandler) move(roomID int32, x int32, y int32, p *game.Player) error {
	h.Lock()
	defer h.Unlock()

	g := h.games[roomID]
	finished, err := g.Move(x, y, p.Color)
	if err != nil {
		return err
	}

	for _, s := range h.client[roomID] {
		err := s.Send(&api.PlayResponse{
			Event: &api.PlayResponse_Move{
				Move: &api.PlayResponse_MoveEvent{
					Player: build.PBPlayer(p),
					Move: &api.Move{
						X: x,
						Y: y,
					},
					Board: build.PBBoard(g.Board),
				},
			},
		})
		if err != nil {
			return err
		}
		if finished {
			err := s.Send(
				&api.PlayResponse{
					Event: &api.PlayResponse_Finished{
						Finished: &api.PlayResponse_FinishedEvent{
							Winner: build.PBColor(g.Winner()),
							Board:  build.PBBoard(g.Board),
						},
					},
				},
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
