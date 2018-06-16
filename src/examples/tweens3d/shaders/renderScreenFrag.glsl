#version 150

in vec2 fragTexCoords;

out vec4 fragColor;

uniform sampler2D BackBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);

	color = texture2D(BackBuffer,fragTexCoords);

	return color;
}

void main()
{
	fragColor = fetchColor();
	if(fragColor.a < 0.1)
	    discard;
}