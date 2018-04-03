#version 410

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
uniform float shadowDistance = 50.0;
uniform float transitionDistance = 5.0;
uniform float bias = 0.005;

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
	if(material.diffuseTextureLoaded)
	{
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
		return texture(material.specularTexture,FragIn.fragTexCoord);
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

	fragColor = finalDiffuseColor + finalSpecularColor  + finalAmbientColor;
}

void calculatePointLights()
{
	// for(uint i = 0;i<numPointLights&&i<MAX_POINT_LIGHTS;i++)
	// {
	// 	calculatePointLight(pointLights[i]);
	// }

	if(numPointLights > 0)
		calculatePointLight(pointLights[0]);
	if(numPointLights > 1)
		calculatePointLight(pointLights[1]);
	if(numPointLights > 2)
		calculatePointLight(pointLights[2]);
	if(numPointLights > 3)
		calculatePointLight(pointLights[3]);
	if(numPointLights > 4)
		calculatePointLight(pointLights[4]);
	if(numPointLights > 5)
		calculatePointLight(pointLights[5]);
	if(numPointLights > 6)
		calculatePointLight(pointLights[6]);
	if(numPointLights > 7)
		calculatePointLight(pointLights[7]);
	if(numPointLights > 8)
		calculatePointLight(pointLights[8]);
	if(numPointLights > 9)
		calculatePointLight(pointLights[9]);
	if(numPointLights > 10)
		calculatePointLight(pointLights[10]);
	if(numPointLights > 11)
		calculatePointLight(pointLights[11]);
	if(numPointLights > 12)
		calculatePointLight(pointLights[12]);
	if(numPointLights > 13)
		calculatePointLight(pointLights[13]);
	if(numPointLights > 14)
		calculatePointLight(pointLights[14]);
	if(numPointLights > 15)
		calculatePointLight(pointLights[15]);
	if(numPointLights > 16)
		calculatePointLight(pointLights[16]);
	if(numPointLights > 17)
		calculatePointLight(pointLights[17]);
	if(numPointLights > 18)
		calculatePointLight(pointLights[18]);
	if(numPointLights > 19)
		calculatePointLight(pointLights[19]);
	// if(numPointLights > 20)
	// 	calculatePointLight(pointLights[20]);
}
void calculateDirectionalLights()
{	
	if(numDirectionalLights > 0)
		calculateDirectionalLight(directionalLights[0]);
	if(numDirectionalLights > 1)
		calculateDirectionalLight(directionalLights[1]);
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
		distance = length(FragIn.fragPos);
		distance = distance - (shadowdistance - transitionDistance);
		distance = distance / transitionDistance;
		distance = clamp(1.0-distance,0.0,1.0);
	}
	vec4 fragPosLightSpace = lightSpaceMatrix*inverse(FragIn.viewMatrix3D)*vec4(FragIn.fragPos,1.0);
	vec3 projCoords = clamp((fragPosLightSpace.xyz / fragPosLightSpace.w)*0.5+0.5,-1.0,1.0);
	float currentDepth = projCoords.z-bias;
	float shadowresult = 0.0;
	float closestDepth = texture(shadowMap, projCoords.xy).r;
	vec2 texelSize = 1.0 / textureSize(shadowMap, 0);
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
	vec3 fragToLight = (inverse(FragIn.viewMatrix3D)*vec4(FragIn.fragPos,1.0)).xyz - pl.position;
	float currentDepth = length(fragToLight)-bias*10.0*pl.farPlane;
	float shadow  = 0.0;
	int samples = 20;
	float viewDistance = length(-FragIn.fragPos);
	float diskRadius = (1.0 + (viewDistance / pl.farPlane)) / 70.0;
	for(int i = 0;i<samples;i++) {
		 float closestDepth = texture(pl.shadowmap, fragToLight + sampleOffsetDirections[i]*diskRadius).r;
	            closestDepth *= pl.farPlane;   // Undo mapping [0;1]
	            if(currentDepth <= closestDepth)
	                shadow += 1.0;
	}
	shadow /= float(samples);
	return shadow;
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

	// Shadow
	float shadow = pl.castsShadows ? calcShadowPointLight(pl) : 1.0;

	diffuse *= attent * shadow;
	specular *= attent * shadow;

	finalDiffuseColor += vec4(diffuse,0.0);
	finalSpecularColor += vec4(specular,0.0);

}
void calculateDirectionalLight(DirectionalLight dl)
{
	vec3 lightDirection = (FragIn.viewMatrix3D*vec4(dl.direction*-1.0,0.0)).xyz;
	vec3 lightDir = normalize(FragIn.fragToTangentSpace*lightDirection);
	
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

	// Shadow
	float shadow = pl.castsShadows ? calcShadow(pl.shadowmap,pl.lightSpaceMatrix,50.0,false) : 1.0;

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

