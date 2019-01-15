package gohome

import (
	"io"
	"strings"
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
		contents, err := ReadAll(reader)
		if err != nil {
			return "", err
		}
		return contents, nil
	}
}

func loadShader(path, name_shader, name string) chan string {
	rv := make(chan string)

	go func() {
		contents, err := loadShaderFile(path)
		if err != nil {
			rv <- "Couldn't load " + name + ": " + err.Error()
		} else {
			rv <- contents
		}
	}()

	return rv
}

func (rsmgr *ResourceManager) LoadShader(name, vertex_path, fragment_path, geometry_path, tesselletion_control_path, eveluation_path, compute_path string) Shader {
	_, already := rsmgr.shaders[name]
	if already {
		ErrorMgr.Message(ERROR_LEVEL_LOG, "Shader", name, "Has already been loaded!")
		return nil
	}

	var contents [6]string
	paths := [6]string{
		vertex_path,
		fragment_path,
		geometry_path,
		tesselletion_control_path,
		eveluation_path,
		compute_path,
	}
	names := [6]string{
		"Vertex File",
		"Fragment File",
		"Geometry File",
		"Tesselletion File",
		"Eveluation File",
		"Compute File",
	}
	var chans [6]chan string

	for i := 0; i < 6; i++ {
		chans[i] = loadShader(paths[i], name, names[i])
	}

	var failed = false

	for i := 0; i < 6; i++ {
		contents[i] = <-chans[i]
		close(chans[i])
		if strings.HasPrefix(contents[i], "Couldn't") {
			failed = true
			ErrorMgr.Error("Shader", name, contents[i])
		}
	}

	if failed {
		return nil
	}

	shader, err := Render.LoadShader(name, contents[VERTEX], contents[FRAGMENT], contents[GEOMETRY], contents[TESSELLETION], contents[EVELUATION], contents[COMPUTE])
	if err != nil {
		ErrorMgr.MessageError(ERROR_LEVEL_ERROR, "Shader", name, err)
		return nil
	}

	if shader != nil {
		rsmgr.shaders[name] = shader
		ErrorMgr.Log("Shader", name, "Finished loading!")
	}
	return shader
}

func (rsmgr *ResourceManager) LoadShaderSource(name, vertex, fragment, geometry, tesselletion_control, eveluation, compute string) Shader {
	if _, ok := rsmgr.shaders[name]; ok {
		ErrorMgr.Error("Shader", name, "Has already been loaded")
		return rsmgr.shaders[name]
	}
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
