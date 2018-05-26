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
}