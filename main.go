package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"studyImage/model"
	"time"
)

func main() {

	begin()

}
func begin() {
	sample, err := gg.LoadJPG("./firefox.jpg")
	if err != nil {
		log.Fatal(err)
	}

	size := model.Size{W: sample.Bounds().Max.X, H: sample.Bounds().Max.Y}
	group := model.Group{}
	group.Init(20, 50, size, sample)
	group.SetMutationRate(0.2)
	fmt.Println("start")
	e := 0
	for i := 1; ; i++ {
		start := time.Now()
		rand.Seed(time.Now().UnixNano())
		group.Kill()
		group.Regenerate()
		e++
		if e >= 10 {
			img := group.GenerateImg(group.MostSimilar.Triangles)
			file, err := os.Create(fmt.Sprintf("./out/%d.png", i))
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
		fmt.Print(group.GetGeneration(),"generations")
		fmt.Print("similarity", group.MSCount)
		fmt.Println(" time", time.Since(start).Milliseconds())
	}
}
