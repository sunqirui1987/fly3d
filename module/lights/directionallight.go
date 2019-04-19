package lights

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
)

type DirectionalLight struct {
	Light
}

func NewDirectionalLight(name string, direction *math32.Vector3, scene *engines.Scene) *DirectionalLight {
	this := &DirectionalLight{}
	this.Init()

	this.Name = name
	this.Id = name

	this._scene = scene

	this.Position = direction.Scale(-1)
	this.Direction = direction

	this._scene.Lights = append(this._scene.Lights, this)

	return this

}

func (this *DirectionalLight) Init() {
	this.Light.Init()
	this.Diffuse = math32.NewColor3(1.0, 1.0, 1.0)
	this.Specular = math32.NewColor3(1.0, 1.0, 1.0)
}

func (this *DirectionalLight) IsSupportShadow() bool {
	return true
}
