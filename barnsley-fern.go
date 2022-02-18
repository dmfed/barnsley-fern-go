package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func createImage(h, v int) *image.RGBA {
	rect := image.Rectangle{
		image.Point{0, 0},
		image.Point{h, v},
	}
	return image.NewRGBA(rect)
}

func fillBackground(img *image.RGBA, c color.Color) {
	rect := img.Bounds()
	for x := 0; x < rect.Max.X; x++ {
		for y := 0; y < rect.Max.Y; y++ {
			img.Set(x, y, c)
		}
	}
}

func min(a, b float64) float64 {
	if a > b {
		a = b
	}
	return a
}

func drawBarnsleyFern(img *image.RGBA, c color.Color, dots int) {
	var (
		x, y, tmpx, tmpy, r, maxy, maxx, scale, yoffset, xoffset float64
	)
	maxy = float64(img.Bounds().Max.Y) // размер картинки по вертикали
	maxx = float64(img.Bounds().Max.X) // размер картиники по горизонтали
	scale = min(maxx, maxy)
	yoffset = scale / 10         // отступы 10% сверху и снизу
	scale = scale - yoffset*2    // масштаб самого папоротника
	xoffset = (maxx - scale) / 2 // равные отступы с двух сторон

	x, y = 0.5, 0.0
	for dots > 0 {
		dots--
		r = rand.Float64()
		if r <= 0.01 {
			// стебель
			tmpx = 0.5
			tmpy = 0.16 * y
		} else if r <= 0.08 {
			// самый большой правый листок
			tmpx = 0.2*x - 0.26*y + 0.400
			tmpy = 0.23*x + 0.22*y - 0.045
		} else if r <= 0.15 {
			// самый большой левый листок
			tmpx = -0.15*x + 0.28*y + 0.575
			tmpy = 0.26*x + 0.24*y - 0.086
		} else {
			// последующие листочки
			tmpx = 0.85*x + 0.04*y + 0.075
			tmpy = -0.04*x + 0.850*y + 0.180
		}
		// рисуем точку
		x, y = tmpx, tmpy
		img.Set(int(xoffset+x*scale), int(maxy-yoffset-y*scale), c)
	}
}

func main() {
	var (
		filename string = "barnsley_fern.png"
		h               = flag.Int("h", 1920, "размер картинки по горизонтали")
		v               = flag.Int("v", 1080, "размер картинки по вертикали")
		dots            = flag.Int("d", 100000, "сколько точек рисовать")
		img      *image.RGBA
		f        *os.File
		err      error
	)
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	img = createImage(*h, *v)

	fillBackground(img, color.White)
	drawBarnsleyFern(img, color.RGBA{0, 153, 0, 255}, *dots)

	if f, err = os.Create(filename); err != nil {
		fmt.Printf("не удалось создать файл %s: %v", filename, err)
	} else if err = png.Encode(f, img); err != nil {
		fmt.Printf("не удалось сохранить изображение: %v", err)
	}
	fmt.Println("done...")
}
