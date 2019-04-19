package windows

import (
	"github.com/suiqirui1987/fly3d/glfw"
)

func KeyString(key glfw.Key) string {
	switch key {
	// Printable keys.
	case glfw.KeyA:
		return "A"
	case glfw.KeyB:
		return "B"
	case glfw.KeyC:
		return "C"
	case glfw.KeyD:
		return "D"
	case glfw.KeyE:
		return "E"
	case glfw.KeyF:
		return "F"
	case glfw.KeyG:
		return "G"
	case glfw.KeyH:
		return "H"
	case glfw.KeyI:
		return "I"
	case glfw.KeyJ:
		return "J"
	case glfw.KeyK:
		return "K"
	case glfw.KeyL:
		return "L"
	case glfw.KeyM:
		return "M"
	case glfw.KeyN:
		return "N"
	case glfw.KeyO:
		return "O"
	case glfw.KeyP:
		return "P"
	case glfw.KeyQ:
		return "Q"
	case glfw.KeyR:
		return "R"
	case glfw.KeyS:
		return "S"
	case glfw.KeyT:
		return "T"
	case glfw.KeyU:
		return "U"
	case glfw.KeyV:
		return "V"
	case glfw.KeyW:
		return "W"
	case glfw.KeyX:
		return "X"
	case glfw.KeyY:
		return "Y"
	case glfw.KeyZ:
		return "Z"
	case glfw.Key1:
		return "1"
	case glfw.Key2:
		return "2"
	case glfw.Key3:
		return "3"
	case glfw.Key4:
		return "4"
	case glfw.Key5:
		return "5"
	case glfw.Key6:
		return "6"
	case glfw.Key7:
		return "7"
	case glfw.Key8:
		return "8"
	case glfw.Key9:
		return "9"
	case glfw.Key0:
		return "0"
	case glfw.KeySpace:
		return "SPACE"
	case glfw.KeyMinus:
		return "MINUS"
	case glfw.KeyEqual:
		return "EQUAL"
	case glfw.KeyLeftBracket:
		return "LEFT BRACKET"
	case glfw.KeyRightBracket:
		return "RIGHT BRACKET"
	case glfw.KeyBackslash:
		return "BACKSLASH"
	case glfw.KeySemicolon:
		return "SEMICOLON"
	case glfw.KeyApostrophe:
		return "APOSTROPHE"
	case glfw.KeyGraveAccent:
		return "GRAVE ACCENT"
	case glfw.KeyComma:
		return "COMMA"
	case glfw.KeyPeriod:
		return "PERIOD"
	case glfw.KeySlash:
		return "SLASH"
	case glfw.KeyWorld1:
		return "WORLD 1"
	case glfw.KeyWorld2:
		return "WORLD 2"
	// Function keys.
	case glfw.KeyEscape:
		return "ESCAPE"
	case glfw.KeyF1:
		return "F1"
	case glfw.KeyF2:
		return "F2"
	case glfw.KeyF3:
		return "F3"
	case glfw.KeyF4:
		return "F4"
	case glfw.KeyF5:
		return "F5"
	case glfw.KeyF6:
		return "F6"
	case glfw.KeyF7:
		return "F7"
	case glfw.KeyF8:
		return "F8"
	case glfw.KeyF9:
		return "F9"
	case glfw.KeyF10:
		return "F10"
	case glfw.KeyF11:
		return "F11"
	case glfw.KeyF12:
		return "F12"
	case glfw.KeyF13:
		return "F13"
	case glfw.KeyF14:
		return "F14"
	case glfw.KeyF15:
		return "F15"
	case glfw.KeyF16:
		return "F16"
	case glfw.KeyF17:
		return "F17"
	case glfw.KeyF18:
		return "F18"
	case glfw.KeyF19:
		return "F19"
	case glfw.KeyF20:
		return "F20"
	case glfw.KeyF21:
		return "F21"
	case glfw.KeyF22:
		return "F22"
	case glfw.KeyF23:
		return "F23"
	case glfw.KeyF24:
		return "F24"
	case glfw.KeyF25:
		return "F25"
	case glfw.KeyUp:
		return "UP"
	case glfw.KeyDown:
		return "DOWN"
	case glfw.KeyLeft:
		return "LEFT"
	case glfw.KeyRight:
		return "RIGHT"
	case glfw.KeyLeftShift:
		return "LEFT SHIFT"
	case glfw.KeyRightShift:
		return "RIGHT SHIFT"
	case glfw.KeyLeftControl:
		return "LEFT CONTROL"
	case glfw.KeyRightControl:
		return "RIGHT CONTROL"
	case glfw.KeyLeftAlt:
		return "LEFT ALT"
	case glfw.KeyRightAlt:
		return "RIGHT ALT"
	case glfw.KeyTab:
		return "TAB"
	case glfw.KeyEnter:
		return "ENTER"
	case glfw.KeyBackspace:
		return "BACKSPACE"
	case glfw.KeyInsert:
		return "INSERT"
	case glfw.KeyDelete:
		return "DELETE"
	case glfw.KeyPageUp:
		return "PAGE UP"
	case glfw.KeyPageDown:
		return "PAGE DOWN"
	case glfw.KeyHome:
		return "HOME"
	case glfw.KeyEnd:
		return "END"
	case glfw.KeyKP0:
		return "KEYPAD 0"
	case glfw.KeyKP1:
		return "KEYPAD 1"
	case glfw.KeyKP2:
		return "KEYPAD 2"
	case glfw.KeyKP3:
		return "KEYPAD 3"
	case glfw.KeyKP4:
		return "KEYPAD 4"
	case glfw.KeyKP5:
		return "KEYPAD 5"
	case glfw.KeyKP6:
		return "KEYPAD 6"
	case glfw.KeyKP7:
		return "KEYPAD 7"
	case glfw.KeyKP8:
		return "KEYPAD 8"
	case glfw.KeyKP9:
		return "KEYPAD 9"
	case glfw.KeyKPDivide:
		return "KEYPAD DIVIDE"
	case glfw.KeyKPMultiply:
		return "KEYPAD MULTPLY"
	case glfw.KeyKPSubtract:
		return "KEYPAD SUBTRACT"
	case glfw.KeyKPAdd:
		return "KEYPAD ADD"
	case glfw.KeyKPDecimal:
		return "KEYPAD DECIMAL"
	case glfw.KeyKPEqual:
		return "KEYPAD EQUAL"
	case glfw.KeyKPEnter:
		return "KEYPAD ENTER"
	case glfw.KeyPrintScreen:
		return "PRINT SCREEN"
	case glfw.KeyNumLock:
		return "NUM LOCK"
	case glfw.KeyCapsLock:
		return "CAPS LOCK"
	case glfw.KeyScrollLock:
		return "SCROLL LOCK"
	case glfw.KeyPause:
		return "PAUSE"
	case glfw.KeyLeftSuper:
		return "LEFT SUPER"
	case glfw.KeyRightSuper:
		return "RIGHT SUPER"
	case glfw.KeyMenu:
		return "MENU"
	default:
		return "UNKNOWN"
	}
}

type IWindow interface {
	GetRenderWidth() int
	GetRenderHeight() int

	StopNewFrame()
	QueueNewFrame(func())
	GetWindowDevicePixelRatio() float32

	GetFullscreen() bool
	ExitFullscreen()
	RequestFullscreen()

	On(name string, fn interface{}) error
	Emit(name string, params ...interface{}) error
	Has(name string) bool
	List() []string
	Remove(names ...string)
}
