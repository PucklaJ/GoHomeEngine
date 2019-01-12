## General

+ Controller Input
+ Simple Loading Screen
+ Add DirectXRenderer
+ WorldEditor (3D, 2D)
+ Add Tween interpreter language
+ Add Parents for TransformableObjects2D
+ Collada loader
+ Parent System for Sprite2D
+ OBJLoader load non-triangulate faces
+ Add more sprites to spriteanimation example
+ Add PostProcessing
    - interface with DoPostProcessing()
    - add GlowShader
    - with GLOW_BIT RenderType
+ Add indirect rendering
    - add mesh/model to IndirectRender struct (?)
    - and then call Render on this struct
+ Add Vertices for NonUV meshes
    - without texCoord and tangent
    - convert Mesh3DVertices to Mesh3DNoUVVertices (with go routines)
    - or rewrite AddVertices -> (vertex,normal,texCoord)
    - and rewrite order of data in Mesh3D
    - change shaders accordingly (remove attributes)
    - in instancedmesh think about index
+ Change order of custom values of InstancedMesh3D
    - to DataOfVal1|DataOfVal2 etc. instead of Val1OfVal1|Val1OfVal2|Val2OfVal1|Val2OfVal2
+ Rewrite PointLight shadows
+ Bump mapping
    - add global fragPos
    - use global fragPos
+ Parallax Mapping / Parallax Occlusion Mapping / Steep Parallax Mapping
+ Add Changeable MAX_POINT_LIGHTS, MAX_DIRECTIONAL_LIGHTS, etc. to LightMgr
    - maybe all makros (APLHA_DISCARD_PADDING)
+ Change obj loading to be concurrent
+ Remove audio manager from framework
	- Add InitAudio method to gohome which takes a AudioManager pointer
+ Fix OpenAL errors
    - device not closed

## 2D

+ Instanced Mesh
+ Particle
+ Add TextureRenderer

## 3D

+ Add Sprite3D
    - Add Option For Billboarding
+ Add Debug Renderer
    - renders AABBs
+ Particle
+ Raycasting
+ Frustum culling
+ Physics
	- [cubez](https://github.com/tbogdala/cubez) maybe
+ Add Short Light add methods
+ Remove light collection and have lights per scene
