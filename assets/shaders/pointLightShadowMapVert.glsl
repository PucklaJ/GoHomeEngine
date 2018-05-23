#version 100

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;


varying	vec2 fragTexCoord;


uniform mat4 transformMatrix3D;

void main()
{
	gl_Position = transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
}