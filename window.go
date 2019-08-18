package astrocyte

//"github.com/go-gl/gl/v2.1/gl"
//"github.com/go-gl/glfw/v3.2/glfw"
import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/singmyr/mainthread"
)

// Window is the general structure for our window objects.
type Window struct {
	window       *glfw.Window
	updateLimit  float64
	previousTime float64
	elapsed      float64
}

// WindowConfig contains all the possible configurations you can use for when creating a new window.
type WindowConfig struct {
	Resizable bool
	FpsLimit  float64
}

var currentWindow *Window

func frameBufferSizeCallback(w *glfw.Window, width int, height int) {
	log.Println("Window has been resized to:", width, height)
	// @todo: is gl.Viewport needed here?
}

// CreateWindow creates a new window and returns the Window object for it.
func CreateWindow(width int, height int, title string, cfg *WindowConfig) (*Window, error) {
	if currentWindow != nil {
		log.Fatalln("A window has already been created.")
	}
	// @todo: Make use of the cfg for these window hints.
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var w *glfw.Window
	var err error
	mainthread.Call(func() {
		w, err = glfw.CreateWindow(width, height, title, nil, nil)
		if err != nil {
			panic(err)
		}

		w.MakeContextCurrent()
		if err := gl.Init(); err != nil {
			log.Fatalln(err)
		}

		version := gl.GoStr(gl.GetString(gl.VERSION))
		log.Println("OpenGL version", version)

		gl.ClearColor(1.0, 0.3, 0.3, 1.0)

		w.SetFramebufferSizeCallback(frameBufferSizeCallback)
	})

	currentWindow = &Window{window: w}
	currentWindow.updateLimit = 0.0
	if cfg.FpsLimit != 0.0 {
		currentWindow.updateLimit = 1.0 / cfg.FpsLimit
	}
	currentWindow.previousTime = glfw.GetTime()

	return currentWindow, err
}

// Run is the first function that has to be called and it has to be in the main function.
func Run(run func()) {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	defer glfw.Terminate()

	mainthread.Run(32, run)
}

var calls int

// Update needs to be called prior to Render inside a loop to cause events to be polled and objects updated.
func (w *Window) Update() {
	mainthread.Call(func() {
		log.Println("Update")
		time := glfw.GetTime()
		delta := time - w.previousTime
		w.previousTime = time
		w.elapsed += delta
		if w.elapsed >= w.updateLimit {
			calls++
			log.Printf("Calls: %d\n", calls)
			w.elapsed -= w.updateLimit
		}

		// fmt.Printf("Delta: %f\n", delta)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	})
}

// Render renders stuff on screen.
func (w *Window) Render() {
	w.window.SwapBuffers()
	mainthread.Call(func() {
		// @todo: Do rendering.

		glfw.PollEvents()
	})
}

// IsOpen returns a boolean indicating whether the window is open or closed.
func (w *Window) IsOpen() bool {
	return !w.window.ShouldClose()
}
