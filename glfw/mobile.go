//+build android darwin,arm darwin,arm64 ios
//+build !mobilebind

package glfw

import (
	"fmt"
	"runtime"

	"github.com/gopherjs/gopherjs/js"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

func init() {
	runtime.LockOSThread()
}

var contextWatcher ContextWatcher

func Init(cw ContextWatcher) error {
	contextWatcher = cw
	return nil
}

func Terminate() error {
	return nil
}

func CreateWindow(_, _ int, title string, monitor *Monitor, share *Window) (*Window, error) {

	w := &Window{
		requestFullscreen: true,
		fullscreen:        true,
	}

	go runLoop(w)

	return w, nil

}

func runLoop(w *Window) {

	var glctx gl.Context
	app.Main(func(a app.App) {

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:

					glctx = e.DrawContext.(gl.Context)

					w.context = glctx

					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					break
				}
				break
			case size.Event:
				windowWidth := int(e.WidthPx)
				windowHeight := int(e.HeightPx)

				if w.framebufferSizeCallback != nil {
					// TODO: Callbacks may be blocking so they need to happen asyncronously. However,
					//       GLFW API promises the callbacks will occur from one thread (i.e., sequentially), so may want to do that.
					go w.framebufferSizeCallback(w, windowWidth, windowHeight)
				}
				if w.sizeCallback != nil {
					go w.sizeCallback(w, windowWidth, windowHeight)
				}

				break
			case key.Event:
				action := Press
				if e.Direction == key.DirPress {
					action = Repeat
				}

				key := toKey(e)

				// Extend slice if needed.
				neededSize := int(key) + 1
				if neededSize > len(w.keys) {
					w.keys = append(w.keys, make([]Action, neededSize-len(w.keys))...)
				}
				w.keys[key] = action

				if w.keyCallback != nil {
					mods := toModifierKey(e)

					go w.keyCallback(w, key, -1, action, mods)
				}

				break
			case paint.Event:
				if e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				a.Publish() // same as SwapBuffers

				// Drive the animation by preparing to paint the next frame
				// after this one is shown. - FPS is ignored here!
				a.Send(paint.Event{})

				break
			case touch.Event:

				switch e.Type {
				case touch.TypeBegin:
					if w.mouseButtonCallback != nil {
						go w.mouseButtonCallback(w, MouseButton(0), Press, 0)
					}
				case touch.TypeMove:

					var movementX, movementY float64
					movementX = float64(e.X) - w.cursorPos[0]
					movementY = float64(e.Y) - w.cursorPos[1]
					w.cursorPos[0], w.cursorPos[1] = float64(e.X), float64(e.Y)
					if w.cursorPosCallback != nil {
						go w.cursorPosCallback(w, w.cursorPos[0], w.cursorPos[1])
					}
					if w.mouseMovementCallback != nil {
						go w.mouseMovementCallback(w, w.cursorPos[0], w.cursorPos[1], movementX, movementY)
					}

				case touch.TypeEnd:
					if w.mouseButtonCallback != nil {
						go w.mouseButtonCallback(w, MouseButton(0), Release, 0)
					}
				}

				break

			}
		}
	})
}

func WaitEvents() {
	// TODO.

	runtime.Gosched()
}

func SwapInterval(interval int) error {
	// TODO: Implement.
	return nil
}
func PollEvents() error {
	return nil
}

func (w *Window) MakeContextCurrent() {
	contextWatcher.OnMakeCurrent(w.context)
}

func DetachCurrentContext() {
	contextWatcher.OnDetach()
}

func GetCurrentContext() *Window {
	fmt.Println("not implemented")
	return nil
}

// Monitor
type Monitor struct{}

func (m *Monitor) GetVideoMode() *VidMode {
	return &VidMode{
		// HACK: Hardcoded sample values.
		// TODO: Try to get real values from browser via some API, if possible.
		Width:       1680,
		Height:      1050,
		RedBits:     8,
		GreenBits:   8,
		BlueBits:    8,
		RefreshRate: 60,
	}
}
func GetPrimaryMonitor() *Monitor {
	// TODO: Implement real functionality.
	return &Monitor{}
}

//Window
type Window struct {
	context           interface{}
	requestFullscreen bool // requestFullscreen is set to true when fullscreen should be entered as soon as possible (in a user input handler).
	fullscreen        bool // fullscreen is true if we're currently in fullscreen mode.

	cursorMode  int
	cursorPos   [2]float64
	mouseButton [3]Action

	keys []Action

	cursorPosCallback       CursorPosCallback
	mouseMovementCallback   MouseMovementCallback
	mouseButtonCallback     MouseButtonCallback
	keyCallback             KeyCallback
	charCallback            CharCallback
	scrollCallback          ScrollCallback
	framebufferSizeCallback FramebufferSizeCallback
	sizeCallback            SizeCallback

	touches *js.Object // Hacky mouse-emulation-via-touch.
}

func (w *Window) SetPos(xpos, ypos int) {
	fmt.Println("not implemented: SetPos:", xpos, ypos)
}

func (w *Window) SetSize(width, height int) {
	fmt.Println("not implemented: SetSize:", width, height)
}

func (w *Window) goFullscreenIfRequested() {
	fmt.Println("not implemented: goFullscreenIfRequested:")
}

func (w *Window) ShouldClose() bool {
	return false
}

func (w *Window) SetTitle(title string) {
	// TODO: Implement.
}

func (w *Window) Show() {
	// TODO: Implement.
}

func (w *Window) Hide() {
	// TODO: Implement.
}

func (w *Window) Destroy() {

}

type CursorPosCallback func(w *Window, xpos float64, ypos float64)

func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback) {
	w.cursorPosCallback = cbfun
	return nil
}

type MouseMovementCallback func(w *Window, xpos float64, ypos float64, xdelta float64, ydelta float64)

func (w *Window) SetMouseMovementCallback(cbfun MouseMovementCallback) (previous MouseMovementCallback) {
	w.mouseMovementCallback = cbfun
	return nil
}

type KeyCallback func(w *Window, key Key, scancode int, action Action, mods ModifierKey)

func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback) {
	w.keyCallback = cbfun
	return nil
}

type CharCallback func(w *Window, char rune)

func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback) {
	w.charCallback = cbfun
	return nil
}

type ScrollCallback func(w *Window, xoff float64, yoff float64)

func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback) {
	w.scrollCallback = cbfun
	return nil
}

type MouseButtonCallback func(w *Window, button MouseButton, action Action, mods ModifierKey)

func (w *Window) SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback) {
	w.mouseButtonCallback = cbfun
	return nil
}

type FramebufferSizeCallback func(w *Window, width int, height int)

func (w *Window) SetFramebufferSizeCallback(cbfun FramebufferSizeCallback) (previous FramebufferSizeCallback) {
	w.framebufferSizeCallback = cbfun
	return nil
}

type CloseCallback func(w *Window)

func (w *Window) SetCloseCallback(cbfun CloseCallback) (previous CloseCallback) {
	return nil
}

type RefreshCallback func(w *Window)

func (w *Window) SetRefreshCallback(cbfun RefreshCallback) (previous RefreshCallback) {
	return nil
}

type SizeCallback func(w *Window, width int, height int)

func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	w.sizeCallback = cbfun
	return nil
}

type CursorEnterCallback func(w *Window, entered bool)

func (w *Window) SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback) {

	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type CharModsCallback func(w *Window, char rune, mods ModifierKey)

func (w *Window) SetCharModsCallback(cbfun CharModsCallback) (previous CharModsCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type PosCallback func(w *Window, xpos int, ypos int)

func (w *Window) SetPosCallback(cbfun PosCallback) (previous PosCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type FocusCallback func(w *Window, focused bool)

func (w *Window) SetFocusCallback(cbfun FocusCallback) (previous FocusCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type IconifyCallback func(w *Window, iconified bool)

func (w *Window) SetIconifyCallback(cbfun IconifyCallback) (previous IconifyCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type DropCallback func(w *Window, names []string)

func (w *Window) SetDropCallback(cbfun DropCallback) (previous DropCallback) {
	// TODO: Implement.

	// TODO: Handle previous.
	return nil
}

type Key int

const (
	KeyA Key = 4
	KeyB Key = 5
	KeyC Key = 6
	KeyD Key = 7
	KeyE Key = 8
	KeyF Key = 9
	KeyG Key = 10
	KeyH Key = 11
	KeyI Key = 12
	KeyJ Key = 13
	KeyK Key = 14
	KeyL Key = 15
	KeyM Key = 16
	KeyN Key = 17
	KeyO Key = 18
	KeyP Key = 19
	KeyQ Key = 20
	KeyR Key = 21
	KeyS Key = 22
	KeyT Key = 23
	KeyU Key = 24
	KeyV Key = 25
	KeyW Key = 26
	KeyX Key = 27
	KeyY Key = 28
	KeyZ Key = 29

	Key0 Key = 30
	Key1 Key = 31
	Key2 Key = 32
	Key3 Key = 33
	Key4 Key = 34
	Key5 Key = 35
	Key6 Key = 36
	Key7 Key = 37
	Key8 Key = 38
	Key9 Key = 39

	KeyEnter  Key = 40
	KeyEscape Key = 41
	KeyDelete Key = 42
	KeyTab    Key = 43
	KeySpace  Key = 44
	KeyMinus  Key = 45

	KeyEqual        Key = 46
	KeyLeftBracket  Key = 47
	KeyRightBracket Key = 48
	KeyBackslash    Key = 49
	KeySemicolon    Key = 51
	KeyApostrophe   Key = 52
	KeyGraveAccent  Key = 53
	KeyComma        Key = 54
	KeyPeriod       Key = 55
	KeySlash        Key = 56

	KeyF1  Key = 58
	KeyF2  Key = 59
	KeyF3  Key = 60
	KeyF4  Key = 61
	KeyF5  Key = 62
	KeyF6  Key = 63
	KeyF7  Key = 64
	KeyF8  Key = 65
	KeyF9  Key = 66
	KeyF10 Key = 67
	KeyF11 Key = 68
	KeyF12 Key = 69

	KeyRight Key = 79
	KeyLeft  Key = 80
	KeyDown  Key = 81
	KeyUp    Key = 82

	KeyWorld1 Key = -iota - 1
	KeyWorld2 Key = -iota - 1
	KeyF13    Key = -iota - 1
	KeyF14    Key = -iota - 1
	KeyF15    Key = -iota - 1
	KeyF16    Key = -iota - 1
	KeyF17    Key = -iota - 1
	KeyF18    Key = -iota - 1
	KeyF19    Key = -iota - 1
	KeyF20    Key = -iota - 1
	KeyF21    Key = -iota - 1
	KeyF22    Key = -iota - 1
	KeyF23    Key = -iota - 1
	KeyF24    Key = -iota - 1
	KeyF25    Key = -iota - 1

	KeyLeftShift   Key = -iota - 1
	KeyRightShift  Key = -iota - 1
	KeyLeftControl Key = -iota - 1

	KeyLeftAlt Key = -iota - 1

	KeyRightAlt Key = -iota - 1

	KeyInsert       Key = -iota - 1
	KeyPageUp       Key = -iota - 1
	KeyPageDown     Key = -iota - 1
	KeyHome         Key = -iota - 1
	KeyRightControl Key = -iota - 1

	KeyBackspace Key = -iota - 1
	KeyEnd       Key = -iota - 1
	KeyKP0       Key = -iota - 1
	KeyKP1       Key = -iota - 1
	KeyKP2       Key = -iota - 1
	KeyKP3       Key = -iota - 1
	KeyKP4       Key = -iota - 1
	KeyKP5       Key = -iota - 1
	KeyKP6       Key = -iota - 1
	KeyKP7       Key = -iota - 1
	KeyKP8       Key = -iota - 1
	KeyKP9       Key = -iota - 1

	KeyKPDivide   Key = -iota - 1
	KeyKPMultiply Key = -iota - 1
	KeyKPSubtract Key = -iota - 1
	KeyKPAdd      Key = -iota - 1
	KeyKPDecimal  Key = -iota - 1
	KeyKPEqual    Key = -iota - 1
	KeyKPEnter    Key = -iota - 1

	KeyPrintScreen Key = -iota - 1
	KeyNumLock     Key = -iota - 1
	KeyCapsLock    Key = -iota - 1
	KeyScrollLock  Key = -iota - 1
	KeyPause       Key = -iota - 1
	KeyLeftSuper   Key = -iota - 1
	KeyRightSuper  Key = -iota - 1
	KeyMenu        Key = -iota - 1
)

// toKey extracts Key from given KeyboardEvent.
func toKey(e key.Event) Key {
	key := Key(e.Code)
	return key
}

// toModifierKey extracts ModifierKey from given KeyboardEvent.
func toModifierKey(e key.Event) ModifierKey {
	mods := ModifierKey(0)
	if e.Modifiers == key.ModShift {
		mods += ModShift
	}
	if e.Modifiers == key.ModControl {
		mods += ModControl
	}
	if e.Modifiers == key.ModAlt {
		mods += ModAlt
	}
	if e.Modifiers == key.ModMeta {
		mods += ModSuper
	}
	return mods
}

type MouseButton int

const (
	MouseButton1 MouseButton = 0
	MouseButton2 MouseButton = 2 // Web MouseEvent has middle and right mouse buttons in reverse order.
	MouseButton3 MouseButton = 1 // Web MouseEvent has middle and right mouse buttons in reverse order.

	MouseButtonLeft   = MouseButton1
	MouseButtonRight  = MouseButton2
	MouseButtonMiddle = MouseButton3
)

type Action int

const (
	Release Action = 0
	Press   Action = 1
	Repeat  Action = 2
)

type InputMode int

const (
	CursorMode InputMode = iota
	StickyKeysMode
	StickyMouseButtonsMode
)

const (
	CursorNormal = iota
	CursorHidden
	CursorDisabled
)

type ModifierKey int

const (
	ModShift ModifierKey = (1 << iota)
	ModControl
	ModAlt
	ModSuper
)
