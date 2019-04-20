package interfaces

import "github.com/suiqirui1987/fly3d/math32"

type IShadowGenerator interface {
	IsReady() bool
	IsUseVarianceShadowMap() bool
	GetShadowMap() IRenderTargetTexture
	Dispose()
}
type ILight interface {
	IsEnabled() bool
	IsSupportShadow() bool

	GetShadowGenerator() IShadowGenerator
	SetShadowGenerator(IShadowGenerator)

	GetIntensity() float32
	GetDiffuse() *math32.Color3
	GetSpecular() *math32.Color3

	GetPosition() *math32.Vector3
	GetDirection() *math32.Vector3
}
