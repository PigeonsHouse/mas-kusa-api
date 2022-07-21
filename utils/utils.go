package utils

import (
	"image"
	"image/color"
)

func GenKusa(tootCount [][7]int) *image.RGBA {
	w, h := 500, 700
	base := image.NewRGBA(image.Rect(0, 0, w, h))
	bg := color.RGBA{255, 255, 255, 255}
	lc := []color.RGBA{
		{235, 237, 240, 255},
		{171, 255, 185, 255},
		{155, 233, 168, 255},
		{64, 196, 99, 255},
		{48, 161, 78, 255},
		{33, 110, 57, 255},
		{23, 74, 18, 255},
	}

	border := []int{5, 15, 40, 70, 100}

	for ox := 0; ox < 5; ox++ {
		for oy := 0; oy < 7; oy++ {
			for x := 0; x < w; x++ {
				for y := 0; y < h; y++ {
					if (x > 10 && y > 10) && (x < 90 && y < 90) {
						if tootCount[ox][oy] == 0 {
							base.Set(x+ox*100, y+oy*100, lc[0])
						} else if tootCount[ox][oy] < border[0] {
							base.Set(x+ox*100, y+oy*100, lc[1])
						} else if tootCount[ox][oy] < border[1] {
							base.Set(x+ox*100, y+oy*100, lc[2])
						} else if tootCount[ox][oy] < border[2] {
							base.Set(x+ox*100, y+oy*100, lc[3])
						} else if tootCount[ox][oy] < border[3] {
							base.Set(x+ox*100, y+oy*100, lc[4])
						} else if tootCount[ox][oy] < border[4] {
							base.Set(x+ox*100, y+oy*100, lc[5])
						} else {
							base.Set(x+ox*100, y+oy*100, lc[6])
						}
					} else {
						base.Set(x+ox*100, y+oy*100, bg)
					}
				}
			}
		}
	}
	return base
}
