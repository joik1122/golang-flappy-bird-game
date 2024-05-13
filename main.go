package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	birdSize     = 20
	pipeWidth    = 50
	pipeGap      = 120
	gravity      = 0.25
	jumpStrength = -5
)

type Bird struct {
	x, y float64
	vy   float64
}

type Pipe struct {
	x      float64
	height float64
}

type Game struct {
	bird           Bird
	pipes          []Pipe
	score          int
	highScore      int
	pipeSpawnTimer int
	gameOver       bool
}

func NewGame() *Game {
	return &Game{
		bird: Bird{x: screenWidth / 4, y: screenHeight / 2},
	}
}

func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Reset()
		}
		return nil
	}

	g.bird.vy += gravity
	g.bird.y += g.bird.vy

	if g.bird.y > screenHeight-birdSize || g.bird.y < 0 {
		g.gameOver = true
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.bird.vy = jumpStrength
	}

	g.pipeSpawnTimer++
	if g.pipeSpawnTimer > 90 {
		g.pipeSpawnTimer = 0
		g.AddPipe()
	}

	for i := range g.pipes {
		g.pipes[i].x -= 3
		if g.pipes[i].x+pipeWidth < 0 {
			g.pipes = append(g.pipes[:i], g.pipes[i+1:]...)
			g.score++
			if g.score > g.highScore {
				g.highScore = g.score
			}
		}
	}

	for _, p := range g.pipes {
		if g.bird.x+birdSize > p.x && g.bird.x < p.x+pipeWidth && (g.bird.y < p.height || g.bird.y+birdSize > p.height+pipeGap) {
			g.gameOver = true
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{135, 206, 235, 255})

	ebitenutil.DrawRect(screen, g.bird.x, g.bird.y, birdSize, birdSize, color.RGBA{255, 255, 0, 255})

	for _, p := range g.pipes {
		ebitenutil.DrawRect(screen, p.x, 0, pipeWidth, p.height, color.RGBA{34, 139, 34, 255})
		ebitenutil.DrawRect(screen, p.x, p.height+pipeGap, pipeWidth, screenHeight, color.RGBA{34, 139, 34, 255})
	}

	msg := "Score: " + fmt.Sprint(g.score)
	if g.gameOver {
		msg += " | Game Over! Press Space to Restart"
	}
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) AddPipe() {
	height := float64(rand.Intn(screenHeight/2) + screenHeight/8)
	g.pipes = append(g.pipes, Pipe{x: screenWidth, height: height})
}

func (g *Game) Reset() {
	g.bird = Bird{x: screenWidth / 4, y: screenHeight / 2}
	g.pipes = nil
	g.score = 0
	g.pipeSpawnTimer = 0
	g.gameOver = false
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Flappy Bird")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
