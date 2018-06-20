#version 110

attribute vec2 vertex;
attribute vec2 texCoord;

varying vec2 fragTexCoords;

uniform float depth;

void main()
{
	fragTexCoords = texCoord;
	gl_Position = vec4(vertex,depth,1.0);
}