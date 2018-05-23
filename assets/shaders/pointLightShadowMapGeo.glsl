#version 320 es

layout(triangles) in;
layout(triangle_strip,max_vertices=18) out;

in VertexOut{
	vec2 fragTexCoord;
} GeoIn[];

out VertexOut{
	vec2 fragTexCoord;
	vec4 fragPos;
} FragIn;

uniform mat4 lightSpaceMatrices[6];
uniform mat4 projectionMatrix3D;

void main()
{
	for(int face = 0;face < 6;++face)
	{
		gl_Layer = face;
		for(int i = 0;i<3;++i)
		{
			FragIn.fragPos = gl_in[i].gl_Position;
			gl_Position = projectionMatrix3D * lightSpaceMatrices[face] * FragIn.fragPos;
			FragIn.fragTexCoord = GeoIn[i].fragTexCoord;
			EmitVertex();
		}
		EndPrimitive();
	}
}