package particles

import (
	"math/rand"
	"strings"

	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/module/effects"
	"github.com/suiqirui1987/fly3d/tools"
)

func randomNumber(min float32, max float32) float32 {
	if min == max {
		return (min)
	}

	random := rand.Float32()

	return ((random * (max - min)) + min)

}

type ParticleSystem struct {
	Name     string
	Id       string
	Capacity int

	// Vectors and colors
	Gravity     *math32.Vector3
	Direction1  *math32.Vector3
	Direction2  *math32.Vector3
	MinEmitBox  *math32.Vector3
	MaxEmitBox  *math32.Vector3
	Color1      *math32.Color4
	Color2      *math32.Color4
	ColorDead   *math32.Color4
	DeadAlpha   float32
	TextureMask *math32.Color4

	Particles       []*Particle
	_stockParticles []*Particle
	_newPartsExcess int

	_vertexDeclaration []int
	_vertexStrideSize  int
	_vertexBuffer      *gl.GLVertexBuffer
	_vertices          []float32

	_indexBuffer       *gl.GLIndexBuffer
	_scene             *engines.Scene
	_scaledUpdateSpeed float32
	_colorDiff         *math32.Color4
	_scaledColorStep   *math32.Color4
	_scaledDirection   *math32.Vector3
	_scaledGravity     *math32.Vector3

	_alive       bool
	_started     bool
	_stopped     bool
	_actualFrame float32

	Emitter            IMesh
	EmitRate           float32
	ManualEmitCount    float32
	UpdateSpeed        float32
	TargetStopDuration float32
	DisposeOnStop      bool
	MinEmitPower       float32
	MaxEmitPower       float32

	MinLifeTime     float32
	MaxLifeTime     float32
	MinSize         float32
	MaxSize         float32
	MinAngularSpeed float32
	MaxAngularSpeed float32
	ParticleTexture ITexture
	OnDispose       func()
	BlendMode       PARTICLEENUM

	//effect
	_effect IEffect
	//cache
	_cachedDefines string
}

func NewParticleSystem(name string, capacity int, scene *engines.Scene) *ParticleSystem {
	this := &ParticleSystem{}

	this.Name = name
	this.Id = name
	this.Capacity = capacity

	this._scene = scene

	this._scene.ParticleSystems = append(this._scene.ParticleSystems, this)

	this.Init()

	return this
}

func (this *ParticleSystem) Init() {

	this._scaledUpdateSpeed = 0.0

	this._alive = false
	this._started = false
	this._stopped = true
	this._actualFrame = 0

	// Vectors and colors
	this.Gravity = math32.NewVector3Zero()
	this.Direction1 = math32.NewVector3(0, 1.0, 0)
	this.Direction2 = math32.NewVector3(0, 1.0, 0)
	this.MinEmitBox = math32.NewVector3(-0.5, -0.5, -0.5)
	this.MaxEmitBox = math32.NewVector3(0.5, 0.5, 0.5)
	this.Color1 = math32.NewColor4(1.0, 1.0, 1.0, 1.0)
	this.Color2 = math32.NewColor4(1.0, 1.0, 1.0, 1.0)
	this.ColorDead = math32.NewColor4(0, 0, 0, 1.0)
	this.DeadAlpha = 0
	this.TextureMask = math32.NewColor4(1.0, 1.0, 1.0, 1.0)

	// Particles
	this.Particles = make([]*Particle, 0)
	this._stockParticles = make([]*Particle, 0)
	this._newPartsExcess = 0

	// VBO
	this._vertexDeclaration = []int{3, 4, 4}
	this._vertexStrideSize = 11 * 4 // 10 floats per particle (x, y, z, r, g, b, a, angle, size, offsetX, offsetY)
	this._vertexBuffer = this._scene.GetEngine().CreateDynamicVertexBuffer(this.Capacity * this._vertexStrideSize * 4)

	indices := make([]uint16, 0)
	var index uint16
	index = 0
	for count := 0; count < this.Capacity; count++ {
		indices = append(indices, index)
		indices = append(indices, index+1)
		indices = append(indices, index+2)
		indices = append(indices, index)
		indices = append(indices, index+2)
		indices = append(indices, index+3)
		index += 4
	}

	this._indexBuffer = this._scene.GetEngine().CreateIndexBuffer(indices, false)
	this._vertices = make([]float32, this.Capacity*this._vertexStrideSize)

	this.Emitter = nil
	this.EmitRate = 10
	this.ManualEmitCount = -1
	this.UpdateSpeed = 0.01
	this.TargetStopDuration = 0
	this.DisposeOnStop = false
	this.MinEmitPower = 1
	this.MaxEmitPower = 1

	this.MinLifeTime = 1
	this.MaxLifeTime = 1
	this.MinSize = 1
	this.MaxSize = 1
	this.MinAngularSpeed = 0
	this.MaxAngularSpeed = 0
	this.ParticleTexture = nil
	this.OnDispose = nil
	this.BlendMode = BLENDMODE_ONEONE

	// Internals
	this._scaledColorStep = math32.NewColor4(0, 0, 0, 0)
	this._colorDiff = math32.NewColor4(0, 0, 0, 0)
	this._scaledDirection = math32.NewVector3Zero()
	this._scaledGravity = math32.NewVector3Zero()

}

func (this *ParticleSystem) IsAlive() bool {
	return this._alive
}
func (this *ParticleSystem) Start() {
	this._started = true
	this._stopped = false
	this._actualFrame = 0
}
func (this *ParticleSystem) Stop() {
	this._stopped = true
}

func (this *ParticleSystem) _appendParticleVertex(index int, particle *Particle, offsetX, offsetY float32) {
	offset := index * 11
	this._vertices[offset] = particle.Position.X
	this._vertices[offset+1] = particle.Position.Y
	this._vertices[offset+2] = particle.Position.Z
	this._vertices[offset+3] = particle.Color.R
	this._vertices[offset+4] = particle.Color.G
	this._vertices[offset+5] = particle.Color.B
	this._vertices[offset+6] = particle.Color.A
	this._vertices[offset+7] = particle.Angle
	this._vertices[offset+8] = particle.Size
	this._vertices[offset+9] = offsetX
	this._vertices[offset+10] = offsetY
}

func (this *ParticleSystem) _update(newParticles int) {

	this._alive = len(this.Particles) > 0
	for index := 0; index < len(this.Particles); index++ {
		particle := this.Particles[index]
		particle.Age += this._scaledUpdateSpeed

		if particle.Age >= particle.LifeTime {
			this.Particles = append(this.Particles[:index], this.Particles[index+1:]...)
			this._stockParticles = append(this._stockParticles, particle)
			index--
			continue
		} else {
			this._scaledColorStep = particle.ColorStep.Scale(this._scaledUpdateSpeed)
			particle.Color = particle.Color.Add(this._scaledColorStep)

			if particle.Color.A < 0 {
				particle.Color.A = 0
			}

			this._scaledDirection = particle.Direction.Scale(this._scaledUpdateSpeed)
			particle.Position = particle.Position.Add(this._scaledDirection)

			particle.Angle += particle.AngularSpeed * this._scaledUpdateSpeed

			this._scaledGravity = this.Gravity.Scale(this._scaledUpdateSpeed)
			particle.Direction = particle.Direction.Add(this._scaledGravity)

		}
	}

	// Add new ones
	worldMatrix := this.Emitter.GetWorldMatrix()

	for index := 0; index < newParticles; index++ {
		if len(this.Particles) == this.Capacity {
			break
		}

		var particle *Particle
		particle = NewParticle()

		this.Particles = append(this.Particles, particle)

		emitPower := randomNumber(this.MinEmitPower, this.MaxEmitPower)

		randX := randomNumber(this.Direction1.X, this.Direction2.X)
		randY := randomNumber(this.Direction1.Y, this.Direction2.Y)
		randZ := randomNumber(this.Direction1.Z, this.Direction2.Z)

		particle.Direction = math32.NewVector3Zero().TransformNormalFromFloats(randX*emitPower, randY*emitPower, randZ*emitPower, worldMatrix)

		particle.LifeTime = randomNumber(this.MinLifeTime, this.MaxLifeTime)

		particle.Size = randomNumber(this.MinSize, this.MaxSize)
		particle.AngularSpeed = randomNumber(this.MinAngularSpeed, this.MaxAngularSpeed)

		randX = randomNumber(this.MinEmitBox.X, this.MaxEmitBox.X)
		randY = randomNumber(this.MinEmitBox.Y, this.MaxEmitBox.Y)
		randZ = randomNumber(this.MinEmitBox.Z, this.MaxEmitBox.Z)

		particle.Position = math32.NewVector3Zero().TransformCoordinatesFromFloats(randX, randY, randZ, worldMatrix)

		step := randomNumber(0, 1.0)

		particle.Color = this.Color1.Lerp(this.Color2, step)

		this._colorDiff = this.ColorDead.Sub(particle.Color)
		particle.ColorStep = this._colorDiff.Scale(1.0 / particle.LifeTime)

	}
}

func (this *ParticleSystem) _getEffect() IEffect {
	defines := make([]string, 0)
	defines = append(defines, "#define EMPTYDEFINED")

	if core.GlobalFly3D.ClipPlane != nil {
		defines = append(defines, "#define CLIPPLANE")
	}

	// Effect
	join := strings.Join(defines, "\n")
	if this._cachedDefines != join {
		this._cachedDefines = join
		this._effect = effects.CreateEffect(this._scene.GetEngine(),
			"particles",
			[]string{"position", "color", "options"},
			[]string{"invView", "view", "projection", "vClipPlane", "textureMask"},
			[]string{"diffuseSampler"}, join)
	}

	return this._effect
}

/*IParticleSystem interface start*/

func (this *ParticleSystem) Animate() {
	if !this._started {
		return
	}

	effect := this._getEffect()
	if effect == nil {
		return
	}

	// Check
	if this.Emitter == nil ||
		!effect.IsReady() ||
		this.ParticleTexture == nil ||
		!this.ParticleTexture.IsReady() {
		return
	}

	this._scaledUpdateSpeed = this.UpdateSpeed * this._scene.GetAnimationRatio()

	// determine the number of particles we need to create
	var emitCout float32

	if this.ManualEmitCount > -1 {
		emitCout = this.ManualEmitCount
		this.ManualEmitCount = 0
	} else {
		emitCout = this.EmitRate
	}

	newParticles := (int(emitCout*this._scaledUpdateSpeed) >> 0)
	this._newPartsExcess = (int)(float32(this._newPartsExcess) + emitCout*this._scaledUpdateSpeed - float32(newParticles))

	if this._newPartsExcess > 1.0 {
		newParticles += this._newPartsExcess >> 0
		this._newPartsExcess -= this._newPartsExcess >> 0
	}

	this._alive = false

	if !this._stopped {
		this._actualFrame += this._scaledUpdateSpeed

		if this.TargetStopDuration != 0.0 && this._actualFrame >= this.TargetStopDuration {
			this.Stop()
		}

	} else {
		newParticles = 0
	}

	this._update(newParticles)

	// Stopped
	if this._stopped {
		if !this._alive {
			this._started = false
			if this.DisposeOnStop {
				this._scene.ToBeDisposed = append(this._scene.ToBeDisposed, this)
			}
		}
	}
	if len(this.Particles) == 0 {
		return
	}

	// Update VBO
	offset := 0
	this._vertices = make([]float32, len(this.Particles)*this._vertexStrideSize)

	for index := 0; index < len(this.Particles); index++ {
		particle := this.Particles[index]

		this._appendParticleVertex(offset, particle, 0, 0)
		offset++

		this._appendParticleVertex(offset, particle, 1, 0)
		offset++

		this._appendParticleVertex(offset, particle, 1, 1)
		offset++
		this._appendParticleVertex(offset, particle, 0, 1)
		offset++
	}
	engine := this._scene.GetEngine()
	engine.UpdateDynamicVertexBuffer(this._vertexBuffer, this._vertices)
}
func (this *ParticleSystem) GetEmitter() IMesh {
	return this.Emitter
}
func (this *ParticleSystem) Render() int {

	effect := this._getEffect()

	// Check
	if this.Emitter == nil ||
		!effect.IsReady() ||
		this.ParticleTexture == nil ||
		!this.ParticleTexture.IsReady() {
		return 0
	}
	engine := this._scene.GetEngine()

	// Render
	engine.EnableEffect(effect)

	viewMatrix := this._scene.GetViewMatrix()
	effect.SetTexture("diffuseSampler", this.ParticleTexture.GetGLTexture())
	effect.SetMatrix("view", viewMatrix)
	effect.SetMatrix("projection", this._scene.GetProjectionMatrix())
	effect.SetFloat4("textureMask", this.TextureMask.R, this.TextureMask.G, this.TextureMask.B, this.TextureMask.A)

	clipplane := core.GlobalFly3D.ClipPlane
	if core.GlobalFly3D.ClipPlane != nil {
		invView := viewMatrix.Clone()
		invView.Invert()

		effect.SetMatrix("invView", invView)
		effect.SetFloat4("vClipPlane", clipplane.Normal.X, clipplane.Normal.Y, clipplane.Normal.Z, clipplane.D)
	}

	// VBOs
	engine.BindBuffers(this._vertexBuffer, this._indexBuffer, this._vertexDeclaration, this._vertexStrideSize, effect)

	// Draw order
	if this.BlendMode == BLENDMODE_ONEONE {
		engine.SetAlphaMode(core.ALPHA_ADD)
	} else {
		engine.SetAlphaMode(core.ALPHA_COMBINE)
	}
	engine.Draw(true, 0, len(this.Particles)*6)
	engine.SetAlphaMode(core.ALPHA_DISABLE)

	return len(this.Particles)
}
func (this *ParticleSystem) Dispose() {
	if this._vertexBuffer != nil {
		this._scene.GetEngine().ReleaseVertexBuffer(this._vertexBuffer)
		this._vertexBuffer = nil
	}

	if this._indexBuffer != nil {
		this._scene.GetEngine().ReleaseIndexBuffer(this._indexBuffer)
		this._indexBuffer = nil
	}

	if this.ParticleTexture != nil {
		this.ParticleTexture.Dispose()
		this.ParticleTexture = nil
	}

	// Remove from scene
	index := tools.IndexOf(this, this._scene.ParticleSystems)
	if index > -1 {
		this._scene.ParticleSystems = append(this._scene.ParticleSystems[:index], this._scene.ParticleSystems[index+1:]...)
	}

	// Callback
	if this.OnDispose != nil {
		this.OnDispose()
	}
}

/*IParticleSystem interface end*/
