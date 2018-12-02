package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK"
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK/gtk"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"log"
	"strconv"
)

type GTKBuilderScene struct {
	cube gohome.Entity3D
	lb   gtk.Label
}

func (this *GTKBuilderScene) Init() {
	builder := gtk.BuilderNew()
	if err := builder.AddFromFile("editscene.ui"); err != nil {
		log.Println("Error:", err.Error())
		return
	}

	window := builder.GetObject("window").ToWidget().ToWindow()
	glarea := builder.GetObject("glarea").ToGLArea()
	lb_assets := builder.GetObject("lb_assets").ToWidget().ToListBox()
	gohome.Framew.(*framework.GTKFramework).InitExternalDefault(&window, &glarea)
	ws, hs := glarea.ToWidget().GetSize()
	gohome.Render.SetNativeResolution(uint32(ws), uint32(hs))

	glarea.ToWidget().SignalConnect("size-allocate", func(widget gtk.Widget) {
		w, h := widget.GetSize()
		gohome.Render.SetNativeResolution(uint32(w), uint32(h))
	})
	this.lb = gtk.LabelNew("Hello Programmer")
	lb_assets.ToContainer().Add(this.lb.ToWidget())

	gohome.Init3DShaders()
	gohome.ResourceMgr.LoadTexture("CubeImage", "cube.png")

	mesh := gohome.Box("Cube", [3]float32{1.0, 1.0, 1.0})
	mesh.GetMaterial().SetTextures("CubeImage", "", "")
	this.cube.InitMesh(mesh)
	this.cube.Transform.Position = [3]float32{0.0, 0.0, -3.0}

	gohome.RenderMgr.AddObject(&this.cube)
	gohome.LightMgr.DisableLighting()

	gohome.RenderMgr.UpdateProjectionWithViewport = true
	gohome.RenderMgr.EnableBackBuffer = true
}

func (this *GTKBuilderScene) Update(delta_time float32) {
	this.cube.Transform.Rotation = this.cube.Transform.Rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{0.0, 1.0, 0.0})).Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{1.0, 0.0, 0.0}))
	this.lb.SetText("FPS: " + strconv.FormatFloat(float64(1.0/delta_time), 'f', 2, 32))
}

func (this *GTKBuilderScene) Terminate() {

}
