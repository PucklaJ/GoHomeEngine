#version 110

varying vec2 fragTexCoord;

uniform sampler2D texture0;

void main()
{
	float depth = texture2D(texture0,fragTexCoord).r;
	gl_FragColor = vec4(depth,depth,depth,1.0);
}