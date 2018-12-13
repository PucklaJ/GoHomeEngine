# GoHomeEngine
[![godoc](https://godoc.org/github.com/PucklaMotzer09/GoHomeEngine/src/gohome?status.svg)](https://godoc.org/github.com/PucklaMotzer09/GoHomeEngine/src/gohome)
[![License: Zlib](https://img.shields.io/badge/License-Zlib-green.svg)](https://github.com/PucklaMotzer09/GoHomeEngine/blob/master/LICENSE.md)
[![GitHub last commit](https://img.shields.io/github/last-commit/PucklaMotzer09/GoHomeEngine.svg)](https://github.com/PucklaMotzer09/GoHomeEngine/commits/master)
<br>
A Game Engine for 2D/3D Games written in go

## Features

##### General
* Loading Shaders
* Dynamic Shader Generation
* Multiple Viewports/Cameras
* Tweens
* Audio (.wav and .mp3)
* GUI with [GTK](https://www.gtk.org)
* Simple OnScreen GUI
* Instancing

##### 3D
* Rendering 3D Models
* Camera
* Loading 3D Models (.obj)
* Materials
* SpecularMaps
* NormalMaps
* PointLights
* DirectionalLights
* SpotLights
* Shadows of all three lighttypes
* Ray Casting
* [Level Editor](https://github.com/PucklaMotzer09/GoHomeEdit) (in development)

##### 2D
* Rendering 2D Sprites
* Camera (Translating, Rotating and Zooming) 
* Sprite Animation
* Rendering 2D Shapes (Point,Line,Rectangle,Polygon, etc.)
* [Physics](https://box2d.org)
* [TiledMaps](https://www.mapeditor.org)

## Dependencies

##### General
+ [mathgl](https://github.com/PucklaMotzer09/mathgl) ([Forked from here](https://github.com/go-gl/mathgl))([License](https://github.com/PucklaMotzer09/mathgl/blob/master/LICENSE))
+ [tga](https://github.com/blezek/tga) ([License](https://github.com/blezek/tga/blob/master/LICENSE.MIT))
+ [go-openal](https://github.com/timshannon/go-openal) ([Forked from here](https://github.com/phf/go-openal)) ([License](https://github.com/timshannon/go-openal/blob/master/LICENSE))
+ [go-wav](https://github.com/PucklaMotzer09/go-wav) ([Forked from here](https://github.com/sdobz/go-wav)) ([License](https://github.com/verdverm/go-wav/blob/master/LICENSE.md))
+ [go-mp3](https://github.com/hajimehoshi/go-mp3) ([License](https://github.com/hajimehoshi/go-mp3/blob/master/LICENSE))
+ [box2d](https://github.com/ByteArena/box2d) ([License](https://github.com/ByteArena/box2d/blob/master/LICENSE.md))
+ [tmx](https://github.com/PucklaMotzer09/tmx) ([Forked from here](https://github.com/elliotmr/tmx)) ([License](https://github.com/PucklaMotzer09/tmx/blob/master/LICENSE))

##### GLFWFramework
+ [glfw](https://github.com/glfw/glfw) ([License](https://github.com/glfw/glfw/blob/master/LICENSE.md))

##### GTKFramework
+ [gtk](https://gtk.org) ([License](http://www.gnu.org/licenses/old-licenses/lgpl-2.1.html))

##### OpenGLRenderer
+ [go-gl/gl](https://github.com/go-gl/gl) ([License](https://github.com/go-gl/gl/blob/master/LICENSE))

## Platforms

|				|Windows| Linux		| Mac		| Android 	| iOS	| Browser |
|---------------|-------|-----------|-----------|-----------|-------|---------|
|Tested 		|	Yes |	Yes		|	No 		|	No		|	No 	|	No    |
|Implemented  	|   Yes |   Yes		|	Yes		|   No		|   No  |   No    |

## Installation
1. Install the c-Dependencies:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// For Desktop (Most of them should already be installed)
	sudo apt-get install libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libpthread-stubs0-dev zlib1g-dev libglfw3-dev libgl1-mesa-dev libxi-dev
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
	go get -u github.com/PucklaMotzer09/mathgl
	go get -u github.com/timshannon/go-openal/openal
	go get -u github.com/PucklaMotzer09/go-wav
	go get -u github.com/hajimehoshi/go-mp3
	go get -u github.com/PucklaMotzer09/tmx
	go get -u github.com/ByteArena/box2d
	go get -u github.com/PucklaMotzer09/GLSLGenerator
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
3. Compile one of the examples to test:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	cd $GOPATH/src/github.com/PucklaMotzer09/GoHomeEngine/src/examples/basic
	go build && ./basic
	// You should see a gopher in the middle
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

## Tutorial

The following code describes what is needed to write a game with GoHome

```go
package main

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/GLFW"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/GoHomeEngine/src/renderers/OpenGL"
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
