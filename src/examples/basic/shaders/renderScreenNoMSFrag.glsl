#version 110

varying vec2 fragTexCoords;

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
	if(gl_FragColor.a < 0.1)
	    discard;
}