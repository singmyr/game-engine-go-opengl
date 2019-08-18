#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;
layout (location = 2) in vec2 aTexCoord;

out vec3 color;
out vec2 texCoord;

uniform vec2 offset;

void main()
{
    gl_Position = vec4(aPos, 1.0) + vec4(offset, 0.0f, 0.0f);
    color = aColor;
    texCoord = aTexCoord;
}
