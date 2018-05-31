package gohome

import "image/color"

type AABBRenderer struct {
	AABB *AxisAlignedBoundingBox
	Lines3D
}

func (this *AABBRenderer) Init(aabb *AxisAlignedBoundingBox,transform TransformableObject,col color.Color) {
	this.AABB = aabb
	this.Lines3D.Init()
	this.SetTransformableObject(transform)

	lines := [12]Line3D{
			{
				{
					aabb.Min.X(), aabb.Min.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Max.X(), aabb.Min.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
			{
				{
					aabb.Max.X(), aabb.Min.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Max.X(), aabb.Max.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
			{
				{
					aabb.Max.X(), aabb.Max.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Min.X(), aabb.Max.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
			{
				{
					aabb.Min.X(), aabb.Max.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Min.X(), aabb.Min.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},



			{
				{
					aabb.Max.X(), aabb.Min.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Max.X(), aabb.Min.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
			{
				{
					aabb.Max.X(), aabb.Max.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Max.X(), aabb.Max.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
			{
				{
					aabb.Min.X(), aabb.Max.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Min.X(), aabb.Max.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
			{
				{
					aabb.Min.X(), aabb.Min.Y(), aabb.Max.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Min.X(), aabb.Min.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},


			{
				{
					aabb.Min.X(), aabb.Min.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Max.X(), aabb.Min.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
			{
				{
					aabb.Max.X(), aabb.Min.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Max.X(), aabb.Max.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
			{
				{
					aabb.Max.X(), aabb.Max.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Min.X(), aabb.Max.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
			{
				{
					aabb.Min.X(), aabb.Max.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
				{
					aabb.Min.X(), aabb.Min.Y(), aabb.Min.Z(), 1.0, 1.0, 1.0, 1.0,
				},
			},
		}

	this.AddLines(lines[:])
	this.SetColor(col)
	this.Load()
}
