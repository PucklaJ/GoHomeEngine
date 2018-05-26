#version 100

precision mediump float;
varying vec2 fragTexCoords;

precision mediump sampler2D;
uniform sampler2D BackBuffer;

vec4 fetchColor()
{
	return texture2D(BackBuffer,fragTexCoords);
}

void main()
{
	gl_FragColor = fetchColor();
	if(gl_FragColor.a < 0.1)
	{
		discard;
	}
}