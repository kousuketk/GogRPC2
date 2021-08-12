package build

import (
	"fmt"

	"github.com/kousuketk/GogRPC2/game"
	"github.com/kousuketk/GogRPC2/gen/api"
)

func Room(r *api.Room) *game.Room {
	return &game.Room{
		ID:    r.GetId(),
		Host:  Player(r.GetHost()),
		Guest: Player(r.GetGuest()),
	}
}

func Player(p *api.Player) *game.Player {
	return &game.Player{
		ID:    p.GetId(),
		Color: Color(p.GetColor()),
	}
}

func Color(c api.Color) game.Color {
	switch c {
	case api.Color_BLACK:
		return game.Black
	case api.Color_WHITE:
		return game.White
	case api.Color_EMPTY:
		return game.Empty
	case api.Color_WALL:
		return game.Wall
	}

	panic(fmt.Sprintf("unkown color=%v", c))
}
