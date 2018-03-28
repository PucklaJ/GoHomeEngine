#version 410

#define MAX_POINT_LIGHTS 20
#define MAX_DIRECTIONAL_LIGHTS 2
#define MAX_SPOT_LIGHTS 10

#define MAX_SPECULAR_EXPONENT 50.0
#define MIN_SPECULAR_EXPONENT 5.0

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

struct Attentuation
{
	float constant;
	float linear;
	float quadratic;
};

struct PointLight
{
	vec3 position;

	vec3 diffuseColor;
	vec3 specularColor;

	Attentuation attentuation;
};

struct DirectionalLight
{
	vec3 direction;

	vec3 diffuseColor;
	vec3 specularColor;
};

struct SpotLight
{
	vec3 position;
	vec3 direction;

	vec3 diffuseColor;
	vec3 specularColor;

	float innerCutOff;
	float outerCutOff;

	Attentuation attentuation;
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
uniform uint numPointLights = 0;
uniform uint numDirectionalLights = 0;
uniform uint numSpotLights = 0;
uniform vec3 ambientLight;
uniform PointLight[MAX_POINT_LIGHTS] pointLights;
uniform DirectionalLight[MAX_DIRECTIONAL_LIGHTS] directionalLights;
uniform SpotLight[MAX_SPOT_LIGHTS] spotLights;

void calculatePointLight(PointLight pl);
void calculateDirectionalLight(DirectionalLight pl);
void calculateSpotLight(SpotLight pl);

void calculatePointLights();
void calculateDirectionalLights();
void calculateSpotLights();

void calculateAllLights();
void calculateColors();

float calculateShinyness(float shinyness);
void setVariables();

vec3 finalDiffuseColor;
vec3 finalSpecularColor;
vec3 finalAmbientColor;
vec3 norm;
vec3 viewDir;

const float levels = 4.0;

void main()
{	
	finalDiffuseColor = vec3(0.0,0.0,0.0);
	finalSpecularColor = vec3(0.0,0.0,0.0);
	finalAmbientColor = ambientLight;
	setVariables();

	calculateAllLights();
	calculateColors();

	// float brightness = (fragColor.r+fragColor.b+fragColor.g)/3.0;
	// float level = floor(brightness * levels);
	// brightness = level / levels;
	// fragColor = vec4(fragColor.rgb*brightness,fragColor.a);
}

void calculateAllLights()
{
	calculatePointLights();
	calculateDirectionalLights();
	calculateSpotLights();
}

vec3 getDiffuseTexture()
{
	if(material.diffuseTextureLoaded)
	{
		// vec3 tex = vec3(texture(material.diffuseTexture,FragIn.fragTexCoord));
		// float brightness = (tex.r+tex.b+tex.g)/3.0;
		// float level = floor(brightness * levels);
		// brightness = level / levels;
		// return tex * brightness;
		return vec3(texture(material.diffuseTexture,FragIn.fragTexCoord));
	}
	else
	{
		return vec3(1.0,1.0,1.0);
	}
}

vec3 getSpecularTexture()
{
	if(material.specularTextureLoaded)
	{
		// vec3 tex = vec3(texture(material.specularTexture,FragIn.fragTexCoord));
		// float brightness = (tex.r+tex.b+tex.g)/3.0;
		// float level = floor(brightness * levels);
		// brightness = level / levels;
		// return tex * brightness;
		return vec3(texture(material.specularTexture,FragIn.fragTexCoord));
	}
	else
	{
		return vec3(1.0,1.0,1.0);
	}
}

void calculateColors()
{
	finalDiffuseColor *= material.diffuseColor * getDiffuseTexture();
	finalSpecularColor *= material.specularColor * getSpecularTexture();
	finalAmbientColor *= material.diffuseColor * getDiffuseTexture();

	fragColor = vec4(finalDiffuseColor + finalSpecularColor + finalAmbientColor,1.0);
}

void calculatePointLights()
{
	for(uint i = 0;i<numPointLights&&i<MAX_POINT_LIGHTS;i++)
	{
		calculatePointLight(pointLights[i]);
	}
}
void calculateDirectionalLights()
{
	for(uint i = 0;i<numDirectionalLights&&i<MAX_DIRECTIONAL_LIGHTS;i++)
	{
		calculateDirectionalLight(directionalLights[i]);
	}
}
void calculateSpotLights()
{
	for(uint i = 0;i<numSpotLights&&i<MAX_SPOT_LIGHTS;i++)
	{
		calculateSpotLight(spotLights[i]);
	}	
}

vec3 diffuseLighting(vec3 lightDir,vec3 diffuse)
{
	float diff = max(dot(norm,lightDir),0.0);
	// float brightness = diff;
	// float level = floor(brightness * levels);
	// brightness = level / levels;
	// diffuse *= diff * brightness;
	diffuse *= diff;
	return diffuse;
}

vec3 specularLighting(vec3 lightDir,vec3 specular)
{
	vec3 reflectDir = reflect(-lightDir, norm);
	vec3 halfwayDir = normalize(lightDir + viewDir);
	float spec = max(pow(max(dot(norm,halfwayDir),0.0),calculateShinyness(material.shinyness)),0.0);
	// float brightness = spec;
	// float level = floor(brightness * levels);
	// brightness = level / levels;
	// specular *= spec*brightness;
	specular *= spec;
	return specular;
}

float calcAttentuation(vec3 lightPosition,Attentuation attentuation)
{
	float distance = distance(lightPosition,FragIn.fragPos);
	float attent = 1.0/(attentuation.quadratic*distance*distance + attentuation.linear*distance + attentuation.constant);
	return attent;
}

void calculatePointLight(PointLight pl)
{
	vec3 lightPosition = (FragIn.viewMatrix3D*vec4(pl.position,1.0)).xyz;
	vec3 lightDir = normalize(FragIn.fragToTangentSpace*(lightPosition - FragIn.fragPos));


	// Diffuse
	vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);

	// Specular
	vec3 specular = specularLighting(lightDir,pl.specularColor);

	// Attentuation
	float attent = calcAttentuation(lightPosition,pl.attentuation);

	diffuse *= attent;
	specular *= attent;

	finalDiffuseColor += diffuse;
	finalSpecularColor += specular;

}
void calculateDirectionalLight(DirectionalLight pl)
{
	vec3 lightDirection = (FragIn.viewMatrix3D*vec4(pl.direction*-1.0,0.0)).xyz;
	vec3 lightDir = normalize(FragIn.fragToTangentSpace*lightDirection);


	// Diffuse
	vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);

	// Specular
	vec3 specular = specularLighting(lightDir,pl.specularColor);

	finalDiffuseColor += diffuse;
	finalSpecularColor += specular;
}

float degToRad(float deg)
{
	return deg / 180.0 * 3.14159265359;
}

float calcSpotAmount(vec3 lightDir,vec3 lightDirection,SpotLight pl)
{
	float theta = dot(lightDir, normalize(FragIn.fragToTangentSpace*lightDirection));
	float spotAmount = 0.0;
	float outerCutOff = cos(degToRad(pl.outerCutOff));
	float innerCutOff = cos(degToRad(pl.innerCutOff));
	float epsilon   = innerCutOff - outerCutOff;
	spotAmount = clamp((theta - outerCutOff) / epsilon, 0.0, 1.0);

	return spotAmount;
}

void calculateSpotLight(SpotLight pl)
{
	vec3 lightPosition = (FragIn.viewMatrix3D*vec4(pl.position,1.0)).xyz;
	vec3 lightDirection = (FragIn.viewMatrix3D*vec4(pl.direction*-1.0,0.0)).xyz;
	vec3 lightDir = normalize(FragIn.fragToTangentSpace*(lightPosition-FragIn.fragPos));

	// Spotamount
	float spotAmount = calcSpotAmount(lightDir,lightDirection,pl);

	// Diffuse
	vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);

	// Specular
	vec3 specular = specularLighting(lightDir,pl.specularColor);

	// Attentuation
	float attent = calcAttentuation(lightPosition,pl.attentuation);

	diffuse *= attent * spotAmount;
	specular *= attent * spotAmount;

	finalDiffuseColor += diffuse;
	finalSpecularColor += specular;
}

float calculateShinyness(float shinyness)
{
	return max(MAX_SPECULAR_EXPONENT*(pow(max(shinyness,0.0),-3.0)-1.0)+MIN_SPECULAR_EXPONENT,0.0);
}

void setVariables()
{
	if(material.normalMapLoaded)
	{
		norm = normalize(2.0*(texture(material.normalMap,FragIn.fragTexCoord)).xyz-1.0);
	}
	else
	{
		norm = normalize(FragIn.fragToTangentSpace*FragIn.fragNormal);
	}
	viewDir = normalize((FragIn.fragToTangentSpace*(FragIn.fragPos*-1.0)));
}

