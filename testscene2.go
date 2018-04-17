package main

import (
	"fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
)

type TestScene2 struct {
}

func (this *TestScene2) Init() {
	fmt.Println("Hello")
	gohome.InitDefaultValues()
	gohome.FPSLimit.MaxFPS = 1000

	gohome.ResourceMgr.LoadTexture("Icon", "icon.png")
}

func (this *TestScene2) Update(delta_time float32) {

}

func (this *TestScene2) Terminate() {

}
