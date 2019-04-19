package animations

import (
	"github.com/suiqirui1987/fly3d/engines"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"reflect"
)

// Animations
// speedRatio 1.0
func BeginAnimation(scene *engines.Scene, target IAnimationTarget, from float32, to float32, loop bool, speedRatio float32) {

	// Local animations
	if target.GetAnimations() != nil {

		StopAnimation(scene, target)

		animatable := NewAnimatable(target, from, to, loop, speedRatio)
		scene.ActiveAnimatables = append(scene.ActiveAnimatables, animatable)
	}

	// Children animations
	if target.GetAnimatables() != nil {
		animatables := target.GetAnimatables()
		for index := 0; index < len(animatables); index++ {
			BeginAnimation(scene, animatables[index].GetTarget(), from, to, loop, speedRatio)
		}
	}
}
func StopAnimation(scene *engines.Scene, target IAnimationTarget) {
	for index := 0; index < len(scene.ActiveAnimatables); index++ {
		if reflect.DeepEqual(scene.ActiveAnimatables[index].GetTarget(), target) {
			scene.ActiveAnimatables = append(scene.ActiveAnimatables[:index], scene.ActiveAnimatables[index+1:]...)
			return
		}
	}
}
