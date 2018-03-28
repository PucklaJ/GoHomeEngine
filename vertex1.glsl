#version 410

layout(location = 0) in vec2 vertex;
layout(location = 1) in vec2 texCoord;

uniform mat3 transformMatrix2D;
uniform mat4 projectionMatrix2D;
uniform mat3 viewMatrix2D;

out vec2 fragTexCoord;

void main()
{
	gl_Position = projectionMatrix2D *vec4(vec2(viewMatrix2D*transformMatrix2D*vec3(vertex,1.0)),0.0,1.0);
	fragTexCoord = texCoord;
}