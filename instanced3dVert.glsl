#version 410

layout(location=0) in vec3 vertex;
layout(location=1) in vec3 normal;
layout(location=2) in vec2 texCoord;
layout(location=3) in vec3 tangent;
layout(location=4) in mat4 transformMatrix3D;

out VertexOut{
	vec2 fragTexCoord;
	vec3 fragPos;
	vec3 fragNormal;
	mat3 fragToTangentSpace;
	mat4 viewMatrix3D;
} FragIn;

uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	FragIn.fragTexCoord = texCoord;
	FragIn.fragPos =  (viewMatrix3D*transformMatrix3D*vec4(vertex,1.0)).xyz;
	// FragIn.fragNormal =  (viewMatrix3D*vec4(mat3(transpose(inverse(transformMatrix3D))) * normal,1.0)).xyz;
	FragIn.fragNormal =  (viewMatrix3D*transformMatrix3D*vec4(normal,0.0)).xyz;

	vec3 norm = normalize(FragIn.fragNormal);
	vec3 tang = normalize((viewMatrix3D*transformMatrix3D*vec4(tangent,0.0)).xyz);
	vec3 bitang = normalize(cross(norm,tang));

	FragIn.fragToTangentSpace = mat3(
		tang.x,bitang.x,norm.x,
		tang.y,bitang.y,norm.y,
		tang.z,bitang.z,norm.z
	);

	FragIn.viewMatrix3D = viewMatrix3D;
}