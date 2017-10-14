package faces

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/napsy/go-opencv/opencv"
)

var cascade *opencv.HaarCascade

func init() {
	cascade = opencv.LoadHaarClassifierCascade("haarcascade_frontalface_alt.xml")
}

// OpenCV has problems if running in parallel
var WorkerN = 1

func drawFace(src image.Image, rect image.Rectangle) image.Image {
	canvas := imaging.Crop(src, rect)
	canvas = imaging.Thumbnail(canvas, 100, 100, imaging.Lanczos)
	return canvas
}

func getFaces(filename string, cutFace bool) ([]image.Image, []image.Rectangle, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer reader.Close()
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, nil, err
	}
	var rects []image.Rectangle
	var faces []image.Image

	// Resize for faster face search
	img = imaging.Resize(img, img.Bounds().Dx()/2, img.Bounds().Dy()/2, imaging.Box)
	cvImg := opencv.FromImage(img)
	cvRects := cascade.DetectObjects(cvImg)
	for _, face := range cvRects {
		r := image.Rectangle{
			image.Point{face.X(), face.Y()},
			image.Point{face.X() + face.Width(),
				face.Y() + face.Height()}}
		if cutFace {
			faces = append(faces, drawFace(img, r))
		}
		rects = append(rects, r)
	}
	return faces, rects, nil
}

func worker(jobIn chan []interface{}, jobOut chan []interface{}) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	for job := range jobIn {
		filename := job[0].(string)
		cutFace := job[1].(bool)
		faces, rects, err := getFaces(filename, cutFace)
		if err != nil {
			jobOut <- nil
			continue
		}
		if len(rects) == 0 {
			jobOut <- nil
			continue
		}
		jobOut <- []interface{}{filename, faces, rects}
	}
}

func ScanDir(dir string, cutFaces bool, fn func(string, []image.Image, []image.Rectangle)) {
	var (
		jobIn  = make(chan []interface{}, WorkerN)
		jobOut = make(chan []interface{}, 4)
		wg     = &sync.WaitGroup{}
	)

	for i := 0; i < cap(jobIn); i++ {
		go worker(jobIn, jobOut)
	}

	go func() {
		for r := range jobOut {
			if r == nil {
				wg.Done()
				continue
			}
			fn(r[0].(string), r[1].([]image.Image), r[2].([]image.Rectangle))
			wg.Done()
		}
	}()

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		wg.Add(1)
		jobIn <- []interface{}{path, cutFaces}
		return nil
	})

	wg.Wait()
	close(jobIn)
	close(jobOut)
}
