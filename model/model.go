package model

import (
	"github.com/fogleman/gg"
	"image"
	"math/rand"
	"studyImage/util"
)

type Point struct {
	X, Y int
}
type Triangle struct {
	Points     [3]Point
	R, G, B, A int //三角形颜色
}

func (this *Triangle) RandInit(size Size) {
	this.R = rand.Intn(256)
	this.G = rand.Intn(256)
	this.B = rand.Intn(256)
	this.A = 127
	for i, _ := range this.Points {
		this.Points[i].X = rand.Intn(size.W)
		this.Points[i].Y = rand.Intn(size.H)
	}
}

type Size struct {
	W, H int
}

type Life struct {
	Triangles []Triangle
}

//为每个点赋随机值，size图片尺寸
func (this *Life) InitRand(numTriangles int, size Size) {
	this.Triangles = make([]Triangle, numTriangles)
	for i, _ := range this.Triangles {
		this.Triangles[i].RandInit(size)
	}
}

type Group struct {
	Lifes        []Life
	Num          int         //种群大小
	NumTriangles int         //每个life拥有三角形的数量
	MostSimilar  Life        //最相似life
	MSCount      uint64      //最相似的life的相似度 , 值越低越相似
	LeastSimilar Life        //最不相似life
	LSCount      uint64      //不相似life的相似度
	MutationRate float32     //0-1突变率
	CrossRate    float32     //交叉率
	Size         Size        //图片尺寸
	Sample       image.Image //样本图片,进化目标
	generation   int
}

//size图片尺寸,num种群大小
func (this *Group) Init(num int, numTriangles int, size Size, sample image.Image) {
	this.MSCount--
	this.Size = size
	this.Sample = sample
	this.Lifes = make([]Life, num)
	this.Num = num
	this.NumTriangles = numTriangles
	this.MutationRate = 0.1 //突变率，一个三角形发生变化的概率(一个life只有可能其中一个三角形发生突变)
	this.CrossRate = 0.86   //交叉率，交叉(繁衍后代)是进行的概率
	for i, _ := range this.Lifes {
		this.Lifes[i].InitRand(numTriangles, size)
	}
}

func (this *Group) SetMutationRate(rate float32) {
	this.MutationRate = rate
}
func (this *Group) SetCrossRate(rate float32) {
	this.CrossRate = rate
}

//两个交叉
func (this *Group) Cross(a, b Life) (Life, Life) {
	if rand.Float32() > this.CrossRate {
		//不交叉,返回原life
		return a, b
	}
	slicingLocation := rand.Intn(this.NumTriangles) //切割位置不大于三角形数量
	aTriangles, bTriangles := make([]Triangle, this.NumTriangles), make([]Triangle, this.NumTriangles)
	//交叉生成aLife
	copy(aTriangles[:slicingLocation], b.Triangles[:slicingLocation])
	copy(aTriangles[slicingLocation:], a.Triangles[slicingLocation:])
	//交叉生成bLife
	copy(bTriangles[:slicingLocation], a.Triangles[:slicingLocation])
	copy(bTriangles[slicingLocation:], b.Triangles[slicingLocation:])
	//(每个个体低概率)变异
	if rand.Float32() < this.MutationRate {
		//随机挑一个三角形变异
		aTriangles[rand.Intn(this.NumTriangles)].RandInit(this.Size)
	}
	if rand.Float32() < this.MutationRate {
		bTriangles[rand.Intn(this.NumTriangles)].RandInit(this.Size)
	}
	return Life{aTriangles}, Life{bTriangles}
}

//淘汰落后者
func (this *Group) Kill() {
	k := 0 //将要淘汰的life的位置
	preSimilarity := uint64(0)
	//fmt.Println("begin kill")
	for i, _ := range this.Lifes {
		img := this.GenerateImg(this.Lifes[i].Triangles)
		//fmt.Println("kill -- in lifes")
		similarity := util.Similarity(img, this.Sample)
		if similarity > preSimilarity {
			preSimilarity = similarity
			k = i
		}
		if similarity < this.MSCount {
			//最相似
			this.MostSimilar = this.Lifes[i]
			this.MSCount = similarity
		}
		if similarity > this.LSCount {
			//最不相似
			this.LeastSimilar = this.Lifes[i]
			this.LSCount = similarity
		}
		this.Fill(k) //相当于淘汰，新替旧
	}

}

//繁衍第二代
func (this *Group) Regenerate() {
	lifes := make([]Life, this.Num)
	this.Lifes = DisorganizeSlice(this.Lifes)
	//fmt.Println("have Disorganized")
	for i := 0; i < (len(lifes) - 1); i++ {
		lifes[i], lifes[i+1] = this.Cross(this.Lifes[i], this.Lifes[i+1])
	}
	this.Lifes = lifes
	this.generation++ //第n代，代数增加
}
func (this *Group) GetGeneration() int {
	return this.generation
}

//填补空缺，淘汰的方式
func (this *Group) Fill(location int) {
	for r := rand.Intn(this.Num); ; r = rand.Intn(this.Num) {
		if r != location {
			a, _ := this.Cross(this.MostSimilar, this.Lifes[r])
			this.Lifes[location] = a
			break
		}
	}

}
func DrawTriangle(context *gg.Context, tri Triangle) {
	context.SetRGBA255(tri.R, tri.G, tri.B, tri.A)
	context.MoveTo(float64(tri.Points[0].X), float64(tri.Points[0].Y))
	context.LineTo(float64(tri.Points[1].X), float64(tri.Points[1].Y))
	context.LineTo(float64(tri.Points[2].X), float64(tri.Points[2].Y))
	context.Fill()

}
func (this *Group) GenerateImg(triangles []Triangle) image.Image {
	//fmt.Print("in function GenerateImg")
	//start := time.Now()
	context := gg.NewContext(this.Size.W, this.Size.H)
	context.SetRGBA255(255, 255, 255, 255)
	//fmt.Println(time.Now().Second())
	context.Clear()
	for i, _ := range triangles {
		DrawTriangle(context, triangles[i])
	}
	//fmt.Println("... end", time.Since(start).Milliseconds())
	return context.Image()

}

//打乱切片
func DisorganizeSlice(a []Life) []Life {
	for i, c := range a {
		for r := rand.Intn(len(a)); ; r = rand.Intn(len(a)) {
			if r != i {
				a[i] = a[r]
				a[r] = c
				break
			}
		}
	}
	return a

}
