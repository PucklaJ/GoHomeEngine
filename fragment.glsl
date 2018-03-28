#version 410

in vec2 fragTexCoord;
out vec4 fragColor;

uniform sampler2D texture0;

void main()
{
	fragColor = texture2D(texture0,fragTexCoord);
	if(fragColor.a < 0.1)
	{
		discard;
	}
}