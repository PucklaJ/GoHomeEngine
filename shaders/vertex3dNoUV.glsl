#version 110

attribute vec3 vertex;
attribute vec3 normal;

varying vec3 fragPos;
varying vec3 fragNormal;
varying mat4 fragViewMatrix3D;
varying mat4 fragInverseViewMatrix3D;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 inverseViewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragPos =  (transformMatrix3D*vec4(vertex,1.0)).xyz;
	// fragNormal =  (viewMatrix3D*vec4(mat3(transpose(inverseMat3(transformMatrix3D))) * normal,1.0)).xyz;
	fragNormal =  normalize((transformMatrix3D*vec4(normal,0.0)).xyz);

	fragViewMatrix3D = viewMatrix3D;
	fragInverseViewMatrix3D = inverseViewMatrix3D;
}