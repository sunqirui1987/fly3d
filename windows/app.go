package windows

import (
	"github.com/suiqirui1987/fly3d/gui/nanogui"

	"github.com/suiqirui1987/fly3d/tools/goevent"

	"github.com/suiqirui1987/fly3d/glfw"
)

var mainloopActive bool = false

type AppOption struct {
	Width  int
	Height int
	Title  string
}

type App struct {
	goevent.Dispatcher

	Width  int
	Height int
	Screen *nanogui.Screen

	_stopFlag  bool
	_renderFun func()

	_posX float32
	_posY float32
}

func NewApp(opt *AppOption) (*App, error) {
	this := &App{}
	this.Width = opt.Width
	this.Height = opt.Height
	this._stopFlag = false

	this.Dispatcher.Initialize()
	nanogui.Init()

	this.Screen = nanogui.NewScreen(this.Width, this.Height, opt.Title, false, false)

	this.Screen.SetResizeEventCallback(this._SizeCallback)

	this.Screen.ScreenCursorPosCallback = this._CursorPosCallback
	this.Screen.ScreenMouseButtonCallback = this._MouseButtonCallback
	this.Screen.ScreenKeyCallback = this._KeyCallback
	this.Screen.ScreenCharCallback = this._CharCallback
	this.Screen.ScreenScrollCallback = this._ScrollCallback

	return this, nil
}

func (this *App) Run() {

	this.Screen.SetDrawContentsCallback(func() {
		if this._renderFun != nil && this._stopFlag == false {
			this._renderFun()
		}
	})

	nanogui.MainLoop()
	defer glfw.Terminate()
}

//interface
func (this *App) GetRenderWidth() int {
	return this.Screen.Fbw
}
func (this *App) GetRenderHeight() int {
	return this.Screen.FbH
}

func (this *App) StopNewFrame() {
	this._stopFlag = true
}
func (this *App) QueueNewFrame(f func()) {
	this._renderFun = f
}
func (this *App) GetWindowDevicePixelRatio() float32 {
	return 1.0
}

func (this *App) GetFullscreen() bool {
	return false
}
func (this *App) ExitFullscreen() {

}
func (this *App) RequestFullscreen() {

}

//events

func (this *App) _SizeCallback(width int, height int) bool {
	this.Width = width
	this.Height = height

	evt := &ResizeEvent{
		WindowWidth:  width,
		WindowHeigth: height,
	}

	this.Emit(Resize, evt)

	return true
}

func (this *App) _CloseCallback(w *glfw.Window) {

}

func (this *App) _RefreshCallback(w *glfw.Window) {

}

func (this *App) _FocusCallback(w *glfw.Window, focused bool) {
	evt := &FocusEvent{
		Focused: focused,
	}

	this.Emit(Focus, evt)
}

func (this *App) _MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	evt := &MouseEvent{
		ClientX: this._posX,
		ClientY: this._posY,
	}

	if action == glfw.Press {
		//mouse down
		this.Emit(MouseDown, evt)
	} else if action == glfw.Release {
		//mouse up
		this.Emit(MouseUp, evt)
	}

}

func (this *App) _CursorPosCallback(w *glfw.Window, x float64, y float64) {

	this._posX = float32(x)
	this._posY = float32(y)
	evt := &MouseEvent{
		ClientX: this._posX,
		ClientY: this._posY,
	}

	this.Emit(MouseMove, evt)

}

func (this *App) _ScrollCallback(w *glfw.Window, x float64, y float64) {
	evt := &WheelEvent{
		DeltaX: float32(x),
		DeltaY: float32(y),
	}
	this.Emit(Wheel, evt)
}

func (this *App) _KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	evt := &KeyboardEvent{
		KeyCode:  int(key),
		CharCode: KeyString(key),
	}

	if action == glfw.Press {
		//key down
		this.Emit(Keydown, evt)
	} else if action == glfw.Release {
		//key up
		this.Emit(Keyup, evt)
	}
}

func (this *App) _CharCallback(w *glfw.Window, char rune) {

}

func (this *App) _CharModsCallback(w *glfw.Window, char rune, mods glfw.ModifierKey) {

}

func (this *App) _DropCallback(w *glfw.Window, names []string) {

}
