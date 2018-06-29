#version 110

#define KEY_COLOR_PADDING 0.1
#define ALPHA_DISCARD_PADDING 0.1

varying vec2 fragTexCoord;

uniform sampler2D texture0;
uniform vec3 keyColor;
uniform vec3 modColor;
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
        return vec4(color.xyz*modColor,color.a);
    }

    return color;
}