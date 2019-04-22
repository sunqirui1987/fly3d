package textures

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/tools"
)

type BaseTexture struct {
	_texture *gl.GLTextureBuffer
	_scene   *engines.Scene

	_hasAlpha bool
	Level     float32

	CoordinatesIndex float32
	CoordinatesMode  int

	OnDispose func()
}

func (this *BaseTexture) Init() {
	this._hasAlpha = false
	this.Level = 1
}

func (this *BaseTexture) GetGLTexture() *gl.GLTextureBuffer {
	return this._texture
}

func (this *BaseTexture) EnableHasAlpha(val bool) {
	this._hasAlpha = val
}
func (this *BaseTexture) HasAlpha() bool {
	return this._hasAlpha
}

func (this *BaseTexture) GetCoordinatesIndex() float32 {
	return this.CoordinatesIndex
}

func (this *BaseTexture) GetCoordinatesMode() int {
	return this.CoordinatesMode
}

func (this *BaseTexture) GetLevel() float32 {
	return this.Level
}

func (this *BaseTexture) IsReady() bool {
	return (this._texture != nil && this._texture.IsReady)
}

func (this *BaseTexture) IsRenderTarget() bool {
	return false
}

func (this *BaseTexture) GetSize() *math32.Vector2 {
	if this._texture.Width > 0 {
		return math32.NewVector2(float32(this._texture.Width), float32(this._texture.Height))
	}

	return math32.NewVector2(0, 0)
}

func (this *BaseTexture) GetBaseSize() *math32.Vector2 {
	if this._texture.BaseWidth > 0 {
		return math32.NewVector2(float32(this._texture.BaseWidth), float32(this._texture.BaseHeight))
	}

	return math32.NewVector2(0, 0)
}

func (this *BaseTexture) _getFromCache(url string, noMipmap bool) *gl.GLTextureBuffer {
	texturesCache := this._scene.GetEngine().GetLoadedTexturesCache()
	for index := 0; index < len(texturesCache); index++ {
		texturesCacheEntry := texturesCache[index]

		if texturesCacheEntry.Url == url && texturesCacheEntry.NoMipmap == noMipmap {
			texturesCacheEntry.References++

			return texturesCacheEntry
		}
	}

	return nil
}

func (this *BaseTexture) ReleaseGLTexture() {
	if this._texture == nil {
		return
	}
	texturesCache := this._scene.GetEngine().GetLoadedTexturesCache()
	this._texture.References--

	// Final reference
	if this._texture.References == 0 {
		index := tools.IndexOf(this._scene.Textures, texturesCache)
		if index > -1 {
			texturesCache = append(texturesCache[:index], texturesCache[index+1:]...)
		}

		this._scene.GetEngine().ReleaseTexture(this._texture)
		this._texture = nil
	}
}

func (this *BaseTexture) ComputeTextureMatrix() *math32.Matrix4 {
	return nil
}
func (this *BaseTexture) ComputeReflectionTextureMatrix() *math32.Matrix4 {
	return nil
}

func (this *BaseTexture) Render() {

}

func (this *BaseTexture) Dispose() {
	if this._texture == nil {
		return
	}
	this.ReleaseGLTexture()

	// Remove from scene
	index := tools.IndexOf(this._scene.Textures, this)
	if index > -1 {
		this._scene.Textures = append(this._scene.Textures[:index], this._scene.Textures[index+1:]...)
	}

	// Callback
	if this.OnDispose != nil {
		this.OnDispose()
	}
}
