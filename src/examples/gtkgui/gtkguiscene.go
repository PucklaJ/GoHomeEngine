package main

import (
	"log"

	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK"
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK/gtk"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

var gtkf *framework.GTKFramework

type GTKGUIScene struct {
	cube gohome.Entity3D
}

func (this *GTKGUIScene) InitGUI() {
	gohome.MainLop.InitWindow()

	var box gtk.Box
	var button gtk.Button
	var button2 gtk.Button
	box = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	button = gtk.ButtonNewWithLabel("Enter Me")
	button2 = gtk.ButtonNewWithLabel("Click Me")
	button.SignalConnect("enter", func(button gtk.Button) {
		log.Println("Entered Button")
	})
	button2.SignalConnect("clicked", func(button gtk.Button) {
		log.Println("Clicked Button2")
	})
	gtk.GetWindow().ToContainer().Add(box.ToWidget())

	box.ToContainer().Add(gtk.GetGLArea().ToWidget())
	box.ToContainer().Add(button.ToWidget())
	box.ToContainer().Add(button2.ToWidget())
	gtk.GetGLArea().ToWidget().SetSizeRequest(640/2, 480/2)
	gtk.GetWindow().ToWidget().ShowAll()

	gohome.Framew.(*framework.GTKFramework).AfterWindowCreation(&gohome.MainLop)
	gohome.RenderMgr.EnableBackBuffer = true
}

func (this *GTKGUIScene) Init() {
	this.InitGUI()
	gohome.Init3DShaders()
	gohome.ResourceMgr.LoadTexture("CubeImage", "cube.png")

	mesh := gohome.Box("Cube", [3]float32{1.0, 1.0, 1.0})
	mesh.GetMaterial().SetTextures("CubeImage", "", "")
	this.cube.InitMesh(mesh)
	this.cube.Transform.Position = [3]float32{0.0, 0.0, -3.0}

	gohome.RenderMgr.AddObject(&this.cube)
	gohome.LightMgr.DisableLighting()
}

func (this *GTKGUIScene) Update(delta_time float32) {
	this.cube.Transform.Rotation = this.cube.Transform.Rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{0.0, 1.0, 0.0})).Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{1.0, 0.0, 0.0}))
}

func (this *GTKGUIScene) Terminate() {

}
