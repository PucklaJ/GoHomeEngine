# GoHomeEngine
[![godoc](https://godoc.org/github.com/PucklaMotzer09/gohomeengine/src/gohome?status.svg)](https://godoc.org/github.com/PucklaMotzer09/gohomeengine/src/gohome)
[![License: Zlib](https://img.shields.io/badge/License-Zlib-green.svg)](https://github.com/PucklaMotzer09/GoHomeEngine/blob/master/LICENSE.md)
[![GitHub last commit](https://img.shields.io/github/last-commit/PucklaMotzer09/GoHomeEngine.svg)](https://github.com/PucklaMotzer09/GoHomeEngine/commits/master)
![GitHub repo size in bytes](https://img.shields.io/github/repo-size/PucklaMotzer09/GoHomeEngine.svg)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/3a84c5eb27bf48099d9e2322b571fff5)](https://www.codacy.com/app/PucklaMotzer09/GoHomeEngine?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=PucklaMotzer09/GoHomeEngine&amp;utm_campaign=Badge_Grade)
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

