package textures

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
)

type CubeTexture struct {
	BaseTexture

	Extensions []string

	_textureMatrix *math32.Matrix4
}

func NewCubeTexture(rootUrl string, scene *engines.Scene, extensions []string) *CubeTexture {
	this := &CubeTexture{}
	this._scene = scene

	this.Init()

	this._texture = this._getFromCache(rootUrl, false)

	if this._texture == nil {
		this._texture = scene.GetEngine().CreateCubeTexture(rootUrl, scene, this.Extensions)
	}

	this._texture.Url = rootUrl
	this._texture.IsCube = true

	this._scene.Textures = append(this._scene.Textures, this)

	this.Extensions = extensions

	this.CoordinatesMode = CUBIC_MODE

	this._textureMatrix = (math32.NewMatrix4()).Identity()

	return this
}

func (this *CubeTexture) Init() {
	this.BaseTexture.Init()
}

func (this *CubeTexture) ComputeReflectionTextureMatrix() *math32.Matrix4 {
	return this._textureMatrix
}
