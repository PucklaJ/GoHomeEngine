#version 100

attribute vec2 vertex;
attribute vec2 texCoord;

varying vec2 fragTexCoord;

void main()
{
	gl_Position = vec4(vertex,0.0,1.0);
	fragTexCoord = texCoord;
}