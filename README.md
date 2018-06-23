# GoHomeEngine
[![godoc](https://godoc.org/github.com/PucklaMotzer09/gohomeengine/src/gohome?status.svg)](https://godoc.org/github.com/PucklaMotzer09/gohomeengine/src/gohome)
[![License: Zlib](https://img.shields.io/badge/License-Zlib-green.svg)](https://github.com/PucklaMotzer09/GoHomeEngine/blob/master/LICENSE.md)
[![GitHub last commit](https://img.shields.io/github/last-commit/PucklaMotzer09/GoHomeEngine.svg)](https://github.com/PucklaMotzer09/GoHomeEngine/commits/master)
<br>
A game engine for 2D and 3D games written in go

## Dependencies

##### General
+ [go-gl/mathgl](https://github.com/go-gl/mathgl)
+ [tga](https://github.com/blezek/tga)

##### GLFWFramework
+ [assimp](https://github.com/assimp/assimp)
+ [glfw](https://github.com/glfw/glfw)

##### AndroidFramework
+ [AndroidSDK](https://developer.android.com/studio/)
+ [gomobile](https://github.com/golang/mobile)

##### OpenGLRenderer
+ [go-gl/gl](https://github.com/go-gl/gl)

##### OpenGLESRenderer
+ [gomobile/gles](https://github.com/golang/mobile)

## Platforms

|				|Windows| Linux		| Mac		| Android 	| iOS	| Browser |
|---------------|-------|-----------|-----------|-----------|-------|---------|
|Tested 		|	Yes |	Yes		|	No 		|	Yes		|	No 	|	No    |
|Implemented  	|   Yes |   Yes		|	Yes		|   Yes		|   Yes |   No    |

## Features

##### General
* Loading Shaders
* Multiple Viewports
* Tweens

##### 3D
* Rendering 3D Models
* Camera
* Loading 3D Models (with [assimp](http://assimp.org/))
* Materials
* SpecularMaps
* NormalMaps
* PointLights
* DirectionalLights
* SpotLights
* Shadows of all three lighttypes

##### 2D
* Rendering 2D Sprites
* Camera (Translating, Rotating and Zooming) 

## Installation
1. Install the c-Dependencies:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// For Desktop (The windows assimp library is provided with this repository)
	sudo apt-get install libassimp-dev libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libpthread-stubs0-dev zlib1g-dev libglfw3-dev libgl1-mesa-dev libxi-dev
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
2. Install the go-Dependencies:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    // For Desktop
	go get -u github.com/go-gl/gl/all-core
	go get -u github.com/go-gl/glfw/v3.2
	go get -u github.com/raedatoui/assimp
	// For Mobile
	go get -u golang.org/x/mobile/cmd/gomobile
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
3. Compile one of the examples to test:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	cd $GOPATH/src/github.com/PucklaMotzer09/gohomeengine/src/examples/basic
	$GOPATH/src/github.com/PucklaMotzer09/gohomeengine/build.sh -linux -run
	// You should see a gopher in the middle
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

## Tutorial

The following code describes what is needed to write a game with GoHome

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~{.go}
package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GLFW"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/renderers/OpenGL"
)

type StartScene struct {

}

func (*StartScene) Init() {

}

func (*StartScene) Update(delta_time float32) {

}

func (*StartScene) Terminate() {

}

func main() {
	gohome.MainLop.Run(&framework.GLFWFramework{},&renderer.OpenGLRenderer{},1280,720,"Example",&StartScene{})
}
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

This program opens a window with a black background. To learn more you can look at the examples in src/examples
