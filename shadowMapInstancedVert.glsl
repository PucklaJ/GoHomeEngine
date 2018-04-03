#version 410

layout(location=0) in vec3 vertex;
layout(location=1) in vec3 normal;
layout(location=2) in vec2 texCoord;
layout(location=3) in vec3 tangent;
layout(location=4) in mat4 transformMatrix3D;

out VertexOut{
	vec2 fragTexCoord;
} FragIn;

uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	FragIn.fragTexCoord = texCoord;
}