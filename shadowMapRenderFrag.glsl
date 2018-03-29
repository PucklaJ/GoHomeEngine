#version 410

in vec2 fragTexCoord;
out vec4 fragColor;

uniform sampler2D texture0;

void main()
{
	float depth = texture(texture0,fragTexCoord).r;
	fragColor = vec4(depth,depth,depth,1.0);
}