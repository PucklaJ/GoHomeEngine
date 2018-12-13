package loader

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"strconv"
)

func loadFile(path string, objLoader *OBJLoader) error {
	reader, err := gohome.Framew.OpenFile(path)
	if err != nil {
		return err
	}
	objLoader.SetDirectory(gohome.GetPathFromFile(path))
	objLoader.SetOpenMaterialFile(gohome.Framew.OpenFile)
	objLoader.SetMaterialPaths(gohome.MATERIAL_PATHS[:])
	err1 := objLoader.LoadReader(reader)
	if err1 != nil {
		return err1
	}
	return nil
}

func loadFileWithPaths(path string, paths []string, objLoader *OBJLoader) error {
	var worked bool = false
	var err error
	for i := 0; i < len(paths); i++ {
		if err = loadFile(paths[i]+path, objLoader); err == nil {
			worked = true
			break
		} else if err = loadFile(paths[i]+gohome.GetFileFromPath(path), objLoader); err == nil {
			worked = true
			break
		}
	}

	if !worked {
		return err
	}

	return nil
}

func getNameForAlreadyLoadedModel(rsmgr *gohome.ResourceManager, name string) string {
	var alreadyLoaded = true
	var count = 1
	var newName string
	for alreadyLoaded {
		newName = name + strconv.FormatInt(int64(count), 10)
		_, alreadyLoaded = rsmgr.Models[newName]
		count++
	}
	return newName
}

func processModel(rsmgr *gohome.ResourceManager, level *gohome.Level, objLoader *OBJLoader, model *OBJModel, preloaded, loadToGPU bool) {
	var alreadyLoaded = false
	if _, alreadyLoaded = rsmgr.Models[model.Name]; !rsmgr.LoadModelsWithSameName && alreadyLoaded {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_LOG, "Model", model.Name, "It has already been loaded!")
		return
	}
	level.LevelObjects = append(level.LevelObjects, gohome.LevelObject{})
	lvlObj := &level.LevelObjects[len(level.LevelObjects)-1]
	var lvlObjTobj gohome.TransformableObject3D
	lvlObjTobj.Position = [3]float32{0.0, 0.0, 0.0}
	lvlObjTobj.Scale = [3]float32{1.0, 1.0, 1.0}
	lvlObjTobj.Rotation.V = [3]float32{0.0, 0.0, -1.0}
	lvlObjTobj.Rotation.W = 0.0
	lvlObjTobj.CalculateTransformMatrix(nil, -1)
	lvlObj.Name = model.Name
	lvlObj.Transform.TransformMatrix = lvlObjTobj.GetTransformMatrix()
	var model3d gohome.Model3D
	if !alreadyLoaded {
		model3d.Name = model.Name
	} else {
		model3d.Name = getNameForAlreadyLoadedModel(rsmgr, model.Name)
	}
	for i := 0; i < len(model.Meshes); i++ {
		mesh3d := gohome.Render.CreateMesh3D(model.Meshes[i].Name)
		processMesh(objLoader, mesh3d, &model.Meshes[i], preloaded, loadToGPU)
		model3d.AddMesh3D(mesh3d)
	}
	lvlObj.Model3D = &model3d
	rsmgr.Models[model3d.Name] = &model3d
	gohome.ErrorMgr.Log("Model", model.Name, "Finished loading!")
}

func toMesh3DVertex(vertex *OBJVertex) gohome.Mesh3DVertex {
	var rv gohome.Mesh3DVertex
	rv[0] = vertex.Position[0]
	rv[1] = vertex.Position[1]
	rv[2] = vertex.Position[2]

	rv[3] = vertex.Normal[0]
	rv[4] = vertex.Normal[1]
	rv[5] = vertex.Normal[2]

	rv[6] = vertex.TextureCoord[0]
	rv[7] = vertex.TextureCoord[1]

	return rv
}

func toGohomeColor(color OBJColor) *gohome.Color {
	var rv gohome.Color
	rv.R = uint8(color[0] * 255.0)
	rv.G = uint8(color[1] * 255.0)
	rv.B = uint8(color[2] * 255.0)
	rv.A = 255
	return &rv
}

func loadMaterialTexture(directory string, path string, preloaded bool) gohome.Texture {
	var rv gohome.Texture
	defer func() {
		if !preloaded && rv != nil {
			rv.SetWrapping(gohome.WRAPPING_REPEAT)
		}
	}()
	if !preloaded {
		gohome.ResourceMgr.LoadTexture(path, directory+path)
		if rv = gohome.ResourceMgr.GetTexture(path); rv == nil {
			gohome.ResourceMgr.LoadTexture(path, directory+gohome.GetFileFromPath(path))
			if rv = gohome.ResourceMgr.GetTexture(path); rv == nil {
				for i := 0; i < len(gohome.TEXTURE_PATHS); i++ {
					gohome.ResourceMgr.LoadTexture(path, gohome.TEXTURE_PATHS[i]+path)
					if rv = gohome.ResourceMgr.GetTexture(path); rv == nil {
						gohome.ResourceMgr.LoadTexture(path, gohome.TEXTURE_PATHS[i]+gohome.GetFileFromPath(path))
						if rv = gohome.ResourceMgr.GetTexture(path); rv != nil {
							return rv
						}
					} else {
						return rv
					}
				}
				for i := 0; i < len(gohome.MATERIAL_PATHS); i++ {
					gohome.ResourceMgr.LoadTexture(path, gohome.MATERIAL_PATHS[i]+path)
					if rv = gohome.ResourceMgr.GetTexture(path); rv == nil {
						gohome.ResourceMgr.LoadTexture(path, gohome.MATERIAL_PATHS[i]+gohome.GetFileFromPath(path))
						if rv = gohome.ResourceMgr.GetTexture(path); rv != nil {
							return rv
						}
					} else {
						return rv
					}
				}
			}
		}
	} else {
		if rv = gohome.ResourceMgr.LoadTextureFunction(path, directory+path, true); rv == nil {
			rv = gohome.ResourceMgr.LoadTextureFunction(path, path, true)
			if rv == nil {
				for i := 0; i < len(gohome.MATERIAL_PATHS); i++ {
					if rv = gohome.ResourceMgr.LoadTextureFunction(path, path, true); rv == nil {
						if rv = gohome.ResourceMgr.LoadTextureFunction(path, gohome.GetFileFromPath(path), true); rv != nil {
							return rv
						}
					} else {
						return rv
					}
				}
			}
		}
	}

	return rv
}

func processMaterial(objLoader *OBJLoader, material *gohome.Material, mat *OBJMaterial, preloaded, loadToGPU bool) {
	if mat == nil {
		return
	}
	material.Name = mat.Name
	material.DiffuseColor = toGohomeColor(mat.DiffuseColor)
	material.SpecularColor = toGohomeColor(mat.SpecularColor)
	material.SetShinyness(mat.SpecularExponent)
	if mat.DiffuseTexture != "" {
		material.DiffuseTexture = loadMaterialTexture(objLoader.directory, mat.DiffuseTexture, preloaded)
	}
	if mat.SpecularTexture != "" {
		material.SpecularTexture = loadMaterialTexture(objLoader.directory, mat.SpecularTexture, preloaded)
	}
	if mat.NormalMap != "" {
		material.NormalMap = loadMaterialTexture(objLoader.directory, mat.NormalMap, preloaded)
	}
	material.Transparency = mat.Transperancy

}

func processMesh(objLoader *OBJLoader, mesh3d gohome.Mesh3D, mesh *OBJMesh, preloaded, loadToGPU bool) {
	var vertices []gohome.Mesh3DVertex
	vertices = make([]gohome.Mesh3DVertex, len(mesh.Vertices))
	for i := 0; i < len(vertices); i++ {
		vertices[i] = toMesh3DVertex(&mesh.Vertices[i])
	}
	var mat gohome.Material
	mat.InitDefault()
	processMaterial(objLoader, &mat, mesh.Material, preloaded, loadToGPU)
	mesh3d.SetMaterial(&mat)
	if len(vertices) != 0 && len(mesh.Indices) != 0 {
		mesh3d.AddVertices(vertices, mesh.Indices)
	}
	if !preloaded {
		if loadToGPU {
			mesh3d.Load()
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_LOG, "Mesh", mesh3d.GetName(), "Finished loading! V: "+strconv.Itoa(int(mesh3d.GetNumVertices()))+" I: "+strconv.Itoa(int(mesh3d.GetNumIndices())))
		}
	} else {
		mesh3d.CalculateTangents()
		gohome.ResourceMgr.PreloadedMeshesChan <- gohome.PreloadedMesh{
			mesh3d,
			loadToGPU,
		}
	}
}

func toGohomeLevel(rsmgr *gohome.ResourceManager, name string, objLoader *OBJLoader, preloaded, loadToGPU bool) *gohome.Level {
	level := &gohome.Level{Name: name}
	for i := 0; i < len(objLoader.Models); i++ {
		processModel(rsmgr, level, objLoader, &objLoader.Models[i], preloaded, loadToGPU)
	}
	return level
}

func getNameForAlreadyLoadedLevel(rsmgr *gohome.ResourceManager, name string) string {
	var alreadyLoaded = true
	var count = 1
	var newName string
	for alreadyLoaded {
		newName = name + strconv.FormatInt(int64(count), 10)
		_, alreadyLoaded = rsmgr.Levels[newName]
		count++
	}
	return newName
}

func LoadLevelOBJ(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	var alreadyLoaded = false
	if _, alreadyLoaded = rsmgr.Levels[name]; alreadyLoaded && !rsmgr.LoadModelsWithSameName {
		gohome.ErrorMgr.Log("Level", name, "Has already been loaded!")
		return nil
	}
	if alreadyLoaded {
		name = getNameForAlreadyLoadedLevel(rsmgr, name)
	}
	var objLoader OBJLoader
	if err := loadFileWithPaths(path, gohome.LEVEL_PATHS[:], &objLoader); err != nil {
		gohome.ErrorMgr.Error("Level", name, "Couldn't load "+path+": "+err.Error())
		return nil
	}
	lvl := toGohomeLevel(rsmgr, name, &objLoader, preloaded, loadToGPU)
	lvl.Name = name
	return lvl
}

func LoadLevelOBJString(rsmgr *gohome.ResourceManager, name, contents, fileName string, preloaded, loadToGPU bool) *gohome.Level {
	var alreadyLoaded = false
	if _, alreadyLoaded = rsmgr.Levels[name]; alreadyLoaded && !rsmgr.LoadModelsWithSameName {
		gohome.ErrorMgr.Log("Level", name, "Has already been loaded!")
		return nil
	}
	if alreadyLoaded {
		name = getNameForAlreadyLoadedLevel(rsmgr, name)
	}
	var objLoader OBJLoader
	objLoader.SetDirectory(gohome.GetPathFromFile(fileName))
	objLoader.SetOpenMaterialFile(gohome.Framew.OpenFile)
	objLoader.SetMaterialPaths(gohome.MATERIAL_PATHS[:])
	if err := objLoader.LoadString(contents); err != nil {
		gohome.ErrorMgr.MessageError(gohome.ERROR_LEVEL_ERROR, "Level", name, err)
		return nil
	}
	lvl := toGohomeLevel(rsmgr, name, &objLoader, preloaded, loadToGPU)
	lvl.Name = name
	return lvl
}
