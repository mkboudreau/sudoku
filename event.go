package main

import (
	"fmt"
	"log"
	"time"
)

type EventHandler interface {
	OnAttemptingCoord(board Board, coord XY)

	OnBeforeClearCoord(board Board, coord XY)
	OnAfterClearCoord(board Board, coord XY)

	OnSuccessfulCoord(board Board, coord XY)
	OnFailedCoord(board Board, coord XY)
}

type LogEventHandler struct {
}

func NewLogEventHandler() *LogEventHandler {
	return &LogEventHandler{}
}
func (e *LogEventHandler) OnAttemptingCoord(board Board, coord XY) {
	log.Printf("Attempting Coordinate %+v", coord)
}
func (e *LogEventHandler) OnBeforeClearCoord(board Board, coord XY) {
	log.Printf("Before Clearing Coordinate %+v", coord)
}
func (e *LogEventHandler) OnAfterClearCoord(board Board, coord XY) {
	log.Printf("After Clearing Coordinate %+v", coord)
}
func (e *LogEventHandler) OnSuccessfulCoord(board Board, coord XY) {
	log.Printf("Successful Coordinate %+v", coord)
}
func (e *LogEventHandler) OnFailedCoord(board Board, coord XY) {
	log.Printf("Failed Coordinate %+v", coord)
}

type ProgressEventHandler struct {
	ShowColors       bool
	Delay            time.Duration
	goodCoords       []*Coord
	failedCoords     []*Coord
	inProgressCoords []*Coord
}

func NewProgressEventHandler(colors bool, delay time.Duration) *ProgressEventHandler {
	return &ProgressEventHandler{ShowColors: colors, Delay: delay}
}

func (e *ProgressEventHandler) OnAttemptingCoord(board Board, coord XY) {
	e.delay()
	e.inProgressCoords = append(e.inProgressCoords, CoordXY(coord))
	inprogressCodeSet := NewColorCoordSet(e.inProgressCoords, StatusInProgressColor)
	successCodeSet := NewColorCoordSet(e.goodCoords, StatusSuccessColor)

	e.clearConsole()
	fmt.Printf("Attempting Coordinate %v\n%v", coord, ColorizeBoard(board, inprogressCodeSet, successCodeSet))
}
func (e *ProgressEventHandler) OnBeforeClearCoord(board Board, coord XY) {}
func (e *ProgressEventHandler) OnAfterClearCoord(board Board, coord XY)  {}
func (e *ProgressEventHandler) OnSuccessfulCoord(board Board, coord XY) {
	e.delay()
	e.removeInProgressCoord(coord)
	e.goodCoords = append(e.goodCoords, CoordXY(coord))
	//failedCodeSet := NewColorCoordSet([]*Coord{CoordXY(coord)}, StatusFailedColor)
	inprogressCodeSet := NewColorCoordSet(e.inProgressCoords, StatusInProgressColor)
	successCodeSet := NewColorCoordSet(e.goodCoords, StatusSuccessColor)

	e.clearConsole()
	fmt.Printf("Successful Coordinate %v\n%v", coord, ColorizeBoard(board, inprogressCodeSet, successCodeSet))
}
func (e *ProgressEventHandler) OnFailedCoord(board Board, coord XY) {
	e.delay()
	e.removeInProgressCoord(coord)
	failedCodeSet := NewColorCoordSet([]*Coord{CoordXY(coord)}, StatusFailedColor)
	inprogressCodeSet := NewColorCoordSet(e.inProgressCoords, StatusInProgressColor)
	successCodeSet := NewColorCoordSet(e.goodCoords, StatusSuccessColor)

	e.clearConsole()
	fmt.Printf("Failed Coordinate %v\n%v", coord, ColorizeBoard(board, inprogressCodeSet, successCodeSet, failedCodeSet))
}

func (e *ProgressEventHandler) clearConsole() {
	fmt.Printf("%v%v", ClearConsole, ResetCursor)
}
func (e *ProgressEventHandler) delay() {
	time.Sleep(e.Delay)
}
func (e *ProgressEventHandler) removeInProgressCoord(coord XY) {
	var index int
	for i, p := range e.inProgressCoords {
		if p.x == coord.X() && p.y == coord.Y() {
			index = i
		}
	}

	if index > 0 {
		e.inProgressCoords = append(e.inProgressCoords[:index], e.inProgressCoords[index+1:]...)
	}
}

type MultiEventHandler struct {
	handlers []EventHandler
}

func NewMultiEventHandler(handlers ...EventHandler) *MultiEventHandler {
	return &MultiEventHandler{handlers: handlers}
}
func (e *MultiEventHandler) OnAttemptingCoord(board Board, coord XY) {
	for _, h := range e.handlers {
		h.OnAttemptingCoord(board, coord)
	}
}
func (e *MultiEventHandler) OnBeforeClearCoord(board Board, coord XY) {
	for _, h := range e.handlers {
		h.OnBeforeClearCoord(board, coord)
	}
}
func (e *MultiEventHandler) OnAfterClearCoord(board Board, coord XY) {
	for _, h := range e.handlers {
		h.OnAfterClearCoord(board, coord)
	}
}
func (e *MultiEventHandler) OnSuccessfulCoord(board Board, coord XY) {
	for _, h := range e.handlers {
		h.OnSuccessfulCoord(board, coord)
	}
}
func (e *MultiEventHandler) OnFailedCoord(board Board, coord XY) {
	for _, h := range e.handlers {
		h.OnFailedCoord(board, coord)
	}
}

type NoEventHandler struct{}

func (e *NoEventHandler) OnAttemptingCoord(board Board, coord XY)  {}
func (e *NoEventHandler) OnBeforeClearCoord(board Board, coord XY) {}
func (e *NoEventHandler) OnAfterClearCoord(board Board, coord XY)  {}
func (e *NoEventHandler) OnSuccessfulCoord(board Board, coord XY)  {}
func (e *NoEventHandler) OnFailedCoord(board Board, coord XY)      {}
