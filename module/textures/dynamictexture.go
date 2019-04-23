package textures

import (
	"image"

	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
)

type DynamicTexture struct {
	Texture
	Name string

	_canvasimg *image.RGBA
}

func NewDynamicTexture(name string, size int, scene *engines.Scene, generateMipMaps bool) *DynamicTexture {
	this := &DynamicTexture{}

	this._scene = scene
	this._scene.Textures = append(this._scene.Textures, this)
	this.Name = name
	this._texture = this._scene.GetEngine().CreateDynamicTexture(size, generateMipMaps)

	this.Init()

	return this
}
func (this *DynamicTexture) Init() {
	this.Texture.Init()
	this._texture.WrapU = gl.CLAMP_ADDRESSMODE
	this._texture.WrapV = gl.CLAMP_ADDRESSMODE
}

func (this *DynamicTexture) Update() {
	if this._canvasimg == nil {
		return
	}
	this._scene.GetEngine().UpdateDynamicTexture(this._texture, this._canvasimg)

}

func (this *DynamicTexture) DrawText(text string, x, y float32) {

}
