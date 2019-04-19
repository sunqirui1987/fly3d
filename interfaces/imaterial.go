package interfaces

import "github.com/suiqirui1987/fly3d/math32"

type IMultiMaterial interface {
	GetId() string
	GetSubMaterial(index int) IMaterial
}
type IMaterial interface {
	IsReady(IMesh) bool
	GetId() string
	GetEffect() IEffect
	PreBind()
	Bind(world *math32.Matrix4, mesh IMesh)
	UnBind()
	HasWireframe() bool

	GetRenderTargetTextures() []ITexture

	NeedAlphaTesting() bool
	NeedAlphaBlending() bool
	GetAlpha() float32

	Dispose()
}
