#version 100

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;
attribute mat4 transformMatrix3D;

varying vec2 fragTexCoord;

void main()
{
	gl_Position = transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
}