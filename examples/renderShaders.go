package main

import (
	"fmt"
	"log"

	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/singmyr/mainthread"

	"github.com/singmyr/astrocyte"
)

// https://learnopengl.com/Getting-started/Hello-Triangle
func run() {
	fmt.Println("Launching astrocyte testing code")
	window, err := astrocyte.CreateWindow(640, 640, "Astrocyte testing", &astrocyte.WindowConfig{FpsLimit: 1.0})
	if err != nil {
		panic(err)
	}

	mainthread.Call(func() {
		gl.Viewport(0, 0, 640, 640)
	})

	vertices := []float32{
		0.5, 0.5, 0.0, 1.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, // bottom left
		-0.5, 0.5, 0.0, 1.0, 0.0, 0.0, // top left
	}

	indices := []uint32{
		// First triangle.
		0, 1, 3,
		// Second triangle.
		1, 2, 3,
	}

	// Create the triangle vertices.
	// vertices := []float32{
	// 	// first triangle
	// 	0.5, 0.5, 0.0, // top right
	// 	0.5, -0.5, 0.0, // bottom right
	// 	-0.5, 0.5, 0.0, // top left
	// 	// second triangle
	// 	0.5, -0.5, 0.0, // bottom right
	// 	-0.5, -0.5, 0.0, // bottom left
	// 	-0.5, 0.5, 0.0, // top left
	// }

	var vbo uint32
	var vao uint32
	var ebo uint32
	var shaderProgram uint32
	mainthread.Call(func() {
		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)

		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		/*
			GL_STATIC_DRAW: the data will most likely not change at all or very rarely.
			GL_DYNAMIC_DRAW: the data is likely to change a lot.
			GL_STREAM_DRAW: the data will change every time it is drawn.
		*/

		gl.GenBuffers(1, &ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		// Vertex shader.
		vShader := `#version 330 core
		layout (location = 0) in vec3 aPos;
		layout (location = 1) in vec3 aColor;

		out vec3 color;

		void main()
		{
			gl_Position = vec4(aPos, 1.0);
			color = aColor;
		}`

		vSrc, vFree := gl.Strs(vShader)
		defer vFree()
		vLen := int32(len(vShader))

		vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
		gl.ShaderSource(vertexShader, 1, vSrc, &vLen)
		gl.CompileShader(vertexShader)

		// @todo: Check if successful or not.
		/*
			int  success;
			char infoLog[512];
			glGetShaderiv(vertexShader, GL_COMPILE_STATUS, &success)
			if(!success)
			{
			    glGetShaderInfoLog(vertexShader, 512, NULL, infoLog);
			    std::cout << "ERROR::SHADER::VERTEX::COMPILATION_FAILED\n" << infoLog << std::endl;
			}
		*/

		var success int32
		gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetShaderInfoLog(vertexShader, logLen, nil, &infoLog[0])
			log.Printf("error compiling vertex shader: %s", string(infoLog))
		}

		defer gl.DeleteShader(vertexShader)

		// Fragment shader.
		fShader := `#version 330 core
		out vec4 FragColor;

		in vec3 color;
		
		void main()
		{
			FragColor = vec4(color, 1.0f);
		}`

		fSrc, fFree := gl.Strs(fShader)
		defer fFree()
		fLen := int32(len(fShader))

		fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
		gl.ShaderSource(fragmentShader, 1, fSrc, &fLen)
		gl.CompileShader(fragmentShader)

		// @todo: Check if successful or not.
		/*
			int  success;
			char infoLog[512];
			glGetShaderiv(fragmentShader, GL_COMPILE_STATUS, &success)
			if(!success)
			{
			    glGetShaderInfoLog(fragmentShader, 512, NULL, infoLog);
			    std::cout << "ERROR::SHADER::VERTEX::COMPILATION_FAILED\n" << infoLog << std::endl;
			}
		*/

		// var success int32 - not required since it was done previously.
		gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(fragmentShader, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetShaderInfoLog(fragmentShader, logLen, nil, &infoLog[0])
			log.Printf("error compiling fragment shader: %s", string(infoLog))
		}

		defer gl.DeleteShader(fragmentShader)

		shaderProgram = gl.CreateProgram()
		gl.AttachShader(shaderProgram, vertexShader)
		gl.AttachShader(shaderProgram, fragmentShader)
		gl.LinkProgram(shaderProgram)
		gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetProgramInfoLog(shaderProgram, logLen, nil, &infoLog[0])
			log.Printf("error linking program: %s", string(infoLog))
		}

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 2*3*4, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 2*3*4, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		/*
			// 0. copy our vertices array in a buffer for OpenGL to use
			glBindBuffer(GL_ARRAY_BUFFER, VBO);
			glBufferData(GL_ARRAY_BUFFER, sizeof(vertices), vertices, GL_STATIC_DRAW);
			// 1. then set the vertex attributes pointers
			glVertexAttribPointer(0, 3, GL_FLOAT, GL_FALSE, 3 * sizeof(float), (void*)0);
			glEnableVertexAttribArray(0);
			// 2. use our shader program when we want to render an object
			glUseProgram(shaderProgram);
			// 3. now draw the object
			someOpenGLFunctionThatDrawsOurTriangle();
		*/
	})

	var ourColor int32
	mainthread.Call(func() {
		// If gl.GetUniformLocation returns -1, it failed to locate it.
		ourColor = gl.GetUniformLocation(shaderProgram, gl.Str("ourColor\x00"))
	})

	for window.IsOpen() {
		// This isn't being run right now because not sure where it fits in.
		// window.Update()
		mainthread.Call(func() {
			// timeValue := glfw.GetTime()
			// redValue := float32((math.Cos(timeValue) / 2.0) + 0.5)
			// greenValue := float32((math.Sin(timeValue) / 2.0) + 0.5)
			// Wireframe mode.
			// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			gl.ClearColor(0.0, 0.0, 0.0, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			gl.UseProgram(shaderProgram)
			// gl.Uniform4f(ourColor, redValue, greenValue, 1.0, 1.0)
			gl.BindVertexArray(vao)
			gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
			gl.BindVertexArray(0)
		})

		if window.IsKeyPressed(glfw.KeyW) {
			log.Println("W is pressed")
		}

		window.Render()
	}
}

func main() {
	astrocyte.Run(run)
}
