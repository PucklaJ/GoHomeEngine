#version 110

#define FLIP_NONE 0
#define FLIP_HORIZONTAL 1
#define FLIP_VERTICAL 2
#define FLIP_DIAGONALLY 3

attribute vec2 vertex;
attribute vec2 texCoord;

uniform mat3 transformMatrix2D;
uniform mat4 projectionMatrix2D;
uniform mat3 viewMatrix2D;
uniform int flip;

varying vec2 fragTexCoord;

vec2 flipTexCoord();

void main()
{
	vec2 flippedTexCoord = flipTexCoord();
	gl_Position = projectionMatrix2D *vec4(vec2(viewMatrix2D*transformMatrix2D*vec3(vertex,1.0)),0.0,1.0);
	fragTexCoord = flippedTexCoord;
}

vec2 flipTexCoord()
{
	vec2 flippedTexCoord;

	
	if(flip == FLIP_NONE)
	{
		flippedTexCoord = texCoord;
	}
	else if(flip == FLIP_HORIZONTAL)
	{
		flippedTexCoord.x = 1.0-texCoord.x;
		flippedTexCoord.y = texCoord.y;
	}
	else if(flip == FLIP_VERTICAL)
	{
		flippedTexCoord.x = texCoord.x;
		flippedTexCoord.y = 1.0-texCoord.y;
	}
	else if(flip == FLIP_DIAGONALLY)
	{
		flippedTexCoord.x = 1.0-texCoord.x;
		flippedTexCoord.y = 1.0-texCoord.y;
	}
	else
	{
		flippedTexCoord = texCoord;
	}

	

	return flippedTexCoord;
}