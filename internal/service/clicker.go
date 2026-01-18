package service

import (
	"sync"
	"time"

	"myproject/internal/domain"
	"myproject/internal/input"

	"github.com/go-vgo/robotgo"
)

type ClickerService struct {
	isRunning      bool
	stopChan       chan struct{}
	wg             sync.WaitGroup
	mu             sync.Mutex
	actions        []domain.ClickAction
	hotkeyEnabled  bool
	onStopCallback func()
}

func NewClickerService() *ClickerService {
	return &ClickerService{
		isRunning: false,
		stopChan:  make(chan struct{}),
		actions:   make([]domain.ClickAction, 0),
	}
}

func (c *ClickerService) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isRunning || len(c.actions) == 0 {
		return
	}

	c.isRunning = true
	c.stopChan = make(chan struct{})

	for _, action := range c.actions {
		c.wg.Add(1)
		go c.clickWorker(action)
	}

	go func() {
		c.wg.Wait()
		c.setStopped()
	}()
}

func (c *ClickerService) clickWorker(action domain.ClickAction) {
	defer c.wg.Done()

	delay := action.Delay
	if delay <= 0 {
		delay = 100
	}
	ticker := time.NewTicker(time.Duration(delay) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			robotgo.Move(action.X, action.Y)
			robotgo.Click("left")
		}
	}
}

func (c *ClickerService) setStopped() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.isRunning = false
}

func (c *ClickerService) Stop() {
	c.mu.Lock()
	if c.isRunning {
		close(c.stopChan)
		c.isRunning = false
		callback := c.onStopCallback
		c.mu.Unlock()

		if callback != nil {
			callback()
		}
		return
	}
	c.mu.Unlock()
}

func (c *ClickerService) SetActions(actions []domain.ClickAction) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.actions = actions
}

func (c *ClickerService) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isRunning
}

func (c *ClickerService) SetStopCallback(callback func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onStopCallback = callback
}

func (c *ClickerService) StartGlobalHotkey() {
	c.mu.Lock()
	if c.hotkeyEnabled {
		c.mu.Unlock()
		return
	}
	c.hotkeyEnabled = true
	c.mu.Unlock()

	input.GetHookManager().EnableHotkey(func() {
		if c.IsRunning() {
			c.Stop()
		}
	})
}

func (c *ClickerService) StopGlobalHotkey() {
	c.mu.Lock()
	if c.hotkeyEnabled {
		c.hotkeyEnabled = false
		c.mu.Unlock()
		input.GetHookManager().DisableHotkey()
		return
	}
	c.mu.Unlock()
}
