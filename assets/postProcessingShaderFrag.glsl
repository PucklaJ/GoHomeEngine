#version 100

precision mediump float;

varying vec2 fragTexCoords;

uniform sampler2D BackBuffer;

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
	return texture2D(backBuffer,texCoord);
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
	gl_FragColor = caclulateKernel(normalKernel);
}