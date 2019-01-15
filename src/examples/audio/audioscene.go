package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/audio"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"runtime"
)

type AudioScene struct {
	sound gohome.Sound
	music gohome.Music
}

func (this *AudioScene) Init() {
	audio.InitAudio()
	gohome.ResourceMgr.LoadSound("Bottle", "assets/bottle.wav")
	gohome.ResourceMgr.LoadMusic("TownTheme", "assets/TownTheme.mp3")

	gohome.ResourceMgr.GetMusic("TownTheme").Play(true)
}

func (this *AudioScene) Update(delta_time float32) {
	if gohome.InputMgr.JustPressed(gohome.Key1) || (runtime.GOOS == "android" && gohome.InputMgr.JustPressed(gohome.KeyBack)) {
		if this.music == nil {
			this.music = gohome.ResourceMgr.GetMusic("TownTheme")
			if this.music != nil && !this.music.IsPlaying() {
				this.music.Play(true)
			} else if this.music != nil {
				this.music.Pause()
			}
		} else {
			if this.music.IsPlaying() {
				this.music.Pause()
			} else {
				this.music.Resume()
			}
		}
	}

	if gohome.InputMgr.JustPressed(gohome.Key2) || (runtime.GOOS == "android" && gohome.InputMgr.JustPressed(gohome.MouseButtonLeft)) {
		if this.sound == nil {
			this.sound = gohome.ResourceMgr.GetSound("Bottle")
			if this.sound != nil {
				this.sound.Play(false)
			}
		} else {
			if this.sound.IsPlaying() {
				this.sound.Pause()
			} else {
				this.sound.Resume()
			}
		}
	}
}

func (this *AudioScene) Terminate() {
	gohome.AudioMgr.Terminate()
}
