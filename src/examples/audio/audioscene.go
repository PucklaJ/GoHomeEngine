package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"fmt"
)

type AudioScene struct {
	sound gohome.Sound
	music gohome.Music
}

func (this *AudioScene) Init() {
	gohome.Framew.GetAudioManager().Init()
	gohome.ResourceMgr.PreloadSound("Bottle","bottle.wav")
	gohome.ResourceMgr.PreloadMusic("TownTheme","TownTheme.mp3")
	gohome.ResourceMgr.LoadPreloadedResources()
}

func (this *AudioScene) Update(delta_time float32) {
	if gohome.InputMgr.JustPressed(gohome.Key1) {
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

	if gohome.InputMgr.JustPressed(gohome.Key2) {
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
		fmt.Println("Music:",this.music.GetPlayingDuration())
	}
}

func (this *AudioScene) Terminate() {
	gohome.Framew.GetAudioManager().Terminate()
}
