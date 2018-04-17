package framework

import (
	// "fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/assimp"
	"log"
	"math"
	"strings"
	"sync"
)

const (
	NUM_GO_ROUTINES_MESH_VERTICES_LOADING uint32 = 10
)

type GLFWFramework struct {
	window           *glfw.Window
	prevMousePos     [2]int16
	prevWindowWidth  int
	prevWindowHeight int
	prevWindowX      int
	prevWindowY      int
}

func (gfw *GLFWFramework) Init(ml *gohome.MainLoop) error {
	gfw.window = nil
	if err := glfw.Init(); err != nil {
		return err
	}
	ml.DoStuff()

	return nil
}
func (GLFWFramework) Update() {
	gohome.InputMgr.Mouse.Wheel[0] = 0
	gohome.InputMgr.Mouse.Wheel[1] = 0
	gohome.InputMgr.Mouse.DPos[0] = 0
	gohome.InputMgr.Mouse.DPos[1] = 0
}

func (gfw *GLFWFramework) Terminate() {
	defer glfw.Terminate()
	defer gfw.window.Destroy()
}

func (gfw *GLFWFramework) CreateWindow(windowWidth, windowHeight uint32, title string) error {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 8)
	var err error
	gfw.window, err = glfw.CreateWindow(int(windowWidth), int(windowHeight), title, nil, nil)
	if err != nil {
		return err
	}
	gfw.window.MakeContextCurrent()
	gfw.window.SetKeyCallback(onKey)
	gfw.window.SetCursorPosCallback(onMousePosChanged)
	gfw.window.SetScrollCallback(onMouseWheelChanged)
	gfw.window.SetMouseButtonCallback(onMouseButton)
	gfw.window.SetFramebufferSizeCallback(onResize)

	glfw.SwapInterval(1)
	return nil
}

func (GLFWFramework) PollEvents() {
	glfw.PollEvents()
}

func (gfw *GLFWFramework) WindowClosed() bool {
	return gfw.window.ShouldClose()
}

func (gfw *GLFWFramework) WindowSwap() {
	gfw.window.SwapBuffers()
}

func (gfw *GLFWFramework) WindowGetSize() mgl32.Vec2 {
	var size mgl32.Vec2
	x, y := gfw.window.GetSize()
	size[0] = float32(x)
	size[1] = float32(y)

	return size
}

func glfwKeysTogohomeKeys(key glfw.Key) gohome.Key {
	switch key {
	case glfw.KeyUnknown:
		return gohome.KeyUnknown
	case glfw.KeySpace:
		return gohome.KeySpace
	case glfw.KeyApostrophe:
		return gohome.KeyApostrophe
	case glfw.KeyComma:
		return gohome.KeyComma
	case glfw.KeyMinus:
		return gohome.KeyMinus
	case glfw.KeyPeriod:
		return gohome.KeyPeriod
	case glfw.KeySlash:
		return gohome.KeySlash
	case glfw.Key0:
		return gohome.Key0
	case glfw.Key1:
		return gohome.Key1
	case glfw.Key2:
		return gohome.Key2
	case glfw.Key3:
		return gohome.Key3
	case glfw.Key4:
		return gohome.Key4
	case glfw.Key5:
		return gohome.Key5
	case glfw.Key6:
		return gohome.Key6
	case glfw.Key7:
		return gohome.Key7
	case glfw.Key8:
		return gohome.Key8
	case glfw.Key9:
		return gohome.Key9
	case glfw.KeySemicolon:
		return gohome.KeySemicolon
	case glfw.KeyEqual:
		return gohome.KeyEqual
	case glfw.KeyA:
		return gohome.KeyA
	case glfw.KeyB:
		return gohome.KeyB
	case glfw.KeyC:
		return gohome.KeyC
	case glfw.KeyD:
		return gohome.KeyD
	case glfw.KeyE:
		return gohome.KeyE
	case glfw.KeyF:
		return gohome.KeyF
	case glfw.KeyG:
		return gohome.KeyG
	case glfw.KeyH:
		return gohome.KeyH
	case glfw.KeyI:
		return gohome.KeyI
	case glfw.KeyJ:
		return gohome.KeyJ
	case glfw.KeyK:
		return gohome.KeyK
	case glfw.KeyL:
		return gohome.KeyL
	case glfw.KeyM:
		return gohome.KeyM
	case glfw.KeyN:
		return gohome.KeyN
	case glfw.KeyO:
		return gohome.KeyO
	case glfw.KeyP:
		return gohome.KeyP
	case glfw.KeyQ:
		return gohome.KeyQ
	case glfw.KeyR:
		return gohome.KeyR
	case glfw.KeyS:
		return gohome.KeyS
	case glfw.KeyT:
		return gohome.KeyT
	case glfw.KeyU:
		return gohome.KeyU
	case glfw.KeyV:
		return gohome.KeyV
	case glfw.KeyW:
		return gohome.KeyW
	case glfw.KeyX:
		return gohome.KeyX
	case glfw.KeyY:
		return gohome.KeyY
	case glfw.KeyZ:
		return gohome.KeyZ
	case glfw.KeyLeftBracket:
		return gohome.KeyLeftBracket
	case glfw.KeyBackslash:
		return gohome.KeyBackslash
	case glfw.KeyRightBracket:
		return gohome.KeyRightBracket
	case glfw.KeyGraveAccent:
		return gohome.KeyGraveAccent
	case glfw.KeyWorld1:
		return gohome.KeyWorld1
	case glfw.KeyWorld2:
		return gohome.KeyWorld2
	case glfw.KeyEscape:
		return gohome.KeyEscape
	case glfw.KeyEnter:
		return gohome.KeyEnter
	case glfw.KeyTab:
		return gohome.KeyTab
	case glfw.KeyBackspace:
		return gohome.KeyBackspace
	case glfw.KeyInsert:
		return gohome.KeyInsert
	case glfw.KeyDelete:
		return gohome.KeyDelete
	case glfw.KeyRight:
		return gohome.KeyRight
	case glfw.KeyLeft:
		return gohome.KeyLeft
	case glfw.KeyDown:
		return gohome.KeyDown
	case glfw.KeyUp:
		return gohome.KeyUp
	case glfw.KeyPageUp:
		return gohome.KeyPageUp
	case glfw.KeyPageDown:
		return gohome.KeyPageDown
	case glfw.KeyHome:
		return gohome.KeyHome
	case glfw.KeyEnd:
		return gohome.KeyEnd
	case glfw.KeyCapsLock:
		return gohome.KeyCapsLock
	case glfw.KeyScrollLock:
		return gohome.KeyScrollLock
	case glfw.KeyNumLock:
		return gohome.KeyNumLock
	case glfw.KeyPrintScreen:
		return gohome.KeyPrintScreen
	case glfw.KeyPause:
		return gohome.KeyPause
	case glfw.KeyF1:
		return gohome.KeyF1
	case glfw.KeyF2:
		return gohome.KeyF2
	case glfw.KeyF3:
		return gohome.KeyF3
	case glfw.KeyF4:
		return gohome.KeyF4
	case glfw.KeyF5:
		return gohome.KeyF5
	case glfw.KeyF6:
		return gohome.KeyF6
	case glfw.KeyF7:
		return gohome.KeyF7
	case glfw.KeyF8:
		return gohome.KeyF8
	case glfw.KeyF9:
		return gohome.KeyF9
	case glfw.KeyF10:
		return gohome.KeyF10
	case glfw.KeyF11:
		return gohome.KeyF11
	case glfw.KeyF12:
		return gohome.KeyF12
	case glfw.KeyF13:
		return gohome.KeyF13
	case glfw.KeyF14:
		return gohome.KeyF14
	case glfw.KeyF15:
		return gohome.KeyF15
	case glfw.KeyF16:
		return gohome.KeyF16
	case glfw.KeyF17:
		return gohome.KeyF17
	case glfw.KeyF18:
		return gohome.KeyF18
	case glfw.KeyF19:
		return gohome.KeyF19
	case glfw.KeyF20:
		return gohome.KeyF20
	case glfw.KeyF21:
		return gohome.KeyF21
	case glfw.KeyF22:
		return gohome.KeyF22
	case glfw.KeyF23:
		return gohome.KeyF23
	case glfw.KeyF24:
		return gohome.KeyF24
	case glfw.KeyF25:
		return gohome.KeyF25
	case glfw.KeyKP0:
		return gohome.KeyKP0
	case glfw.KeyKP1:
		return gohome.KeyKP1
	case glfw.KeyKP2:
		return gohome.KeyKP2
	case glfw.KeyKP3:
		return gohome.KeyKP3
	case glfw.KeyKP4:
		return gohome.KeyKP4
	case glfw.KeyKP5:
		return gohome.KeyKP5
	case glfw.KeyKP6:
		return gohome.KeyKP6
	case glfw.KeyKP7:
		return gohome.KeyKP7
	case glfw.KeyKP8:
		return gohome.KeyKP8
	case glfw.KeyKP9:
		return gohome.KeyKP9
	case glfw.KeyKPDecimal:
		return gohome.KeyKPDecimal
	case glfw.KeyKPDivide:
		return gohome.KeyKPDivide
	case glfw.KeyKPMultiply:
		return gohome.KeyKPMultiply
	case glfw.KeyKPSubtract:
		return gohome.KeyKPSubtract
	case glfw.KeyKPAdd:
		return gohome.KeyKPAdd
	case glfw.KeyKPEnter:
		return gohome.KeyKPEnter
	case glfw.KeyKPEqual:
		return gohome.KeyKPEqual
	case glfw.KeyLeftShift:
		return gohome.KeyLeftShift
	case glfw.KeyLeftControl:
		return gohome.KeyLeftControl
	case glfw.KeyLeftAlt:
		return gohome.KeyLeftAlt
	case glfw.KeyLeftSuper:
		return gohome.KeyLeftSuper
	case glfw.KeyRightShift:
		return gohome.KeyRightShift
	case glfw.KeyRightControl:
		return gohome.KeyRightControl
	case glfw.KeyRightAlt:
		return gohome.KeyRightAlt
	case glfw.KeyRightSuper:
		return gohome.KeyRightSuper
	case glfw.KeyMenu:
		return gohome.KeyMenu
	}

	return gohome.KeyUnknown
}

func glfwMouseButtonTogohomeKeys(mb glfw.MouseButton) gohome.Key {
	switch mb {
	// case glfw.MouseButton1:
	// 	return MouseButton1
	// case glfw.MouseButton2:
	// 	return MouseButton2
	// case glfw.MouseButton3:
	// 	return MouseButton3
	// case glfw.MouseButton4:
	// 	return MouseButton4
	// case glfw.MouseButton5:
	// 	return MouseButton5
	// case glfw.MouseButton6:
	// 	return MouseButton6
	// case glfw.MouseButton7:
	// 	return MouseButton7
	// case glfw.MouseButton8:
	// 	return MouseButton8
	case glfw.MouseButtonLast:
		return gohome.MouseButtonLast
	case glfw.MouseButtonLeft:
		return gohome.MouseButtonLeft
	case glfw.MouseButtonRight:
		return gohome.MouseButtonRight
	case glfw.MouseButtonMiddle:
		return gohome.MouseButtonMiddle
	}

	return gohome.MouseButtonLast
}

func (gfw *GLFWFramework) CurserShow() {
	gfw.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
}
func (gfw *GLFWFramework) CursorHide() {
	gfw.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
}
func (gfw *GLFWFramework) CursorDisable() {
	gfw.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
}
func (gfw *GLFWFramework) CursorShown() bool {
	return gfw.window.GetInputMode(glfw.CursorMode) == glfw.CursorNormal
}
func (gfw *GLFWFramework) CursorHidden() bool {
	return gfw.window.GetInputMode(glfw.CursorMode) == glfw.CursorHidden
}
func (gfw *GLFWFramework) CursorDisabled() bool {
	return gfw.window.GetInputMode(glfw.CursorMode) == glfw.CursorDisabled
}
func (gfw *GLFWFramework) WindowSetFullscreen(b bool) {
	var monitor *glfw.Monitor
	var refreshRate int
	var x, y int
	var width, height int
	if b {
		monitor = getFocusedMonitor(gfw.window)
		if monitor == nil {
			monitor = glfw.GetPrimaryMonitor()
		}
		refreshRate = monitor.GetVideoMode().RefreshRate
		x = 0
		y = 0
		gfw.prevWindowWidth, gfw.prevWindowHeight = gfw.window.GetSize()
		gfw.prevWindowX, gfw.prevWindowY = gfw.window.GetPos()
	} else {
		monitor = nil
		temp := getFocusedMonitor(gfw.window)
		if temp == nil {
			temp = glfw.GetPrimaryMonitor()
		}
		refreshRate = temp.GetVideoMode().RefreshRate
		x, y = gfw.prevWindowX, gfw.prevWindowY
	}
	width, height = gfw.prevWindowWidth, gfw.prevWindowHeight

	gfw.window.SetMonitor(monitor, x, y, width, height, refreshRate)
}

func getFocusedMonitor(window *glfw.Window) *glfw.Monitor {
	wx, wy := window.GetPos()
	ww, wh := window.GetSize()

	var mx, my, mw, mh int
	var monitor *glfw.Monitor
	var max, maxIndex int = -1, -1
	var axmin, axmax, aymin, aymax int
	var bxmin, bxmax, bymin, bymax int

	for i := 0; i < len(glfw.GetMonitors()); i++ {
		monitor = glfw.GetMonitors()[i]
		mx, my = monitor.GetPos()
		mw, mh = monitor.GetVideoMode().Width, monitor.GetVideoMode().Height

		if wx+ww > mx && wx < mx+mw &&
			wy+wh > my && wy < my+mh {

			axmin, axmax = wx, wx+ww
			aymin, aymax = wy, wy+wh
			bxmin, bxmax = mx, mx+mw
			bymin, bymax = my, my+mh

			dx := int(math.Min(float64(axmax), float64(bxmax)) - math.Max(float64(axmin), float64(bxmin)))
			dy := int(math.Min(float64(aymax), float64(bymax)) - math.Max(float64(aymin), float64(bymin)))

			mean := dx * dy
			if mean > max {
				max = mean
				maxIndex = i
			}
		}
	}

	if maxIndex == -1 {
		return nil
	} else {
		return glfw.GetMonitors()[maxIndex]
	}

}

func (gfw *GLFWFramework) WindowIsFullscreen() bool {
	return gfw.window.GetMonitor() != nil
}

func (gfw *GLFWFramework) LoadLevel(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	if _, ok := rsmgr.Levels[name]; ok {
		log.Println("The level with the name", name, "has already been loaded!")
		return nil
	}
	level := &gohome.Level{Name: name}
	var scene *assimp.Scene
	if scene = assimp.ImportFile(path, uint(assimp.Process_Triangulate|assimp.Process_FlipUVs|assimp.Process_GenNormals|assimp.Process_OptimizeMeshes)); scene == nil || (scene.Flags()&assimp.SceneFlags_Incomplete) != 0 || scene.RootNode() == nil {
		log.Println("Couldn't load level", name, "with path", path, ":", assimp.GetErrorString())
		return nil
	}

	directory := path
	if index := strings.LastIndex(directory, "/"); index != -1 {
		directory = directory[index:]
	} else {
		directory = ""
	}

	gfw.ProcessNode(rsmgr, scene.RootNode(), scene, level, directory, preloaded, loadToGPU)

	return level
}

func (gfw *GLFWFramework) ProcessNode(rsmgr *gohome.ResourceManager, node *assimp.Node, scene *assimp.Scene, level *gohome.Level, directory string, preloaded, loadToGPU bool) {
	if node != scene.RootNode() {
		model := &gohome.Model3D{}
		gfw.initModel(model, node, scene, level, directory, preloaded, loadToGPU)
		if !preloaded {
			if _, ok := rsmgr.Models[model.Name]; ok {
				log.Println("Model", model.Name, "has already been loaded! Overwritting ...")
			}
			rsmgr.Models[model.Name] = model
			log.Println("Finished loading Model", model.Name, "!")
		} else {
			rsmgr.PreloadedModelsChan <- model
		}

	}
	for i := 0; i < node.NumChildren(); i++ {
		gfw.ProcessNode(rsmgr, node.Children()[i], scene, level, directory, preloaded, loadToGPU)
	}
}

func (gfw *GLFWFramework) initModel(model *gohome.Model3D, node *assimp.Node, scene *assimp.Scene, level *gohome.Level, directory string, preloaded, loadToGPU bool) {
	level.LevelObjects = append(level.LevelObjects, gohome.LevelObject{
		Name: node.Name(),
	})
	gfw.setTransformLevelObject(&level.LevelObjects[len(level.LevelObjects)-1], node.Transformation())

	model.Name = node.Name()
	for i := 0; i < node.NumMeshes(); i++ {
		aiMesh := scene.Meshes()[node.Meshes()[i]]
		mesh := gohome.Render.CreateMesh3D(aiMesh.Name())
		gfw.addVerticesAssimpMesh3D(mesh, aiMesh, node, scene, level, directory, preloaded)
		model.AddMesh3D(mesh)
		if !preloaded {
			if loadToGPU {
				mesh.Load()
			}
			log.Println("Finished loading mesh", mesh.GetName(), "V:", mesh.GetNumVertices(), "I:", mesh.GetNumIndices(), "!")
		} else {
			mesh.CalculateTangents()
			gohome.ResourceMgr.PreloadedMeshesChan <- gohome.PreloadedMesh{
				mesh,
				loadToGPU,
			}
		}
	}
}

func loadVertices(vertices *[]gohome.Mesh3DVertex, mesh *assimp.Mesh, start_index, end_index, max_index uint32, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start_index; i < uint32(mesh.NumVertices()) && i < end_index; i++ {
		var texCoords mgl32.Vec2
		if len(mesh.TextureCoords(0)) > 0 {
			texCoords[0] = mesh.TextureCoords(0)[i].X()
			texCoords[1] = mesh.TextureCoords(0)[i].Y()
		} else {
			texCoords[0] = 0.0
			texCoords[1] = 0.0
		}
		vertex := gohome.Mesh3DVertex{
			/* X,Y,Z,
			   NX,NY,NZ,
			   U,V,
			   TX,TY,TZ,
			*/
			mesh.Vertices()[i].X(), mesh.Vertices()[i].Y(), mesh.Vertices()[i].Z(),
			mesh.Normals()[i].X(), mesh.Normals()[i].Y(), mesh.Vertices()[i].Z(),
			texCoords[0], texCoords[1],
			0.0, 0.0, 0.0,
		}
		(*vertices)[i] = vertex
	}
}

func (gfw *GLFWFramework) addVerticesAssimpMesh3D(oglm gohome.Mesh3D, mesh *assimp.Mesh, node *assimp.Node, scene *assimp.Scene, level *gohome.Level, directory string, preloaded bool) {
	vertices := make([]gohome.Mesh3DVertex, mesh.NumVertices())
	var indices []uint32
	var wg sync.WaitGroup
	var i uint32
	deltaIndex := uint32(mesh.NumVertices()) / NUM_GO_ROUTINES_MESH_VERTICES_LOADING
	wg.Add(int(NUM_GO_ROUTINES_MESH_VERTICES_LOADING))
	for i = 0; i < NUM_GO_ROUTINES_MESH_VERTICES_LOADING; i++ {
		go loadVertices(&vertices, mesh, i*deltaIndex, (i+1)*deltaIndex, uint32(mesh.NumVertices()), &wg)
	}

	for i = 0; i < uint32(mesh.NumFaces()); i++ {
		face := mesh.Faces()[i]
		faceIndices := face.CopyIndices()
		indices = append(indices, faceIndices...)
	}

	wg.Wait()

	oglm.AddVertices(vertices, indices)

	mat := &gohome.Material{}
	gfw.initMaterial(mat, scene.Materials()[mesh.MaterialIndex()], scene, directory, preloaded)
	oglm.SetMaterial(mat)
}

func (gfw *GLFWFramework) initMaterial(mat *gohome.Material, material *assimp.Material, scene *assimp.Scene, directory string, preloaded bool) {
	var ret assimp.Return
	var matDifColor assimp.Color4
	var matSpecColor assimp.Color4
	var matShininess float32

	matDifColor, ret = material.GetMaterialColor(assimp.MatKey_ColorDiffuse, 0, 0)
	if ret == assimp.Return_Failure {
		mat.DiffuseColor = &gohome.Color{255, 255, 255, 255}
	} else {
		mat.DiffuseColor = convertAssimpColor(matDifColor)
	}
	matSpecColor, ret = material.GetMaterialColor(assimp.MatKey_ColorSpecular, 0, 0)
	if ret == assimp.Return_Failure {
		mat.SpecularColor = &gohome.Color{255, 255, 255, 255}
	} else {
		mat.SpecularColor = convertAssimpColor(matSpecColor)
	}
	matShininess, ret = material.GetMaterialFloat(assimp.MatKey_Shininess, 0, 0)
	if ret == assimp.Return_Failure {
		mat.Shinyness = 0.0
	} else {
		mat.Shinyness = matShininess
	}

	diffuseTextures := material.GetMaterialTextureCount(1)
	specularTextures := material.GetMaterialTextureCount(2)
	normalMaps := material.GetMaterialTextureCount(6)
	for i := 0; i < diffuseTextures; i++ {
		texPath, _, _, _, _, _, _, ret := material.GetMaterialTexture(1, i)
		if ret == assimp.Return_Failure {
			log.Println("Couldn't return diffuse Texture")
		} else {
			if !preloaded {
				gohome.ResourceMgr.LoadTexture(texPath, directory+texPath)
				mat.DiffuseTexture = gohome.ResourceMgr.GetTexture(texPath)
			} else {
				mat.DiffuseTexture = gohome.ResourceMgr.LoadTextureFunction(texPath, directory+texPath, true)
			}

			break
		}
	}
	for i := 0; i < specularTextures; i++ {
		texPath, _, _, _, _, _, _, ret := material.GetMaterialTexture(2, i)
		if ret == assimp.Return_Failure {
			log.Println("Couldn't return specular Texture")
		} else {
			if !preloaded {
				gohome.ResourceMgr.LoadTexture(texPath, directory+texPath)
				mat.SpecularTexture = gohome.ResourceMgr.GetTexture(texPath)
			} else {
				mat.SpecularTexture = gohome.ResourceMgr.LoadTextureFunction(texPath, directory+texPath, true)
			}
			break
		}
	}
	for i := 0; i < normalMaps; i++ {
		texPath, _, _, _, _, _, _, ret := material.GetMaterialTexture(6, i)
		if ret == assimp.Return_Failure {
			log.Println("Couldn't return normal map")
		} else {
			if !preloaded {
				gohome.ResourceMgr.LoadTexture(texPath, directory+texPath)
				mat.NormalMap = gohome.ResourceMgr.GetTexture(texPath)
			} else {
				mat.NormalMap = gohome.ResourceMgr.LoadTextureFunction(texPath, directory+texPath, true)
			}
			break
		}
	}
}

func convertAssimpColor(color assimp.Color4) *gohome.Color {
	return &gohome.Color{uint8(color.R() * 255.0), uint8(color.G() * 255.0), uint8(color.B() * 255.0), uint8(color.A() * 255.0)}
}

func (gfw *GLFWFramework) setTransformLevelObject(this *gohome.LevelObject, mat assimp.Matrix4x4) {
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			this.Transform.TransformMatrix[this.Transform.TransformMatrix.Index(r, c)] = mat.Values()[c][r]
		}
	}
}
