package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"sync"

	"faces"
)

func face2base64(img image.Image) string {
	b := &bytes.Buffer{}
	if err := png.Encode(b, img); err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func saveFace(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	result := &sync.Map{}
	i := 0
	faces.ScanDir(flag.Arg(0), true, func(filename string, faces []image.Image, rects []image.Rectangle) {
		fmt.Printf("image %s has %d faces in it\n", filename, len(rects))
		for _, face := range faces {
			saveFace(fmt.Sprintf("%s/face_%d.png", flag.Arg(1), i), face)
			i++
		}
		result.Store(filename, faces)
	})
}

var exitChan = make(chan bool, 1)

func showGui() {
	go serveHTTP()
	<-exitChan
}
