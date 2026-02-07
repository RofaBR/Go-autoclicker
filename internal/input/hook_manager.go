package input

import (
	"sync"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

type HookManager struct {
	mu         sync.Mutex
	started    bool
	recording  bool
	recordChan chan Coordinates
	hotkeys    map[uint16]*hotkeyEntry
}

type hotkeyEntry struct {
	enabled  bool
	callback func()
}

var (
	instance *HookManager
	once     sync.Once
)

func GetHookManager() *HookManager {
	once.Do(func() {
		instance = &HookManager{
			hotkeys: make(map[uint16]*hotkeyEntry),
		}
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
	evChan := hook.Start()
	defer hook.End()

	for e := range evChan {
		if e.Kind == hook.KeyDown {
			h.mu.Lock()
			entry, found := h.hotkeys[e.Keycode]
			var action func()
			if found && entry.enabled {
				action = entry.callback
			}
			h.mu.Unlock()

			if action != nil {
				go action()
			}
		}
		if e.Kind == hook.MouseDown {
			h.mu.Lock()
			if h.recording && e.Button == hook.MouseMap["left"] {
				x, y := robotgo.Location()
				ch := h.recordChan
				h.recording = false
				h.mu.Unlock()

				if ch != nil {
					ch <- Coordinates{X: x, Y: y}
				}
				continue
			}
			h.mu.Unlock()
		}
	}
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

func (h *HookManager) RegisterHotkey(key string, callback func()) {
	h.mu.Lock()
	defer h.mu.Unlock()
	kc := hook.Keycode[key]
	h.hotkeys[kc] = &hotkeyEntry{
		enabled:  false,
		callback: callback,
	}
}

func (h *HookManager) EnableHotkey(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	kc := hook.Keycode[key]
	if _, ok := h.hotkeys[kc]; ok {
		h.hotkeys[kc].enabled = true
	}
}
func (h *HookManager) EnableHotkeys() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, entry := range h.hotkeys {
		entry.enabled = true
	}
}

func (h *HookManager) DisableHotkey(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	kc := hook.Keycode[key]
	if _, ok := h.hotkeys[kc]; ok {
		h.hotkeys[kc].enabled = false
	}
}

func (h *HookManager) DisableHotkeys() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, entry := range h.hotkeys {
		entry.enabled = false
	}
}

func (h *HookManager) IsRecording() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.recording
}
