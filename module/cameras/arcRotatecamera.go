package cameras

import (
	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/windows"
	"math"
)

type IArcRotateCameraTarget interface {
	GetPostion() *math32.Vector3
}

type ArcRotateCamera struct {
	Camera

	Alpha  float32
	Beta   float32
	Radius float32
	Target IArcRotateCameraTarget

	InertialAlphaOffset float32
	InertialBetaOffset  float32

	_keys []int
}

func NewArcRotateCamera(name string, alpha, beta, radius float32, target IArcRotateCameraTarget, scene *engines.Scene) *ArcRotateCamera {
	this := &ArcRotateCamera{}

	this.Init(name, alpha, beta, radius, target, scene)

	this._scene.Cameras = append(this._scene.Cameras, this)

	if this._scene.ActiveCamera == nil {
		this._scene.ActiveCamera = this
	}

	this.GetViewMatrix()
	return this
}

func (this *ArcRotateCamera) Init(name string, alpha, beta, radius float32, target IArcRotateCameraTarget, scene *engines.Scene) {

	this.Camera.Init(name, target.GetPostion(), scene)

	this.Alpha = alpha
	this.Beta = beta
	this.Radius = radius
	this.Target = target

	this._keys = make([]int, 0)

	this._scene = scene

}

func (this *ArcRotateCamera) AttachControl(win windows.IWindow) {

	var previousPosition *math32.Vector2
	that := this

	_onPointerDown := func(evt *windows.MouseEvent) error {

		previousPosition = math32.NewVector2(evt.ClientX, evt.ClientY)

		evt.StopPropagation()
		return nil
	}

	_onPointerUp := func(evt *windows.MouseEvent) error {
		previousPosition = nil
		evt.StopPropagation()
		return nil
	}

	_onPointerMove := func(evt *windows.MouseEvent) error {
		if previousPosition == nil {
			return nil
		}

		offsetX := evt.ClientX - previousPosition.X
		offsetY := evt.ClientY - previousPosition.Y

		that.InertialAlphaOffset = that.InertialAlphaOffset - offsetX/1000.0
		that.InertialBetaOffset = that.InertialBetaOffset - offsetY/1000.0

		previousPosition = math32.NewVector2(evt.ClientX, evt.ClientY)

		evt.StopPropagation()

		return nil
	}

	_wheel := func(evt *windows.WheelEvent) error {
		var delta float32
		delta = 0.0
		if evt.DeltaX != 0.0 {
			delta = evt.DeltaX / 3.0
		} else if evt.DeltaY != 0.0 {
			delta = evt.DeltaY / 3.0
		}

		that.Radius -= delta

		evt.StopPropagation()

		return nil
	}

	_onKeyDown := func(evt *windows.KeyboardEvent) error {
		if evt.KeyCode == KEYS_UP ||
			evt.KeyCode == KEYS_DOWN ||
			evt.KeyCode == KEYS_LEFT ||
			evt.KeyCode == KEYS_RIGHT {

			index := -1
			for i, code := range that._keys {
				if code == evt.KeyCode {
					index = i
				}

			}
			if index == -1 {
				that._keys = append(that._keys, evt.KeyCode)
			}
			evt.StopPropagation()
		}

		return nil
	}

	_onKeyUp := func(evt *windows.KeyboardEvent) error {
		if evt.KeyCode == KEYS_UP ||
			evt.KeyCode == KEYS_DOWN ||
			evt.KeyCode == KEYS_LEFT ||
			evt.KeyCode == KEYS_RIGHT {

			index := -1
			for i, code := range that._keys {
				if code == evt.KeyCode {
					index = i
				}

			}

			if index >= 0 {
				that._keys = append(that._keys[:index], that._keys[index+1:]...)
			}
			evt.StopPropagation()

		}

		return nil
	}

	_onLostFocus := func(evt *windows.FocusEvent) error {
		if evt.Focused == false {
			that._keys = make([]int, 0)
		}
		return nil
	}

	// Subscribe to events
	win.On(windows.MouseDown, _onPointerDown)
	win.On(windows.MouseUp, _onPointerUp)
	win.On(windows.MouseOut, _onPointerUp)
	win.On(windows.MouseMove, _onPointerMove)

	win.On(windows.Wheel, _wheel)

	win.On(windows.Keydown, _onKeyDown)
	win.On(windows.Keyup, _onKeyUp)
	win.On(windows.Focus, _onLostFocus)
}

func (this *ArcRotateCamera) DetachControl(win windows.IWindow) {
	win.Remove(windows.MouseDown)
	win.Remove(windows.MouseUp)
	win.Remove(windows.MouseOut)
	win.Remove(windows.MouseMove)

	win.Remove(windows.Wheel)

	win.Remove(windows.Keydown)
	win.Remove(windows.Keyup)
	win.Remove(windows.Focus)
}

func (this *ArcRotateCamera) Update() {
	// Keyboard
	for index := 0; index < len(this._keys); index++ {
		keyCode := this._keys[index]

		if keyCode == KEYS_LEFT {
			this.InertialAlphaOffset -= 0.01
		} else if keyCode == KEYS_UP {
			this.InertialBetaOffset -= 0.01
		} else if keyCode == KEYS_RIGHT {
			this.InertialAlphaOffset += 0.01
		} else if keyCode == KEYS_DOWN {
			this.InertialBetaOffset += 0.01
		}
	}

	// Inertia
	if this.InertialAlphaOffset != 0 || this.InertialBetaOffset != 0 {

		this.Alpha += this.InertialAlphaOffset
		this.Beta += this.InertialBetaOffset

		this.InertialAlphaOffset *= this.Inertia
		this.InertialBetaOffset *= this.Inertia

		if math32.Abs(this.InertialAlphaOffset) < core.Epsilon {
			this.InertialAlphaOffset = 0
		}

		if math32.Abs(this.InertialBetaOffset) < core.Epsilon {
			this.InertialBetaOffset = 0
		}

	}
}

func (this *ArcRotateCamera) SetPosition(position *math32.Vector3) {
	radiusv3 := position.Sub(this.Target.GetPostion())
	this.Radius = radiusv3.Length()

	this.Alpha = math32.Atan(radiusv3.Z / radiusv3.X)
	this.Beta = math32.Acos(radiusv3.Y / this.Radius)
}

func (this *ArcRotateCamera) GetViewMatrix() *math32.Matrix4 {

	// Compute
	if this.Beta > math.Pi {
		this.Beta = math.Pi
	}

	if this.Beta <= 0 {
		this.Beta = 0.01
	}

	cosa := math32.Cos(this.Alpha)
	sina := math32.Sin(this.Alpha)
	cosb := math32.Cos(this.Beta)
	sinb := math32.Sin(this.Beta)

	this.Position = this.Target.GetPostion().Add(math32.NewVector3(this.Radius*cosa*sinb, this.Radius*cosb, this.Radius*sina*sinb))

	m := math32.NewMatrix4()
	viewm := m.LookAtLH(this.Position, this.Target.GetPostion(), math32.NewVector3Up())

	return viewm

}
