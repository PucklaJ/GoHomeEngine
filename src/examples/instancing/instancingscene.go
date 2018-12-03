package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"golang.org/x/image/colornames"
	"log"
	"sync"
)

const SIZE int = 20
const USE_INSTANCING = true

type InstancingScene struct {
	ent gohome.InstancedEntity3D
	cam gohome.Camera3D
}

func (this *InstancingScene) Init() {
	cubeMesh := gohome.Box("Box", [3]float32{1.0, 1.0, 1.0}, false)
	if USE_INSTANCING {
		cubeInstanced := gohome.InstancedMesh3DFromMesh3D(cubeMesh)
		this.ent.InitMesh(cubeInstanced, uint32(SIZE*SIZE*SIZE))
		for x := 0; x < SIZE; x++ {
			for y := 0; y < SIZE; y++ {
				for z := 0; z < SIZE; z++ {
					this.ent.Transforms[x+y*SIZE+z*SIZE*SIZE].Position = mgl32.Vec3{float32(x) * 2.0, float32(y) * 2.0, float32(z) * 2.0}
				}
			}
		}

		this.ent.UpdateInstancedValues()
		this.ent.StopUpdatingInstancedValues = true
		gohome.RenderMgr.AddObject(&this.ent)
	} else {
		gohome.Init3DShaders()
		cubeMesh.Load()
		for x := 0; x < SIZE; x++ {
			for y := 0; y < SIZE; y++ {
				for z := 0; z < SIZE; z++ {
					var ent gohome.Entity3D
					ent.InitMesh(cubeMesh)
					ent.Transform.Position = [3]float32{float32(x) * 2.0, float32(y) * 2.0, float32(z) * 2.0}
					gohome.RenderMgr.AddObject(&ent)
				}
			}
		}
	}

	gohome.LightMgr.SetAmbientLight(colornames.Gray, 0)
	gohome.LightMgr.AddDirectionalLight(
		&gohome.DirectionalLight{
			Direction:     mgl32.Vec3{1.0, -1.0, -1.0},
			DiffuseColor:  colornames.White,
			SpecularColor: colornames.Black,
			CastsShadows:  0,
		}, 0,
	)

	var move gohome.TestCameraMovement3D
	this.cam.Init()
	move.Init(&this.cam)
	gohome.UpdateMgr.AddObject(&move)

	this.cam.Position[2] = 3.0
	gohome.RenderMgr.SetCamera3D(&this.cam, 0)
}

var addedRows int = 1
var removedRows int = 0

func (this *InstancingScene) addBoxes() {
	this.ent.SetNumInstances(uint32((SIZE + addedRows) * (SIZE + addedRows) * (SIZE + addedRows)))
	var wg sync.WaitGroup
	wg.Add((SIZE + addedRows) * (SIZE + addedRows) * (SIZE + addedRows))
	for x := 0; x < SIZE+addedRows; x++ {
		for y := 0; y < SIZE+addedRows; y++ {
			for z := 0; z < SIZE+addedRows; z++ {
				go func(_x, _y, _z int) {
					this.ent.Transforms[_x+_y*(SIZE+addedRows)+_z*(SIZE+addedRows)*(SIZE+addedRows)].Position = mgl32.Vec3{float32(_x) * 2.0, float32(_y) * 2.0, float32(_z) * 2.0}
					wg.Done()
				}(x, y, z)
			}
		}
	}
	wg.Wait()
	this.ent.UpdateInstancedValues()
	addedRows++
	removedRows = 0
}

func (this *InstancingScene) Update(delta_time float32) {
	if !USE_INSTANCING {
		return
	}
	if gohome.InputMgr.JustPressed(gohome.KeyH) {
		this.addBoxes()
		log.Println("Size:", SIZE+addedRows)
	} else if gohome.InputMgr.JustPressed(gohome.KeyJ) {
		if addedRows < 2 {
			return
		}
		addedRows -= 2
		this.addBoxes()
		log.Println("Size:", SIZE+addedRows)
	} else if gohome.InputMgr.JustPressed(gohome.KeyK) {
		this.ent.SetNumUsedInstances(this.ent.Model3D.GetNumUsedInstances() - 1)
	}
}

func (this *InstancingScene) Terminate() {
	this.ent.Terminate()
}
