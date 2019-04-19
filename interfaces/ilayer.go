package interfaces

type ILayer interface {
	IsBackground() bool
	Render()
	Dispose()
}
