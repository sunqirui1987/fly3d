package interfaces

type IAnimatable interface {
	GetTarget() IAnimationTarget
	Animate() bool
}

type IAnimation interface {
	Animate(target IAnimationTarget, delay float32, from float32, to float32, loop bool, speedRatio float32) bool
}

type IAnimationTarget interface {
	GetAnimations() []IAnimation
	GetAnimatables() []IAnimatable
}
