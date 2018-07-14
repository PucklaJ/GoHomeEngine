package main

import "github.com/PucklaMotzer09/gohomeengine/src/gohome"

type AudioScene struct {
	sound gohome.Sound
}

func (this *AudioScene) Init() {
	gohome.Framew.GetAudioManager().Init()
	gohome.ResourceMgr.LoadSound("Bach","bach.wav")
	this.sound = gohome.ResourceMgr.GetSound("Bach")
	if this.sound != nil {
		this.sound.Play()
	}
}

func (this *AudioScene) Update(delta_time float32) {
	if gohome.InputMgr.JustPressed(gohome.MouseButtonLeft) {
		if this.sound.IsPlaying() {
			this.sound.Pause()
		} else {
			this.sound.Resume()
		}
	}
}

func (this *AudioScene) Terminate() {
	gohome.Framew.GetAudioManager().Terminate()
}
