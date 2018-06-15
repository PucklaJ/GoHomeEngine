#version 150

layout(triangles) in;
layout(triangle_strip,max_vertices=18) out;

in vec2 geoTexCoord[];

out vec2 fragTexCoord;
out	vec4 fragPos;

uniform mat4 lightSpaceMatrices[6];
uniform mat4 projectionMatrix3D;

void main()
{
	for(int face = 0;face < 6;++face)
	{
		gl_Layer = face;
		for(int i = 0;i<3;++i)
		{
			fragPos = gl_in[i].gl_Position;
			gl_Position = projectionMatrix3D * lightSpaceMatrices[face] * fragPos;
			fragTexCoord = geoTexCoord[i];
			EmitVertex();
		}
		EndPrimitive();
	}
}