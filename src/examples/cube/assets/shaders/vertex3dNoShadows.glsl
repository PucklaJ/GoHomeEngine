#version 100

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;

varying vec2 fragTexCoord;
varying vec3 fragPos;
varying vec3 fragNormal;
varying mat3 fragToTangentSpace;
varying mat4 fragViewMatrix3D;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
	fragPos =  (viewMatrix3D*transformMatrix3D*vec4(vertex,1.0)).xyz;
	// fragNormal =  (viewMatrix3D*vec4(mat3(transpose(inverse(transformMatrix3D))) * normal,1.0)).xyz;
	fragNormal =  (viewMatrix3D*transformMatrix3D*vec4(normal,0.0)).xyz;

	vec3 norm = normalize(fragNormal);
	vec3 tang = normalize((viewMatrix3D*transformMatrix3D*vec4(tangent,0.0)).xyz);
	vec3 bitang = normalize(cross(norm,tang));

	fragToTangentSpace = mat3(
		tang.x,bitang.x,norm.x,
		tang.y,bitang.y,norm.y,
		tang.z,bitang.z,norm.z
	);

	fragViewMatrix3D = viewMatrix3D;
}