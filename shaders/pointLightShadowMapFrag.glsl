#version 150

in vec2 fragTexCoord;
in vec4 fragPos;

uniform vec3 lightPos;
uniform float farPlane;
uniform struct Material
{
	bool DiffuseTextureLoaded;
} material;
uniform	sampler2D materialdiffuseTexture;


vec4 fetchColor()
{
	if(material.DiffuseTextureLoaded)
	{
		return texture2D(materialdiffuseTexture,fragTexCoord);
	}
	else
	{
		return vec4(1.0,1.0,1.0,1.0);
	}
}

void main()
{
	vec4 color = fetchColor();
	if(color.a < 0.1)
		discard;
	float lightDistance = length(fragPos.xyz - lightPos);
	lightDistance = lightDistance / farPlane;
	gl_FragDepth = lightDistance;
}