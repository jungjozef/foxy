package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth      = 800
	screenHeight     = 450
	spriteSize       = 32
	platformWidth    = 500
	platformHeight   = 450 // Csökkentett érték
	initialPlatformY = 300 // Az első platform Y pozíciója
	platformGapX     = 50  // Növelt vízszintes távolság
	gravity          = 0.5
	jumpForce        = -12
	scrollSpeed      = 3.
	fallLimit        = screenHeight + spriteSize*2
	maxHeightOffset  = 50
	platformCount    = 3 // Csökkentett platform szám
)

type Platform struct {
	position rl.Vector2
	width    float32
	height   float32
}

type Player struct {
	texture         rl.Texture2D
	position        rl.Vector2
	velocityY       float32
	isJumping       bool
	currentFrame    int
	frameCounter    int
	framesPerSecond int
	row             int
	maxFrames       int
	flipX           bool
}

func NewPlayer(texture rl.Texture2D, posX, posY float32, framesPerSecond, row int) *Player {
	return &Player{
		texture:         texture,
		position:        rl.Vector2{X: posX, Y: posY},
		velocityY:       0,
		isJumping:       false,
		currentFrame:    0,
		frameCounter:    0,
		framesPerSecond: framesPerSecond,
		row:             row,
		maxFrames:       8,
		flipX:           false,
	}
}

func (p *Player) Update(platforms []Platform, landSnd rl.Sound) {
	// Animation update
	p.frameCounter++
	if p.frameCounter >= (60 / p.framesPerSecond) {
		p.frameCounter = 0
		p.currentFrame++
		if p.currentFrame >= p.maxFrames {
			p.currentFrame = 0
		}
	}

	// Apply gravity with terminal velocity
	p.velocityY += gravity
	if p.velocityY > 12 { // Add terminal velocity
		p.velocityY = 12
	}

	// Store the next position for collision checking
	nextY := p.position.Y + p.velocityY

	// Reset jumping state for new collision checks
	wasOnPlatform := false

	// Collision detection with platforms
	for _, platform := range platforms {
		// Calculate player bounds
		playerBottom := nextY + float32(spriteSize*2)
		playerTop := nextY
		playerLeft := p.position.X + float32(spriteSize*0.5)  // Add offset for more precise collision
		playerRight := p.position.X + float32(spriteSize*1.5) // Adjust collision box

		// Calculate platform bounds
		platformTop := platform.position.Y
		platformBottom := platform.position.Y + platform.height
		platformLeft := platform.position.X
		platformRight := platform.position.X + platform.width

		// Check for collision
		if playerRight > platformLeft &&
			playerLeft < platformRight &&
			playerBottom > platformTop &&
			playerTop < platformBottom {

			// Landing on top of platform
			if p.velocityY > 0 && p.position.Y+float32(spriteSize*2) <= platformTop+5 {
				p.position.Y = platformTop - float32(spriteSize*2)
				p.velocityY = 0

				wasOnPlatform = true
			}
		}
	}

	// Update jumping state
	if p.isJumping && wasOnPlatform {
		rl.PlaySound(landSnd)
	}
	p.isJumping = !wasOnPlatform

	// Update position
	p.position.Y += p.velocityY
}

func (p *Player) Draw() {
	sourceRec := rl.Rectangle{
		X:      float32(p.currentFrame * spriteSize),
		Y:      float32(p.row * spriteSize),
		Width:  spriteSize,
		Height: spriteSize,
	}

	// Tükrözi a textúrát, ha szükséges
	if p.flipX {
		sourceRec.Width = -spriteSize
	}

	destRec := rl.Rectangle{
		X:      p.position.X,
		Y:      p.position.Y,
		Width:  spriteSize * 2, // Nagyítás a jobb láthatóságért
		Height: spriteSize * 2,
	}

	rl.DrawTexturePro(p.texture, sourceRec, destRec, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
}

func GeneratePlatform(xPos float32, prevY float32) Platform {
	// Calculate new Y position with smoother transitions
	yVariation := float32(rand.Intn(2*maxHeightOffset+1)) - float32(maxHeightOffset)
	newY := prevY + yVariation

	// Enforce platform height limits
	minY := float32(screenHeight - 300)
	maxY := float32(screenHeight - 100)

	if newY > maxY {
		newY = maxY
	}
	if newY < minY {
		newY = minY
	}

	// Ensure the height difference isn't too extreme
	maxDiff := float32(60)
	if newY > prevY+maxDiff {
		newY = prevY + maxDiff
	} else if newY < prevY-maxDiff {
		newY = prevY - maxDiff
	}

	return Platform{
		position: rl.Vector2{
			X: xPos,
			Y: newY,
		},
		width:  platformWidth,
		height: platformHeight,
	}
}

type Sun struct {
	radius   float32
	position rl.Vector2
}

func NewSun(radius float32, position rl.Vector2) Sun {
	return Sun{
		radius:   radius,
		position: position,
	}
}

func (s *Sun) Draw() {
	rl.DrawCircleGradient(int32(s.position.X), int32(s.position.Y), s.radius, rl.Yellow, rl.Orange)
}

func (s *Sun) Update() {
	//s.position.X -= s.radius / 60
	//if s.position.X < -2*s.radius {
	//	s.position.X = screenWidth + s.radius
	//}
	s.position.Y = 70 + 5*float32(math.Sin(rl.GetTime()))
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Endless Runner")
	defer rl.CloseWindow()
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	// Cél FPS beállítása
	rl.SetTargetFPS(60)

	// Véletlenszám generátor inicializálása
	rand.Seed(time.Now().UnixNano())

	// Betölti a sprite sheet-et
	spriteSheet := rl.LoadTexture("fox.png")
	defer rl.UnloadTexture(spriteSheet)

	jmpSnd := rl.LoadSound("jump.wav")
	defer rl.UnloadSound(jmpSnd)
	landSnd := rl.LoadSound("landing.wav")
	defer rl.UnloadSound(landSnd)

	// Játékos inicializálása az első platformon
	playerStartX := float32(100)
	playerStartY := float32(initialPlatformY - spriteSize*3)
	player := NewPlayer(spriteSheet, playerStartX, playerStartY, 8, 2)

	// Platformok inicializálása
	platforms := make([]Platform, platformCount)
	platforms[0] = Platform{
		position: rl.Vector2{X: 0, Y: float32(initialPlatformY)},
		width:    platformWidth,
		height:   platformHeight,
	}
	for i := 1; i < platformCount; i++ {
		platforms[i] = GeneratePlatform((float32(i)+platforms[i-1].position.X+platforms[i-1].width)+platformGapX, platforms[i-1].position.Y)
	}

	sunSpr := NewSunSprite(0, 10)

	gameOver := false
	score := 0

	// music
	musToPlay := rand.Intn(12)
	inGameMus := NewInGameMusic(fmt.Sprintf("music%02d.ogg", musToPlay+1))

	gameOverMus := NewInGameMusic("gameover.ogg")

	for !rl.WindowShouldClose() {
		// Bemenet kezelése
		if !gameOver {
			gameOverMus.Stop()
			inGameMus.Update()
			// Ugrás kezelése
			if rl.IsKeyPressed(rl.KeySpace) && !player.isJumping {
				player.velocityY = jumpForce
				player.isJumping = true
				rl.PlaySound(jmpSnd)
			}

			// Mozgás kezelése
			// Ha szükséges, engedélyezheted a balra és jobbra mozgást
			// Például:
			/*
				if rl.IsKeyDown(rl.KeyRight) {
					player.position.X += 2.0
					player.row = 1 // Járás animáció
					player.flipX = false
				} else if rl.IsKeyDown(rl.KeyLeft) {
					player.position.X -= 2.0
					player.row = 1 // Járás animáció
					player.flipX = true
				} else {
					player.row = 0 // Állás animáció
				}
			*/

			// Játékos frissítése (gravitáció és ütközés)
			player.Update(platforms, landSnd)

			// Platformok görgetése
			for i := range platforms {
				platforms[i].position.X -= scrollSpeed
			}

			// Platformok újragenerálása, ha kilépnek a képernyőről
			for i := 0; i < len(platforms); i++ {
				if platforms[i].position.X+platforms[i].width < 0 {
					// A legjobban jobbra lévő platform meghatározása
					rightmostPlatform := platforms[0]
					for _, p := range platforms {
						if p.position.X > rightmostPlatform.position.X {
							rightmostPlatform = p
						}
					}
					// Új platform generálása a legjobban jobbra lévő platform után
					newPlatform := GeneratePlatform(rightmostPlatform.position.X+platformGapX+rightmostPlatform.width, rightmostPlatform.position.Y)
					platforms[i] = newPlatform
					score++ // Növeli a pontszámot
				}
			}

			// Ellenőrzi, hogy a játékos lehullott-e
			if player.position.Y > fallLimit {
				gameOver = true
			}
		} else {
			gameOverMus.Update()
			inGameMus.Stop()
			//inGameMus.Close()
			// Játék vége állapot kezelése
			if rl.IsKeyPressed(rl.KeyR) {
				musToPlay = rand.Intn(12)
				inGameMus = NewInGameMusic(fmt.Sprintf("music%02d.ogg", musToPlay+1))
				// Játék visszaállítása
				player.position = rl.Vector2{X: playerStartX, Y: playerStartY}
				player.velocityY = 0
				player.isJumping = false
				gameOver = false
				score = 0

				// Platformok visszaállítása
				platforms[0] = Platform{
					position: rl.Vector2{X: 0, Y: float32(initialPlatformY)},
					width:    platformWidth,
					height:   platformHeight,
				}
				for i := 1; i < platformCount; i++ {
					platforms[i] = GeneratePlatform((float32(i)+platforms[i-1].position.X+platforms[i-1].width)+platformGapX, platforms[i-1].position.Y)
				}
			}
		}
		sunSpr.Update()

		// Rajzolás
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawRectangleGradientV(0, 0, screenWidth, screenHeight-40, rl.SkyBlue, rl.LightGray)
		rl.DrawRectangleGradientV(0, screenHeight-40, screenWidth, screenHeight, rl.Brown, rl.Beige)

		sunSpr.Draw()

		// Platformok rajzolása
		for _, platform := range platforms {
			rl.DrawRectangle(int32(platform.position.X), int32(platform.position.Y), int32(platform.width), int32(platform.height), rl.Gray)
			//rl.DrawRectangleLines(int32(platform.position.X), int32(platform.position.Y), int32(platform.width), int32(platform.height), rl.Red)

		}

		// Játékos rajzolása
		player.Draw()

		// Pontszám megjelenítése
		rl.DrawText("Score: "+strconv.Itoa(score), 10, 10, 20, rl.Black)

		// "Game Over" üzenet megjelenítése
		if gameOver {
			rl.DrawText("Game Over! Press R to Restart", screenWidth/2-150, screenHeight/2-10, 20, rl.Red)
		}

		// Debugging: Rajzolj bounding box-okat
		//rl.DrawRectangleLines(int32(player.position.X), int32(player.position.Y), int32(spriteSize*2), int32(spriteSize*2), rl.Blue)

		rl.EndDrawing()
	}
}
