#version 320 es

layout(location=0) in vec3 vertex;
layout(location=1) in vec3 normal;
layout(location=2) in vec2 texCoord;
layout(location=3) in vec3 tangent;
layout(location=4) in mat4 transformMatrix3D;

out VertexOut{
	vec2 fragTexCoord;
} GeoIn;

void main()
{
	gl_Position = transformMatrix3D*vec4(vertex,1.0);
	GeoIn.fragTexCoord = texCoord;
}