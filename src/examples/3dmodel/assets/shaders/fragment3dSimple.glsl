#version 100

precision mediump float;
precision mediump sampler2D;

#define MAX_POINT_LIGHTS 5
#define MAX_DIRECTIONAL_LIGHTS 2
#define MAX_SPOT_LIGHTS 1

#define MAX_SPECULAR_EXPONENT 50.0
#define MIN_SPECULAR_EXPONENT 5.0

struct Attentuation
{
	float constant;
	float linear;
	float quadratic;
};

uniform int numPointLights;
uniform int numDirectionalLights;
uniform int numSpotLights;
uniform vec3 ambientLight;
uniform struct PointLight
{
	vec3 position;

	vec3 diffuseColor;
	vec3 specularColor;

	Attentuation attentuation;
} pointLights[MAX_POINT_LIGHTS];
uniform struct DirectionalLight
{
	vec3 direction;

	vec3 diffuseColor;
	vec3 specularColor;
} directionalLights[MAX_DIRECTIONAL_LIGHTS];
uniform struct SpotLight
{
	vec3 position;
	vec3 direction;

	vec3 diffuseColor;
	vec3 specularColor;

	float innerCutOff;
	float outerCutOff;

	Attentuation attentuation;
} spotLights[MAX_SPOT_LIGHTS];
uniform struct Material
{
	vec3 diffuseColor;
	vec3 specularColor;

	bool DiffuseTextureLoaded;

	float shinyness;
} material;
uniform sampler2D materialdiffuseTexture;

void calculatePointLight(PointLight pl);
void calculateDirectionalLight(DirectionalLight pl);
void calculateSpotLight(SpotLight pl);

void calculatePointLights();
void calculateDirectionalLights();
void calculateSpotLights();

void calculateAllLights();
void calculateLightColors();

float calculateShinyness(float shinyness);
void setVariables();

vec4 getDiffuseTexture();
vec4 getSpecularTexture();

vec4 finalDiffuseColor;
vec4 finalSpecularColor;
vec4 finalAmbientColor;
vec3 norm = vec3(0.0,0.0,0.0);
vec3 viewDir = vec3(0.0,0.0,0.0);

varying vec2 fragTexCoord;
varying vec3 fragPos;
varying vec3 fragNormal;
varying mat3 fragToTangentSpace;
varying mat4 fragViewMatrix3D;

void main()
{	
	finalDiffuseColor = vec4(0.0,0.0,0.0,0.0);
	finalSpecularColor = vec4(0.0,0.0,0.0,0.0);
	finalAmbientColor = vec4(ambientLight,1.0);
	setVariables();

	calculateAllLights();
	calculateLightColors();
}

void calculateAllLights()
{
	calculatePointLights();
	calculateDirectionalLights();
	calculateSpotLights();
}

vec4 getDiffuseTexture()
{
	return mix(vec4(1.0,1.0,1.0,1.0),texture2D(materialdiffuseTexture,fragTexCoord),material.DiffuseTextureLoaded ? 1.0 : 0.0);
}

vec4 getSpecularTexture()
{
	return vec4(1.0,1.0,1.0,1.0);
}

void calculateLightColors()
{	
	vec4 texDifCol = getDiffuseTexture();
	vec4 texSpecCol = getSpecularTexture();

	if(texDifCol.a < 0.1)
	{
		discard;
	}

	finalDiffuseColor *= vec4(material.diffuseColor,1.0) * texDifCol;
	finalSpecularColor *= vec4(material.specularColor,1.0) * texSpecCol;
	finalAmbientColor *= vec4(material.diffuseColor,1.0) * texDifCol;

	gl_FragColor = finalDiffuseColor + finalSpecularColor  + finalAmbientColor;
}

void calculatePointLights()
{
	// for(uint i = 0;i<numPointLights&&i<MAX_POINT_LIGHTS;i++)
	// {
	// 	calculatePointLight(pointLights[i]);
	// }

	#if MAX_POINT_LIGHTS > 0
	if(int(numPointLights) > 0)
		calculatePointLight(pointLights[0]);
	#endif
	#if MAX_POINT_LIGHTS > 1
	if(int(numPointLights) > 1)
		calculatePointLight(pointLights[1]);
	#endif
	#if MAX_POINT_LIGHTS > 2
	if(int(numPointLights) > 2)
		calculatePointLight(pointLights[2]);
	#endif
	#if MAX_POINT_LIGHTS > 3
	if(int(numPointLights) > 3)
		calculatePointLight(pointLights[3]);
	#endif
	#if MAX_POINT_LIGHTS > 4
	if(int(numPointLights) > 4)
		calculatePointLight(pointLights[4]);
	#endif
	#if MAX_POINT_LIGHTS > 5
	if(int(numPointLights) > 5)
		calculatePointLight(pointLights[5]);
	#endif
	#if MAX_POINT_LIGHTS > 6
	if(int(numPointLights) > 6)
		calculatePointLight(pointLights[6]);
	#endif
	#if MAX_POINT_LIGHTS > 7
	if(int(numPointLights) > 7)
		calculatePointLight(pointLights[7]);
	#endif
	#if MAX_POINT_LIGHTS > 8
	if(int(numPointLights) > 8)
		calculatePointLight(pointLights[8]);
	#endif
	#if MAX_POINT_LIGHTS > 9
	if(int(numPointLights) > 9)
		calculatePointLight(pointLights[9]);
	#endif
	#if MAX_POINT_LIGHTS > 10
	if(int(numPointLights) > 10)
		calculatePointLight(pointLights[10]);
	#endif
	#if MAX_POINT_LIGHTS > 11
	if(int(numPointLights) > 11)
		calculatePointLight(pointLights[11]);
	#endif
	#if MAX_POINT_LIGHTS > 12
	if(int(numPointLights) > 12)
		calculatePointLight(pointLights[12]);
	#endif
	#if MAX_POINT_LIGHTS > 13
	if(int(numPointLights) > 13)
		calculatePointLight(pointLights[13]);
	#endif
	#if MAX_POINT_LIGHTS > 14
	if(int(numPointLights) > 14)
		calculatePointLight(pointLights[14]);
	#endif
	#if MAX_POINT_LIGHTS > 15
	if(int(numPointLights) > 15)
		calculatePointLight(pointLights[15]);
	#endif
	#if MAX_POINT_LIGHTS > 16
	if(int(numPointLights) > 16)
		calculatePointLight(pointLights[16]);
	#endif
	#if MAX_POINT_LIGHTS > 17
	if(int(numPointLights) > 17)
		calculatePointLight(pointLights[17]);
	#endif
	#if MAX_POINT_LIGHTS > 18
	if(int(numPointLights) > 18)
		calculatePointLight(pointLights[18]);
	#endif
	#if MAX_POINT_LIGHTS > 19
	if(int(numPointLights) > 19)
		calculatePointLight(pointLights[19]);
	#endif
	#if MAX_POINT_LIGHTS > 20
	if(int(numPointLights) > 20)
		calculatePointLight(pointLights[20]);
	#endif
}
void calculateDirectionalLights()
{	
	#if MAX_DIRECTIONAL_LIGHTS > 0
	if(int(numDirectionalLights) > 0)
		calculateDirectionalLight(directionalLights[0]);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 1
	if(int(numDirectionalLights) > 1)
		calculateDirectionalLight(directionalLights[1]);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 2
	if(int(numDirectionalLights) > 2)
		calculateDirectionalLight(directionalLights[2]);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 3
	if(int(numDirectionalLights) > 3)
		calculateDirectionalLight(directionalLights[3]);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 4
	if(int(numDirectionalLights) > 4)
		calculateDirectionalLight(directionalLights[4]);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 5
	if(int(numDirectionalLights) > 5)
		calculateDirectionalLight(directionalLights[5]);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 6
	if(int(numDirectionalLights) > 6)
		calculateDirectionalLight(directionalLights[6]);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 7
	if(int(numDirectionalLights) > 7)
		calculateDirectionalLight(directionalLights[7]);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 8
	if(int(numDirectionalLights) > 8)
		calculateDirectionalLight(directionalLights[8]);
	#endif
}
void calculateSpotLights()
{
	// for(int i=0; i<numSpotLights && i<MAX_SPOT_LIGHTS ; i++)
	// {
	// 	calculateSpotLight(spotLights[i]);
	// }
	#if MAX_SPOT_LIGHTS > 0
	if(int(numSpotLights) > 0)
		calculateSpotLight(spotLights[0]);
	#endif
	#if MAX_SPOT_LIGHTS > 1
	if(int(numSpotLights) > 1)
		calculateSpotLight(spotLights[1]);
	#endif
	#if MAX_SPOT_LIGHTS > 2
	if(int(numSpotLights) > 2)
		calculateSpotLight(spotLights[2]);
	#endif
	#if MAX_SPOT_LIGHTS > 3
	if(int(numSpotLights) > 3)
		calculateSpotLight(spotLights[3]);
	#endif
	#if MAX_SPOT_LIGHTS > 4
	if(int(numSpotLights) > 4)
		calculateSpotLight(spotLights[4]);
	#endif
	#if MAX_SPOT_LIGHTS > 5
	if(int(numSpotLights) > 5)
		calculateSpotLight(spotLights[5]);
	#endif
	#if MAX_SPOT_LIGHTS > 6
	if(int(numSpotLights) > 6)
		calculateSpotLight(spotLights[6]);
	#endif
	#if MAX_SPOT_LIGHTS > 7
	if(int(numSpotLights) > 7)
		calculateSpotLight(spotLights[7]);
	#endif
	#if MAX_SPOT_LIGHTS > 8
	if(int(numSpotLights) > 8)
		calculateSpotLight(spotLights[8]);
	#endif
	#if MAX_SPOT_LIGHTS > 9
	if(int(numSpotLights) > 9)
		calculateSpotLight(spotLights[9]);
	#endif
	#if MAX_SPOT_LIGHTS > 10
	if(int(numSpotLights) > 10)
		calculateSpotLight(spotLights[10]);
	#endif
	#if MAX_SPOT_LIGHTS > 11
	if(int(numSpotLights) > 11)
		calculateSpotLight(spotLights[11]);
	#endif
	#if MAX_SPOT_LIGHTS > 12
	if(int(numSpotLights) > 12)
		calculateSpotLight(spotLights[12]);
	#endif
	#if MAX_SPOT_LIGHTS > 13
	if(int(numSpotLights) > 13)
		calculateSpotLight(spotLights[13]);
	#endif
	#if MAX_SPOT_LIGHTS > 14
	if(int(numSpotLights) > 14)
		calculateSpotLight(spotLights[14]);
	#endif
	#if MAX_SPOT_LIGHTS > 15
	if(int(numSpotLights) > 15)
		calculateSpotLight(spotLights[15]);
	#endif
	#if MAX_SPOT_LIGHTS > 16
	if(int(numSpotLights) > 16)
		calculateSpotLight(spotLights[16]);
	#endif
}

vec3 diffuseLighting(vec3 lightDir,vec3 diffuse)
{
	float diff = max(dot(norm,lightDir),0.0);
	diffuse *= diff;
	return diffuse;
}

vec3 specularLighting(vec3 lightDir,vec3 specular)
{
	vec3 reflectDir = reflect(-lightDir, norm);
	vec3 halfwayDir = normalize(lightDir + viewDir);
	float spec = max(pow(max(dot(norm,halfwayDir),0.0),calculateShinyness(material.shinyness)),0.0);
	specular *= spec;
	return specular;
}

float calcAttentuation(vec3 lightPosition,Attentuation attentuation)
{
	float distance = distance(lightPosition,fragPos);
	float attent = 1.0/(attentuation.quadratic*distance*distance + attentuation.linear*distance + attentuation.constant);
	return attent;
}

void calculatePointLight(PointLight pl)
{
	vec3 lightPosition = (fragViewMatrix3D*vec4(pl.position,1.0)).xyz;
	vec3 lightDir = normalize(fragToTangentSpace*(lightPosition - fragPos));


	// Diffuse
	vec3 diffuse = diffuseLighting(lightDir,pl.diffuseColor);

	// Specular
	vec3 specular = specularLighting(lightDir,pl.specularColor);

	// Attentuation
	float attent = calcAttentuation(lightPosition,pl.attentuation);

	diffuse *= attent;
	specular *= attent;

	finalDiffuseColor += vec4(diffuse,0.0);
	finalSpecularColor += vec4(specular,0.0);

}
void calculateDirectionalLight(DirectionalLight dl)
{
	vec3 lightDirection = (fragViewMatrix3D*vec4(dl.direction*-1.0,0.0)).xyz;
	vec3 lightDir = normalize(fragToTangentSpace*lightDirection);
	
	// Diffuse
	vec3 diffuse = diffuseLighting(lightDir,dl.diffuseColor);
	
	// Specular
	vec3 specular = specularLighting(lightDir,dl.specularColor);

	finalDiffuseColor += vec4(diffuse,0.0);
	finalSpecularColor += vec4(specular,0.0);
}

float degToRad(float deg)
{
	return deg / 180.0 * 3.14159265359;
}

float calcSpotAmount(vec3 lightDir,vec3 lightDirection,SpotLight pl)
{
	float theta = dot(lightDir, normalize(fragToTangentSpace*lightDirection));
	float spotAmount = 0.0;
	float outerCutOff = cos(degToRad(pl.outerCutOff));
	float innerCutOff = cos(degToRad(pl.innerCutOff));
	float epsilon   = innerCutOff - outerCutOff;
	spotAmount = clamp((theta - outerCutOff) / epsilon, 0.0, 1.0);

	return spotAmount;
}

void calculateSpotLight(SpotLight pl)
{
	vec3 lightPosition = (fragViewMatrix3D*vec4(pl.position,1.0)).xyz;
	vec3 lightDirection = (fragViewMatrix3D*vec4(pl.direction*-1.0,0.0)).xyz;
	vec3 lightDir = normalize(fragToTangentSpace*(lightPosition-fragPos));

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

	finalDiffuseColor += vec4(diffuse,0.0);
	finalSpecularColor += vec4(specular,0.0);
}

float calculateShinyness(float shinyness)
{
	return max(MAX_SPECULAR_EXPONENT*(pow(max(shinyness,0.0),-3.0)-1.0)+MIN_SPECULAR_EXPONENT,0.0);
}

void setVariables()
{
	norm = normalize(fragToTangentSpace*fragNormal);
	viewDir = normalize((fragToTangentSpace*(fragPos*-1.0)));
}

