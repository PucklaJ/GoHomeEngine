#version 100

precision mediump float;
precision mediump sampler2D;

varying vec2 fragTexCoord;

uniform sampler2D texture0;

void main()
{
	gl_FragColor = texture2D(texture0,fragTexCoord);
	if(gl_FragColor.a < 0.1)
	{
		discard;
	}
}