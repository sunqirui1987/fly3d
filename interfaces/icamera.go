package interfaces

import (
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/windows"
)

type ICamera interface {
	GetId() string
	GetPosition() *math32.Vector3
	GetMinZ() float32
	GetMaxZ() float32

	AttachControl(win windows.IWindow)
	DetachControl(win windows.IWindow)

	GetViewMatrix() *math32.Matrix4
	GetProjectionMatrix() *math32.Matrix4

	Update()
}
