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
		// vec3 tex = vec3(texture(material.diffuseTexture,FragIn.fragTexCoord));
		// float brightness = (tex.r+tex.b+tex.g)/3.0;
		// float level = floor(brightness * levels);
		// brightness = level / levels;
		// return tex * brightness;
		return texture(material.diffuseTexture,FragIn.fragTexCoord);
	}
	else
	{
		return vec4(1.0,1.0,1.0,1.0);
	}
}

vec4 getSpecularTexture()
{
	if(material.specularTextureLoaded)
	{
		// vec3 tex = vec3(texture(material.specularTexture,FragIn.fragTexCoord));
		// float brightness = (tex.r+tex.b+tex.g)/3.0;
		// float level = floor(brightness * levels);
		// brightness = level / levels;
		// return tex * brightness;
		return texture(material.specularTexture,FragIn.fragTexCoord);
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

	gl_FragDepth = gl_FragCoord.z;
}