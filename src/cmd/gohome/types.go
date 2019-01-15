package main

type Build interface {
	Build() bool
	Install() bool
	Generate()
	IsGenerated() bool
	Run() bool
	Export()
	Clean()
	Env()
}

type DesktopBuild struct {
	title          string
	width          int
	height         int
	gtkwholewindow bool
	gtkmenubar     bool
}

type AndroidBuild struct {
}

type JSBuild struct {
}
