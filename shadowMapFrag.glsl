#version 410

struct Material
{
	vec3 diffuseColor;
	vec3 specularColor;

	sampler2D diffuseTexture;
	sampler2D specularTexture;
	sampler2D normalMap;

	bool diffuseTextureLoaded;
	bool specularTextureLoaded;
	bool normalMapLoaded;

	float shinyness;
};

in VertexOut{
	vec2 fragTexCoord;
} FragIn;

uniform Material material;

vec4 getDiffuseTexture()
{
	if(material.diffuseTextureLoaded)
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
	vec4 texDifCol = getDiffuseTexture();

	if(texDifCol.a < 0.1)
	{
		discard;
	}
}