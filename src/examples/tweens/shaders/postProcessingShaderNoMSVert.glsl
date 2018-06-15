#version 120

attribute vec2 vertex;
attribute vec2 texCoord;

varying vec2 fragTexCoords;

void main()
{
	fragTexCoords = texCoord;
	gl_Position = vec4(vertex,0.0,1.0);
}