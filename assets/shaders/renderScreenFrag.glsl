#version 100

precision mediump float;
varying vec2 fragTexCoords;

precision mediump sampler2D;
uniform sampler2D BackBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);

	color = texture2D(BackBuffer,fragTexCoords);

	return color;
}

void main()
{
	gl_FragColor = fetchColor();
}