#version 150

in vec3 vertex;
in vec3 normal;
in vec2 texCoord;
in vec3 tangent;

out vec2 geoTexCoord;

uniform mat4 transformMatrix3D;

void main()
{
	gl_Position = transformMatrix3D*vec4(vertex,1.0);
	geoTexCoord = texCoord;
}