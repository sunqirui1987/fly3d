package lights

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
)

type PointLight struct {
	Light
}

func NewPointLight(name string, position *math32.Vector3, scene *engines.Scene) *PointLight {
	this := &PointLight{}

	this.Init()

	this.Name = name
	this.Id = name
	this._scene = scene

	this.Position = position
	this.ShadowGenerator = nil

	this._scene.Lights = append(this._scene.Lights, this)

	return this
}

func (this *PointLight) Init() {
	this.Light.Init()
	this.Diffuse = math32.NewColor3(1.0, 1.0, 1.0)
	this.Specular = math32.NewColor3(1.0, 1.0, 1.0)
}

func (this *PointLight) IsSupportShadow() bool {
	return false
}
