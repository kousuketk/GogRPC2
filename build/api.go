package build

import (
	"github.com/kousuketk/GogRPC2/game"
	"github.com/kousuketk/GogRPC2/gen/api"
)

func PBRoom(r *game.Room) *api.Room {
	return &api.Room{
		Id:    r.ID,
		Host:  PBPlayer(r.Host),
		Guest: PBPlayer(r.Guest),
	}
}

func PBPlayer(p *game.Player) *api.Player {
	if p == nil {
		return nil
	}
	return &api.Player{
		Id:    p.ID,
		Color: PBColor(p.Color),
	}
}

func PBColor(c game.Color) api.Color {
	switch c {
	case game.Black:
		return api.Color_BLACK
	case game.White:
		return api.Color_WHITE
	case game.Empty:
		return api.Color_EMPTY
	case game.Wall:
		return api.Color_WALL
	}

	return api.Color_UNKNOWN
}

func PBBoard(b *game.Board) *api.Board {
	pbCols := make([]*api.Board_Col, 0, 10)

	for _, col := range b.Cells {
		pbCells := make([]api.Color, 0, 10)
		for _, c := range col {
			pbCells = append(pbCells, PBColor(c))
		}
		pbCols = append(pbCols, &api.Board_Col{
			Cells: pbCells,
		})
	}

	return &api.Board{
		Cols: pbCols,
	}
}
