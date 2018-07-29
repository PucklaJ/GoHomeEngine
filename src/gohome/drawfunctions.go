package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
)

func toLine3D(pos1 mgl32.Vec3, pos2 mgl32.Vec3, col color.Color) Line3D {
	var line Line3D
	line[0][0] = pos1[0]
	line[0][1] = pos1[1]
	line[0][2] = pos1[2]

	vecCol := ColorToVec4(col)

	line[0][3] = vecCol[0]
	line[0][4] = vecCol[1]
	line[0][5] = vecCol[2]
	line[0][6] = vecCol[3]

	line[1][0] = pos2[0]
	line[1][1] = pos2[1]
	line[1][2] = pos2[2]

	line[1][3] = vecCol[0]
	line[1][4] = vecCol[1]
	line[1][5] = vecCol[2]
	line[1][6] = vecCol[3]

	return line
}

func DrawLine3D(pos1 mgl32.Vec3, pos2 mgl32.Vec3, col color.Color) {
	line := toLine3D(pos1, pos2, col)
	var robj Lines3D
	robj.Init()
	robj.AddLine(line)
	robj.Load()
	RenderMgr.RenderRenderObject(&robj)
	robj.Terminate()
}
