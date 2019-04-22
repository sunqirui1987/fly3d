package animations

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	log "github.com/suiqirui1987/fly3d/tools/logrus"
	"github.com/suiqirui1987/fly3d/tools/reflections"
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

func getFloat(unk interface{}) (float32, bool) {
	switch i := unk.(type) {
	case float64:
		return float32(i), true
	case float32:
		return float32(i), true
	case int64:
		return float32(i), true
	case int32:
		return float32(i), true
	case int16:
		return float32(i), true
	case int8:
		return float32(i), true
	case uint64:
		return float32(i), true
	case uint32:
		return float32(i), true
	case uint16:
		return float32(i), true
	case uint8:
		return float32(i), true
	case int:
		return float32(i), true
	case uint:
		return float32(i), true
	case string:
		f, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return 0, false
		}
		return float32(f), true
	default:
		return 0, false
	}
}

type AnimationKeyFrame struct {
	Frame float32
	Value interface{} // Vector3 or Quaternion or Matrix or Float
}

type Animation struct {
	Name string

	FramePerSecond float32
	DataType       int
	LoopMode       int
	CurrentFrame   float32

	_targetProperty     string
	_targetPropertyPath []string
	_keys               []*AnimationKeyFrame
	_offsetsCache       map[string]interface{}
	_highLimitsCache    map[string]interface{}
}

//loopmodel = -1
func NewAnimation(name string, targetProperty string, framePerSecond float32, dataType int, loopMode int) *Animation {
	this := &Animation{}

	this.Init(name, targetProperty, framePerSecond, dataType, loopMode)
	return this
}
func (this *Animation) Init(name string, targetProperty string, framePerSecond float32, dataType int, loopMode int) {
	this.Name = name

	this.FramePerSecond = framePerSecond
	this.DataType = dataType

	this._targetProperty = targetProperty
	this._targetPropertyPath = strings.Split(targetProperty, ".")

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
	clone := NewAnimation(this.Name, this._targetProperty, this.FramePerSecond, this.DataType, this.LoopMode)

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
		if this._keys[key+1].Frame >= currentFrame {
			startValue_obj := this._keys[key].Value
			endValue_obj := this._keys[key+1].Value
			gradient := (float32)(currentFrame-this._keys[key].Frame) / (float32)(this._keys[key+1].Frame-this._keys[key].Frame)
			/*
				if (this._easingFunction != nil) {
							gradient = this._easingFunction.ease(gradient);
					}
			*/
			switch this.DataType {
			// Float
			case ANIMATIONTYPE_FLOAT:
				startValue, ok := getFloat(startValue_obj)
				if !ok {
					log.Printf("_interpolate The interface type is incorrect, request float32 ")
					return 0.0
				}
				endValue, ok := getFloat(endValue_obj)
				if !ok {
					log.Printf("_interpolate The interface type is incorrect, request float32 ")
					return 0.0
				}

				var offsetValue float32
				if offsetValue_Obj != nil {
					offsetValue, _ = offsetValue_Obj.(float32)
				} else {
					offsetValue = 0.0
				}

				switch loopMode {
				case ANIMATIONLOOPMODE_CYCLE:
					return startValue + (endValue-startValue)*gradient
				case ANIMATIONLOOPMODE_CONSTANT:
					return startValue + (endValue-startValue)*gradient
				case ANIMATIONLOOPMODE_RELATIVE:
					return offsetValue*float32(repeatCount) + (startValue + (endValue-startValue)*gradient)
				}
				break
				// Quaternion
			case ANIMATIONTYPE_QUATERNION:
				var quaternion *math32.Quaternion
				startValue, ok := startValue_obj.(*math32.Quaternion)
				if !ok {
					log.Printf("_interpolate The interface type is incorrect, request Quaternion ")
					return math32.NewQuaternionZero()
				}
				endValue, ok := endValue_obj.(*math32.Quaternion)
				if !ok {
					log.Printf("_interpolate The interface type is incorrect, request Quaternion ")
					return math32.NewQuaternionZero()
				}

				var offsetValue *math32.Quaternion
				if offsetValue_Obj != nil {
					offsetValue, _ = offsetValue_Obj.(*math32.Quaternion)
				} else {
					offsetValue = math32.NewQuaternion(0, 0, 0, 0)
				}
				switch loopMode {
				case ANIMATIONLOOPMODE_CYCLE:
					quaternion = startValue.Slerp(endValue, gradient)
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

				startValue, ok := startValue_obj.(*math32.Vector3)
				if !ok {
					log.Printf("_interpolate The interface type is incorrect, request Quaternion ")
					return math32.NewVector3Zero()
				}
				endValue, ok := endValue_obj.(*math32.Vector3)
				if !ok {
					log.Printf("_interpolate The interface type is incorrect, request Quaternion ")
					return math32.NewVector3Zero()
				}

				var offsetValue *math32.Vector3
				if offsetValue_Obj != nil {
					offsetValue, _ = offsetValue_Obj.(*math32.Vector3)
				} else {
					offsetValue = math32.NewVector3(0, 0, 0)
				}
				switch loopMode {
				case ANIMATIONLOOPMODE_CYCLE:
					return startValue.Lerp(endValue, gradient)
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

	return this._keys[len(this._keys)-1].Value
}

//interface
/*
type IAnimation interface {
	Animate(target IAnimationTarget, delay float32, from float32, to float32, loop bool, speedRatio float32) bool
}
*/
func (this *Animation) Animate(target IAnimationTarget, delay float32, from float32, to float32, loop bool, speedRatio float32) bool {

	if this._targetProperty == "" || len(this._targetPropertyPath) < 1 {
		return false
	}

	if len(this._keys) == 0 {
		return false
	}

	// Check limits
	if from < this._keys[0].Frame || from > this._keys[len(this._keys)-1].Frame {
		from = this._keys[0].Frame
	}
	if to < this._keys[0].Frame || to > this._keys[len(this._keys)-1].Frame {
		to = this._keys[len(this._keys)-1].Frame
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
	if len(this._targetPropertyPath) > 1 {
		property, err := reflections.GetField(target, this._targetPropertyPath[0])
		if err != nil {
			log.Printf("Animate reflections.GetField %s", property)
			return false
		}

		for index := 1; index < len(this._targetPropertyPath)-1; index++ {
			property, err = reflections.GetField(property, this._targetPropertyPath[index])
			if err != nil {
				log.Printf("Animate reflections.GetField %s", property)
				return false
			}
		}

		valname := this._targetPropertyPath[len(this._targetPropertyPath)-1]
		reflections.SetField(property, valname, currentValue)
	} else {
		reflections.SetField(target, this._targetProperty, currentValue)
	}

	return true
}
