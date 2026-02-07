package service

import "myproject/internal/input"

const stopHotkeyKey = "f10"
const changeVisibilityKey = "f9"

func (c *ClickerService) RegisterHotkeys() {
	input.GetHookManager().RegisterHotkey(stopHotkeyKey, func() {
		if c.IsRunning() {
			c.Stop()
		}
	})
	input.GetHookManager().RegisterHotkey(changeVisibilityKey, func() {
		c.mu.Lock()
		callback := c.onVisibilityCallback
		c.mu.Unlock()
		if callback != nil {
			callback()
		}
	})
}

func (c *ClickerService) StartGlobalHotkey() {
	c.mu.Lock()
	if c.hotkeyEnabled {
		c.mu.Unlock()
		return
	}
	c.hotkeyEnabled = true
	c.mu.Unlock()

	input.GetHookManager().EnableHotkey(stopHotkeyKey)
}

func (c *ClickerService) StopGlobalHotkey() {
	c.mu.Lock()
	if c.hotkeyEnabled {
		c.hotkeyEnabled = false
		c.mu.Unlock()
		input.GetHookManager().DisableHotkey(stopHotkeyKey)
		return
	}
	c.mu.Unlock()
}

func (c *ClickerService) EnableVisibilityHotkey() {
	input.GetHookManager().EnableHotkey(changeVisibilityKey)
}
