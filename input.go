package astrocyte

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// IsKeyPressed checks whether a given key is pressed or not.
func (w *Window) IsKeyPressed(k glfw.Key) bool {
	return w.window.GetKey(k) == glfw.Press
}
