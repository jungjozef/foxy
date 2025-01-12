package main

import rl "github.com/gen2brain/raylib-go/raylib"

type InGameMusic struct {
	stream rl.Music
}

func NewInGameMusic(fn string) *InGameMusic {
	mus := rl.LoadMusicStream(fn)
	return &InGameMusic{
		stream: mus,
	}
}

func (music *InGameMusic) Play() {
	rl.PlayMusicStream(music.stream)
}

func (music *InGameMusic) Pause() {
	rl.PauseMusicStream(music.stream)
}

func (music *InGameMusic) Stop() {
	if rl.IsMusicStreamPlaying(music.stream) {
		rl.StopMusicStream(music.stream)
	}
}

func (music *InGameMusic) Update() {
	if !rl.IsMusicStreamPlaying(music.stream) {
		music.Play()
	}
	rl.UpdateMusicStream(music.stream)
}
func (music *InGameMusic) Close() {
	rl.UnloadMusicStream(music.stream)
}
