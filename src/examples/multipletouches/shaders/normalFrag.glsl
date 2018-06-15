#version 120

varying vec2 fragTexCoord;
varying vec3 fragPos;
varying vec3 fragNormal;
varying mat3 fragToTangentSpace;
varying mat4 fragViewMatrix3D;
varying mat4 fragInverseViewMatrix3D;

uniform struct Material
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
} material;

float detMat3(mat3 m)
{
	return m[0][0]*m[1][1]*m[2][2] + m[1][0]*m[2][1]*m[0][2] + m[2][0]*m[0][1]*m[1][2] - m[2][0]*m[1][1]*m[0][2] - m[1][0]*m[0][1]*m[2][2] - m[0][0]*m[2][1]*m[1][2];
}

mat3 inverseMat3(mat3 m)
{
	float det = detMat3(m);
	if(det == 0.0)
		return m;

	mat3 retMat; 
	retMat[0][0] = m[1][1]*m[2][2] - m[1][2]*m[2][1];
	retMat[0][1] = m[0][2]*m[2][1] - m[0][1]*m[2][2];
	retMat[0][2] = m[0][1]*m[1][2] - m[0][2]*m[1][1];
	retMat[1][0] = m[1][2]*m[2][0] - m[1][0]*m[2][2];
	retMat[1][1] = m[0][0]*m[2][2] - m[0][2]*m[2][0];
	retMat[1][2] = m[0][2]*m[1][0] - m[0][0]*m[1][2];
	retMat[2][0] = m[1][0]*m[2][1] - m[1][1]*m[2][0];
	retMat[2][1] = m[0][1]*m[2][0] - m[0][0]*m[2][1];
	retMat[2][2] = m[0][0]*m[1][1] - m[0][1]*m[1][0];

	return retMat * (1.0 / det);
}

void main()
{	
	if(material.diffuseColor.r < 100.0 || material.specularColor.r < 100.0 || material.shinyness < 100.0 || material.diffuseTextureLoaded || material.specularTextureLoaded)
	{
		vec3 normal = vec3(0.0,0.0,0.0);
		if(material.normalMapLoaded)
		{
			normal = 2.0*(texture2D(material.normalMap,fragTexCoord).xyz)-1.0;
			normal = normalize((fragInverseViewMatrix3D*vec4(inverseMat3(fragToTangentSpace)*normal,0.0)).xyz);
		}
		else 
		{
			normal = fragNormal;
			normal = normalize((fragInverseViewMatrix3D*vec4(normal,0.0)).xyz);
		}

		gl_FragColor = vec4(normal,1.0);
		if (gl_FragColor.r < 0.0) {
			gl_FragColor.r = abs(gl_FragColor.r) * 0.1;
		}
		if (gl_FragColor.g < 0.0) {
			gl_FragColor.g = abs(gl_FragColor.g) * 0.1;
		}
		if (gl_FragColor.b < 0.0) {
			gl_FragColor.b = abs(gl_FragColor.b) * 0.1;
		}
		gl_FragColor = vec4(1.0,0.0,0.0,1.0);
	}
}