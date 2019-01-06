package main

import (
	"fmt"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"runtime"
)

type AudioScene struct {
	sound gohome.Sound
	music gohome.Music
}

func (this *AudioScene) Init() {
	gohome.Framew.GetAudioManager().Init()
	gohome.ResourceMgr.LoadSound("Bottle", "bottle.wav")
	gohome.ResourceMgr.LoadMusic("TownTheme", "TownTheme.mp3")
}

func (this *AudioScene) Update(delta_time float32) {
	if gohome.InputMgr.JustPressed(gohome.Key1) || (runtime.GOOS == "android" && gohome.InputMgr.JustPressed(gohome.KeyBack)) {
		if this.music == nil {
			this.music = gohome.ResourceMgr.GetMusic("TownTheme")
			if this.music != nil {
				this.music.Play(true)
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

	if this.music != nil {
		fmt.Println("Music:", this.music.GetPlayingDuration())
	}
}

func (this *AudioScene) Terminate() {
	gohome.Framew.GetAudioManager().Terminate()
}
