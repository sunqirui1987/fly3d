package gl

const (
	EXPLICIT_MODE   Enum = 0
	SPHERICAL_MODE  Enum = 1
	PLANAR_MODE     Enum = 2
	CUBIC_MODE      Enum = 3
	PROJECTION_MODE Enum = 4
	SKYBOX_MODE     Enum = 5

	CLAMP_ADDRESSMODE  Enum = 0
	WRAP_ADDRESSMODE   Enum = 1
	MIRROR_ADDRESSMODE Enum = 2
)

type GLCullState struct {
	Culling bool
}

type GLTextureBuffer struct {
	Tex             Texture
	FrameBuf        Framebuffer
	DepthBuf        Renderbuffer
	GenerateMipMaps bool

	BaseWidth  int
	BaseHeight int
	Width      int
	Height     int
	Url        string
	NoMipmap   bool
	References int

	IsReady bool
	IsCube  bool

	//更新接口
	UpdateFunc func() bool

	//cache
	CoordinatesMode      Enum
	CacheCoordinatesMode Enum

	WrapU       Enum
	CachedWrapU Enum

	WrapV       Enum
	CachedWrapV Enum
}

func NewGLTextureBuffer() *GLTextureBuffer {
	tex := &GLTextureBuffer{
		GenerateMipMaps: false,
		BaseWidth:       0,
		BaseHeight:      0,
		Width:           0,
		Height:          0,
		Url:             "",
		NoMipmap:        false,
		References:      0,
		IsReady:         false,
		IsCube:          false,
		UpdateFunc:      nil,
	}
	tex.Tex = Texture{}
	tex.FrameBuf = Framebuffer{}
	tex.DepthBuf = Renderbuffer{}
	return tex
}

type GLVertexBuffer struct {
	Vbo        Buffer
	References int
	StrideSize int
}
type GLIndexBuffer struct {
	Vbo        Buffer
	References int
	Is32Bits   bool
}
