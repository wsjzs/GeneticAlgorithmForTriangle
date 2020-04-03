package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"studyImage/model"
	"time"
)

func main() {

	test1()

}

func Create() {

}

func Test() {

	size := model.Size{256, 256}
	context := gg.NewContext(size.W, size.H)
	context.SetRGBA255(255, 255, 255, 255)
	context.Clear()
	//随机3个三角形
	rand.Seed(time.Now().UnixNano())

	group := model.Group{}
	sample, _ := gg.LoadJPG("test.jpg")
	group.Init(16, 100, size, sample)

	//画第一个life
	triangles := group.Lifes[0].Triangles
	for i, _ := range triangles {
		model.DrawTriangle(context, triangles[i])
	}
	context.SavePNG("out.png")
	fmt.Println(triangles)

}

func test2() {
	img, _ := gg.LoadJPG("firefox.jpg")

	//img:=model.Point{}
	c := img.At(0, 0)
	d := color.RGBAModel.Convert(c)
	fmt.Println(d)

}

func test1() {
	sample, err := gg.LoadJPG("./qqsmall.jpg")
	if err != nil {
		log.Fatal(err)
	}

	size := model.Size{W: sample.Bounds().Max.X, H: sample.Bounds().Max.Y}
	group := model.Group{}
	group.Init(20, 50, size, sample)
	group.SetMutationRate(0.2)
	fmt.Println(sample)
	fmt.Println("start")
	e := 0
	for i := 1; ; i++ {
		start := time.Now()
		rand.Seed(time.Now().UnixNano())
		group.Kill()
		//fmt.Println("kill end")
		group.Regenerate()
		//context := gg.NewContextForImage(group.GenerateImg(group.MostSimilar.Triangles))
		//context.SavePNG(fmt.Sprintf("\\out\\%d.png", i))
		e++
		if e >= 10 {
			img := group.GenerateImg(group.MostSimilar.Triangles)
			file, err := os.Create(fmt.Sprintf("./out2/%d.png", i))
			if err != nil {
				log.Fatal(err)
			}
			err = png.Encode(file, img)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()
			e = 0
		}
		fmt.Print("第", group.GetGeneration(), "代 ")
		fmt.Print("相似度 ", group.MSCount)
		fmt.Println(" 耗时", time.Since(start).Milliseconds())
	}
}
