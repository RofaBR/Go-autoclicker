package main

import (
	"context"
	"myproject/internal/domain"
	"myproject/internal/input"
	"myproject/internal/service"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx      context.Context
	points   []domain.ClickAction
	nextID   int
	mu       sync.Mutex
	recorder *input.RecorderService
	clicker  *service.ClickerService
}

func NewApp() *App {
	return &App{
		points:   []domain.ClickAction{},
		nextID:   1,
		recorder: input.NewRecorderService(),
		clicker:  service.NewClickerService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.clicker.SetStopCallback(func() {
		runtime.EventsEmit(a.ctx, "clicker:stopped")
	})
}

func (a *App) shutdown(ctx context.Context) {
	a.clicker.StopGlobalHotkey()
}

func (a *App) AddPoint() domain.ClickAction {
	a.mu.Lock()
	defer a.mu.Unlock()

	point := domain.ClickAction{
		ID:    a.nextID,
		X:     0,
		Y:     0,
		Delay: 1000,
	}
	a.nextID++
	a.points = append(a.points, point)
	return point
}

func (a *App) GetPoints() []domain.ClickAction {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.points
}

func (a *App) UpdatePointCoordinates(id int, x int, y int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i := range a.points {
		if a.points[i].ID == id {
			a.points[i].X = x
			a.points[i].Y = y
			return true
		}
	}
	return false
}

func (a *App) UpdatePointDelay(id int, delay int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i := range a.points {
		if a.points[i].ID == id {
			a.points[i].Delay = delay
			return true
		}
	}
	return false
}

func (a *App) RemovePoint(id int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i, point := range a.points {
		if point.ID == id {
			a.points = append(a.points[:i], a.points[i+1:]...)
			return true
		}
	}
	return false
}

func (a *App) RecordCoordinates(id int) bool {
	runtime.WindowHide(a.ctx)

	x, y := a.recorder.StartRecording()
	success := a.UpdatePointCoordinates(id, x, y)

	runtime.WindowShow(a.ctx)

	return success
}

func (a *App) StartClicker() {
	a.mu.Lock()

	currentPoints := make([]domain.ClickAction, len(a.points))
	copy(currentPoints, a.points)
	a.mu.Unlock()
	a.clicker.SetActions(currentPoints)
	a.clicker.Start()

	a.clicker.StartGlobalHotkey()
}

func (a *App) StopClicker() {
	a.clicker.Stop()

	a.clicker.StopGlobalHotkey()
}
