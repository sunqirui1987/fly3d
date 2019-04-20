package interfaces

import (
	"github.com/suiqirui1987/fly3d/gl"
	"github.com/suiqirui1987/fly3d/math32"
)

type ITexture interface {
	IsReady() bool
	IsRenderTarget() bool
	HasAlpha() bool
	GetCoordinatesIndex() float32
	GetCoordinatesMode() int
	GetLevel() float32

	GetGLTexture() *gl.GLTextureBuffer
	ComputeTextureMatrix() *math32.Matrix4
	ComputeReflectionTextureMatrix() *math32.Matrix4

	Render()
	Dispose()
}

type IRenderTargetTexture interface {
	ITexture

	GetRenderList() []IMesh
}
