package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SunSprite struct {
	texture         rl.Texture2D
	position        rl.Vector2
	currentFrame    int
	frameCounter    int
	framesPerSecond int
	maxFrames       int
	spriteSize      float32
	sizeMultiplier  float32
}

func (s *SunSprite) Draw() {
	sourceRec := rl.Rectangle{
		X:      float32(s.currentFrame) * s.spriteSize,
		Y:      s.spriteSize,
		Width:  s.spriteSize,
		Height: s.spriteSize,
	}

	destRec := rl.Rectangle{
		X:      s.position.X,
		Y:      s.position.Y,
		Width:  s.spriteSize * s.sizeMultiplier,
		Height: s.spriteSize * s.sizeMultiplier,
	}
	rl.DrawTexturePro(s.texture, sourceRec, destRec, rl.Vector2{X: 0, Y: 0}, 0, rl.White)
}

func (s *SunSprite) Update() {
	s.frameCounter++
	s.position.Y = s.position.Y + 0.1*float32(math.Sin(rl.GetTime()))
	if s.frameCounter >= (60 / s.framesPerSecond) {
		s.frameCounter = 0
		s.currentFrame++
		if s.currentFrame >= s.maxFrames {
			s.currentFrame = 0
		}
	}
}

func NewSunSprite(x, y float32) *SunSprite {
	rl.LoadTexture("sun.png")
	return &SunSprite{
		texture:         rl.LoadTexture("sun.png"),
		position:        rl.NewVector2(x, y),
		currentFrame:    0,
		frameCounter:    0,
		maxFrames:       64,
		framesPerSecond: 10,
		spriteSize:      200,
		sizeMultiplier:  0.5,
	}
}
