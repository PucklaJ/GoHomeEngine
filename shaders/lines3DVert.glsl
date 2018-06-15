#version 110

attribute vec3 vertex;
attribute vec4 color;

varying vec4 fragColor;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragColor = color;
}