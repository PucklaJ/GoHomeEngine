#version 410

in VertexOut{
	vec2 fragTexCoord;
	vec4 fragPos;
} FragIn;

uniform vec3 lightPos;
uniform float farPlane;

void main()
{
	float lightDistance = length(FragIn.fragPos.xyz - lightPos);
	lightDistance = lightDistance / farPlane;
	gl_FragDepth = lightDistance;
}