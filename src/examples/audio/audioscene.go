package main

import "github.com/PucklaMotzer09/gohomeengine/src/gohome"

type AudioScene struct {

}

func (this *AudioScene) Init() {
	gohome.Framew.GetAudioManager().Init()
	gohome.ResourceMgr.LoadSound("Explode","explode.wav")
	sound := gohome.ResourceMgr.GetSound("Explode")
	sound.Play()
}

func (this *AudioScene) Update(delta_time float32) {

}

func (this *AudioScene) Terminate() {
	gohome.Framew.GetAudioManager().Terminate()
}
