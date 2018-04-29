#version 150

in vec2 fragTexCoords;

out vec4 fragColor;

uniform sampler2DMS backBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);
	ivec2 texCoords = ivec2(fragTexCoords * textureSize(backBuffer));

	for(int i = 0;i<8;i++)
	{
		color += texelFetch(backBuffer,texCoords,i);
	}
	color /= 8.0;

	return color;
}

void main()
{
	fragColor = fetchColor();
	if(fragColor.a < 0.1)
	{
		discard;
	}
}