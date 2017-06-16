package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
    "math/rand"
	"image/color"
    "time"
    "fmt"
)

type myScene struct {}

const scaling = (400/15)

var head Square = Square{BasicEntity: ecs.NewBasic(), xv: 1}
var food Square = Square{
    BasicEntity: ecs.NewBasic(),
}
var tail = make([]Square, 0, 224)
var gameOver bool = false

type MouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type Square struct {
    ecs.BasicEntity
    common.RenderComponent
    common.SpaceComponent
    x int
    y int
    xv int
    yv int
}

type SnakeSystem struct {
    world *ecs.World
}


func (cb *SnakeSystem) New(w *ecs.World) {
    cb.world = w
}

func (*SnakeSystem) Remove(ecs.BasicEntity) {}

func checkFood() {
    for (head.x == food.x && head.y == food.y) {
        s := rand.NewSource(time.Now().UnixNano())
        r := rand.New(s)
        food.x = r.Intn(15)
        s = rand.NewSource(time.Now().UnixNano())
        r = rand.New(s)
        food.y = r.Intn(15)
        tail = append(tail, Square{BasicEntity: ecs.NewBasic(),
        x: head.x,
        y: head.y,
    })
}
}

func checkCollision() {
    for i := 0; i < len(tail); i++ {
        if head.x == tail[i].x && head.y == tail[i].y {
            gameOver = true
            fmt.Println("Game Over!")
            fmt.Println("Press space to start again")
        }
    }
}

func (cb *SnakeSystem) Update(dt float32) {
    if !gameOver {
        if engo.Input.Button("Up").JustPressed() {
            head.yv = -1
            head.xv = 0
        } else if engo.Input.Button("Down").JustPressed() {
            head.yv = 1
            head.xv = 0
        } else if engo.Input.Button("Right").JustPressed() {
            head.xv = 1
            head.yv = 0
        } else if engo.Input.Button("Left").JustPressed() {
            head.xv = -1
            head.yv = 0
        }

        pointx := (head.x + head.xv) %15
        pointy := (head.y + head.yv) %15
        if pointx < 0 {
            pointx = 14
        }
        if pointy < 0 {
            pointy = 14
        }


        if len(tail) > 0 {
            for i := len(tail)-1; i>0; i-- {
                tail[i].x = tail[i-1].x
                tail[i].y = tail[i-1].y
            }
            tail[0].x = head.x
            tail[0].y = head.y
        }

        head.x = pointx
        head.y = pointy


        checkCollision()
        checkFood()

        textureBlue, _ := common.LoadedSprite("blue-square.png")
        textureVeryBlue, _ := common.LoadedSprite("veryblue-square.png")
        textureGreen, _ := common.LoadedSprite("green-square.png")


        head.SpaceComponent = common.SpaceComponent {
            Position: engo.Point{float32(head.x*scaling), float32(head.y*scaling)},
            Width:  1,
            Height: 1,
        }
        head.RenderComponent = common.RenderComponent{
            Drawable: textureBlue,
            Scale:    engo.Point{0.08, 0.08},
        }


        food.SpaceComponent = common.SpaceComponent {
            Position: engo.Point{float32(food.x*scaling), float32(food.y*scaling)},
            Width:  1,
            Height: 1,
        }

        food.RenderComponent = common.RenderComponent{
            Drawable: textureGreen,
            Scale:    engo.Point{0.08, 0.08},
        }

        for i := 0; i<len(tail); i++ {
            tail[i].SpaceComponent = common.SpaceComponent {
                Position: engo.Point{float32(tail[i].x*scaling), float32(tail[i].y*scaling)},
                Width:  1,
                Height: 1,
            }
            tail[i].RenderComponent = common.RenderComponent{
                Drawable: textureVeryBlue,
                Scale:    engo.Point{0.08, 0.08},
            }
        }

        for _, system := range cb.world.Systems() {
            switch sys := system.(type) {
            case *common.RenderSystem:
                sys.Add(&head.BasicEntity, &head.RenderComponent, &head.SpaceComponent)
                sys.Add(&food.BasicEntity, &food.RenderComponent, &food.SpaceComponent)
                for i := 0; i<len(tail); i++ {
                    sys.Add(&tail[i].BasicEntity, &tail[i].RenderComponent, &tail[i].SpaceComponent)
                }
            }
        }
    } else {
        if engo.Input.Button("Space").JustPressed() {
            for i := 0; i<len(tail); i++ {
                tail[i].x = 500
                tail[i].y = 500
            }
            for i := 0; i<len(tail); i++ {
                tail[i].SpaceComponent = common.SpaceComponent {
                    Position: engo.Point{float32(tail[i].x*scaling), float32(tail[i].y*scaling)},
                    Width:  1,
                    Height: 1,
                }
            }

            tail = tail[:0]
            head.x = 0
            head.y = 0
            head.yv = 0
            head.xv = 1
            s := rand.NewSource(time.Now().UnixNano())
            r := rand.New(s)
            food.x = r.Intn(15)
            s = rand.NewSource(time.Now().UnixNano())
            r = rand.New(s)
            food.y = r.Intn(15)
            common.SetBackground(color.White)
            gameOver = false
        }
    }
}

type System interface {
	Update(dt float32)
	Remove(ecs.BasicEntity)
}

func (*myScene) Type() string { return "myGame" }

func (*myScene) Preload() {
    engo.Files.Load("blue-square.png")
    engo.Files.Load("veryblue-square.png")
    engo.Files.Load("green-square.png")
}

func (*myScene) Setup(world *ecs.World) {
    engo.Input.RegisterButton("Up", engo.ArrowUp)
    engo.Input.RegisterButton("Down", engo.ArrowDown)
    engo.Input.RegisterButton("Left", engo.ArrowLeft)
    engo.Input.RegisterButton("Right", engo.ArrowRight)
    engo.Input.RegisterButton("Space", engo.Space)
    common.SetBackground(color.White)
    world.AddSystem(&common.RenderSystem{})
    world.AddSystem(&SnakeSystem{})
    s := rand.NewSource(time.Now().UnixNano())
    r := rand.New(s)
    food.x = r.Intn(15)
    s = rand.NewSource(time.Now().UnixNano())
    r = rand.New(s)
    food.y = r.Intn(15)
}

func main() {
	opts := engo.RunOptions{
		Title: "Snake",
		Width:  390,
		Height: 390,
        Fullscreen : false,
        FPSLimit : 10,
	}
	engo.Run(opts, &myScene{})
}
