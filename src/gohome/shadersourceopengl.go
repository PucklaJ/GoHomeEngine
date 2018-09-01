package gohome

// BackBuffer Shader
const (
	BACKBUFFER_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 150

out vec2 fragTexCoords;

uniform float depth;

vec2 vertices[6];
vec2 texCoords[6];

void setValues()
{
	vertices[0] = vec2(-1.0,-1.0);
	vertices[1] = vec2(1.0,-1.0);
	vertices[2] = vec2(1.0,1.0);
	vertices[3] = vec2(1.0,1.0);
	vertices[4] = vec2(-1.0,1.0);
	vertices[5] = vec2(-1.0,-1.0);

	texCoords[0] = vec2(0.0,0.0);
	texCoords[1] = vec2(1.0,0.0);
	texCoords[2] = vec2(1.0,1.0);
	texCoords[3] = vec2(1.0,1.0);
	texCoords[4] = vec2(0.0,1.0);
	texCoords[5] = vec2(0.0,0.0);
}

void main()
{
	setValues();
	fragTexCoords = texCoords[gl_VertexID];
	gl_Position = vec4(vertices[gl_VertexID],depth,1.0);
}`
	BACKBUFFER_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 150

in vec2 fragTexCoords;

out vec4 fragColor;

uniform sampler2DMS BackBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);
	ivec2 texCoords = ivec2(fragTexCoords * textureSize(BackBuffer));

	for(int i = 0;i<8;i++)
	{
		color += texelFetch(BackBuffer,texCoords,i);
	}
	color /= 8.0;

	return color;
}

void main()
{
	fragColor = fetchColor();
	if(fragColor.a < 0.1)
		discard;
}`
	BACKBUFFER_NOMS_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec2 vertex;
attribute vec2 texCoord;

varying vec2 fragTexCoords;

uniform float depth;

void main()
{
	fragTexCoords = texCoord;
	gl_Position = vec4(vertex,depth,1.0);
}`
	BACKBUFFER_NOMS_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

varying vec2 fragTexCoords;

uniform sampler2D BackBuffer;

vec4 fetchColor()
{
	return texture2D(BackBuffer,fragTexCoords);
}

void main()
{
	gl_FragColor = fetchColor();
	if(gl_FragColor.a < 0.1)
		discard;
}`
)

// PostProcessingShader
const (
	POST_PROCESSING_SHADER_VERTEX_SOURCE_OPENGL = `
#version 150

out vec2 fragTexCoords;

vec2 vertices[6];
vec2 texCoords[6];

void setValues()
{
	vertices[0] = vec2(-1.0,-1.0);
	vertices[1] = vec2(1.0,-1.0);
	vertices[2] = vec2(1.0,1.0);
	vertices[3] = vec2(1.0,1.0);
	vertices[4] = vec2(-1.0,1.0);
	vertices[5] = vec2(-1.0,-1.0);

	texCoords[0] = vec2(0.0,0.0);
	texCoords[1] = vec2(1.0,0.0);
	texCoords[2] = vec2(1.0,1.0);
	texCoords[3] = vec2(1.0,1.0);
	texCoords[4] = vec2(0.0,1.0);
	texCoords[5] = vec2(0.0,0.0);
}

void main()
{
	setValues();
	fragTexCoords = texCoords[gl_VertexID];
	gl_Position = vec4(vertices[gl_VertexID],0.0,1.0);
}`
	POST_PROCESSING_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 150

in vec2 fragTexCoords;

out vec4 fragColor;

uniform sampler2DMS BackBuffer;

const float offset = 1.0 / 300.0;
vec2 offsets[9] = vec2[](
        vec2(-offset,  offset), // top-left
        vec2( 0.0f,    offset), // top-center
        vec2( offset,  offset), // top-right
        vec2(-offset,  0.0f),   // center-left
        vec2( 0.0f,    0.0f),   // center-center
        vec2( offset,  0.0f),   // center-right
        vec2(-offset, -offset), // bottom-left
        vec2( 0.0f,   -offset), // bottom-center
        vec2( offset, -offset)  // bottom-right    
    );

float blurKernel[9] = float[](
        1.0/16.0, 2.0/16.0, 1.0/16.0,
        2.0/16.0,  4.0/16.0, 2.0/16.0,
        1.0/16.0, 2.0/16.0, 1.0/16.0
    );

float normalKernel[9] = float[](
	0.0,0.0,0.0,
	0.0,1.0,0.0,
	0.0,0.0,0.0
	);

vec4 fetchColor(vec2 texCoord)
{
	vec4 color = vec4(0.0);
	ivec2 texCoords = ivec2(texCoord * textureSize(BackBuffer));

	for(int i = 0;i<8;i++)
	{
		color += texelFetch(BackBuffer,texCoords,i);
	}
	color /= 8.0;

	return color;
}

vec4 caclulateKernel(float _kernel[9])
{
	vec3 sampleTex[9];
    for(int i = 0; i < 9; i++)
    {
        sampleTex[i] = vec3(fetchColor(fragTexCoords + offsets[i]));
    }
    vec3 col = vec3(0.0);
    for(int i = 0; i < 9; i++)
        col += sampleTex[i] * _kernel[i];

    return vec4(col,1.0);
}

void main()
{
	fragColor = caclulateKernel(normalKernel);
	if(fragColor.a < 0.1)
	    discard;
}`
	POST_PROCESSING_SHADER_NOMS_VERTEX_SOURCE_OPENGL string = `
#version 120

attribute vec2 vertex;
attribute vec2 texCoord;

varying vec2 fragTexCoords;

void main()
{
	fragTexCoords = texCoord;
	gl_Position = vec4(vertex,0.0,1.0);
}`
	POST_PROCESSING_SHADER_NOMS_FRAGMENT_SOURCE_OPENGL string = `
#version 120

varying vec2 fragTexCoords;

uniform sampler2D BackBuffer;

const float offset = 1.0 / 300.0;
vec2 offsets[9];

float blurKernel[9];

float normalKernel[9];

void setValues()
{
        offsets[0] = vec2(-offset,  offset); // top-left
        offsets[1] = vec2( 0.0,    offset); // top-center
        offsets[2] = vec2( offset,  offset); // top-right
        offsets[3] = vec2(-offset,  0.0);   // center-left
        offsets[4] = vec2( 0.0,    0.0);   // center-center
        offsets[5] = vec2( offset,  0.0);   // center-right
        offsets[6] = vec2(-offset, -offset); // bottom-left
        offsets[7] = vec2( 0.0,   -offset); // bottom-center
        offsets[8] = vec2( offset, -offset); // bottom-right

        blurKernel[0] = 1.0/16.0; blurKernel[1] = 2.0/16.0; blurKernel[2] = 1.0/16.0;
        blurKernel[3] = 2.0/16.0;  blurKernel[4] = 4.0/16.0; blurKernel[5] = 2.0/16.0;
        blurKernel[6] = 1.0/16.0; blurKernel[7] = 2.0/16.0; blurKernel[8] = 1.0/16.0;

        normalKernel[0] = 0.0;normalKernel[1] = 0.0;normalKernel[2] = 0.0;
        normalKernel[3] = 0.0;normalKernel[4] = 1.0;normalKernel[5] = 0.0;
        normalKernel[6] = 0.0;normalKernel[7] = 0.0;normalKernel[8] = 0.0;
}

vec4 fetchColor(vec2 texCoord)
{
	return texture2D(BackBuffer,texCoord);
}

vec4 caclulateKernel(float _kernel[9])
{
	vec3 sampleTex[9];
    for(int i = 0; i < 9; i++)
    {
        sampleTex[i] = vec3(fetchColor(fragTexCoords + offsets[i]));
    }
    vec3 col = vec3(0.0);
    for(int i = 0; i < 9; i++)
        col += sampleTex[i] * _kernel[i];

    return vec4(col,1.0);
}

void main()
{
    setValues();
	gl_FragColor = caclulateKernel(normalKernel);
	if(gl_FragColor.a < 0.1)
	    discard;
}`
)

// RenderScreenShader
const (
	RENDER_SCREEN_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 150

in vec2 fragTexCoords;

out vec4 fragColor;

uniform sampler2D BackBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);

	color = texture2D(BackBuffer,fragTexCoords);

	return color;
}

void main()
{
	fragColor = fetchColor();
	if(fragColor.a < 0.1)
	    discard;
}`
	RENDER_SCREEN_NOMS_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

varying vec2 fragTexCoords;

uniform sampler2D BackBuffer;

vec4 fetchColor()
{
	vec4 color = vec4(0.0);

	color = texture2D(BackBuffer,fragTexCoords);

	return color;
}

void main()
{
	gl_FragColor = fetchColor();
	if(gl_FragColor.a < 0.1)
	    discard;
}`
)

// Entity3D Shader
const (
	ENTITY_3D_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;

varying vec2 fragTexCoord;
varying vec3 fragPos;
varying vec3 fragNormal;
varying mat3 fragToTangentSpace;
varying mat4 fragViewMatrix3D;
varying mat4 fragInverseViewMatrix3D;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 inverseViewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
	fragPos =  (viewMatrix3D*transformMatrix3D*vec4(vertex,1.0)).xyz;
	// fragNormal =  (viewMatrix3D*vec4(mat3(transpose(inverseMat3(transformMatrix3D))) * normal,1.0)).xyz;
	fragNormal =  (viewMatrix3D*transformMatrix3D*vec4(normal,0.0)).xyz;

	vec3 norm = normalize(fragNormal);
	vec3 tang = normalize((viewMatrix3D*transformMatrix3D*vec4(tangent,0.0)).xyz);
	vec3 bitang = normalize(cross(norm,tang));

	fragToTangentSpace = mat3(
		tang.x,bitang.x,norm.x,
		tang.y,bitang.y,norm.y,
		tang.z,bitang.z,norm.z
	);

	fragViewMatrix3D = viewMatrix3D;
	fragInverseViewMatrix3D = inverseViewMatrix3D;
}`
	ENTITY_3D_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

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
vec3 norm;
vec3 viewDir;

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
	if(material.DiffuseTextureLoaded)
	{
		return texture2D(materialdiffuseTexture,fragTexCoord);
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
		return texture2D(materialspecularTexture,fragTexCoord);
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
	float closestDepth = texture2D(shadowMap, projCoords.xy).r;
	vec2 texelSize = 1.0 / vec2(shadowMapSize);
	for(int x = -1; x <= 1; ++x)
	{
	    for(int y = -1; y <= 1; ++y)
	    {
	        float pcfDepth = texture2D(shadowMap, projCoords.xy + vec2(x, y) * texelSize).r; 
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
	int samples = 20;
	float viewDistance = length(-fragPos);
	float diskRadius = (1.0 + (viewDistance / pl.farPlane)) / 70.0;
	for(int i = 0;i<samples;i++) {
		 float closestDepth = textureCube(shadowmap, fragToLight + sampleOffsetDirections[i]*diskRadius).r;
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
	if(material.NormalMapLoaded)
	{
		norm = normalize(2.0*(texture2D(materialnormalMap,fragTexCoord)).xyz-1.0);
	}
	else
	{
		norm = normalize(fragToTangentSpace*fragNormal);
	}
	viewDir = normalize((fragToTangentSpace*(fragPos*-1.0)));
}`
	ENTITY_3D_NO_SHADOWS_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;

varying vec2 fragTexCoord;
varying vec3 fragPos;
varying vec3 fragNormal;
varying mat3 fragToTangentSpace;
varying mat4 fragViewMatrix3D;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
	fragPos =  (viewMatrix3D*transformMatrix3D*vec4(vertex,1.0)).xyz;
	// fragNormal =  (viewMatrix3D*vec4(mat3(transpose(inverseMat3(transformMatrix3D))) * normal,1.0)).xyz;
	fragNormal =  (viewMatrix3D*transformMatrix3D*vec4(normal,0.0)).xyz;

	vec3 norm = normalize(fragNormal);
	vec3 tang = normalize((viewMatrix3D*transformMatrix3D*vec4(tangent,0.0)).xyz);
	vec3 bitang = normalize(cross(norm,tang));

	fragToTangentSpace = mat3(
		tang.x,bitang.x,norm.x,
		tang.y,bitang.y,norm.y,
		tang.z,bitang.z,norm.z
	);

	fragViewMatrix3D = viewMatrix3D;
}`
	ENTITY_3D_NO_SHADOWS_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

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
		return texture2D(materialdiffuseTexture,fragTexCoord);
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
		return texture2D(materialspecularTexture,fragTexCoord);
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

	diffuse *= attent;
	specular *= attent;

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
	if(material.NormalMapLoaded)
	{
		norm = normalize(2.0*(texture2D(materialnormalMap,fragTexCoord)).xyz-1.0);
	}
	else
	{
		norm = normalize(fragToTangentSpace*fragNormal);
	}
	viewDir = normalize((fragToTangentSpace*(fragPos*-1.0)));
}`
	ENTITY_3D_NOUV_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;

varying vec3 fragPos;
varying vec3 fragNormal;
varying mat4 fragViewMatrix3D;
varying mat4 fragInverseViewMatrix3D;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 inverseViewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragPos =  (transformMatrix3D*vec4(vertex,1.0)).xyz;
	// fragNormal =  (viewMatrix3D*vec4(mat3(transpose(inverseMat3(transformMatrix3D))) * normal,1.0)).xyz;
	fragNormal =  normalize((transformMatrix3D*vec4(normal,0.0)).xyz);

	fragViewMatrix3D = viewMatrix3D;
	fragInverseViewMatrix3D = inverseViewMatrix3D;
}`
	ENTITY_3D_NOUV_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

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

	float shinyness;
} material;

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

vec4 finalDiffuseColor;
vec4 finalSpecularColor;
vec4 finalAmbientColor;
vec3 norm;
vec3 viewDir;

varying vec3 fragPos;
varying vec3 fragNormal;
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

void calculateLightColors()
{
	finalDiffuseColor *= vec4(material.diffuseColor,1.0);
	finalSpecularColor *= vec4(material.specularColor,1.0);
	finalAmbientColor *= vec4(material.diffuseColor,1.0);

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
	vec4 fragPosLightSpace = lightSpaceMatrix*vec4(fragPos,1.0);
	vec3 projCoords = clamp((fragPosLightSpace.xyz / fragPosLightSpace.w)*0.5+0.5,-1.0,1.0);
	float currentDepth = projCoords.z-bias;
	float shadowresult = 0.0;
	float closestDepth = texture2D(shadowMap, projCoords.xy).r;
	vec2 texelSize = 1.0 / vec2(shadowMapSize);
	for(int x = -1; x <= 1; ++x)
	{
	    for(int y = -1; y <= 1; ++y)
	    {
	        float pcfDepth = texture2D(shadowMap, projCoords.xy + vec2(x, y) * texelSize).r; 
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

	vec3 fragToLight = fragPos - pl.position;
	float currentDepth = length(fragToLight)-bias*10.0*pl.farPlane;
	float shadow  = 0.0;
	int samples = 20;
	float viewDistance = length(-((fragViewMatrix3D*vec4(fragPos,1.0)).xyz));
	float diskRadius = (1.0 + (viewDistance / pl.farPlane)) / 70.0;
	for(int i = 0;i<samples;i++) {
		 float closestDepth = textureCube(shadowmap, fragToLight + sampleOffsetDirections[i]*diskRadius).r;
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

void calculatePointLight(PointLight pl,int index)
{
	vec3 lightPosition = pl.position;
	vec3 lightDir = normalize(lightPosition - fragPos);


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
	vec3 lightDirection = -dl.direction;
	vec3 lightDir = normalize(lightDirection);
	
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
	float theta = dot(lightDir, lightDirection);
	float spotAmount = 0.0;
	float outerCutOff = cos(degToRad(pl.outerCutOff));
	float innerCutOff = cos(degToRad(pl.innerCutOff));
	float epsilon   = innerCutOff - outerCutOff;
	spotAmount = clamp((theta - outerCutOff) / epsilon, 0.0, 1.0);

	return spotAmount;
}

void calculateSpotLight(SpotLight pl,int index)
{
	vec3 lightPosition = pl.position;
	vec3 lightDirection = -pl.direction;
	vec3 lightDir = normalize(lightPosition-fragPos);

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
	norm = fragNormal;
	vec3 camPos = (fragInverseViewMatrix3D*vec4(0.0,0.0,0.0,1.0)).xyz;
	viewDir = camPos - fragPos;
}`
	ENTITY_3D_NOUV_NO_SHADOWS_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

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

	float shinyness;
} material;

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

vec4 finalDiffuseColor;
vec4 finalSpecularColor;
vec4 finalAmbientColor;
vec3 norm;
vec3 viewDir;

varying vec3 fragPos;
varying vec3 fragNormal;
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

void calculateLightColors()
{
	finalDiffuseColor *= vec4(material.diffuseColor,1.0);
	finalSpecularColor *= vec4(material.specularColor,1.0);
	finalAmbientColor *= vec4(material.diffuseColor,1.0);

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

float calcAttentuation(vec3 lightPosition,Attentuation attentuation)
{
	float distance = distance(lightPosition,fragPos);
	float attent = 1.0/(attentuation.quadratic*distance*distance + attentuation.linear*distance + attentuation.constant);
	return attent;
}

void calculatePointLight(PointLight pl,int index)
{
	vec3 lightPosition = pl.position;
	vec3 lightDir = normalize(lightPosition - fragPos);


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
void calculateDirectionalLight(DirectionalLight dl,int index)
{
	vec3 lightDirection = -dl.direction;
	vec3 lightDir = normalize(lightDirection);
	
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
	float theta = dot(lightDir, lightDirection);
	float spotAmount = 0.0;
	float outerCutOff = cos(degToRad(pl.outerCutOff));
	float innerCutOff = cos(degToRad(pl.innerCutOff));
	float epsilon   = innerCutOff - outerCutOff;
	spotAmount = clamp((theta - outerCutOff) / epsilon, 0.0, 1.0);

	return spotAmount;
}

void calculateSpotLight(SpotLight pl,int index)
{
	vec3 lightPosition = pl.position;
	vec3 lightDirection = -pl.direction;
	vec3 lightDir = normalize(lightPosition-fragPos);

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
	norm = fragNormal;
	vec3 camPos = (fragInverseViewMatrix3D*vec4(0.0,0.0,0.0,1.0)).xyz;
	viewDir = camPos - fragPos;
}`
	ENTITY_3D_SIMPLE_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

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
		return texture2D(materialdiffuseTexture,fragTexCoord);
	}
	else
	{
		return vec4(1.0,1.0,1.0,1.0);
	}
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

	diffuse *= attent;
	specular *= attent;

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
}`
	ENTITY_3D_INSTANCED_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;
attribute mat4 transformMatrix3D;

varying vec2 fragTexCoord;
varying vec3 fragPos;
varying vec3 fragNormal;
varying mat3 fragToTangentSpace;
varying mat4 fragViewMatrix3D;
varying mat4 fragInverseViewMatrix3D;

uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;
uniform mat4 inverseViewMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
	fragPos =  (viewMatrix3D*transformMatrix3D*vec4(vertex,1.0)).xyz;
	// fragNormal =  (viewMatrix3D*vec4(mat3(transpose(inverse(transformMatrix3D))) * normal,1.0)).xyz;
	fragNormal =  (viewMatrix3D*transformMatrix3D*vec4(normal,0.0)).xyz;

	vec3 norm = normalize(fragNormal);
	vec3 tang = normalize((viewMatrix3D*transformMatrix3D*vec4(tangent,0.0)).xyz);
	vec3 bitang = normalize(cross(norm,tang));

	fragToTangentSpace = mat3(
		tang.x,bitang.x,norm.x,
		tang.y,bitang.y,norm.y,
		tang.z,bitang.z,norm.z
	);

	fragViewMatrix3D = viewMatrix3D;
	fragInverseViewMatrix3D = inverseViewMatrix3D;
}`
)

// Lines3D Shader
const (
	LINES_3D_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec4 color;

varying vec4 fragColor;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragColor = color;
}`
	LINES_3D_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

varying vec4 fragColor;

void main()
{
    gl_FragColor = fragColor;
}`
)

// PointLight Shadowmap Shader
const (
	POINTLIGHT_SHADOWMAP_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 150

in vec3 vertex;
in vec3 normal;
in vec2 texCoord;
in vec3 tangent;

out vec2 geoTexCoord;

uniform mat4 transformMatrix3D;

void main()
{
	gl_Position = transformMatrix3D*vec4(vertex,1.0);
	geoTexCoord = texCoord;
}`
	POINTLIGHT_SHADOWMAP_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 150

in vec2 fragTexCoord;
in vec4 fragPos;

uniform vec3 lightPos;
uniform float farPlane;
uniform struct Material
{
	bool DiffuseTextureLoaded;
} material;
uniform	sampler2D materialdiffuseTexture;


vec4 fetchColor()
{
	if(material.DiffuseTextureLoaded)
	{
		return texture2D(materialdiffuseTexture,fragTexCoord);
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
	float lightDistance = length(fragPos.xyz - lightPos);
	lightDistance = lightDistance / farPlane;
	gl_FragDepth = lightDistance;
}`
	POINTLIGHT_SHADOWMAP_SHADER_GEOMETRY_SOURCE_OPENGL string = `
#version 150

layout(triangles) in;
layout(triangle_strip,max_vertices=18) out;

in vec2 geoTexCoord[];

out vec2 fragTexCoord;
out	vec4 fragPos;

uniform mat4 lightSpaceMatrices[6];
uniform mat4 projectionMatrix3D;

void main()
{
	for(int face = 0;face < 6;++face)
	{
		gl_Layer = face;
		for(int i = 0;i<3;++i)
		{
			fragPos = gl_in[i].gl_Position;
			gl_Position = projectionMatrix3D * lightSpaceMatrices[face] * fragPos;
			switch(i)
			{
			    case 0:
			        fragTexCoord = geoTexCoord[0];
			        break;
			    case 1:
                	fragTexCoord = geoTexCoord[1];
                	break;
                case 2:
                    fragTexCoord = geoTexCoord[2];
                	break;
                default:
                    break;
			}
			EmitVertex();
		}
		EndPrimitive();
	}
}`
	POINTLIGHT_SHADOWMAP_INSTANCED_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 150

in vec3 vertex;
in vec3 normal;
in vec2 texCoord;
in vec3 tangent;
in mat4 transformMatrix3D;

out	vec2 geoTexCoord;


void main()
{
	gl_Position = transformMatrix3D*vec4(vertex,1.0);
	geoTexCoord = texCoord;
}`
)

// Shadowmap Shader
const (
	SHADOWMAP_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;


varying vec2 fragTexCoord;

uniform mat4 transformMatrix3D;
uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
}`
	SHADOWMAP_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

varying	vec2 fragTexCoord;

uniform struct Material
{
	bool DiffuseTextureLoaded;
} material;
uniform sampler2D materialdiffuseTexture;

vec4 getDiffuseTexture()
{
	if(material.DiffuseTextureLoaded)
	{
		return texture2D(materialdiffuseTexture,fragTexCoord);
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
}`
	SHADOWMAP_INSTANCED_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec3 vertex;
attribute vec3 normal;
attribute vec2 texCoord;
attribute vec3 tangent;
attribute mat4 transformMatrix3D;

varying vec2 fragTexCoord;


uniform mat4 viewMatrix3D;
uniform mat4 projectionMatrix3D;

void main()
{
	gl_Position = projectionMatrix3D*viewMatrix3D*transformMatrix3D*vec4(vertex,1.0);
	fragTexCoord = texCoord;
}`
	SHADOWMAP_RENDER_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

varying vec2 fragTexCoord;

uniform sampler2D texture0;

void main()
{
	float depth = texture2D(texture0,fragTexCoord).r;
	gl_FragColor = vec4(depth,depth,depth,1.0);
}`
)

// Sprite2D Shader
const (
	SPRITE_2D_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

#define FLIP_NONE 0
#define FLIP_HORIZONTAL 1
#define FLIP_VERTICAL 2
#define FLIP_DIAGONALLY 3

attribute vec2 vertex;
attribute vec2 texCoord;

uniform mat3 transformMatrix2D;
uniform mat4 projectionMatrix2D;
uniform mat3 viewMatrix2D;
uniform vec4 textureRegion;
uniform bool enableTextureRegion;
uniform int flip;
uniform float depth;

varying vec2 fragTexCoord;

vec2 textureRegionToTexCoord(vec2 tc);
vec2 flipTexCoord(vec2 tc);

void main()
{
	gl_Position = projectionMatrix2D *vec4(vec2(viewMatrix2D*transformMatrix2D*vec3(vertex,1.0)),0.0,1.0);
	gl_Position.z = depth;
	fragTexCoord = textureRegionToTexCoord(flipTexCoord(texCoord));
}

vec2 textureRegionToTexCoord(vec2 tc)
{
	if(!enableTextureRegion)
		return tc;

    // X: 0 ->      0 -> Min X
    // Y: 0 ->      0 -> Min Y
    // Z: WIDTH ->  1 -> Max X
    // W: HEIGHT -> 1 -> Max Y

    vec2 newTexCoord = tc;
    newTexCoord.x = newTexCoord.x * (textureRegion.z-textureRegion.x)+textureRegion.x;
    newTexCoord.y = newTexCoord.y * (textureRegion.w-textureRegion.y)+textureRegion.y;

    return newTexCoord;
}

vec2 flipTexCoord(vec2 tc)
{
	vec2 flippedTexCoord;


	if(flip == FLIP_NONE)
	{
		flippedTexCoord = tc;
	}
	else if(flip == FLIP_HORIZONTAL)
	{
		flippedTexCoord.x = 1.0-tc.x;
		flippedTexCoord.y = tc.y;
	}
	else if(flip == FLIP_VERTICAL)
	{
		flippedTexCoord.x = tc.x;
		flippedTexCoord.y = 1.0-tc.y;
	}
	else if(flip == FLIP_DIAGONALLY)
	{
		flippedTexCoord.x = 1.0-tc.x;
		flippedTexCoord.y = 1.0-tc.y;
	}
	else
	{
		flippedTexCoord = tc;
	}

	return flippedTexCoord;
}`
	SPRITE_2D_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

#define KEY_COLOR_PADDING 0.1
#define ALPHA_DISCARD_PADDING 0.1

varying vec2 fragTexCoord;

uniform sampler2D texture0;
uniform vec3 keyColor;
uniform vec4 modColor;
uniform bool enableKey;
uniform bool enableMod;

vec4 applyKeyColor(vec4 color);
vec4 applyModColor(vec4 color);

void main()
{
	gl_FragColor = applyModColor(applyKeyColor(texture2D(texture0,fragTexCoord)));
	if(gl_FragColor.a < ALPHA_DISCARD_PADDING)
	{
		discard;
	}
}

vec4 applyKeyColor(vec4 color)
{
    if(enableKey)
    {
        if(color.r >= keyColor.r - KEY_COLOR_PADDING && color.r <= keyColor.r + KEY_COLOR_PADDING &&
           color.g >= keyColor.g - KEY_COLOR_PADDING && color.g <= keyColor.g + KEY_COLOR_PADDING &&
           color.b >= keyColor.b - KEY_COLOR_PADDING && color.b <= keyColor.b + KEY_COLOR_PADDING)
        {
           discard;
        }
    }

    return color;
}

vec4 applyModColor(vec4 color)
{
    if(enableMod)
    {
    	color *= modColor;
        return color;
    }

    return color;
}`
)

// Text2D Shader
const (
	TEXT_2D_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

#define ALPHA_DISCARD_PADDING 0.3

varying vec2 fragTexCoord;

uniform sampler2D texture0;
uniform vec4 color;

void main()
{
	gl_FragColor = texture2D(texture0,fragTexCoord) * color;
	if(gl_FragColor.a < ALPHA_DISCARD_PADDING)
	{
		discard;
	}
}`
)

// Shape2D Shader
const (
	SHAPE_2D_SHADER_VERTEX_SOURCE_OPENGL string = `
#version 110

attribute vec2 vertex;
attribute vec4 color;

uniform mat3 transformMatrix2D;
uniform mat4 projectionMatrix2D;
uniform mat3 viewMatrix2D;
uniform float depth;

varying vec4 fragColor;

void main()
{
	gl_Position = projectionMatrix2D *vec4(vec2(viewMatrix2D*transformMatrix2D*vec3(vertex,1.0)),0.0,1.0);
	gl_Position.z = depth;
	fragColor = color;
}
`
	SHAPE_2D_SHADER_FRAGMENT_SOURCE_OPENGL string = `
#version 110

#define ALPHA_DISCARD_PADDING 0.1

varying vec4 fragColor;

void main()
{
	gl_FragColor = fragColor;
	if(gl_FragColor.a < ALPHA_DISCARD_PADDING)
		discard;
}
`
)
