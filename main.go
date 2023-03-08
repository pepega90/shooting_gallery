package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	WIDTH  = 840
	HEIGHT = 640
)

type Target struct {
	x, y           float64
	img            *ebiten.Image
	speedX, speedY float64
}

func (t *Target) Draw(screen *ebiten.Image) {
	tp := &ebiten.DrawImageOptions{}
	tp.GeoM.Translate(t.x, t.y)
	screen.DrawImage(t.img, tp)
}

type Game struct {
	bg            *ebiten.Image
	crossHairImg  *ebiten.Image
	targetImg     *ebiten.Image
	shellImg      *ebiten.Image
	mousePosition struct {
		x, y int
	}
	listTarget []Target
}

func CheckPointCollision(mx int, my int, rect Target) bool {
	width, height := rect.img.Size()
	if mx >= int(rect.x) && // right of the left edge AND
		mx <= int(rect.x+float64(width)) && // left of the right edge AND
		my >= int(rect.y) && // below the top AND
		my <= int(rect.y+float64(height)) { // above the bottom
		return true
	}
	return false
}

// * Global variabel
var (
	lastUpdate   time.Time
	countdown    = 25
	f            font.Face
	titleF       font.Face
	subTitleF    font.Face
	creditF      font.Face
	gameOverF    font.Face
	ammo         int
	lastTime     time.Time
	currentScene = MENU
	score        = 0
	playCount    = 6
	lastPlay     time.Time
	play         bool
	trigger      = true
)

const (
	MENU = iota
	PLAY
	GAMEOVER
)

var (
	ORANGE = color.RGBA{250, 246, 2, 255}
)

func (g *Game) Update() error {
	g.mousePosition.x, g.mousePosition.y = ebiten.CursorPosition()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	switch currentScene {
	case MENU:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			play = true
		}

		if time.Since(lastPlay) > time.Second && playCount > 0 && play {
			playCount -= 1
			lastPlay = time.Now()
		}

		if playCount == 0 {

			currentScene = PLAY
		}
	case PLAY:

		// coutndown
		if time.Since(lastUpdate) > time.Second && countdown > 0 {
			countdown -= 1
			lastUpdate = time.Now()
		}

		if trigger {
			g.listTarget = append(g.listTarget[:1], g.listTarget[1+1:]...)
		}

		if countdown == 23 {
			trigger = false
		}

		// if countdown < 25 {
		if trigger {
			for len(g.listTarget) < 10 {
				rand.Seed(time.Now().UnixNano())
				g.listTarget = append(g.listTarget, Target{
					img:    g.targetImg,
					speedX: float64(rand.Intn(6-3) + 3),
					speedY: float64(rand.Intn(6-3) + 3),
					x:      float64(rand.Intn((WIDTH-33)-33) + 33),
					y:      float64(rand.Intn((HEIGHT-33)-0) + 0),
				})
			}
		}
		// }

		// update target
		for i := 0; i < len(g.listTarget); i++ {
			g.listTarget[i].x += g.listTarget[i].speedX
			g.listTarget[i].y += g.listTarget[i].speedY
		}

		// bikin target stay in window
		for i := 0; i < len(g.listTarget); i++ {
			targetWidth, targetHeight := g.listTarget[i].img.Size()
			if int(g.listTarget[i].x) > WIDTH-targetWidth || g.listTarget[i].x < 0 {
				g.listTarget[i].speedX *= -1
			}

			if int(g.listTarget[i].y) > HEIGHT-targetHeight || g.listTarget[i].y < 0 {
				g.listTarget[i].speedY *= -1
			}
		}

		// jika player klik kiri maka ammonya berkurang
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ammo > 0 {
			if time.Since(lastTime) > time.Millisecond*300 {
				ammo -= 1
				lastTime = time.Now()
			}
		}

		// check collision antara mouse position dengan target
		for i := 0; i < len(g.listTarget); i++ {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && CheckPointCollision(g.mousePosition.x, g.mousePosition.y, g.listTarget[i]) {
				g.listTarget = append(g.listTarget[:i], g.listTarget[i+1:]...)
				score += 1
			}
		}

		if ammo == 0 || countdown == 0 || len(g.listTarget) < 1 {
			currentScene = GAMEOVER
		}
	case GAMEOVER:
		// reset game
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			// g.listTarget = []Target{}
			for i := 0; i < 5; i++ {
				rand.Seed(time.Now().UnixNano())
				g.listTarget = append(g.listTarget, Target{
					img:    g.targetImg,
					speedX: float64(rand.Intn(6-3) + 3),
					speedY: float64(rand.Intn(6-3) + 3),
					x:      float64(rand.Intn((WIDTH-33)-33) + 33),
					y:      float64(rand.Intn((HEIGHT-33)-0) + 0),
				})
			}
			trigger = true
			ammo = 15
			countdown = 25
			score = 0
			currentScene = PLAY
		}

	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background
	for y := 0; y < 3; y++ {
		for x := 0; x < 4; x++ {
			bp := &ebiten.DrawImageOptions{}
			bp.GeoM.Translate(float64(x*250), float64(y*200))
			screen.DrawImage(g.bg, bp)
		}
	}

	switch currentScene {
	case MENU:
		if play {
			text.Draw(screen, fmt.Sprintf("Ready?\n\n  %v", playCount), titleF, WIDTH/2-60, HEIGHT/2, ORANGE)
		} else {
			text.Draw(screen, "KKona Shooting Gallery", titleF, 100, HEIGHT/4, ORANGE)
			text.Draw(screen, "Tekan \"Space\" untuk play!", subTitleF, 180, HEIGHT/2, color.White)
			text.Draw(screen, "created by aji mustofa @pepega90", creditF, 187, 617, color.RGBA{0, 0, 0, 255})
		}
	case PLAY:
		// draw target
		for i := 0; i < len(g.listTarget); i++ {
			g.listTarget[i].Draw(screen)
		}

		// draw crosshair
		cp := &ebiten.DrawImageOptions{}
		cWidth, cHeight := g.crossHairImg.Size()
		cp.GeoM.Translate(float64(g.mousePosition.x)-float64(cWidth)/2, float64(g.mousePosition.y)-float64(cHeight)/2)
		screen.DrawImage(g.crossHairImg, cp)

		// draw indicator
		ip := &ebiten.DrawImageOptions{}
		ip.GeoM.Scale(1.1, 1.1)
		ip.GeoM.Translate(738, 25)
		screen.DrawImage(g.targetImg, ip)
		text.Draw(screen, fmt.Sprintf("%v", len(g.listTarget)), f, 785, 55, ORANGE)

		// draw shotgun shell
		for i := 0; i < ammo; i++ {
			sp := &ebiten.DrawImageOptions{}
			sp.GeoM.Translate(float64(20+i*21), 20)
			screen.DrawImage(g.shellImg, sp)
		}

		// draw countdown
		text.Draw(screen, fmt.Sprintf("%v", countdown), f, WIDTH/2, 40, color.White)
	case GAMEOVER:
		text.Draw(screen, "GAME OVER", f, WIDTH/2-100, HEIGHT/4, color.RGBA{255, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("Score Kamu: %v", score), f, WIDTH/2-120, HEIGHT/2, ORANGE)
		text.Draw(screen, "Tekan \"R\" untuk restart!", gameOverF, WIDTH/2-150, HEIGHT/2+150, color.White)
	}
	// draw mouse position, untuk debugging
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("last update: %v", lastUpdate))
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Mouse X: %v\nMouse Y: %v", g.mousePosition.x, g.mousePosition.y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Kkona Shooting gallery")

	// load font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatalf("error parse font: %v", err)
	}

	tt2, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatalf("error parse font: %v", err)
	}

	titleF, _ = opentype.NewFace(tt2, &opentype.FaceOptions{
		Size:    30,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	subTitleF, _ = opentype.NewFace(tt2, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	creditF, _ = opentype.NewFace(tt2, &opentype.FaceOptions{
		Size:    15,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	f, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	gameOverF, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    30,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// load assets
	shotgunShellImage, _, _ := ebitenutil.NewImageFromFile("./assets/icon_bullet_gold_long.png")
	bgImage, _, _ := ebitenutil.NewImageFromFile("./assets/bg.png")
	crossHair, _, _ := ebitenutil.NewImageFromFile("./assets/crosshair.png")
	targetImg, _, _ := ebitenutil.NewImageFromFile("./assets/target.png")

	g := new(Game)
	g.bg = bgImage
	g.crossHairImg = crossHair
	g.targetImg = targetImg
	g.shellImg = shotgunShellImage

	// sizeTarget := 10
	for i := 0; i < 5; i++ {
		rand.Seed(time.Now().UnixNano())
		td := Target{
			img:    g.targetImg,
			speedX: float64(rand.Intn(6-3) + 3),
			speedY: float64(rand.Intn(6-3) + 3),
			x:      float64(rand.Intn((WIDTH-33)-33) + 33),
			y:      float64(rand.Intn((HEIGHT-33)-0) + 0),
		}

		g.listTarget = append(g.listTarget, td)
	}

	ammo = 15

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
