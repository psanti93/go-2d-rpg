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
	player      *Player
	enemies     []*Enemy
	potions     []*Potion
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
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

	opts := ebiten.DrawImageOptions{}

	// drawing the map layers
	for _, layer := range g.tilemapJSON.Layers {
		// index is the position in the array and id is the value
		for index, id := range layer.Data {
			// get the position of the tile
			x := index % layer.Width
			y := index / layer.Width

			// give the pixle position of the tile
			x *= 16
			y *= 16

			// grabbing the image from the tileset_floor.png
			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)
			opts.GeoM.Reset()
		}
	}
	// draw our player
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
	ebiten.SetWindowSize(500, 500)
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

	// loading the json file for the tile map
	tileMapJson, err := NewTilemapJSON("assets/maps/spawn.json")

	if err != nil {
		log.Fatal(err)
	}

	//next we want to load the tile set image
	tileSetImg, _, err := ebitenutil.NewImageFromFile("assets/images/tileset_floor.png")
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
		tilemapJSON: tileMapJson,
		tilemapImg:  tileSetImg,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
