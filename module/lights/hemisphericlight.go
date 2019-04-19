package lights

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
)

type HemisphericLight struct {
	Light

	GroundColor *math32.Color3

	// Animations
}

func NewHemisphericLight(name string, direction *math32.Vector3, scene *engines.Scene) *HemisphericLight {
	this := &HemisphericLight{}
	this.Init()
	this.Name = name
	this.Id = name

	this._scene = scene
	this.Direction = direction

	this._scene.Lights = append(this._scene.Lights, this)

	return this
}
func (this *HemisphericLight) Init() {
	this.Light.Init()
	this.Diffuse = math32.NewColor3(1.0, 1.0, 1.0)
	this.Specular = math32.NewColor3(1.0, 1.0, 1.0)
	this.GroundColor = math32.NewColor3(0.0, 0.0, 0.0)
}

func (this *HemisphericLight) IsSupportShadow() bool {
	return false
}
