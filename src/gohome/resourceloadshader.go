package gohome

import (
	"io"
	"io/ioutil"
)

var (
	SHADER_PATHS = [4]string{
		"",
		"shaders/",
		"assets/",
		"assets/shaders/",
	}
)

func loadShaderFile(path string) (string, error) {
	if path == "" {
		return "", nil
	} else {
		var reader io.Reader
		var err1 error
		reader, _, err1 = OpenFileWithPaths(path, SHADER_PATHS[:])
		if err1 != nil {
			return "", err1
		}
		contents, err := ioutil.ReadAll(reader)
		if err != nil {
			return "", err
		} else {
			return string(contents[:len(contents)]), nil
		}
	}
}

func loadShader(path, name_shader, name string) (string, bool) {
	path = Render.FilterShaderFiles(name_shader, path, name)
	contents, err := loadShaderFile(path)
	if err != nil {
		ErrorMgr.Error("Shader", name_shader, "Couldn't load "+name+": "+err.Error())
		return "", true
	}
	return contents, false
}

func (rsmgr *ResourceManager) LoadShader(name, vertex_path, fragment_path, geometry_path, tesselletion_control_path, eveluation_path, compute_path string) Shader {
	_, already := rsmgr.shaders[name]
	if already {
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Shader", name, "Has already been loaded!")
		return nil
	}

	var contents [6]string
	var err bool
	var erro error

	contents[VERTEX], err = loadShader(vertex_path, name, "Vertex File")
	if err {
		return nil
	}
	contents[FRAGMENT], err = loadShader(fragment_path, name, "Fragment File")
	if err {
		return nil
	}
	contents[GEOMETRY], err = loadShader(geometry_path, name, "Geometry File")
	if err {
		return nil
	}
	contents[TESSELLETION], err = loadShader(tesselletion_control_path, name, "Tesselletion File")
	if err {
		return nil
	}
	contents[EVELUATION], err = loadShader(eveluation_path, name, "Eveluation File")
	if err {
		return nil
	}
	contents[COMPUTE], err = loadShader(compute_path, name, "Compute File")
	if err {
		return nil
	}

	var shader Shader

	shader, erro = Render.LoadShader(name, contents[VERTEX], contents[FRAGMENT], contents[GEOMETRY], contents[TESSELLETION], contents[EVELUATION], contents[COMPUTE])
	if erro != nil {
		ErrorMgr.MessageError(ERROR_LEVEL_ERROR, "Shader", name, erro)
		return nil
	}

	if shader != nil {
		rsmgr.shaders[name] = shader
		ErrorMgr.Log("Shader", name, "Finished loading!")
	}
	return shader
}

func filterShaderSource(name, vertex, fragment, geometry, tesselletion_control, eveluation, compute string) (string, string, string, string, string, string) {
	vertex = Render.FilterShaderSource(name, vertex, "Vertex File")
	fragment = Render.FilterShaderSource(name, fragment, "Fragment File")
	geometry = Render.FilterShaderSource(name, geometry, "Geometry File")
	tesselletion_control = Render.FilterShaderSource(name, tesselletion_control, "Tesselation Control File")
	eveluation = Render.FilterShaderSource(name, eveluation, "Eveluation File")
	compute = Render.FilterShaderSource(name, compute, "Compute File")

	return vertex, fragment, geometry, tesselletion_control, eveluation, compute
}

func (rsmgr *ResourceManager) LoadShaderSource(name, vertex, fragment, geometry, tesselletion_control, eveluation, compute string) Shader {
	if _, ok := rsmgr.shaders[name]; ok {
		ErrorMgr.Error("Shader", name, "Has already been loaded")
		return rsmgr.shaders[name]
	}
	vertex, fragment, geometry, tesselletion_control, eveluation, compute = filterShaderSource(name, vertex, fragment, geometry, tesselletion_control, eveluation, compute)
	shader, err := Render.LoadShader(name, vertex, fragment, geometry, tesselletion_control, eveluation, compute)
	if err != nil {
		ErrorMgr.Error("Shader", name, "Loading source: "+err.Error())
		return nil
	}
	rsmgr.shaders[name] = shader
	ErrorMgr.Log("Shader", name, "Finished Loading!")
	return shader
}

func (rsmgr *ResourceManager) GetShader(name string) Shader {
	s := rsmgr.shaders[name]
	return s
}

func (rsmgr *ResourceManager) SetShader(name string, name1 string) {
	s := rsmgr.shaders[name1]
	if s == nil {
		ErrorMgr.Message(ERROR_LEVEL_ERROR, "Shader", name, "Couldn't set to "+name1+" (It is nil)")
		return
	}
	rsmgr.shaders[name] = s
	ErrorMgr.Message(ERROR_LEVEL_LOG, "Shader", name, "Set to "+name1)
}

func (rsmgr *ResourceManager) DeleteShader(name string) {
	if _, ok := rsmgr.shaders[name]; ok {
		rsmgr.shaders[name].Terminate()
		delete(rsmgr.shaders, name)
		rsmgr.deleteResourceFileName(name)
		ErrorMgr.Log("Shader", name, "Deleted!")
	} else {
		ErrorMgr.Warning("Shader", name, "Couldn't delete! It hasn't been loaded!")
	}
}
