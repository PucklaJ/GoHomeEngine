#version 320 es

precision mediump float;
in vec2 fragTexCoords;

out vec4 fragColor;

precision mediump sampler2D;
uniform sampler2D backBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);

	color = texture(backBuffer,fragTexCoords);

	return color;
}

void main()
{
	fragColor = fetchColor();
}