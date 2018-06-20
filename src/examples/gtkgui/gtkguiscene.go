package main

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK"
	"github.com/PucklaMotzer09/gohomeengine/src/frameworks/GTK/gtk"
)

var gtkf *framework.GTKFramework

type GTKGUIScene struct {
	cube gohome.Entity3D
}

func (this *GTKGUIScene) InitGUI() {
	var box gtk.Box
	var button gtk.Button
	box = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL,0)
	button = gtk.ButtonNewWithLabel("Click Me")
	gtk.GetWindow().ToContainer().Add(box.ToWidget())

	box.ToContainer().Add(gtk.GetGLArea().ToWidget())
	box.ToContainer().Add(button.ToWidget())

	gtk.GetGLArea().ToWidget().SetSizeRequest(640/2,480/2)
}

func (this *GTKGUIScene) Init() {
	gtkf = gohome.Framew.(*framework.GTKFramework)
	if !gtkf.UseWholeWindowAsGLArea {
		gohome.MainLop.InitWindow()
		this.InitGUI()
		gohome.MainLop.InitRenderer()
		gohome.MainLop.InitManagers()
		gohome.Render.AfterInit()
		gohome.RenderMgr.RenderToScreenFirst = true
	}

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
	this.cube.Transform.Rotation[0] += 30.0 * delta_time
	this.cube.Transform.Rotation[1] += 30.0 * delta_time
}

func (this *GTKGUIScene) Terminate() {

}
