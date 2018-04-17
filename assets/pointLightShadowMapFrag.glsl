#version 320 es

precision mediump float;
precision mediump sampler2D;

struct Material
{
	vec3 diffuseColor;
	vec3 specularColor;

	sampler2D diffuseTexture;
	sampler2D specularTexture;
	sampler2D normalMap;

	bool DiffuseTextureLoaded;
	bool SpecularTextureLoaded;
	bool NormalMapLoaded;

	float shinyness;
};

in VertexOut{
	vec2 fragTexCoord;
	vec4 fragPos;
} FragIn;

uniform vec3 lightPos;
uniform float farPlane;
uniform Material material;

vec4 fetchColor()
{
	if(material.DiffuseTextureLoaded)
	{
		return texture(material.diffuseTexture,FragIn.fragTexCoord);
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
	float lightDistance = length(FragIn.fragPos.xyz - lightPos);
	lightDistance = lightDistance / farPlane;
	gl_FragDepth = lightDistance;
}