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
}