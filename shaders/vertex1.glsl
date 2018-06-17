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
uniform vec4 textureRegion;
uniform int flip;
uniform float depth;

varying vec2 fragTexCoord;

vec2 textureRegionToTexCoord(vec2 tc);
vec2 flipTexCoord(vec2 tc);

void main()
{
	gl_Position = projectionMatrix2D *vec4(vec2(viewMatrix2D*transformMatrix2D*vec3(vertex,1.0)),depth,1.0);
	fragTexCoord = textureRegionToTexCoord(flipTexCoord(texCoord));
}

vec2 textureRegionToTexCoord(vec2 tc)
{
    // X: 0 ->      0 -> Min X
    // Y: 0 ->      0 -> Min Y
    // Z: WIDTH ->  1 -> Max X
    // W: HEIGHT -> 1 -> Max Y

    vec2 newTexCoord = tc;
    newTexCoord.x = newTexCoord.x * (textureRegion.z-textureRegion.x)+textureRegion.x;
    newTexCoord.y = newTexCoord.y * (textureRegion.w-textureRegion.y)+textureRegion.y;

    return newTexCoord;
}

vec2 flipTexCoord(vec2 tc)
{
	vec2 flippedTexCoord;


	if(flip == FLIP_NONE)
	{
		flippedTexCoord = tc;
	}
	else if(flip == FLIP_HORIZONTAL)
	{
		flippedTexCoord.x = 1.0-tc.x;
		flippedTexCoord.y = tc.y;
	}
	else if(flip == FLIP_VERTICAL)
	{
		flippedTexCoord.x = tc.x;
		flippedTexCoord.y = 1.0-tc.y;
	}
	else if(flip == FLIP_DIAGONALLY)
	{
		flippedTexCoord.x = 1.0-tc.x;
		flippedTexCoord.y = 1.0-tc.y;
	}
	else
	{
		flippedTexCoord = tc;
	}



	return flippedTexCoord;
}