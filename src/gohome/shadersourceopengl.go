package gohome

// Shadowmap Shader
const (
	SHADOWMAP_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;


varying vec2 fragTexCoord;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
}`
	SHADOWMAP_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

varying	vec2 fragTexCoord;

uniform struct Material
{
	bool DiffuseTextureLoaded;
	float transparency;
} material;
uniform sampler2D materialdiffuseTexture;

vec4 getDiffuseTexture()
{
	if(material.DiffuseTextureLoaded)
	{
		return texture2D(materialdiffuseTexture,fragTexCoord);
	}
	else
	{
		return vec4(1.0,1.0,1.0,1.0);
	}
}

void main()
{
	vec4 texDifCol = getDiffuseTexture();
	texDifCol.w *= material.transparency;

	if(texDifCol.a < 0.1)
	{
		discard;
	}
}`
	SHADOWMAP_INSTANCED_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;
attribute mat4 transformMatrix3D;

varying vec2 fragTexCoord;


uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
}`
)
