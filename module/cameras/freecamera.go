package cameras

import (
	"math"

	. "github.com/suiqirui1987/fly3d/interfaces"

	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/module/collisions"
	"github.com/suiqirui1987/fly3d/tools"
	"github.com/suiqirui1987/fly3d/windows"

	log "github.com/suiqirui1987/fly3d/tools/logrus"
)

type FreeCamera struct {
	Camera

	CameraDirection    *math32.Vector3
	CameraRotation     *math32.Vector2
	Rotation           *math32.Vector3
	Ellipsoid          *math32.Vector3
	AngularSensibility float32
	MoveSensibility    float32
	OnCollide          func(IMesh)
	_keys              []string

	// Collisions
	_collider           *collisions.Collider
	_needMoveForGravity bool

	// Offset
	_offsetX        float32
	_offsetY        float32
	_pointerCount   int
	_pointerPressed []int

	//speed
	Speed           float32
	CheckCollisions bool
	ApplyGravity    bool

	KeysUp    []string
	KeysDown  []string
	KeysLeft  []string
	KeysRight []string
}

func NewFreeCamera(name string, position *math32.Vector3, scene *engines.Scene) *FreeCamera {

	this := &FreeCamera{}

	this.Init(name, position, scene)
	this._scene.Cameras = append(this._scene.Cameras, this)
	if this._scene.ActiveCamera == nil {
		this._scene.ActiveCamera = this
	}

	this.KeysUp = []string{"UP", "W"}
	this.KeysDown = []string{"DOWN", "S"}
	this.KeysLeft = []string{"LEFT", "A"}
	this.KeysRight = []string{"RIGHT", "D"}

	return this

}

func (this *FreeCamera) Init(name string, position *math32.Vector3, scene *engines.Scene) {

	this.Camera.Init(name, position, scene)

	this.CameraDirection = math32.NewVector3(0, 0, 0)
	this.CameraRotation = math32.NewVector2(0.0, 0.0)
	this.Rotation = math32.NewVector3(0, 0, 0)
	this.Ellipsoid = math32.NewVector3(0.5, 1, 0.5)

	this._keys = make([]string, 0)

	// Collisions
	this._collider = collisions.NewCollider()
	this._needMoveForGravity = true

	this.Speed = 2.0
	this.CheckCollisions = false
	this.ApplyGravity = false

}

func (this *FreeCamera) _computeLocalCameraSpeed() float32 {
	return this.Speed * (tools.GetDeltaTime() / (tools.GetFps() * 10.0))
}

//

func (this *FreeCamera) SetTarget(target *math32.Vector3) {
	m := math32.NewMatrix4()
	camMatrix := m.LookAtLH(this.Position, target, math32.NewVector3Up())
	camMatrix.Invert()

	this.Rotation.X = math32.Atan(camMatrix[6] / camMatrix[10])
	vDir := target.Sub(this.Position)

	if vDir.X >= 0.0 {
		this.Rotation.Y = (-math32.Atan(vDir.Z/vDir.X) + math.Pi/2.0)
	} else {
		this.Rotation.Y = (-math32.Atan(vDir.Z/vDir.X) - math.Pi/2.0)
	}

	v := math32.NewVector3(0.0, 1.0, 0.0)
	this.Rotation.Z = -math32.Acos(v.Dot(math32.NewVector3Up()))

}

func (this *FreeCamera) AttachControl(win windows.IWindow) {
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

	_onPointerOut := func(evt *windows.MouseEvent) error {
		previousPosition = nil
		that._keys = make([]string, 0)
		evt.StopPropagation()
		return nil
	}

	_onPointerMove := func(evt *windows.MouseEvent) error {
		if previousPosition == nil {
			return nil
		}

		offsetX := evt.ClientX - previousPosition.X
		offsetY := evt.ClientY - previousPosition.Y

		that.CameraRotation.Y += offsetX / 2000.0
		that.CameraRotation.X += offsetY / 2000.0

		previousPosition = math32.NewVector2(evt.ClientX, evt.ClientY)

		evt.StopPropagation()

		return nil
	}

	_onKeyDown := func(evt *windows.KeyboardEvent) error {
		log.Printf("FreeCamera onKeyDown %d, %s", evt.KeyCode, evt.CharCode)
		if tools.IndexOf(evt.CharCode, this.KeysUp) != -1 ||
			tools.IndexOf(evt.CharCode, this.KeysDown) != -1 ||
			tools.IndexOf(evt.CharCode, this.KeysLeft) != -1 ||
			tools.IndexOf(evt.CharCode, this.KeysRight) != -1 {

			index := tools.IndexOf(evt.CharCode, that._keys)
			if index == -1 {
				that._keys = append(that._keys, evt.CharCode)
			}
			evt.StopPropagation()
		}

		return nil
	}

	_onKeyUp := func(evt *windows.KeyboardEvent) error {
		log.Printf("FreeCamera onKeyUp %d, %s", evt.KeyCode, evt.CharCode)
		if tools.IndexOf(evt.CharCode, this.KeysUp) != -1 ||
			tools.IndexOf(evt.CharCode, this.KeysDown) != -1 ||
			tools.IndexOf(evt.CharCode, this.KeysLeft) != -1 ||
			tools.IndexOf(evt.CharCode, this.KeysRight) != -1 {

			index := tools.IndexOf(evt.CharCode, that._keys)
			if index >= 0 {
				that._keys = append(that._keys[:index], that._keys[index+1:]...)
			}
			evt.StopPropagation()

		}

		return nil
	}

	_onLostFocus := func(evt *windows.FocusEvent) error {
		if evt.Focused == false {
			that._keys = make([]string, 0)
		}
		return nil
	}

	// Subscribe to events
	win.On(windows.MouseDown, _onPointerDown)
	win.On(windows.MouseUp, _onPointerUp)
	win.On(windows.MouseOut, _onPointerOut)
	win.On(windows.MouseMove, _onPointerMove)

	win.On(windows.Keydown, _onKeyDown)
	win.On(windows.Keyup, _onKeyUp)
	win.On(windows.Focus, _onLostFocus)
}

func (this *FreeCamera) DetachControl(win windows.IWindow) {
	win.Remove(windows.MouseDown)
	win.Remove(windows.MouseUp)
	win.Remove(windows.MouseOut)
	win.Remove(windows.MouseMove)

	win.Remove(windows.Keydown)
	win.Remove(windows.Keyup)
	win.Remove(windows.Focus)
}

func (this *FreeCamera) _collideWithWorld(velocity *math32.Vector3) {
	oldPosition := this.Position.Sub(math32.NewVector3(0, this.Ellipsoid.Y, 0))
	this._collider.Radius = this.Ellipsoid

	if this.ApplyGravity {
		velocity = velocity.Add(this._scene.Gravity)
	}

	newPosition := collisions.GetNewPosition(this._scene, oldPosition, velocity, this._collider, 3)
	diffPosition := newPosition.Sub(oldPosition)

	if diffPosition.Length() > core.CollisionsEpsilon {
		this.Position = this.Position.Add(diffPosition)

		collidedMesh := this._collider.GetMesh()
		if this.OnCollide != nil && !tools.IsNil(collidedMesh) {
			this.OnCollide(collidedMesh)
		}
	}
}

func (this *FreeCamera) _checkInputs() {
	// Keyboard
	for index := 0; index < len(this._keys); index++ {
		charCode := this._keys[index]

		var direction *math32.Vector3
		speed := this._computeLocalCameraSpeed()

		if tools.IndexOf(charCode, this.KeysLeft) != -1 {
			direction = math32.NewVector3(-speed, 0, 0)
		} else if tools.IndexOf(charCode, this.KeysUp) != -1 {
			direction = math32.NewVector3(0, 0, speed)
		} else if tools.IndexOf(charCode, this.KeysRight) != -1 {
			direction = math32.NewVector3(speed, 0, 0)
		} else if tools.IndexOf(charCode, this.KeysDown) != -1 {
			direction = math32.NewVector3(0, 0, -speed)
		}

		//log.Printf("_checkInputs %s direction %f ", charCode, direction)

		cameraTransformMatrix := math32.NewMatrix4().RotationYawPitchRoll(this.Rotation.Y, this.Rotation.X, 0)
		transformedDirection := direction.TransformCoordinates(cameraTransformMatrix)
		this.CameraDirection = this.CameraDirection.Add(transformedDirection)
	}
}
func (this *FreeCamera) Update() {
	this._checkInputs()

	needToMove := this._needMoveForGravity || math32.Abs(this.CameraDirection.X) > 0 || math32.Abs(this.CameraDirection.Y) > 0 || math32.Abs(this.CameraDirection.Z) > 0
	needToRotate := math32.Abs(this.CameraRotation.X) > 0 || math32.Abs(this.CameraRotation.Y) > 0

	// Move
	if needToMove {
		if this.CheckCollisions && this._scene.CollisionsEnabled {
			this._collideWithWorld(this.CameraDirection)
		} else {
			this.Position = this.Position.Add(this.CameraDirection)
		}
	}

	// Rotate
	if needToRotate {
		this.Rotation.X += this.CameraRotation.X
		this.Rotation.Y += this.CameraRotation.Y

		limit := (float32)((math32.Pi / 2) * 0.95)

		if this.Rotation.X > limit {
			this.Rotation.X = limit
		}

		if this.Rotation.X < -limit {
			this.Rotation.X = -limit
		}

	}

	// Inertia
	if needToMove {
		this.CameraDirection = this.CameraDirection.Scale(this.Inertia)
	}
	if needToRotate {
		this.CameraRotation = this.CameraRotation.Scale(this.Inertia)
	}
}
func (this *FreeCamera) GetViewMatrix() *math32.Matrix4 {

	// Compute
	referencePoint := math32.NewVector3(0, 0, 1)

	mt := math32.NewMatrix4()
	transform := mt.RotationYawPitchRoll(this.Rotation.Y, this.Rotation.X, this.Rotation.Z)

	currentTarget := this.Position.Add(referencePoint.TransformCoordinates(transform))

	m := math32.NewMatrix4()
	return m.LookAtLH(this.Position, currentTarget, math32.NewVector3Up())
}
