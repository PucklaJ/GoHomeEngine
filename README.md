# GoHomeEngine
[![godoc](https://godoc.org/github.com/PucklaJ/GoHomeEngine/src/gohome?status.svg)](https://godoc.org/github.com/PucklaJ/GoHomeEngine/src/gohome)
[![License: Zlib](https://img.shields.io/badge/License-Zlib-green.svg)](https://github.com/PucklaJ/GoHomeEngine/blob/master/LICENSE.md)
[![GitHub last commit](https://img.shields.io/github/last-commit/PucklaJ/GoHomeEngine.svg)](https://github.com/PucklaJ/GoHomeEngine/commits/master)
<br>
A Game Engine for 2D/3D Games written in go

## Features

##### General
* Multiple Platform support: Windows, Linux, Mac, Android, Browser
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
* Shadows of DirectionalLight and SpotLight
* Ray Casting
* [Level Editor](https://github.com/PucklaJ/GoHomeEdit) (in development)

##### 2D
* Rendering 2D Sprites
* Camera (Translating, Rotating and Zooming)
* Sprite Animation
* Rendering 2D Shapes (Point,Line,Rectangle,Polygon, etc.)
* [Physics](https://box2d.org)
* [TiledMaps](https://www.mapeditor.org)

## Platforms

|             | Windows | Linux | Mac | Android | iOS | Browser |
| ----------- | ------- | ----- | --- | ------- | --- | ------- |
| Tested      | Yes     | Yes   | No  | Yes     | No  | Yes     |
| Implemented | Yes     | Yes   | Yes | Yes     | No  | Yes     |

## Dependencies

+ [mathgl](https://github.com/PucklaJ/mathgl) ([Forked from here](https://github.com/go-gl/mathgl))([License](https://github.com/PucklaJ/mathgl/blob/master/LICENSE))
+ [tga](https://github.com/blezek/tga) ([License](https://github.com/blezek/tga/blob/master/LICENSE.MIT))
+ [go-openal](https://github.com/PucklaJ/go-openal) ([Forked from here](https://github.com/timshannon/go-openal)) ([License](https://github.com/PucklaJ/go-openal/blob/master/LICENSE))
+ [go-wav](https://github.com/PucklaJ/go-wav) ([Forked from here](https://github.com/sdobz/go-wav)) ([License](https://github.com/verdverm/go-wav/blob/master/LICENSE.md))
+ [go-mp3](https://github.com/hajimehoshi/go-mp3) ([License](https://github.com/hajimehoshi/go-mp3/blob/master/LICENSE))
+ [box2d](https://github.com/ByteArena/box2d) ([License](https://github.com/ByteArena/box2d/blob/master/LICENSE.md))
+ [tmx](https://github.com/PucklaJ/tmx) ([Forked from here](https://github.com/elliotmr/tmx)) ([License](https://github.com/PucklaJ/tmx/blob/master/LICENSE))
+ [glfw](https://github.com/glfw/glfw) ([License](https://github.com/glfw/glfw/blob/master/LICENSE.md))
+ [gtk](https://gtk.org) ([License](http://www.gnu.org/licenses/old-licenses/lgpl-2.1.html))
+ [go-sdl2](https://github.com/PucklaJ/go-sdl2) ([Forked from here](https://github.com/veandco/go-sdl2)) ([License](https://github.com/PucklaJ/go-sdl2/blob/master/LICENSE))
+ [gopherjs](https://github.com/gopherjs/gopherjs) ([License](https://github.com/gopherjs/gopherjs/blob/master/LICENSE))
+ [go-gl/gl](https://github.com/go-gl/gl) ([License](https://github.com/go-gl/gl/blob/master/LICENSE))
+ [android-go](https://github.com/PucklaJ/android-go) ([Forked from here](https://github.com/xlab/android-go)) ([License](https://github.com/PucklaJ/android-go/blob/master/LICENSE.txt))
+ [webgl](https://github.com/PucklaJ/webgl) ([Forked from here](https://github.com/gopherjs/webgl)) ([License](https://github.com/PucklaJ/webgl/blob/master/LICENSE))

## Installation
1. Install the c-Dependencies:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// For Desktop (Most of them should already be installed)
	sudo apt-get install libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libpthread-stubs0-dev zlib1g-dev libgl1-mesa-dev libxi-dev
	// For GTK
	sudo apt-get install libgtk-3-dev
	// On Windows use msys and execute one of the following commands
	pacman -S mingw-w64-x86_64-gtk3 // for 64-Bit
	pacman -S mingw-w64-i686-gtk3   // for 32-Bit
	// For SDL2
	sudo apt-get install libsdl2-dev
	// On Windows use msys and execute one of the following commands
	pacman -S mingw-w64-x86_64-sdl2 // for 64-Bit
	pacman -S mingw-w64-i686-sdl2   // for 32-Bit

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
2. Install the go-Dependencies:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    // Desktop
	// Default (GLFW,OpenGL) if you only want to build desktop applications only execute this command
	go get -u github.com/PucklaJ/mathgl/mgl32 github.com/PucklaJ/tmx github.com/ByteArena/box2d github.com/PucklaJ/GLSLGenerator github.com/go-gl/glfw/v3.2 github.com/PucklaJ/go-openal/openal github.com/PucklaJ/go-wav github.com/hajimehoshi/go-mp3 github.com/go-gl/gl/all-core/gl

	// Use some of the following commands to build for a different platform or if you want to use
	// a different Framework or Renderer

	// Always Needed
	go get -u github.com/PucklaJ/mathgl/mgl32 github.com/PucklaJ/tmx github.com/ByteArena/box2d github.com/PucklaJ/GLSLGenerator

	// For GLFW
	go get -u github.com/go-gl/glfw/v3.2 github.com/PucklaJ/go-openal/openal github.com/PucklaJ/go-wav github.com/hajimehoshi/go-mp3

	// For SDL2
	go get -u github.com/PucklaJ/go-sdl2/sdl github.com/PucklaJ/go-openal/openal github.com/PucklaJ/go-wav github.com/hajimehoshi/go-mp3

	// For GTK
	go get -u github.com/PucklaJ/go-openal/openal github.com/PucklaJ/go-wav github.com/hajimehoshi/go-mp3

	// For OpenGL
	go get -u github.com/go-gl/gl/all-core/gl

	// Browser
	go get -u github.com/gopherjs/gopherjs github.com/PucklaJ/webgl

	// Android
	go get -u github.com/PucklaJ/go-sdl2/sdl github.com/PucklaJ/android-go/gles2

	// For OpenGLES2
	go get -u github.com/PucklaJ/android-go/gles2

	// For OpenGLES3
	go get -u github.com/PucklaJ/android-go/gles3

	// For OpenGLES31
	go get -u github.com/PucklaJ/android-go/gles31
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
3. Compile one of the examples to test:<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	cd $GOPATH/src/github.com/PucklaJ/GoHomeEngine/src/examples/basic
	go build && ./basic
	// You should see a gopher in the middle
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Optional: Install libraries for audio on android<br>
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    // If you want to use audio on android you first have to install SDL_mixer with mp3 and ogg support
	// This guide is only for linux ubuntu and debian. For other platforms it should be similiar or the
	// default precompiled binaries already include mp3 and ogg support.

	// go to some directory where you want to store the respositories
	cd $GIT_HOME
	// clone and install mpg123
	git clone https://github.com/kepstin/mpg123.git
	cd mpg123
	./configure
	make -j4
	sudo make install
	cd ..
	// install libvorbis
	sudo apt-get install libvorbis-dev
	// clone and install SDL_mixer
	git clone https://github.com/SDL-mirror/SDL_mixer.git
	cd SDL_mixer
	./configure --enable-music-mp3 --enable-music-ogg
	make -j4
	sudo make install
	// Now you are ready to use audio on android
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

## Tutorial

The following code describes what is needed to write a game with GoHome

```go
package main

import (
	"github.com/PucklaJ/GoHomeEngine/src/frameworks/GLFW"
	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	"github.com/PucklaJ/GoHomeEngine/src/renderers/OpenGL"
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

## Build tool

There is a build tool available with which you can build your games easily for different platforms
Read more [here](https://github.com/PucklaJ/GoHomeEngine/blob/master/src/cmd/gohome/README.md).