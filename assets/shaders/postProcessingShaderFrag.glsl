#version 100

precision mediump float;

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
}