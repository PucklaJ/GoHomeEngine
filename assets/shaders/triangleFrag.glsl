#version 100

precision mediump float;
varying vec2 fragTexCoord;

void main()
{
	gl_FragColor = vec4(fragTexCoord.r,1.0,1.0,1.0);
}