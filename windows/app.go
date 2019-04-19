package windows

import (
	log "github.com/sirupsen/logrus"
	"github.com/suiqirui1987/fly3d/tools/goevent"

	"github.com/suiqirui1987/fly3d/gl"
	"github.com/suiqirui1987/fly3d/glfw"
)

type AppOption struct {
	Width  int
	Height int
	HInt   int
	Title  string
}

type App struct {
	goevent.Dispatcher

	Width        int
	Height       int
	RenderWidth  int
	RenderHeight int
	HInt         int

	_stopFlag  bool
	_window    *glfw.Window
	_renderFun func()

	_posX float32
	_posY float32
}

func NewApp(opt *AppOption) (*App, error) {
	this := &App{}
	this.Width = opt.Width
	this.Height = opt.Height
	this.HInt = opt.HInt
	this._stopFlag = false

	this.Dispatcher.Initialize()

	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		log.Printf("glfw.Init Failed %s", err)
		return nil, err
	}

	var windowSize = [2]int{this.Width, this.Height}
	glfw.WindowHint(glfw.Samples, this.HInt) // Anti-aliasing.

	this._window, err = glfw.CreateWindow(windowSize[0], windowSize[1], opt.Title, nil, nil)
	if err != nil {
		log.Printf("glfw.CreateWindow Failed %s", err)
		return nil, err
	}
	this._window.MakeContextCurrent()
	glfw.SwapInterval(1)

	log.Printf("OpenGL: %s %s %s; %v samples.\n", gl.GetString(gl.VENDOR), gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION), gl.GetInteger(gl.SAMPLES))
	log.Printf("GLSL: %s.\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	this._window.SetSizeCallback(this._SizeCallback)
	this._window.SetCloseCallback(this._CloseCallback)
	this._window.SetRefreshCallback(this._RefreshCallback)
	this._window.SetFocusCallback(this._FocusCallback)
	this._window.SetMouseButtonCallback(this._MouseButtonCallback)
	this._window.SetCursorPosCallback(this._CursorPosCallback)
	this._window.SetScrollCallback(this._ScrollCallback)
	this._window.SetKeyCallback(this._KeyCallback)
	this._window.SetCharCallback(this._CharCallback)
	this._window.SetCharModsCallback(this._CharModsCallback)
	this._window.SetDropCallback(this._DropCallback)

	framebufferSizeCallback := func(w *glfw.Window, framebufferSize0, framebufferSize1 int) {
		this.Width, this.Height = w.GetSize()
		log.Printf("NewApp render %d ,%d size %d,%d \n", framebufferSize0, framebufferSize1, this.Width, this.Height)

		this.RenderWidth = framebufferSize0
		this.RenderHeight = framebufferSize1

	}
	this._window.SetFramebufferSizeCallback(framebufferSizeCallback)
	{
		var framebufferSize [2]int
		framebufferSize[0], framebufferSize[1] = this._window.GetFramebufferSize()
		framebufferSizeCallback(this._window, framebufferSize[0], framebufferSize[1])
	}

	return this, nil
}

func (this *App) Run() {
	for !this._window.ShouldClose() {

		if this._stopFlag {
			log.Print("App Exit")
			break
		}
		if this._renderFun != nil {
			this._renderFun()
		}

		this._window.SwapBuffers()
		glfw.PollEvents()
	}

	defer glfw.Terminate()
}

func (this *App) _resize() {
	this._window.SetSize(this.Width, this.Height)
}

//interface
func (this *App) GetRenderWidth() int {
	return this.RenderWidth
}
func (this *App) GetRenderHeight() int {
	return this.RenderHeight
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

func (this *App) _SizeCallback(w *glfw.Window, width int, height int) {
	this.Width = width
	this.Height = height

	evt := &ResizeEvent{
		WindowWidth:  width,
		WindowHeigth: height,
	}

	this.Emit(Resize, evt)

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
		this.Emit(Keydown, evt)
	}
}

func (this *App) _CharCallback(w *glfw.Window, char rune) {

}

func (this *App) _CharModsCallback(w *glfw.Window, char rune, mods glfw.ModifierKey) {

}

func (this *App) _DropCallback(w *glfw.Window, names []string) {

}
