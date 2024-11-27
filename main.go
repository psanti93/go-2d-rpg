package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64 // position of the sprite
}

type Player struct {
	*Sprite
	Health uint // the health can never be negative
}

type Enemy struct {
	*Sprite       // struct embedding
	FollowsPlayer bool
}

type Potion struct {
	*Sprite
	AmountHeal uint
}

type Game struct {
	player  *Player
	enemies []*Enemy
	potions []*Potion
}

// note for the key cortesen the y axis is inverted
func (g *Game) Update() error {
	// react to key presses to move the player
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}

	// algorithm for sprites to follow player
	for _, enemy := range g.enemies {

		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.X += 1
			} else if enemy.X > g.player.X {
				enemy.X -= 1
			}

			if enemy.Y < g.player.Y {
				enemy.Y += 1
			} else if enemy.Y > g.player.Y {
				enemy.Y -= 1
			}
		}
	}

	for _, potion := range g.potions {
		if g.player.X > potion.X {
			g.player.Health += potion.AmountHeal
			fmt.Printf("Player picked up Potion! Health is now %d\n", g.player.Health)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	// draw our player
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.player.X, g.player.Y)
	// the coordinates of the image of the player that you want to draw
	screen.DrawImage(g.player.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
	opts.GeoM.Reset()

	// drawing all the enemies
	for _, enemy := range g.enemies {
		opts.GeoM.Translate(enemy.X, enemy.Y)
		screen.DrawImage(enemy.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}

	// draw the potions
	for _, potion := range g.potions {
		opts.GeoM.Translate(potion.X, potion.Y)
		screen.DrawImage(potion.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), &opts)
		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	// allowing the window to resize on the screen
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// loading up the player images
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/lionboy.png")
	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	potionsImg, _, err := ebitenutil.NewImageFromFile("assets/images/potion.png")
	if err != nil {
		log.Fatal(err)
	}
	game := &Game{
		player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				X:   100.0,
				Y:   100.0,
			},
			Health: 10,
		},
		enemies: []*Enemy{
			{
				&Sprite{
					Img: skeletonImg,
					X:   100.0,
					Y:   100.0,
				},
				true,
			},
			{
				&Sprite{
					Img: skeletonImg,
					X:   150.0,
					Y:   50.0,
				},
				false,
			},
		},
		potions: []*Potion{
			{
				&Sprite{
					Img: potionsImg,
					X:   210.0,
					Y:   100,
				},
				1.0,
			},
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
