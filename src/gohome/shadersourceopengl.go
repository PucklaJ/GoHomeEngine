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
	float transparency;
} material;
uniform	sampler2D materialdiffuseTexture;


vec4 fetchColor()
{	
	vec4 col;
	if(material.DiffuseTextureLoaded)
	{
		col = texture2D(materialdiffuseTexture,fragTexCoord);
	}
	else
	{
		col = vec4(1.0,1.0,1.0,1.0);
	}

	col.w *= material.transparency;
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
	float transparency;
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
	texDifCol.w *= material.transparency;

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
