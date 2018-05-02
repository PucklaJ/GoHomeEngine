# GoHomeEngine
[![godoc](https://godoc.org/github.com/PucklaMotzer09/gohomeengine/src/gohome?status.svg)](https://godoc.org/github.com/PucklaMotzer09/gohomeengine/src/gohome)
[![AUR](https://img.shields.io/aur/license/yaourt.svg)](https://github.com/PucklaMotzer09/GoHomeEngine/blob/master/LICENSE.md)
<br>
A game engine for 2D and 3D games written in go

## Dependencies
- [go-gl/gl](https://github.com/go-gl/gl)
- [go-gl/mathgl](https://github.com/go-gl/mathgl)
- [assimp](https://github.com/assimp/assimp)
  - [bindings](https://github.com/raedatoui/assimp)
- [glfw](https://github.com/glfw/glfw)
  - [bindings](https://github.com/go-gl/glfw)
- [tga](https://github.com/blezek/tga)
- [gomobile](https://github.com/golang/mobile)

## Platforms

* Windows
* Linux (Ubuntu)
* Mac (not tested) 
* Android
* iOS (not tested) 


## Features

#### General
* Loading Shaders
* Multiple Viewports

#### 3D
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

#### 2D
* Rendering 2D Sprites
* Camera (Translating, Rotating and Zooming) 

