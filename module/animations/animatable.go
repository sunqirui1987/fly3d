package animations

import (
	. "github.com/suiqirui1987/fly3d/interfaces"
	"time"
)

type Animatable struct {
	_target IAnimationTarget

	FromFrame float32
	ToFrame   float32

	LoopAnimation        bool
	AnimationStartedDate int64

	SpeedRatio float32
}

func NewAnimatable(target IAnimationTarget, from float32, to float32, loop bool, speedRatio float32) *Animatable {

	this := &Animatable{
		SpeedRatio: 1.0,
	}
	this.Init(target, from, to, loop, speedRatio)
	return this
}

func (this *Animatable) Init(target IAnimationTarget, from float32, to float32, loop bool, speedRatio float32) {

	this._target = target
	this.FromFrame = from
	this.ToFrame = to
	this.LoopAnimation = loop
	this.SpeedRatio = speedRatio

	this.AnimationStartedDate = time.Now().Unix()

}

//interface
/*
type IAnimatable interface {
	GetTarget() IAnimationTarget
	Animate() bool
}
*/
func (this *Animatable) GetTarget() IAnimationTarget {
	return this._target
}

func (this *Animatable) Animate() bool {

	//Getting time
	var delay float32
	delay = (float32)(time.Now().Unix() - this.AnimationStartedDate)

	// Animating
	running := false
	animations := this._target.GetAnimations()
	for i := 0; i < len(animations); i++ {
		animation, ok := animations[i].(IAnimation)
		if ok {
			isRunning := animation.Animate(this._target, delay, this.FromFrame, this.ToFrame, this.LoopAnimation, this.SpeedRatio)
			running = running || isRunning
		}

	}

	return running
}
