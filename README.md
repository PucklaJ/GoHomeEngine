# GoHomeEngine
[![godoc](https://godoc.org/github.com/PucklaMotzer09/GoHomeEngine/src/gohome?status.svg)](https://godoc.org/github.com/PucklaMotzer09/GoHomeEngine/src/gohome)
[![License: Zlib](https://img.shields.io/badge/License-Zlib-green.svg)](https://github.com/PucklaMotzer09/GoHomeEngine/blob/master/LICENSE.md)
[![GitHub last commit](https://img.shields.io/github/last-commit/PucklaMotzer09/GoHomeEngine.svg)](https://github.com/PucklaMotzer09/GoHomeEngine/commits/master)
<br>
A game engine for 2D and 3D games written in go

## Dependencies

##### General
+ [mathgl](https://github.com/PucklaMotzer09/mathgl) ([Forked from here](https://github.com/go-gl/mathgl))([License](https://github.com/PucklaMotzer09/mathgl/blob/master/LICENSE))
+ [tga](https://github.com/blezek/tga) ([License](https://github.com/blezek/tga/blob/master/LICENSE.MIT))
+ [go-openal](https://github.com/timshannon/go-openal) ([License](https://github.com/timshannon/go-openal/blob/master/LICENSE))
+ [go-wav](https://github.com/PucklaMotzer09/go-wav) ([Forked from here](https://github.com/sdobz/go-wav)) ([License](https://github.com/verdverm/go-wav/blob/master/LICENSE.md))
+ [go-mp3](https://github.com/hajimehoshi/go-mp3) ([License](https://github.com/hajimehoshi/go-mp3/blob/master/LICENSE))
+ [box2d](https://github.com/BxteArena/box2d)([License](https://github.com/ByteArena/box2d/blob/master/LICENSE.md))
+ [tmx](https://github.com/elliotmr/tmx)([License](https://github.com/elliotmr/tmx/blob/master/LICENSE))

##### GLFWFramework
+ [assimp](https://github.com/assimp/assimp) ([License](https://github.com/assimp/assimp/blob/master/LICENSE))
+ [glfw](https://github.com/glfw/glfw) ([License](https://github.com/glfw/glfw/blob/master/LICENSE.md))

##### GTKFramework
+ [gtk](https://gtk.org) ([License](http://www.gnu.org/licenses/old-licenses/lgpl-2.1.html))

##### AndroidFramework
+ [AndroidSDK](https://developer.android.com/studio/) ([License](https://developer.android.com/studio/terms))
+ [gomobile](https://github.com/golang/mobile) ([License](https://github.com/golang/mobile/blob/master/LICENSE))

##### OpenGLRenderer
+ [go-gl/gl](https://github.com/go-gl/gl) ([License](https://github.com/go-gl/gl/blob/master/LICENSE))

##### OpenGLESRenderer
+ [gomobile/gles](https://github.com/golang/mobile/tree/master/gl) ([License](https://github.com/golang/mobile/blob/master/LICENSE))

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
* Audio (.wav and .mp3)

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
* Sprite Animation

## Installation
1. Install the c-Dependencies:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// For Desktop (The windows assimp library is provided with this repository)
	sudo apt-get install libassimp-dev libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libpthread-stubs0-dev zlib1g-dev libglfw3-dev libgl1-mesa-dev libxi-dev
	// For GTK
	sudo apt-get install libgtk-3-dev
	// On Windows use msys and execute one of the following commands
	pacman -S mingw-w64-x86_64-gtk3 // for 64-Bit
	pacman -S mingw-w64-i686-gtk3   // for 32-Bit
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
2. Install the go-Dependencies:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// For Desktop
	go get -u github.com/go-gl/gl/all-core
	go get -u github.com/go-gl/glfw/v3.2
	go get -u github.com/raedatoui/assimp
	go get -u github.com/phf/go-openal
	go get -u github.com/PucklaMotzer09/go-wav
	go get -u github.com/hajimehoshi/go-mp3
	// For Mobile
	go get -u golang.org/x/mobile/cmd/gomobile
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
3. Compile one of the examples to test:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	cd $GOPATH/src/github.com/PucklaMotzer09/gohomeengine/src/examples/basic
	$GOPATH/src/github.com/PucklaMotzer09/gohomeengine/build.sh -linux -run
	// You should see a gopher in the middle (use -windows instead of -linux on windows)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

## Tutorial

The following code describes what is needed to write a game with GoHome

```go
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
```

This program opens a window with a black background. To learn more you can look at the examples in src/examples
