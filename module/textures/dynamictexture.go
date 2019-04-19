package textures

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
)

type DynamicTexture struct {
	BaseTexture
	Name string
}

func NewDynamicTexture(name string, size int, scene *engines.Scene, generateMipMaps bool) *DynamicTexture {
	this := &DynamicTexture{}

	this._scene = scene
	this._scene.Textures = append(this._scene.Textures, this)
	this.Name = name
	this._texture = this._scene.GetEngine().CreateDynamicTexture(size, generateMipMaps)
	this._texture.WrapU = gl.CLAMP_ADDRESSMODE
	this._texture.WrapV = gl.CLAMP_ADDRESSMODE

	return this
}

func (this *DynamicTexture) Update() {

}
