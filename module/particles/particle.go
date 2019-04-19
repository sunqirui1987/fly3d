package particles

import (
	"github.com/suiqirui1987/fly3d/math32"
)

type Particle struct {
	Position     *math32.Vector3
	Direction    *math32.Vector3
	LifeTime     float32
	Age          float32
	Size         float32
	Angle        float32
	AngularSpeed float32
	Color        *math32.Color4
	ColorStep    *math32.Color4
}

func NewParticle() *Particle {
	this := &Particle{}
	return this
}

func (this *Particle) Init() {
	this.Position = nil
	this.Direction = nil
	this.LifeTime = 1.0
	this.Age = 0.0
	this.Size = 0.0
	this.Angle = 0.0
	this.AngularSpeed = 0.0
	this.Color = nil
	this.ColorStep = nil
}
