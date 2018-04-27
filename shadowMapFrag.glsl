#version 120

varying	vec2 fragTexCoord;

uniform struct Material
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
} material;

vec4 getDiffuseTexture()
{
	if(material.DiffuseTextureLoaded)
	{
		return texture2D(material.diffuseTexture,fragTexCoord);
	}
	else
	{
		return vec4(1.0,1.0,1.0,1.0);
	}
}

void main()
{
	vec4 texDifCol = getDiffuseTexture();

	if(texDifCol.a < 0.1)
	{
		discard;
	}
}