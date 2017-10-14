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
	flagWeb := flag.Bool("web", false, "show faces on web page localhost:4000")
	flag.Parse()
	i := 0
	imgChan := make(chan []interface{}, 1)
	go faces.ScanDir(flag.Arg(0), true, func(filename string, faces []image.Image, rects []image.Rectangle) {
		fmt.Printf("image %s has %d faces in it\n", filename, len(rects))
		imgs := []string{}
		for _, face := range faces {
			if *flagWeb {
				img := face2base64(face)
				imgs = append(imgs, img)
			} else {
				saveFace(fmt.Sprintf("%s/face_%d", flag.Arg(1), i), face)
				i++
			}
		}
		if *flagWeb {
			imgChan <- []interface{}{filename, imgs}
		}
	})
	serveHTTP(imgChan)
}
