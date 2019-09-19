package main

import (
	"log"

	framework "github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/GTK"
	"github.com/PucklaMotzer09/GoHomeEngine/src/frameworks/GTK/gtk"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

var gtkf *framework.GTKFramework

type GTKGUIScene struct {
	cube gohome.Entity3D
	lb   gtk.ListBox
}

func (this *GTKGUIScene) InitGUI() {
	gohome.MainLop.InitWindow()

	var box gtk.Box
	var button gtk.Button
	var button2 gtk.Button
	lbl := gtk.LabelNew("I am a label")
	this.lb = gtk.ListBoxNew()
	box = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	button = gtk.ButtonNewWithLabel("Enter Me")
	button2 = gtk.ButtonNewWithLabel("Click Me")
	button.SignalConnect("enter", func(button gtk.Button) {
		log.Println("Entered Button")
	})
	this.lb.ToWidget().SignalConnect("button-press-event", func(widget gtk.Widget, event gtk.Event) {
		log.Println("Button Pressed Event")
	})
	button2.SignalConnect("clicked", func(button gtk.Button) {
		log.Println("Clicked Button2")
	})
	this.lb.SignalConnect("row-selected", func(listBox gtk.ListBox, listBoxRow gtk.ListBoxRow) {
		lbr := listBoxRow
		if !lbr.IsNULL() {
			cont := lbr.ToContainer()
			lbl := cont.GetChildren().Data().ToWidget().ToLabel()
			log.Println("Selected row:", lbl.GetText())
		}
	})
	gtk.GetGLArea().ToWidget().SignalConnect("size-allocate", func(widget gtk.Widget) {
		w, h := widget.GetSize()
		gohome.Render.SetNativeResolution(int(w), int(h))
	})
	gtk.GetWindow().ToContainer().Add(box.ToWidget())

	box.ToContainer().Add(gtk.GetGLArea().ToWidget())
	box.ToContainer().Add(button.ToWidget())
	box.ToContainer().Add(button2.ToWidget())
	box.ToContainer().Add(this.lb.ToWidget())
	this.lb.ToContainer().Add(lbl.ToWidget())
	this.lb.ToContainer().Add(gtk.LabelNew("You are a programmer").ToWidget())
	this.lb.ToContainer().Add(gtk.LabelNew("This is GTK").ToWidget())
	gtk.GetGLArea().ToWidget().SetSizeRequest(640/2, 480/2)
	gtk.GetGLArea().ToWidget().SetCanFocus(true)
	gtk.GetWindow().ToWidget().ShowAll()

	gohome.Framew.(*framework.GTKFramework).AfterWindowCreation(&gohome.MainLop)
}

func (this *GTKGUIScene) Init() {
	this.InitGUI()
	gohome.LightMgr.DisableLighting()
	gohome.ResourceMgr.LoadTexture("CubeImage", "cube.png")

	mesh := gohome.Box("Cube", [3]float32{1.0, 1.0, 1.0}, true)
	mesh.GetMaterial().SetTextures("CubeImage", "", "")
	this.cube.InitMesh(mesh)
	this.cube.Transform.Position = [3]float32{0.0, 0.0, -3.0}

	gohome.RenderMgr.AddObject(&this.cube)
	gohome.RenderMgr.UpdateProjectionWithViewport = true
}

func (this *GTKGUIScene) Update(delta_time float32) {
	this.cube.Transform.Rotation = this.cube.Transform.Rotation.Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{0.0, 1.0, 0.0})).Mul(mgl32.QuatRotate(mgl32.DegToRad(30.0)*delta_time, mgl32.Vec3{1.0, 0.0, 0.0}))
	if gohome.InputMgr.IsPressed(gohome.KeySpace) {
		log.Println("Space is pressed")
	}
}

func (this *GTKGUIScene) Terminate() {

}
