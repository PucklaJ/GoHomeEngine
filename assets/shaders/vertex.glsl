#version 320 es

uniform vec2 vertex1 = vec2(-0.5,-0.5);
uniform vec2 vertex2 = vec2(0.5,-0.5);
uniform vec2 vertex3 = vec2(0.0,0.5);

void main()
{
	// vertices[0] = vec2(-0.5,-0.5);
	// vertices[1] = vec2(0.5,-0.5);
	// vertices[2] = vec2(0.0,0.5);

	// gl_Position = vec4(vertices[gl_VertexID],0.0,1.0);

	if(gl_VertexID == 0)
	{
		gl_Position = vec4(vertex1,0.0,1.0);
	}
	else if(gl_VertexID == 1)
	{
		gl_Position = vec4(vertex2,0.0,1.0);
	}
	else if(gl_VertexID == 2)
	{
		gl_Position = vec4(vertex3,0.0,1.0);
	}
}