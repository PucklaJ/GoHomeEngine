package gohome

func (rsmgr *ResourceManager) CheckLevel(name, path string) (l, q bool, n string) {

	if _, ok := rsmgr.Levels[name]; ok {
		l = true
		if !rsmgr.LoadModelsWithSameName {
			q = true
			ErrorMgr.Log("Level", name, "Has already been loaded!")
		}

	}

	if name, ok := rsmgr.resourceFileNames[path]; ok {
		n = name
	}

	return
}

func (rsmgr *ResourceManager) CheckModel(name string) (l, q bool) {
	if _, ok := rsmgr.Models[name]; ok {
		l = true
		if !rsmgr.LoadModelsWithSameName {
			q = true
			ErrorMgr.Log("Model", name, "Has already been loaded!")
		}
	}
	return
}
