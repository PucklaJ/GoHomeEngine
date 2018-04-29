#version 110

varying vec2 fragTexCoords;

uniform sampler2D backBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);

	color = texture2D(backBuffer,fragTexCoords);

	return color;
}

void main()
{
	gl_FragColor = fetchColor();
}