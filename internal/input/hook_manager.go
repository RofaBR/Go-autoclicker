package input

import (
	"sync"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

type HookManager struct {
	mu             sync.Mutex
	started        bool
	recording      bool
	recordChan     chan Coordinates
	hotkeyEnabled  bool
	hotkeyCallback func()
}

var (
	instance *HookManager
	once     sync.Once
)

func GetHookManager() *HookManager {
	once.Do(func() {
		instance = &HookManager{}
	})
	return instance
}

func (h *HookManager) Start() {
	h.mu.Lock()
	if h.started {
		h.mu.Unlock()
		return
	}
	h.started = true
	h.mu.Unlock()

	go h.runEventLoop()
}

func (h *HookManager) runEventLoop() {
	hook.Register(hook.KeyDown, []string{"f10"}, func(e hook.Event) {
		h.mu.Lock()
		enabled := h.hotkeyEnabled
		callback := h.hotkeyCallback
		h.mu.Unlock()

		if enabled && callback != nil {
			callback()
		}
	})

	hook.Register(hook.MouseDown, []string{}, func(e hook.Event) {
		h.mu.Lock()
		if h.recording && e.Button == hook.MouseMap["left"] {
			x, y := robotgo.Location()
			ch := h.recordChan
			h.recording = false
			h.mu.Unlock()

			if ch != nil {
				ch <- Coordinates{X: x, Y: y}
			}
			return
		}
		h.mu.Unlock()
	})

	s := hook.Start()
	<-hook.Process(s)
}

func (h *HookManager) StartRecording() (int, int) {
	h.mu.Lock()
	h.recordChan = make(chan Coordinates, 1)
	h.recording = true
	h.mu.Unlock()

	result := <-h.recordChan

	h.mu.Lock()
	h.recordChan = nil
	h.mu.Unlock()

	return result.X, result.Y
}

func (h *HookManager) EnableHotkey(callback func()) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.hotkeyCallback = callback
	h.hotkeyEnabled = true
}

func (h *HookManager) DisableHotkey() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.hotkeyEnabled = false
}

func (h *HookManager) IsRecording() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.recording
}
