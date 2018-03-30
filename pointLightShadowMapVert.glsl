#version 410

layout(location=0) in vec3 vertex;
layout(location=1) in vec3 normal;
layout(location=2) in vec2 texCoord;
layout(location=3) in vec3 tangent;

out VertexOut{
	vec2 fragTexCoord;
} GeoIn;

uniform mat4 transformMatrix3D;

void main()
{
	gl_Position = transformMatrix3D*vec4(vertex,1.0);
	GeoIn.fragTexCoord = texCoord;
}