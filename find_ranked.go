package main

import (
	"sort"
)

type RankedCoordFinder struct { // slow implementation. using only for visualiza
	board     Board
	remaining []*Coord
}

func (f *RankedCoordFinder) NextOpenCoordinate(board Board, coord XY) (*Coord, bool) {
	f.board = board
	f.RefreshCoordinates()

	if len(f.remaining) == 0 {
		return nil, false
	}

	c := f.remaining[0]
	f.remaining = f.remaining[0:]

	return c, true
}

func (f *RankedCoordFinder) RefreshCoordinates() {
	f.remaining = []*Coord{}
	for x := 0; x < len(f.board); x++ {
		for y := 0; y < len(f.board); y++ {
			if f.board[x][y] == 0 {
				f.remaining = append(f.remaining, &Coord{x, y})
			}
		}
	}
	sort.Sort(f)
}

func (f *RankedCoordFinder) Len() int {
	return len(f.remaining)
}
func (f *RankedCoordFinder) Less(i, j int) bool {
	return len(f.board.AvailableValuesAtCoordinate(f.remaining[i])) <
		len(f.board.AvailableValuesAtCoordinate(f.remaining[j]))
}
func (f *RankedCoordFinder) Swap(i, j int) {
	f.remaining[i], f.remaining[j] = f.remaining[j], f.remaining[i]
}
