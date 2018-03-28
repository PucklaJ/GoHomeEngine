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

in VertexOut {
	vec2 fragTexCoord;
	vec3 fragPos;
	vec3 fragNormal;
	mat3 fragToTangentSpace;
	mat4 viewMatrix3D;
} FragIn;

out vec4 fragColor;

uniform Material material;

void main()
{	
	if(material.diffuseColor.r < 100.0 || material.specularColor.r < 100.0 || material.shinyness < 100.0 || material.diffuseTextureLoaded || material.specularTextureLoaded)
	{
		vec3 normal = vec3(0.0,0.0,0.0);
		if(material.normalMapLoaded)
		{
			normal = 2.0*(texture(material.normalMap,FragIn.fragTexCoord).xyz)-1.0;
			normal = normalize((inverse(FragIn.viewMatrix3D)*vec4(inverse(FragIn.fragToTangentSpace)*normal,0.0)).xyz);
		}
		else 
		{
			normal = FragIn.fragNormal;
			normal = normalize((inverse(FragIn.viewMatrix3D)*vec4(normal,0.0)).xyz);
		}

		fragColor = vec4(normal,1.0);
		if (fragColor.r < 0.0) {
			fragColor.r = abs(fragColor.r) * 0.1;
		}
		if (fragColor.g < 0.0) {
			fragColor.g = abs(fragColor.g) * 0.1;
		}
		if (fragColor.b < 0.0) {
			fragColor.b = abs(fragColor.b) * 0.1;
		}
	}
}