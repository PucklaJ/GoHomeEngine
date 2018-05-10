package loader

import (
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/assimp"
	"log"
	"strconv"
	"strings"
	"sync"
)

const (
	NUM_GO_ROUTINES_MESH_VERTICES_LOADING uint32 = 10
)

func importFile(path string) *assimp.Scene {
	var scene *assimp.Scene
	scene = assimp.ImportFile(path, uint(assimp.Process_Triangulate|assimp.Process_FlipUVs|assimp.Process_GenNormals|assimp.Process_OptimizeMeshes))
	if scene == nil || (scene.Flags()&assimp.SceneFlags_Incomplete) != 0 || scene.RootNode() == nil {
		return nil
	}
	return scene
}

func importFileWithPaths(path string, paths []string) (*assimp.Scene, string) {
	var scene *assimp.Scene
	var worked bool = false
	var workingPath string
	for i := 0; i < len(paths); i++ {
		if scene = importFile(paths[i] + path); scene != nil {
			worked = true
			workingPath = paths[i] + path
			break
		} else if scene = importFile(paths[i] + gohome.GetFileFromPath(path)); scene != nil {
			worked = true
			workingPath = paths[i] + gohome.GetFileFromPath(path)
			break
		}
	}

	if worked {
		return scene, workingPath
	} else {
		return nil, ""
	}
}

func LoadLevelAssimp(rsmgr *gohome.ResourceManager, name, path string, preloaded, loadToGPU bool) *gohome.Level {
	if _, ok := rsmgr.Levels[name]; ok {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_LOG, "Level", name, "Has already been loaded!")
		return nil
	}
	level := &gohome.Level{Name: name}
	var scene *assimp.Scene
	var workingPath string
	if scene, workingPath = importFileWithPaths(path, gohome.LEVEL_PATHS[:]); scene == nil {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Level", name, "Couldn't load file "+path+": "+assimp.GetErrorString())
		return nil
	}

	directory := workingPath
	if index := strings.LastIndex(directory, "/"); index != -1 {
		directory = directory[:index+1]
	} else {
		directory = ""
	}

	processNode(rsmgr, scene.RootNode(), scene, level, directory, preloaded, loadToGPU)

	return level
}

func processNode(rsmgr *gohome.ResourceManager, node *assimp.Node, scene *assimp.Scene, level *gohome.Level, directory string, preloaded, loadToGPU bool) {
	if node != scene.RootNode() {
		model := &gohome.Model3D{}
		initModel(model, node, scene, level, directory, preloaded, loadToGPU)
		if !preloaded {
			if _, ok := rsmgr.Models[model.Name]; ok {
				gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_LOG, "Model", model.Name, "Has already been loaded! Overwritting ...")
			}
			rsmgr.Models[model.Name] = model
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_LOG, "Model", model.Name, "Finished loading!")
		} else {
			rsmgr.PreloadedModelsChan <- model
		}

	}
	for i := 0; i < node.NumChildren(); i++ {
		processNode(rsmgr, node.Children()[i], scene, level, directory, preloaded, loadToGPU)
	}
}

func initModel(model *gohome.Model3D, node *assimp.Node, scene *assimp.Scene, level *gohome.Level, directory string, preloaded, loadToGPU bool) {
	level.LevelObjects = append(level.LevelObjects, gohome.LevelObject{
		Name: node.Name(),
	})
	setTransformLevelObject(&level.LevelObjects[len(level.LevelObjects)-1], node.Transformation())

	model.Name = node.Name()
	for i := 0; i < node.NumMeshes(); i++ {
		aiMesh := scene.Meshes()[node.Meshes()[i]]
		mesh := gohome.Render.CreateMesh3D(aiMesh.Name())
		addVerticesAssimpMesh3D(mesh, aiMesh, node, scene, level, directory, preloaded)
		model.AddMesh3D(mesh)
		if !preloaded {
			if loadToGPU {
				mesh.Load()
			}
			log.Println("Finished loading mesh", mesh.GetName(), "V:", mesh.GetNumVertices(), "I:", mesh.GetNumIndices(), "!")
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_LOG, "Mesh", mesh.GetName(), "Finished loading! V:"+strconv.Itoa(int(mesh.GetNumVertices()))+" I: "+strconv.Itoa(int(mesh.GetNumIndices())))
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

func addVerticesAssimpMesh3D(oglm gohome.Mesh3D, mesh *assimp.Mesh, node *assimp.Node, scene *assimp.Scene, level *gohome.Level, directory string, preloaded bool) {
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
	initMaterial(mat, scene.Materials()[mesh.MaterialIndex()], scene, directory, preloaded)
	oglm.SetMaterial(mat)
}

func initMaterial(mat *gohome.Material, material *assimp.Material, scene *assimp.Scene, directory string, preloaded bool) {
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
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Material", "assimp", "Couldn't return diffuse Texture")
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
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Material", "assimp", "Couldn't return specular Texture")
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
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "Material", "assimp", "Couldn't return normal map")
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

func setTransformLevelObject(this *gohome.LevelObject, mat assimp.Matrix4x4) {
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			this.Transform.TransformMatrix[this.Transform.TransformMatrix.Index(r, c)] = mat.Values()[c][r]
		}
	}
}
