package gohome

import "image/color"

// This class is used for debug purposes
// It renders an AxisAlignedBoundingBox
type AABBRenderer struct {
	AABB *AxisAlignedBoundingBox
	Shape3D
}

// Init takes a pointer to the AxisAlignedBoundingBox which should be drawn
// It also takes the color in which it should be drawn
func (this *AABBRenderer) Init(aabb *AxisAlignedBoundingBox, transform TransformableObject, col color.Color) {
	this.AABB = aabb
	this.Shape3D.Init()
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
