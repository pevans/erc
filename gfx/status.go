package gfx

// StatusOverlay is the global overlay instance for status messages. Only one
// overlay can be shown at a time.
var StatusOverlay *Overlay

func init() {
	StatusOverlay = NewOverlay()
}

// ShowStatus is a convenience function to show a status overlay with fade.
func ShowStatus(data []byte) error {
	return StatusOverlay.Show(data)
}

// ShowPersistentStatus shows a status overlay without fade.
func ShowPersistentStatus(data []byte) error {
	return StatusOverlay.ShowPersistent(data)
}

// HideStatus hides the current status overlay.
func HideStatus() {
	StatusOverlay.Hide()
}
