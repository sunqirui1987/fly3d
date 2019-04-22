package tools

import (
	"time"
)

var (
	fpsRange               float32 = 60.0
	previousFramesDuration []int64 = make([]int64, 0)
	fps                    float32 = 60.0
	deltaTime              float32 = 0.0
)

func GetFps() float32 {
	return fps
}

func GetDeltaTime() float32 {
	return deltaTime
}

func MeasureFps() {

	millis := time.Now().UnixNano() / 1000000
	previousFramesDuration = append(previousFramesDuration, millis)
	length := len(previousFramesDuration)
	if length >= 2 {
		deltaTime = float32(previousFramesDuration[length-1] - previousFramesDuration[length-2])
	}
	if float32(length) >= fpsRange {
		if float32(length) > fpsRange {
			previousFramesDuration = previousFramesDuration[1:]
			length = len(previousFramesDuration)
		}

		var sum int64
		sum = 0
		for id := 0; id < length-1; id++ {
			sum += previousFramesDuration[id+1] - previousFramesDuration[id]
		}

		fps = 1000.0 / float32(sum/int64(length-1))
	}

}
