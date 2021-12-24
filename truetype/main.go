package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"time"

	//"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "/usr/share/fonts/truetype/liberation/LiberationSerif-Regular.ttf", "filename of the ttf font")
	hinting  = flag.Int("hinting", 0, "font.HintingNone | font.HintingFull")
	size     = flag.Float64("size", 12, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
	text     = string("JOJO")
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
	const (
		width        = 72
		height       = 36
		startingDotX = 6
		startingDotY = 28
	)

	flag.Parse()
	fmt.Printf("Loading fontfile %q\n", *fontfile)
	b, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(b)
	if err != nil {
		log.Println(err)
		return
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size:              *size,
		DPI:               *dpi,
		Hinting:           font.Hinting(*hinting),
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	})

	defer timeTrack(time.Now(), "Timer")
	dst := image.NewGray(image.Rect(0, 0, width, height))
	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	fmt.Printf("The dot is at %v\n", d.Dot)
	d.DrawString("jel")
	fmt.Printf("The dot is at %v\n", d.Dot)
	d.Src = image.NewUniform(color.Gray{0x7F})
	d.DrawString("ly")
	fmt.Printf("The dot is at %v\n", d.Dot)

	const asciiArt = " +*#"
	buf := make([]byte, 0, height*(width+1))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := asciiArt[dst.GrayAt(x, y).Y>>6]
			if c != '.' {
				// No-op.
			} else if x == startingDotX-1 {
				c = ']'
			} else if y == startingDotY-1 {
				c = '_'
			}
			buf = append(buf, c)
		}
		buf = append(buf, '\n')
	}
	os.Stdout.Write(buf)
}
