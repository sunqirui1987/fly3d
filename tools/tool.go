package tools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/suiqirui1987/fly3d/math32"
)

//implement interval multi-timer
func SetInterval(d time.Duration, f func(args ...interface{}), args ...interface{}) *Timer {
	t := &Timer{
		D:     d,
		F:     f,
		Args:  args,
		max:   -1,
		count: 0,
	}

	t.Start()

	return t
}

//implement timeout multi-timer
func SetTimeout(d time.Duration, f func(args ...interface{}), args ...interface{}) *Timer {
	t := &Timer{
		D:     d,
		F:     f,
		Args:  args,
		max:   1,
		count: 0,
	}

	t.Start()

	return t
}

// Snow build gives an error that haxe.Timer has no delay method...
func Delay(f func(args ...interface{}), time_ms int) {

	d := time.Duration(time_ms) * time.Millisecond
	t := &Timer{
		D:     d,
		F:     f,
		Args:  nil,
		max:   1,
		count: 0,
	}

	t.Start()

}
func GetCurrentTimeMs() int {
	return (int)(time.Now().UnixNano() / 1e6)
}

func DecodeImage(content []byte) (*image.RGBA, error) {

	// Decodes image
	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	// Converts image to RGBA format
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba, nil
}

func IndexOf(params ...interface{}) int {
	v := reflect.ValueOf(params[0])
	arr := reflect.ValueOf(params[1])

	var t = reflect.TypeOf(params[1]).Kind()

	if t != reflect.Slice && t != reflect.Array {
		panic("Type Error! Second argument must be an array or a slice.")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == v.Interface() {
			return i
		}
	}
	return -1
}

func BytesUint16(byteOrder binary.ByteOrder, values ...uint16) []byte {
	le := false
	switch byteOrder {
	case binary.BigEndian:
	case binary.LittleEndian:
		le = true
	default:
		panic(fmt.Sprintf("invalid byte order %v", byteOrder))
	}

	b := make([]byte, 2*len(values))
	for i, v := range values {
		u := v
		if le {
			b[2*i+0] = byte(u >> 0)
			b[2*i+1] = byte(u >> 8)
		} else {
			b[2*i+0] = byte(u >> 8)
			b[2*i+1] = byte(u >> 0)
		}
	}
	return b
}

func OpenGeneralFile(file string) ([]byte, error) {
	if strings.Index(file, "http") == 0 {
		//http
		return DownHttpFile(file)
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func DownHttpFile(url string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ExtractMinAndMax(positions []float32, start int, count int) (*math32.Vector3, *math32.Vector3) {
	minimum := math32.NewVector3(math.MaxFloat32, math.MaxFloat32, math.MaxFloat32)
	maximum := math32.NewVector3(-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32)

	for index := start; index < start+count; index++ {
		current := math32.NewVector3(positions[index*3], positions[index*3+1], positions[index*3+2])

		minimum = current.Min(minimum)
		maximum = current.Max(maximum)
	}

	return minimum, maximum
}

// Clean returns b with the 3 byte BOM stripped off the front if it is present.
// If the BOM is not present, then b is returned.
func Clean(b []byte) []byte {
	if len(b) >= 3 &&
		b[0] == 0xef &&
		b[1] == 0xbb &&
		b[2] == 0xbf {
		return b[3:]
	}
	return b
}
