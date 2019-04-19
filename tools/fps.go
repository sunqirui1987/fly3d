package tools

import (
	"time"
)

var (
	fpsRange               float32   = 60.0
	previousFramesDuration []float32 = make([]float32, 0)
	fps                    float32   = 60.0
	deltaTime              float32   = 0.0
)

func GetFps() float32 {
	return fps
}

func GetDeltaTime() float32 {
	return deltaTime
}

func MeasureFps() {

	previousFramesDuration = append(previousFramesDuration, float32(time.Now().UnixNano())/1e6)
	length := len(previousFramesDuration)
	if length >= 2 {
		deltaTime = previousFramesDuration[length-1] - previousFramesDuration[length-2]
	}
	if float32(length) >= fpsRange {
		if float32(length) > fpsRange {
			previousFramesDuration = previousFramesDuration[1:]
			length = len(previousFramesDuration)
		}

		var sum float32
		sum = 0
		for id := 0; id < length-1; id++ {
			sum += previousFramesDuration[id+1] - previousFramesDuration[id]
		}

		fps = 1000.0 / (sum / float32(length-1))
	}

}
