#version 100

precision mediump float;
precision mediump sampler2D;
precision mediump samplerCube;

#define MAX_POINT_LIGHTS 20
#define MAX_DIRECTIONAL_LIGHTS 2
#define MAX_SPOT_LIGHTS 1

#define MAX_SPECULAR_EXPONENT 50.0
#define MIN_SPECULAR_EXPONENT 5.0

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

	mat4 lightSpaceMatrix[6];
	samplerCube shadowmap;
	bool castsShadows;

	float farPlane;
};

struct DirectionalLight
{
	vec3 direction;

	vec3 diffuseColor;
	vec3 specularColor;

	mat4 lightSpaceMatrix;
	sampler2D shadowmap;
	bool castsShadows;

	float shadowDistance;
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

	mat4 lightSpaceMatrix;
	sampler2D shadowmap;
	bool castsShadows;
};



const float shadowDistance = 50.0;
const float transitionDistance = 5.0;
const float bias = 0.005;

uniform int numPointLights;
uniform int numDirectionalLights;
uniform int numSpotLights;
uniform vec3 ambientLight;
uniform PointLight[MAX_POINT_LIGHTS] pointLights;
uniform DirectionalLight[MAX_DIRECTIONAL_LIGHTS] directionalLights;
uniform SpotLight[MAX_SPOT_LIGHTS] spotLights;
uniform Material material;

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
vec3 norm;
vec3 viewDir;

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
	if(material.DiffuseTextureLoaded)
	{
		return texture2D(material.diffuseTexture,fragTexCoord);
	}
	else
	{
		return vec4(1.0,1.0,1.0,1.0);
	}
}

vec4 getSpecularTexture()
{
	if(material.SpecularTextureLoaded)
	{
		return texture2D(material.specularTexture,fragTexCoord);
	}
	else
	{
		return vec4(1.0,1.0,1.0,1.0);
	}
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

	if(int(numPointLights) > 0)
		calculatePointLight(pointLights[0]);
	if(int(numPointLights) > 1)
		calculatePointLight(pointLights[1]);
	if(int(numPointLights) > 2)
		calculatePointLight(pointLights[2]);
	if(int(numPointLights) > 3)
		calculatePointLight(pointLights[3]);
	if(int(numPointLights) > 4)
		calculatePointLight(pointLights[4]);
	if(int(numPointLights) > 5)
		calculatePointLight(pointLights[5]);
	if(int(numPointLights) > 6)
		calculatePointLight(pointLights[6]);
	if(int(numPointLights) > 7)
		calculatePointLight(pointLights[7]);
	if(int(numPointLights) > 8)
		calculatePointLight(pointLights[8]);
	if(int(numPointLights) > 9)
		calculatePointLight(pointLights[9]);
	if(int(numPointLights) > 10)
		calculatePointLight(pointLights[10]);
	if(int(numPointLights) > 11)
		calculatePointLight(pointLights[11]);
	if(int(numPointLights) > 12)
		calculatePointLight(pointLights[12]);
	if(int(numPointLights) > 13)
		calculatePointLight(pointLights[13]);
	if(int(numPointLights) > 14)
		calculatePointLight(pointLights[14]);
	if(int(numPointLights) > 15)
		calculatePointLight(pointLights[15]);
	if(int(numPointLights) > 16)
		calculatePointLight(pointLights[16]);
	if(int(numPointLights) > 17)
		calculatePointLight(pointLights[17]);
	if(int(numPointLights) > 18)
		calculatePointLight(pointLights[18]);
	if(int(numPointLights) > 19)
		calculatePointLight(pointLights[19]);
	// if(int(numPointLights) > 20)
	// 	calculatePointLight(pointLights[20]);
}
void calculateDirectionalLights()
{	
	if(int(numDirectionalLights) > 0)
		calculateDirectionalLight(directionalLights[0]);
	if(int(numDirectionalLights) > 1)
		calculateDirectionalLight(directionalLights[1]);
}
void calculateSpotLights()
{
	uint i;
	for(i=uint(0);i<numSpotLights&&i<uint(MAX_SPOT_LIGHTS);i++)
	{
		calculateSpotLight(spotLights[i]);
	}	
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

float calcShadow(sampler2D shadowMap,mat4 lightSpaceMatrix,float shadowdistance,bool distanceTransition)
{	
	float distance = 0.0;
	if(distanceTransition)
	{
		distance = length(fragPos);
		distance = distance - (shadowdistance - transitionDistance);
		distance = distance / transitionDistance;
		distance = clamp(1.0-distance,0.0,1.0);
	}
	vec4 fragPosLightSpace = lightSpaceMatrix*inverse(fragViewMatrix3D)*vec4(fragPos,1.0);
	vec3 projCoords = clamp((fragPosLightSpace.xyz / fragPosLightSpace.w)*0.5+0.5,-1.0,1.0);
	float currentDepth = projCoords.z-bias;
	float shadowresult = 0.0;
	float closestDepth = texture2D(shadowMap, projCoords.xy).r;
	vec2 texelSize = 1.0 / vec2(textureSize(shadowMap, 0));
	for(int x = -1; x <= 1; ++x)
	{
	    for(int y = -1; y <= 1; ++y)
	    {
	        float pcfDepth = texture(shadowMap, projCoords.xy + vec2(x, y) * texelSize).r; 
	        shadowresult += currentDepth > pcfDepth ? 0.0 : 1.0;        
	    }    
	}
	shadowresult /= 9.0;
	if(distanceTransition)
	{
		shadowresult = 1.0 - (1.0-shadowresult)*distance;
	}
	return shadowresult;
}

vec3 sampleOffsetDirections[20] = vec3[]
(
   vec3( 1,  1,  1), vec3( 1, -1,  1), vec3(-1, -1,  1), vec3(-1,  1,  1), 
   vec3( 1,  1, -1), vec3( 1, -1, -1), vec3(-1, -1, -1), vec3(-1,  1, -1),
   vec3( 1,  1,  0), vec3( 1, -1,  0), vec3(-1, -1,  0), vec3(-1,  1,  0),
   vec3( 1,  0,  1), vec3(-1,  0,  1), vec3( 1,  0, -1), vec3(-1,  0, -1),
   vec3( 0,  1,  1), vec3( 0, -1,  1), vec3( 0, -1, -1), vec3( 0,  1, -1)
);   

float calcShadowPointLight(PointLight pl)
{
	vec3 fragToLight = (inverse(fragViewMatrix3D)*vec4(fragPos,1.0)).xyz - pl.position;
	float currentDepth = length(fragToLight)-bias*10.0*pl.farPlane;
	float shadow  = 0.0;
	int samples = 20;
	float viewDistance = length(-fragPos);
	float diskRadius = (1.0 + (viewDistance / pl.farPlane)) / 70.0;
	for(int i = 0;i<samples;i++) {
		 float closestDepth = textureCube(pl.shadowmap, fragToLight + sampleOffsetDirections[i]*diskRadius).r;
	            closestDepth *= pl.farPlane;   // Undo mapping [0;1]
	            if(currentDepth <= closestDepth)
	                shadow += 1.0;
	}
	shadow /= float(samples);
	return shadow;
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

	// Shadow
	float shadow = pl.castsShadows ? calcShadowPointLight(pl) : 1.0;

	diffuse *= attent * shadow;
	specular *= attent * shadow;

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
	
	// Shadow
	float shadow = dl.castsShadows ? calcShadow(dl.shadowmap,dl.lightSpaceMatrix,dl.shadowDistance,true) : 1.0;
	
	diffuse *= shadow;
	specular *= shadow;

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

	// Shadow
	float shadow = pl.castsShadows ? calcShadow(pl.shadowmap,pl.lightSpaceMatrix,50.0,false) : 1.0;
	// float shadow = 1.0;

	diffuse *= attent * spotAmount * shadow;
	specular *= attent * spotAmount * shadow;

	finalDiffuseColor += vec4(diffuse,0.0);
	finalSpecularColor += vec4(specular,0.0);
}

float calculateShinyness(float shinyness)
{
	return max(MAX_SPECULAR_EXPONENT*(pow(max(shinyness,0.0),-3.0)-1.0)+MIN_SPECULAR_EXPONENT,0.0);
}

void setVariables()
{
	if(material.NormalMapLoaded)
	{
		norm = normalize(2.0*(texture2D(material.normalMap,fragTexCoord)).xyz-1.0);
	}
	else
	{
		norm = normalize(fragToTangentSpace*fragNormal);
	}
	viewDir = normalize((fragToTangentSpace*(fragPos*-1.0)));
}

