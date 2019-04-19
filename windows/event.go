package windows

import (
	"github.com/suiqirui1987/fly3d/tools/goevent"
)

const (
	Resize           = "resize"
	Keydown          = "keydown"
	Keyup            = "keyup"
	MouseDown        = "mousedown"
	MouseUp          = "mouseup"
	MouseMove        = "mousemove"
	MouseOut         = "mouseout"
	Wheel            = "wheel"
	Focus            = "focus"
	FullScreenChange = "fullscreenchange"
)

type ResizeEvent struct {
	goevent.Event

	WindowWidth  int
	WindowHeigth int
}
type MouseEvent struct {
	goevent.Event

	ClientX float32
	ClientY float32
}

type WheelEvent struct {
	goevent.Event

	DeltaX float32
	DeltaY float32
}

type KeyboardEvent struct {
	goevent.Event

	KeyCode  int
	CharCode string
}
type FocusEvent struct {
	goevent.Event

	Focused bool
}
