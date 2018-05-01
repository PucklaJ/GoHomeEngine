package loader

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"log"
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

func processModel(rsmgr *gohome.ResourceManager, level *gohome.Level, model *OBJModel, preloaded, loadToGPU bool) {
	level.LevelObjects = append(level.LevelObjects, gohome.LevelObject{})
	lvlObj := &level.LevelObjects[len(level.LevelObjects)-1]
	var lvlObjTobj gohome.TransformableObject3D
	lvlObjTobj.Position = [3]float32{0.0, 0.0, 0.0}
	lvlObjTobj.Scale = [3]float32{1.0, 1.0, 1.0}
	lvlObjTobj.Rotation = [3]float32{0.0, 0.0, 0.0}
	lvlObjTobj.CalculateTransformMatrix(nil, -1)
	lvlObj.Name = model.Name
	lvlObj.Transform.TransformMatrix = lvlObjTobj.GetTransformMatrix()
	var model3d gohome.Model3D
	model3d.Name = model.Name
	for i := 0; i < len(model.Meshes); i++ {
		mesh3d := gohome.Render.CreateMesh3D(model.Meshes[i].Name)
		processMesh(mesh3d, &model.Meshes[i], preloaded, loadToGPU)
		model3d.AddMesh3D(mesh3d)
	}
	lvlObj.Entity3D.InitModel(&model3d, nil)
	rsmgr.Models[model3d.Name] = &model3d
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

func loadMaterialTexture(path string, preloaded bool) gohome.Texture {
	var rv gohome.Texture
	if !preloaded {
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
	} else {
		rv = gohome.ResourceMgr.LoadTextureFunction(path, path, true)
		if rv == nil {
			for i := 0; i < len(gohome.MATERIAL_PATHS); i++ {
				if rv = gohome.ResourceMgr.LoadTextureFunction(path, path, true); rv == nil {
					if rv = gohome.ResourceMgr.LoadTextureFunction(path, path, true); rv != nil {
						return rv
					}
				} else {
					return rv
				}
			}
		}
	}

	return rv
}

func processMaterial(material *gohome.Material, mat *OBJMaterial, preloaded, loadToGPU bool) {
	material.Name = mat.Name
	material.DiffuseColor = toGohomeColor(mat.DiffuseColor)
	material.SpecularColor = toGohomeColor(mat.SpecularColor)
	material.SetShinyness(mat.SpecularExponent)
	if loadToGPU {
		if mat.DiffuseTexture != "" {
			material.DiffuseTexture = loadMaterialTexture(mat.DiffuseTexture, preloaded)
		}
		if mat.SpecularTexture != "" {
			material.SpecularTexture = loadMaterialTexture(mat.SpecularTexture, preloaded)
		}
		if mat.NormalMap != "" {
			material.NormalMap = loadMaterialTexture(mat.NormalMap, preloaded)
		}
	}

}

func processMesh(mesh3d gohome.Mesh3D, mesh *OBJMesh, preloaded, loadToGPU bool) {
	var vertices []gohome.Mesh3DVertex
	vertices = make([]gohome.Mesh3DVertex, len(mesh.Vertices))
	for i := 0; i < len(vertices); i++ {
		vertices[i] = toMesh3DVertex(&mesh.Vertices[i])
	}
	var mat gohome.Material
	mat.InitDefault()
	processMaterial(&mat, mesh.Material, preloaded, loadToGPU)
	mesh3d.SetMaterial(&mat)
	mesh3d.AddVertices(vertices, mesh.Indices)
	if !preloaded {
		if loadToGPU {
			mesh3d.Load()
			log.Println("Finished loading mesh", mesh3d.GetName(), "V:", mesh3d.GetNumVertices(), "I:", mesh3d.GetNumIndices())
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
		processModel(rsmgr, level, &objLoader.Models[i], preloaded, loadToGPU)
	}
	return level
}

func LoadLevelOBJ(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	if _, ok := rsmgr.Levels[name]; ok {
		log.Println("The level with the name", name, "has already been loaded!")
		return nil
	}
	var objLoader OBJLoader
	if err := loadFileWithPaths(path, gohome.LEVEL_PATHS[:], &objLoader); err != nil {
		log.Println("Couldn't load level", name, "with path", path, ":", err.Error())
		return nil
	}
	lvl := toGohomeLevel(rsmgr, name, &objLoader, preloaded, loadToGPU)
	return lvl
}
