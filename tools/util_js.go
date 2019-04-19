// +build js
package tools

import (
	"image"
	"log"
)

func LoadImage(url string, onload func(*image.RGBA), onfail func(error)) {

	content, err := DownHttpFile(url)
	if err != nil {
		log.Printf("LoadImage  DownHttpFile Failed %s", err)
		onfail(err)
		return
	}
	rgba, err := DecodeImage(content)
	if err != nil {
		log.Printf("LoadImage DecodeImage Failed %s", err)
		onfail(err)
		return
	}
	onload(rgba)
	return
}

func LoadFile(url string, callback func(string), progressCallBack func(int)) {
	content, err := DownHttpFile(url)
	if err != nil {
		log.Printf("LoadFile Failed %s", err)
		return
	}
	content_str := string(Clean(content))
	callback(content_str)
}
