package service

import (
	"sync"
	"time"

	"myproject/internal/domain"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

type ClickerService struct {
	isRunning      bool
	stopChan       chan struct{}
	wg             sync.WaitGroup
	mu             sync.Mutex
	actions        []domain.ClickAction
	hookRunning    bool
	hookStopChan   chan struct{}
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

	for {
		select {
		case <-c.stopChan:
			return
		default:
			robotgo.Move(action.X, action.Y, action.Delay)
			robotgo.Click("left")
			time.Sleep(time.Duration(action.Delay) * time.Millisecond)
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
	if c.hookRunning {
		c.mu.Unlock()
		return
	}
	c.hookRunning = true
	c.hookStopChan = make(chan struct{})
	c.mu.Unlock()

	go c.globalKeyListener()
}

func (c *ClickerService) StopGlobalHotkey() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.hookRunning {
		close(c.hookStopChan)
		c.hookRunning = false
		hook.End()
	}
}

func (c *ClickerService) globalKeyListener() {
	hook.Register(hook.KeyDown, []string{"f10"}, func(e hook.Event) {
		if c.IsRunning() {
			c.Stop()
		}
	})

	s := hook.Start()

	eventChan := hook.Process(s)

	select {
	case <-c.hookStopChan:
		hook.End()
		return
	case <-eventChan:
		return
	}
}
