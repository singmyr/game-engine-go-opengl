package astrocyte

import (
	"io/ioutil"
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/singmyr/mainthread"
)

// Shader is the base type for our shaders.
type Shader struct {
	program uint32
}

// CreateShader compiles and links shaders.
func CreateShader(vertexShaderPath string, fragmentShaderPath string) *Shader {
	// Load the Vertex Shader.
	vsFile, err := ioutil.ReadFile(vertexShaderPath)
	if err != nil {
		log.Fatal(err)
	}

	fsFile, err := ioutil.ReadFile(fragmentShaderPath)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize empty Shader object.
	shader := &Shader{}

	mainthread.Call(func() {
		shader.program = gl.CreateProgram()

		// Vertex Shader.
		vSrc, vFree := gl.Strs(string(vsFile))
		defer vFree()
		vLen := int32(len(vsFile))

		vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
		gl.ShaderSource(vertexShader, 1, vSrc, &vLen)
		gl.CompileShader(vertexShader)

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

		// Fragment Shader.
		fSrc, fFree := gl.Strs(string(fsFile))
		defer fFree()
		fLen := int32(len(fsFile))

		fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
		gl.ShaderSource(fragmentShader, 1, fSrc, &fLen)
		gl.CompileShader(fragmentShader)

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

		// var success int32
		gl.AttachShader(shader.program, vertexShader)
		gl.AttachShader(shader.program, fragmentShader)
		gl.LinkProgram(shader.program)
		gl.GetProgramiv(shader.program, gl.LINK_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(shader.program, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetProgramInfoLog(shader.program, logLen, nil, &infoLog[0])
			log.Printf("error linking program: %s", string(infoLog))
		}
	})

	return shader
}

// Use activates the shader.
func (s *Shader) Use() {
	mainthread.Call(func() {
		gl.UseProgram(s.program)
	})
}
