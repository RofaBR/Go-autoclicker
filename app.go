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
		points:   make([]domain.ClickAction, 0),
		nextID:   1,
		recorder: input.NewRecorderService(),
		clicker:  service.NewClickerService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	input.GetHookManager().Start()

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
		Delay: 5000,
	}
	a.points = append(a.points, point)
	a.nextID++
	return point
}

func (a *App) GetPoints() []domain.ClickAction {
	a.mu.Lock()
	defer a.mu.Unlock()

	result := make([]domain.ClickAction, len(a.points))
	copy(result, a.points)
	return result
}

func (a *App) findPointIndex(id int) int {
	for i, p := range a.points {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func (a *App) UpdatePointCoordinates(id int, x int, y int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	idx := a.findPointIndex(id)
	if idx == -1 {
		return false
	}
	a.points[idx].X = x
	a.points[idx].Y = y
	return true
}

func (a *App) UpdatePointDelay(id int, delay int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	idx := a.findPointIndex(id)
	if idx == -1 {
		return false
	}
	a.points[idx].Delay = delay
	return true
}

func (a *App) RemovePoint(id int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	idx := a.findPointIndex(id)
	if idx == -1 {
		return false
	}
	a.points = append(a.points[:idx], a.points[idx+1:]...)
	return true
}

func (a *App) RecordCoordinates(id int) bool {
	//runtime.WindowHide(a.ctx)
	//runtime.WindowShow(a.ctx)
	x, y := a.recorder.StartRecording()
	return a.UpdatePointCoordinates(id, x, y)
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
