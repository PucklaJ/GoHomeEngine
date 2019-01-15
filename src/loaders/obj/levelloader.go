package loader

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"io/ioutil"
	"strconv"
	"strings"
)

func LoadLevelOBJ(name, path string, loadToGPU bool) *gohome.Level {
	alreadyLoaded, quit, _ := gohome.ResourceMgr.CheckLevel(name, path)
	if quit {
		return nil
	}

	if alreadyLoaded {
		name = getNameForAlreadyLoadedLevel(name)
	}
	var objLoader OBJLoader
	objLoader.DisableGoRoutines = true
	if err := objLoader.Load(path); err != nil {
		gohome.ErrorMgr.Error("Level", name, "Couldn't load "+path+": "+err.Error())
		return nil
	}
	lvl := toGohomeLevel(name, &objLoader, loadToGPU)
	lvl.Name = name
	return lvl
}

func LoadLevelOBJString(name, contents, fileName string, loadToGPU bool) *gohome.Level {
	alreadyLoaded, quit, _ := gohome.ResourceMgr.CheckLevel(name, fileName)
	if quit {
		return nil
	}
	if alreadyLoaded {
		name = getNameForAlreadyLoadedLevel(name)
	}

	reader := ioutil.NopCloser(strings.NewReader(contents))

	var objLoader OBJLoader
	objLoader.SetDirectory(gohome.GetPathFromFile(fileName))
	objLoader.DisableGoRoutines = true
	if err := objLoader.LoadReader(reader); err != nil {
		gohome.ErrorMgr.Error("Level", name, err.Error())
		return nil
	}
	lvl := toGohomeLevel(name, &objLoader, loadToGPU)
	lvl.Name = name
	return lvl
}

func processModel(level *gohome.Level, objLoader *OBJLoader, model *OBJModel, loadToGPU bool) {
	alreadyLoaded, quit := gohome.ResourceMgr.CheckModel(model.Name)
	if quit {
		return
	}

	level.LevelObjects = append(level.LevelObjects, gohome.LevelObject{})
	lvlObj := &level.LevelObjects[len(level.LevelObjects)-1]
	lvlObj.Name = model.Name

	var model3d gohome.Model3D
	if !alreadyLoaded {
		model3d.Name = model.Name
	} else {
		model3d.Name = getNameForAlreadyLoadedModel(model.Name)
	}
	for i := 0; i < len(model.Meshes); i++ {
		mesh3d := gohome.Render.CreateMesh3D(model.Meshes[i].Name)
		processMesh(objLoader, mesh3d, &model.Meshes[i], loadToGPU)
		model3d.AddMesh3D(mesh3d)
	}
	lvlObj.Model3D = &model3d
	gohome.ResourceMgr.Models[model3d.Name] = &model3d
	gohome.ErrorMgr.Log("Model", model.Name, "Finished loading!")
}

func toGohomeColor(color [3]float32) *gohome.Color {
	var rv gohome.Color
	rv.R = uint8(color[0] * 255.0)
	rv.G = uint8(color[1] * 255.0)
	rv.B = uint8(color[2] * 255.0)
	rv.A = 255
	return &rv
}

func loadMaterialTexture(directory string, path string) (rv gohome.Texture) {
	if rv = gohome.ResourceMgr.LoadTexture(path, directory+path); rv == nil {
		rv = gohome.ResourceMgr.LoadTexture(path, directory+gohome.GetFileFromPath(path))
	}
	if rv != nil {
		rv.SetWrapping(gohome.WRAPPING_REPEAT)
	}

	return
}

func processMaterial(objLoader *OBJLoader, material *gohome.Material, mat *OBJMaterial, loadToGPU bool) {
	if mat == nil {
		return
	}
	material.Name = mat.Name
	material.DiffuseColor = mat.DiffuseColor
	material.SpecularColor = mat.SpecularColor
	material.SetShinyness(mat.SpecularExponent)
	if mat.DiffuseTexture != "" {
		material.DiffuseTexture = loadMaterialTexture(objLoader.directory, mat.DiffuseTexture)
	}
	if mat.SpecularTexture != "" {
		material.SpecularTexture = loadMaterialTexture(objLoader.directory, mat.SpecularTexture)
	}
	if mat.NormalMap != "" {
		material.NormalMap = loadMaterialTexture(objLoader.directory, mat.NormalMap)
	}
	material.Transparency = mat.Transperancy

}

func processMesh(objLoader *OBJLoader, mesh3d gohome.Mesh3D, mesh *OBJMesh, loadToGPU bool) {
	mesh3d.AddVertices(mesh.Vertices, mesh.Indices)
	var mat gohome.Material
	mat.InitDefault()
	processMaterial(objLoader, &mat, mesh.Material, loadToGPU)
	mesh3d.SetMaterial(&mat)

	if loadToGPU {
		mesh3d.Load()
	}
	gohome.ErrorMgr.Log("Mesh", mesh3d.GetName(), "Finished loading! V: "+strconv.Itoa(int(mesh3d.GetNumVertices()))+" I: "+strconv.Itoa(int(mesh3d.GetNumIndices())))
}

func toGohomeLevel(name string, objLoader *OBJLoader, loadToGPU bool) *gohome.Level {
	level := &gohome.Level{Name: name}
	for i := 0; i < len(objLoader.Models); i++ {
		processModel(level, objLoader, &objLoader.Models[i], loadToGPU)
	}
	return level
}

func getNameForAlreadyLoadedLevel(name string) string {
	var alreadyLoaded = true
	var count = 1
	var newName string
	for alreadyLoaded {
		newName = name + strconv.FormatInt(int64(count), 10)
		_, alreadyLoaded = gohome.ResourceMgr.Levels[newName]
		count++
	}
	return newName
}

func getNameForAlreadyLoadedModel(name string) string {
	var alreadyLoaded = true
	var count = 1
	var newName string
	for alreadyLoaded {
		newName = name + strconv.FormatInt(int64(count), 10)
		_, alreadyLoaded = gohome.ResourceMgr.Models[newName]
		count++
	}
	return newName
}
