#version 410

layout (triangles) in;
layout (line_strip, max_vertices = 6) out;

in VertexOut {
	vec2 fragTexCoord;
	vec3 fragPos;
	vec3 fragNormal;
} GeoIn[];

out VertexOut {
	vec2 fragTexCoord;
	vec3 fragPos;
	vec3 fragNormal;
} FragIn;

void makeLine(vec4 pos1,vec4 pos2, int index1,int index2);

void main()
{
	makeLine(gl_in[0].gl_Position,gl_in[1].gl_Position,0,1);
	makeLine(gl_in[1].gl_Position,gl_in[2].gl_Position,1,2);
	makeLine(gl_in[2].gl_Position,gl_in[0].gl_Position,2,0);
}

void makeLine(vec4 pos1,vec4 pos2,int index1, int index2)
{
	gl_Position = pos1;
	FragIn.fragTexCoord = GeoIn[index1].fragTexCoord;
	FragIn.fragPos = GeoIn[index1].fragPos;
	FragIn.fragNormal = GeoIn[index1].fragNormal;
	EmitVertex();

	gl_Position = pos2;
	FragIn.fragTexCoord = GeoIn[index2].fragTexCoord;
	FragIn.fragPos = GeoIn[index2].fragPos;
	FragIn.fragNormal = GeoIn[index2].fragNormal;
	EmitVertex();

	EndPrimitive();
}