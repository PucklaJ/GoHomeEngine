#version 100

precision mediump float;
precision mediump sampler2D;
precision mediump samplerCube;

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

const float shadowDistance = 50.0;
const float transitionDistance = 5.0;
const float bias = 0.005;

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

	mat4 lightSpaceMatrix[6];
	bool castsShadows;

	float farPlane;
} pointLights[MAX_POINT_LIGHTS];
uniform samplerCube pointLightsshadowmap[MAX_POINT_LIGHTS];
uniform struct DirectionalLight
{
	vec3 direction;

	vec3 diffuseColor;
	vec3 specularColor;

	mat4 lightSpaceMatrix;
	bool castsShadows;
	ivec2 shadowMapSize;

	float shadowDistance;
} directionalLights[MAX_DIRECTIONAL_LIGHTS];
uniform sampler2D directionalLightsshadowmap[MAX_DIRECTIONAL_LIGHTS];
uniform struct SpotLight
{
	vec3 position;
	vec3 direction;

	vec3 diffuseColor;
	vec3 specularColor;

	float innerCutOff;
	float outerCutOff;

	Attentuation attentuation;

	mat4 lightSpaceMatrix;
	bool castsShadows;
	ivec2 shadowMapSize;
} spotLights[MAX_SPOT_LIGHTS];
uniform sampler2D spotLightsshadowmap[MAX_SPOT_LIGHTS];
uniform struct Material
{
	vec3 diffuseColor;
	vec3 specularColor;

	bool DiffuseTextureLoaded;
	bool SpecularTextureLoaded;
	bool NormalMapLoaded;

	float shinyness;
} material;
uniform sampler2D materialdiffuseTexture;
uniform sampler2D materialspecularTexture;
uniform sampler2D materialnormalMap;

void calculatePointLight(PointLight pl,int index);
void calculateDirectionalLight(DirectionalLight pl,int index);
void calculateSpotLight(SpotLight pl,int index);

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
varying mat4 fragInverseViewMatrix3D;

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
	return mix(vec4(1.0,1.0,1.0,1.0),texture2D(materialspecularTexture,fragTexCoord),material.SpecularTextureLoaded ? 1.0 : 0.0);
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
	if(numPointLights > 0)
		calculatePointLight(pointLights[0],0);
	#endif
	#if MAX_POINT_LIGHTS > 1
	if(numPointLights > 1)
		calculatePointLight(pointLights[1],1);
	#endif
	#if MAX_POINT_LIGHTS > 2
	if(numPointLights > 2)
		calculatePointLight(pointLights[2],2);
	#endif
	#if MAX_POINT_LIGHTS > 3
	if(numPointLights > 3)
		calculatePointLight(pointLights[3],3);
	#endif
	#if MAX_POINT_LIGHTS > 4
	if(numPointLights > 4)
		calculatePointLight(pointLights[4],4);
	#endif
	#if MAX_POINT_LIGHTS > 5
	if(numPointLights > 5)
		calculatePointLight(pointLights[5],5);
	#endif
	#if MAX_POINT_LIGHTS > 6
	if(numPointLights > 6)
		calculatePointLight(pointLights[6],6);
	#endif
	#if MAX_POINT_LIGHTS > 7
	if(numPointLights > 7)
		calculatePointLight(pointLights[7],7);
	#endif
	#if MAX_POINT_LIGHTS > 8
	if(numPointLights > 8)
		calculatePointLight(pointLights[8],8);
	#endif
	#if MAX_POINT_LIGHTS > 9
	if(numPointLights > 9)
		calculatePointLight(pointLights[9],9);
	#endif
	#if MAX_POINT_LIGHTS > 10
	if(numPointLights > 10)
		calculatePointLight(pointLights[10],10);
	#endif
	#if MAX_POINT_LIGHTS > 11
	if(numPointLights > 11)
		calculatePointLight(pointLights[11],11);
	#endif
	#if MAX_POINT_LIGHTS > 12
	if(numPointLights > 12)
		calculatePointLight(pointLights[12],12);
	#endif
	#if MAX_POINT_LIGHTS > 13
	if(numPointLights > 13)
		calculatePointLight(pointLights[13],13);
	#endif
	#if MAX_POINT_LIGHTS > 14
	if(numPointLights > 14)
		calculatePointLight(pointLights[14],14);
	#endif
	#if MAX_POINT_LIGHTS > 15
	if(numPointLights > 15)
		calculatePointLight(pointLights[15],15);
	#endif
	#if MAX_POINT_LIGHTS > 16
	if(numPointLights > 16)
		calculatePointLight(pointLights[16],16);
	#endif
	#if MAX_POINT_LIGHTS > 17
	if(numPointLights > 17)
		calculatePointLight(pointLights[17],17);
	#endif
	#if MAX_POINT_LIGHTS > 18
	if(numPointLights > 18)
		calculatePointLight(pointLights[18],18);
	#endif
	#if MAX_POINT_LIGHTS > 19
	if(numPointLights > 19)
		calculatePointLight(pointLights[19],19);
	#endif
	#if MAX_POINT_LIGHTS > 20
	if(numPointLights > 20)
		calculatePointLight(pointLights[20],20);
	#endif
}
void calculateDirectionalLights()
{	
	#if MAX_DIRECTIONAL_LIGHTS > 0
	if(int(numDirectionalLights) > 0)
		calculateDirectionalLight(directionalLights[0],0);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 1
	if(int(numDirectionalLights) > 1)
		calculateDirectionalLight(directionalLights[1],1);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 2
	if(int(numDirectionalLights) > 2)
		calculateDirectionalLight(directionalLights[2],2);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 3
	if(int(numDirectionalLights) > 3)
		calculateDirectionalLight(directionalLights[3],3);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 4
	if(int(numDirectionalLights) > 4)
		calculateDirectionalLight(directionalLights[4],4);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 5
	if(int(numDirectionalLights) > 5)
		calculateDirectionalLight(directionalLights[5],5);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 6
	if(int(numDirectionalLights) > 6)
		calculateDirectionalLight(directionalLights[6],6);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 7
	if(int(numDirectionalLights) > 7)
		calculateDirectionalLight(directionalLights[7],7);
	#endif
	#if MAX_DIRECTIONAL_LIGHTS > 8
	if(int(numDirectionalLights) > 8)
		calculateDirectionalLight(directionalLights[8],8);
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
		calculateSpotLight(spotLights[0],0);
	#endif
	#if MAX_SPOT_LIGHTS > 1
	if(int(numSpotLights) > 1)
		calculateSpotLight(spotLights[1],1);
	#endif
	#if MAX_SPOT_LIGHTS > 2
	if(int(numSpotLights) > 2)
		calculateSpotLight(spotLights[2],2);
	#endif
	#if MAX_SPOT_LIGHTS > 3
	if(int(numSpotLights) > 3)
		calculateSpotLight(spotLights[3],3);
	#endif
	#if MAX_SPOT_LIGHTS > 4
	if(int(numSpotLights) > 4)
		calculateSpotLight(spotLights[4],4);
	#endif
	#if MAX_SPOT_LIGHTS > 5
	if(int(numSpotLights) > 5)
		calculateSpotLight(spotLights[5],5);
	#endif
	#if MAX_SPOT_LIGHTS > 6
	if(int(numSpotLights) > 6)
		calculateSpotLight(spotLights[6],6);
	#endif
	#if MAX_SPOT_LIGHTS > 7
	if(int(numSpotLights) > 7)
		calculateSpotLight(spotLights[7],7);
	#endif
	#if MAX_SPOT_LIGHTS > 8
	if(int(numSpotLights) > 8)
		calculateSpotLight(spotLights[8],8);
	#endif
	#if MAX_SPOT_LIGHTS > 9
	if(int(numSpotLights) > 9)
		calculateSpotLight(spotLights[9],9);
	#endif
	#if MAX_SPOT_LIGHTS > 10
	if(int(numSpotLights) > 10)
		calculateSpotLight(spotLights[10],10);
	#endif
	#if MAX_SPOT_LIGHTS > 11
	if(int(numSpotLights) > 11)
		calculateSpotLight(spotLights[11],11);
	#endif
	#if MAX_SPOT_LIGHTS > 12
	if(int(numSpotLights) > 12)
		calculateSpotLight(spotLights[12],12);
	#endif
	#if MAX_SPOT_LIGHTS > 13
	if(int(numSpotLights) > 13)
		calculateSpotLight(spotLights[13],13);
	#endif
	#if MAX_SPOT_LIGHTS > 14
	if(int(numSpotLights) > 14)
		calculateSpotLight(spotLights[14],14);
	#endif
	#if MAX_SPOT_LIGHTS > 15
	if(int(numSpotLights) > 15)
		calculateSpotLight(spotLights[15],15);
	#endif
	#if MAX_SPOT_LIGHTS > 16
	if(int(numSpotLights) > 16)
		calculateSpotLight(spotLights[16],16);
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

float calcShadow(sampler2D shadowMap,mat4 lightSpaceMatrix,float shadowdistance,bool distanceTransition,ivec2 shadowMapSize)
{	
	float distance = 0.0;
	if(distanceTransition)
	{
		distance = length(fragPos);
		distance = distance - (shadowdistance - transitionDistance);
		distance = distance / transitionDistance;
		distance = clamp(1.0-distance,0.0,1.0);
	}
	vec4 fragPosLightSpace = lightSpaceMatrix*fragInverseViewMatrix3D*vec4(fragPos,1.0);
	vec3 projCoords = clamp((fragPosLightSpace.xyz / fragPosLightSpace.w)*0.5+0.5,-1.0,1.0);
	float currentDepth = projCoords.z-bias;
	float shadowresult = 0.0;
	// vec2 texelSize = 1.0 / vec2(shadowMapSize);
	// for(int x = -1; x <= 1; ++x)
	// {
	//     for(int y = -1; y <= 1; ++y)
	//     {
	//         float pcfDepth = texture2D(shadowMap, projCoords.xy + vec2(x, y) * texelSize).r; 
	//         shadowresult += currentDepth > pcfDepth ? 0.0 : 1.0;        
	//     }    
	// }
	// shadowresult /= 9.0;
	float closestDepth = texture2D(shadowMap, projCoords.xy).r;
	shadowresult = currentDepth > closestDepth ? 0.0 : 1.0;
	if(distanceTransition)
	{
		shadowresult = 1.0 - (1.0-shadowresult)*distance;
	}
	return shadowresult;
}

vec3 sampleOffsetDirections[20];

void setOffsetDirections()
{
   sampleOffsetDirections[1] = vec3( 1,  1,  1); sampleOffsetDirections[2] = vec3( 1, -1,  1); sampleOffsetDirections[3] = vec3(-1, -1,  1); sampleOffsetDirections[4] = vec3(-1,  1,  1); 
   sampleOffsetDirections[5] = vec3( 1,  1, -1); sampleOffsetDirections[6] = vec3( 1, -1, -1); sampleOffsetDirections[7] = vec3(-1, -1, -1); sampleOffsetDirections[8] = vec3(-1,  1, -1);
   sampleOffsetDirections[9] = vec3( 1,  1,  0); sampleOffsetDirections[10] = vec3( 1, -1,  0); sampleOffsetDirections[11] = vec3(-1, -1,  0); sampleOffsetDirections[12] = vec3(-1,  1,  0);
   sampleOffsetDirections[13] = vec3( 1,  0,  1); sampleOffsetDirections[14] = vec3(-1,  0,  1); sampleOffsetDirections[15] = vec3( 1,  0, -1); sampleOffsetDirections[15] = vec3(-1,  0, -1);
   sampleOffsetDirections[17] = vec3( 0,  1,  1); sampleOffsetDirections[18] = vec3( 0, -1,  1); sampleOffsetDirections[19] = vec3( 0, -1, -1); sampleOffsetDirections[19] = vec3( 0,  1, -1);
}

float calcShadowPointLight(PointLight pl,samplerCube shadowmap)
{
	setOffsetDirections();

	vec3 fragToLight = (fragInverseViewMatrix3D*vec4(fragPos,1.0)).xyz - pl.position;
	float currentDepth = length(fragToLight)-bias*10.0*pl.farPlane;
	float shadow  = 0.0;
	// int samples = 20;
	// float viewDistance = length(-fragPos);
	// float diskRadius = (1.0 + (viewDistance / pl.farPlane)) / 70.0;
	// for(int i = 0;i<samples;i++) {
	// 	 float closestDepth = textureCube(shadowmap, fragToLight + sampleOffsetDirections[i]*diskRadius).r;
	//             closestDepth *= pl.farPlane;   // Undo mapping [0;1]
	//             if(currentDepth <= closestDepth)
	//                 shadow += 1.0;
	// }
	// shadow /= float(samples);
	float closestDepth = textureCube(shadowmap,fragToLight).r;
	closestDepth *= pl.farPlane;
	if(currentDepth <= closestDepth)
		shadow = 1.0;
	return shadow;
}

float calcAttentuation(vec3 lightPosition,Attentuation attentuation)
{
	float distance = distance(lightPosition,fragPos);
	float attent = 1.0/(attentuation.quadratic*distance*distance + attentuation.linear*distance + attentuation.constant);
	return attent;
}

void calculatePointLight(PointLight pl,int index)
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
	float shadow = pl.castsShadows ? calcShadowPointLight(pl,pointLightsshadowmap[index]) : 1.0;

	diffuse *= attent * shadow;
	specular *= attent * shadow;

	finalDiffuseColor += vec4(diffuse,0.0);
	finalSpecularColor += vec4(specular,0.0);

}
void calculateDirectionalLight(DirectionalLight dl,int index)
{
	vec3 lightDirection = (fragViewMatrix3D*vec4(dl.direction*-1.0,0.0)).xyz;
	vec3 lightDir = normalize(fragToTangentSpace*lightDirection);
	
	// Diffuse
	vec3 diffuse = diffuseLighting(lightDir,dl.diffuseColor);
	
	// Specular
	vec3 specular = specularLighting(lightDir,dl.specularColor);
	
	// Shadow
	float shadow = dl.castsShadows ? calcShadow(directionalLightsshadowmap[index],dl.lightSpaceMatrix,dl.shadowDistance,true,dl.shadowMapSize) : 1.0;
	
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

void calculateSpotLight(SpotLight pl,int index)
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
	float shadow = pl.castsShadows ? calcShadow(spotLightsshadowmap[index],pl.lightSpaceMatrix,50.0,false,pl.shadowMapSize) : 1.0;
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
	norm = mix(normalize(fragToTangentSpace*fragNormal),normalize(2.0*(texture2D(materialnormalMap,fragTexCoord)).xyz-1.0),material.NormalMapLoaded ? 1.0 : 0.0);
	viewDir = normalize((fragToTangentSpace*(fragPos*-1.0)));
}

