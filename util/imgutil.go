package util

import (
	"image"
	"image/color"
	"math/rand"
)

//计算相同尺寸图像的相似度
func Similarity(imgA, imgB image.Image) uint64 {
	simi := uint64(0)
	//fmt.Println("w:", imgA.Bounds().Dx(), "h:", imgA.Bounds().Dy())
	for x := 0; x < imgA.Bounds().Dx(); x++ {
		for y := 0; y < imgA.Bounds().Dy(); y++ {
			//s1 := imgA.At(x, y).(color.RGBA).R + imgA.At(x, y).(color.RGBA).G + imgA.At(x, y).(color.RGBA).B + imgA.At(x, y).(color.RGBA).A
			//s2 := imgB.At(x, y).(color.RGBA).R + imgB.At(x, y).(color.RGBA).G + imgB.At(x, y).(color.RGBA).B + imgB.At(x, y).(color.RGBA).A
			//sampleColor := color.RGBAModel.Convert(imgB.At(x, y)).(color.RGBA)
			//s2 := sampleColor.(color.RGBA).R + sampleColor.(color.RGBA).G + sampleColor.(color.RGBA).B + sampleColor.(color.RGBA).A
			//c1 := imgA.At(x, y).(color.RGBA).R - imgA.At(x, y).(color.RGBA).R
			//c2 := imgA.At(x, y).(color.RGBA).G - imgA.At(x, y).(color.RGBA).G
			//c1 := imgA.At(x, y).(color.RGBA).B - imgA.At(x, y).(color.RGBA).B
			//c2 := imgA.At(x, y).(color.RGBA).A - imgA.At(x, y).(color.RGBA).A
			//simi += uint64(math.Abs(float64(s1) - float64(s2)))
			sampleColor := color.RGBAModel.Convert(imgB.At(x, y)).(color.RGBA)

			//R
			c1 := imgA.At(x, y).(color.RGBA).R
			c2 := sampleColor.R
			if c1 > c2 {
				simi += uint64(c1 - c2)
			} else {
				simi += uint64(c2 - c1)
			}
			//G
			c1 = imgA.At(x, y).(color.RGBA).G
			c2 = sampleColor.G
			if c1 > c2 {
				simi += uint64(c1 - c2)
			} else {
				simi += uint64(c2 - c1)
			}

			//B
			c1 = imgA.At(x, y).(color.RGBA).B
			c2 = sampleColor.B
			if c1 > c2 {
				simi += uint64(c1 - c2)
			} else {
				simi += uint64(c2 - c1)
			}

			//A
			c1 = imgA.At(x, y).(color.RGBA).A
			c2 = sampleColor.A
			if c1 > c2 {
				simi += uint64(c1 - c2)
			} else {
				simi += uint64(c2 - c1)
			}

		}
	}
	return simi
}

func GetRand255() int {
	return rand.Intn(256)
}
