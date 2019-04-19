package animations

import (
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/tools/reflections"
	"fmt"
)

const (
	ANIMATIONTYPE_FLOAT      = 0
	ANIMATIONTYPE_VECTOR3    = 1
	ANIMATIONTYPE_QUATERNION = 2
	ANIMATIONTYPE_MATRIX     = 3

	ANIMATIONLOOPMODE_RELATIVE = 0
	ANIMATIONLOOPMODE_CYCLE    = 1
	ANIMATIONLOOPMODE_CONSTANT = 2
)

type AnimationKeyFrame struct {
	frame float32
	value interface{} // Vector3 or Quaternion or Matrix or Float
}

type Animation struct {
	Name             string
	TargetProperty   string
	FramePerSecond   float32
	DataType         int
	LoopMode         int
	CurrentFrame     float32
	_keys            []*AnimationKeyFrame
	_offsetsCache    map[string]interface{}
	_highLimitsCache map[string]interface{}
}

//loopmodel = -1
func NewAnimation(name string, targetProperty string, framePerSecond float32, dataType int, loopMode int) *Animation {
	this := &Animation{}

	this.Init(name, targetProperty, framePerSecond, dataType, loopMode)
	return this
}
func (this *Animation) Init(name string, targetProperty string, framePerSecond float32, dataType int, loopMode int) {
	this.Name = name
	this.TargetProperty = targetProperty
	this.FramePerSecond = framePerSecond
	this.DataType = dataType

	if loopMode == -1 {
		this.LoopMode = ANIMATIONLOOPMODE_CYCLE
	} else {
		this.LoopMode = loopMode
	}

	this._keys = make([]*AnimationKeyFrame, 0)

	// Cache
	this._offsetsCache = map[string]interface{}{}
	this._highLimitsCache = map[string]interface{}{}
}

// Methods
func (this *Animation) Clone() *Animation {
	clone := NewAnimation(this.Name, this.TargetProperty, this.FramePerSecond, this.DataType, this.LoopMode)

	clone.SetKeys(this._keys)

	return clone
}

func (this *Animation) SetKeys(values []*AnimationKeyFrame) {
	this._keys = values
	this._offsetsCache = map[string]interface{}{}
	this._highLimitsCache = map[string]interface{}{}
}

func (this *Animation) _interpolate(currentFrame float32, repeatCount int, loopMode int, offsetValue_Obj interface{}, highLimitValue_Obj interface{}) interface{} {

	if loopMode == ANIMATIONLOOPMODE_CONSTANT && repeatCount > 0 {
		return highLimitValue_Obj
	}
	for key := 0; key < len(this._keys); key++ {
		if this._keys[key+1].frame >= currentFrame {
			startValue_obj := this._keys[key].value
			endValue_obj := this._keys[key+1].value
			gradient := (float32)(currentFrame-this._keys[key].frame) / (float32)(this._keys[key+1].frame-this._keys[key].frame)
			/*
				if (this._easingFunction != nil) {
							gradient = this._easingFunction.ease(gradient);
					}
			*/
			switch this.DataType {
			// Float
			case ANIMATIONTYPE_FLOAT:
				startValue, _ := startValue_obj.(float32)
				endValue, _ := endValue_obj.(float32)

				var offsetValue float32
				if offsetValue_Obj != nil {
					offsetValue, _ = offsetValue_Obj.(float32)
				} else {
					offsetValue = 0.0
				}

				switch loopMode {
				case ANIMATIONLOOPMODE_CYCLE:
				case ANIMATIONLOOPMODE_CONSTANT:
					return startValue + (endValue-startValue)*gradient
				case ANIMATIONLOOPMODE_RELATIVE:
					return offsetValue*float32(repeatCount) + (startValue + (endValue-startValue)*gradient)
				}
				break
				// Quaternion
			case ANIMATIONTYPE_QUATERNION:
				var quaternion *math32.Quaternion
				startValue, _ := startValue_obj.(*math32.Quaternion)
				endValue, _ := endValue_obj.(*math32.Quaternion)
				var offsetValue *math32.Quaternion
				if offsetValue_Obj != nil {
					offsetValue, _ = offsetValue_Obj.(*math32.Quaternion)
				} else {
					offsetValue = math32.NewQuaternion(0, 0, 0, 0)
				}
				switch loopMode {
				case ANIMATIONLOOPMODE_CYCLE:
				case ANIMATIONLOOPMODE_CONSTANT:
					quaternion = startValue.Slerp(endValue, gradient)
					break
				case ANIMATIONLOOPMODE_RELATIVE:
					quaternion = startValue.Slerp(endValue, gradient).Add(offsetValue.Scale(float32(repeatCount)))
					break
				}

				return quaternion
			// Vector3
			case ANIMATIONTYPE_VECTOR3:

				startValue, _ := startValue_obj.(*math32.Vector3)
				endValue, _ := endValue_obj.(*math32.Vector3)

				var offsetValue *math32.Vector3
				if offsetValue_Obj != nil {
					offsetValue, _ = offsetValue_Obj.(*math32.Vector3)
				} else {
					offsetValue = math32.NewVector3(0, 0, 0)
				}
				switch loopMode {
				case ANIMATIONLOOPMODE_CYCLE:
				case ANIMATIONLOOPMODE_CONSTANT:
					return startValue.Lerp(endValue, gradient)
				case ANIMATIONLOOPMODE_RELATIVE:
					return startValue.Lerp(endValue, gradient).Add(offsetValue.Scale(float32(repeatCount)))
				}
			default:
				break
			}
			break

		}
	}

	return this._keys[len(this._keys)-1].value
}

//interface
/*
type IAnimation interface {
	Animate(target IAnimationTarget, delay float32, from float32, to float32, loop bool, speedRatio float32) bool
}
*/
func (this *Animation) Animate(target IAnimationTarget, delay float32, from float32, to float32, loop bool, speedRatio float32) bool {

	if this.TargetProperty == "" {
		return false
	}

	// Check limits
	if from < this._keys[0].frame || from > this._keys[len(this._keys)-1].frame {
		from = this._keys[0].frame
	}
	if to < this._keys[0].frame || to > this._keys[len(this._keys)-1].frame {
		to = this._keys[len(this._keys)-1].frame
	}

	// Compute ratio
	rangeval := to - from
	ratio := delay * float32(this.FramePerSecond*speedRatio) / 1000.0

	if ratio > rangeval && !loop { // If we are out of range and not looping get back to caller
		return false
	}
	var offsetValue interface{}
	var highLimitValue interface{}
	if this.LoopMode != ANIMATIONLOOPMODE_CYCLE {
		keyOffset := fmt.Sprintf("form[%f]to[%f]", from, to)
		_, ok := this._offsetsCache[keyOffset]
		if !ok {

			fromValue_obj := this._interpolate(from, 0, ANIMATIONLOOPMODE_CYCLE, nil, nil)
			toValue_obj := this._interpolate(to, 0, ANIMATIONLOOPMODE_CYCLE, nil, nil)

			switch this.DataType {
			// Float
			case ANIMATIONTYPE_FLOAT:
				toValue, _ := toValue_obj.(float32)
				fromValue, _ := fromValue_obj.(float32)
				this._offsetsCache[keyOffset] = toValue - fromValue
				break
			// Quaternion
			case ANIMATIONTYPE_QUATERNION:
				toValue, _ := toValue_obj.(*math32.Quaternion)
				fromValue, _ := fromValue_obj.(*math32.Quaternion)
				this._offsetsCache[keyOffset] = toValue.Sub(fromValue)
				break
			// Vector3
			case ANIMATIONTYPE_VECTOR3:
				toValue, _ := toValue_obj.(*math32.Vector3)
				fromValue, _ := fromValue_obj.(*math32.Vector3)
				this._offsetsCache[keyOffset] = toValue.Sub(fromValue)
			default:
				break
			}

			this._highLimitsCache[keyOffset] = toValue_obj

		}

		highLimitValue, _ = this._highLimitsCache[keyOffset]
		offsetValue, _ = this._offsetsCache[keyOffset]
	}
	// Compute value
	repeatCount := int(ratio/rangeval) >> 0
	currentFrame := from + (float32)(int(ratio)%int(rangeval))
	currentValue := this._interpolate(currentFrame, repeatCount, this.LoopMode, offsetValue, highLimitValue)

	// Set value
	reflections.SetField(target, this.TargetProperty, currentValue)

	return true
}
