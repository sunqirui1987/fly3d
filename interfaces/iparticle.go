package interfaces

type IParticleSystem interface {
	Animate()
	GetEmitter() IMesh
	Render() int
	Dispose()
}
