package engines

import (
	"encoding/binary"
	"image"
	"reflect"

	log "github.com/suiqirui1987/fly3d/tools/logrus"

	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/gl"
	"github.com/suiqirui1987/fly3d/gl/glutil"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/tools"
	"github.com/suiqirui1987/fly3d/tools/resize"
	"github.com/suiqirui1987/fly3d/windows"

	"golang.org/x/mobile/exp/f32"
)

type EngineBuffersCache struct {
	_cachedVertexBuffer           *gl.GLVertexBuffer
	_cachedVertexBuffers          map[string]*gl.GLVertexBuffer
	_cachedIndexBuffer            *gl.GLIndexBuffer
	_cachedEffectForVertexBuffers IEffect
}

//engine

type EngineCaps struct {
	MaxTexturesImageUnits int
	MaxTextureSize        int
	MaxCubemapTextureSize int
	MaxRenderTextureSize  int

	StandardDerivatives bool
}

type Engine struct {
	ForceWireframe bool
	CullBackFaces  bool
	Scenes         []*Scene

	_caps *EngineCaps

	//Cache
	_loadedTexturesCache []*gl.GLTextureBuffer
	_activeTexturesCache []*gl.GLTextureBuffer
	_buffersCache        *EngineBuffersCache
	_currentEffect       IEffect
	_currentState        *gl.GLCullState

	//window
	IsFullscreen bool

	//get Properties
	_aspectRatio          float32
	_renderingCanvas      windows.IWindow
	_hardwareScalingLevel float32
	_alphaTest            bool
	_runningLoop          bool

	//
	CompiledEffects map[string]IEffect

	//render
	_renderFunction func()
}

func NewEngine(canvas windows.IWindow, antialias bool) *Engine {
	this := &Engine{}
	this._renderingCanvas = canvas
	this._alphaTest = false

	// Options
	this.ForceWireframe = false
	this.CullBackFaces = true

	//Viewport
	this._hardwareScalingLevel = 1.0 / canvas.GetWindowDevicePixelRatio()

	// Caps
	this._caps = &EngineCaps{}
	this._caps.MaxTexturesImageUnits = gl.GetInteger(gl.MAX_TEXTURE_IMAGE_UNITS)
	this._caps.MaxTextureSize = gl.GetInteger(gl.MAX_TEXTURE_SIZE)
	this._caps.MaxCubemapTextureSize = gl.GetInteger(gl.MAX_CUBE_MAP_TEXTURE_SIZE)
	this._caps.MaxRenderTextureSize = gl.GetInteger(gl.MAX_RENDERBUFFER_SIZE)

	// Extensions
	//derivatives := gl.GetExtension("OES_standard_derivatives")
	this._caps.StandardDerivatives = true

	// Cache
	this._loadedTexturesCache = make([]*gl.GLTextureBuffer, 0)
	this._activeTexturesCache = make([]*gl.GLTextureBuffer, 0)
	this._buffersCache = &EngineBuffersCache{}
	this._currentState = &gl.GLCullState{
		Culling: false,
	}

	this.CompiledEffects = map[string]IEffect{}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	this.IsFullscreen = false

	that := this
	this._renderingCanvas.On(windows.FullScreenChange, func() error {
		that.IsFullscreen = that._renderingCanvas.GetFullscreen()
		return nil
	})
	this._renderingCanvas.On(windows.Resize, func(evt windows.ResizeEvent) error {
		this._aspectRatio = (float32)(this._renderingCanvas.GetRenderWidth()) / (float32)(this._renderingCanvas.GetRenderHeight())
		return nil
	})
	this._aspectRatio = (float32)(this._renderingCanvas.GetRenderWidth()) / (float32)(this._renderingCanvas.GetRenderHeight())

	return this
}

func (this *Engine) GetAspectRatio() float32 {
	return this._aspectRatio
}

func (this *Engine) GetRenderWidth() int {
	return this._renderingCanvas.GetRenderWidth()
}

func (this *Engine) GetRenderHeight() int {
	return this._renderingCanvas.GetRenderHeight()
}

func (this *Engine) GetRenderingCanvas() windows.IWindow {
	return this._renderingCanvas
}

func (this *Engine) SetHardwareScalingLevel(level float32) {
	this._hardwareScalingLevel = level
}

func (this *Engine) GetHardwareScalingLevel() float32 {
	return this._hardwareScalingLevel
}

func (this *Engine) GetLoadedTexturesCache() []*gl.GLTextureBuffer {
	return this._loadedTexturesCache
}
func (this *Engine) GetCaps() *EngineCaps {
	return this._caps
}

//
func (this *Engine) StopRenderLoop() {
	this._renderFunction = nil
	this._runningLoop = false
}

func (this *Engine) _renderLoop() {
	this.BeginFrame()

	if this._renderFunction != nil {
		this._renderFunction()
	}

	// Present
	this.EndFrame()

	//是否停止
	if !this._runningLoop {
		this._renderingCanvas.StopNewFrame()
	}
}
func (this *Engine) RunRenderLoop(renderFunction func()) {
	this._runningLoop = true
	this._renderFunction = renderFunction

	that := this
	loop := func() {
		that._renderLoop()
	}

	that._renderingCanvas.QueueNewFrame(loop)

}
func (this *Engine) SwitchFullscreen() {
	if this.IsFullscreen {
		this._renderingCanvas.ExitFullscreen()
	} else {
		this._renderingCanvas.RequestFullscreen()
	}
}

func (this *Engine) Clear(color *math32.Color3, backBuffer, depthStencil bool) {

	gl.ClearColor(color.R, color.G, color.B, 1.0)
	gl.ClearDepthf(1.0)

	var mode gl.Enum
	mode = 0
	if backBuffer || this.ForceWireframe {
		mode |= gl.COLOR_BUFFER_BIT
	}

	if depthStencil {
		mode |= gl.DEPTH_BUFFER_BIT
	}

	gl.Clear(mode)

}

func (this *Engine) BeginFrame() {
	tools.MeasureFps()

	gl.Viewport(0, 0, this._renderingCanvas.GetRenderWidth(), this._renderingCanvas.GetRenderHeight())

}
func (this *Engine) EndFrame() {
	this.FlushFramebuffer()
}

func (this *Engine) BindFramebuffer(texture *gl.GLTextureBuffer) {

	gl.BindFramebuffer(gl.FRAMEBUFFER, texture.FrameBuf)
	gl.Viewport(0, 0, texture.Width, texture.Height)

	this.WipeCaches()
}

func (this *Engine) UnBindFramebuffer(texture *gl.GLTextureBuffer) {
	if texture.GenerateMipMaps {
		gl.BindTexture(gl.TEXTURE_2D, texture.Tex)
		gl.GenerateMipmap(gl.TEXTURE_2D)
		gl.BindTexture(gl.TEXTURE_2D, gl.Texture{})
	}
}

func (this *Engine) FlushFramebuffer() {
	gl.Flush()
}

func (this *Engine) RestoreDefaultFramebuffer() {

	gl.BindFramebuffer(gl.FRAMEBUFFER, gl.Framebuffer{})
	gl.Viewport(0, 0, this._renderingCanvas.GetRenderWidth(), this._renderingCanvas.GetRenderHeight())

	this.WipeCaches()
}

// VBOs
func (this *Engine) CreateVertexBuffer(vertices []float32) *gl.GLVertexBuffer {
	vertices_data := f32.Bytes(binary.LittleEndian, vertices...)

	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, vertices_data, gl.STATIC_DRAW)
	this._buffersCache._cachedVertexBuffer = nil

	vbobuf := &gl.GLVertexBuffer{}
	vbobuf.Vbo = vbo
	vbobuf.References = 1

	return vbobuf

}

func (this *Engine) CreateDynamicVertexBuffer(capacity int) *gl.GLVertexBuffer {

	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferInit(gl.ARRAY_BUFFER, capacity, gl.DYNAMIC_DRAW)
	this._buffersCache._cachedVertexBuffer = nil

	vbobuf := &gl.GLVertexBuffer{}
	vbobuf.Vbo = vbo
	vbobuf.References = 1

	return vbobuf
}

func (this *Engine) UpdateDynamicVertexBuffer(vertexBuffer *gl.GLVertexBuffer, vertices []float32) {
	vertices_data := f32.Bytes(binary.LittleEndian, vertices...)

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer.Vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, vertices_data)

}

func (this *Engine) CreateIndexBuffer(indices []uint16, is32Bits bool) *gl.GLIndexBuffer {
	indices_data := tools.BytesUint16(binary.LittleEndian, indices...)

	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices_data, gl.STATIC_DRAW)
	this._buffersCache._cachedIndexBuffer = nil

	vbobuf := &gl.GLIndexBuffer{}
	vbobuf.Vbo = vbo
	vbobuf.References = 1
	vbobuf.Is32Bits = is32Bits
	return vbobuf
}

func (this *Engine) BindBuffers(vertexBuffer *gl.GLVertexBuffer, indexBuffer *gl.GLIndexBuffer, vertexDeclaration []int, vertexStrideSize int, effect IEffect) {
	//VertexBuffer
	if !reflect.DeepEqual(this._buffersCache._cachedVertexBuffer, vertexBuffer) {
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer.Vbo)
		this._buffersCache._cachedVertexBuffer = vertexBuffer

		offset := 0
		for index := 0; index < len(vertexDeclaration); index++ {
			order := effect.GetAttribute(index)
			if order.Valid() {
				gl.VertexAttribPointer(order, vertexDeclaration[index], gl.FLOAT, false, vertexStrideSize, offset)
			}

			offset += vertexDeclaration[index] * 4
		}
	}

	//IndexBuffer
	if !reflect.DeepEqual(this._buffersCache._cachedIndexBuffer, indexBuffer) {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer.Vbo)
		this._buffersCache._cachedIndexBuffer = indexBuffer
	}

	if err := gl.GetError(); err != 0 {
		log.Printf("Draw gl error: %v \r\n", err)
	}
}

func (this *Engine) BindMultiBuffers(vertexBuffers map[string]*gl.GLVertexBuffer, indexBuffer *gl.GLIndexBuffer, effect IEffect) {

	if !reflect.DeepEqual(this._buffersCache._cachedVertexBuffers, vertexBuffers) ||
		!reflect.DeepEqual(this._buffersCache._cachedEffectForVertexBuffers, effect) {
		this._buffersCache._cachedVertexBuffers = vertexBuffers
		this._buffersCache._cachedEffectForVertexBuffers = effect

		attributes := effect.GetAttributesNames()

		for index := 0; index < len(attributes); index++ {
			order := effect.GetAttribute(index)

			if order.Valid() {
				vertexBuffer := vertexBuffers[attributes[index]]
				stride := vertexBuffer.StrideSize
				gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer.Vbo)
				gl.VertexAttribPointer(order, stride, gl.FLOAT, false, stride*4, 0)
			}
		}
	}

	if !reflect.DeepEqual(this._buffersCache._cachedIndexBuffer, indexBuffer) {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer.Vbo)
		this._buffersCache._cachedIndexBuffer = indexBuffer
	}

	if err := gl.GetError(); err != 0 {
		log.Printf("BindMultiBuffers gl error: %v \r\n", err)
	}
}

func (this *Engine) ReleaseIndexBuffer(buffer *gl.GLIndexBuffer) {
	buffer.References--

	if buffer.References == 0 {
		gl.DeleteBuffer(buffer.Vbo)
	}
}

func (this *Engine) ReleaseVertexBuffer(buffer *gl.GLVertexBuffer) {
	buffer.References--

	if buffer.References == 0 {
		gl.DeleteBuffer(buffer.Vbo)
	}
}

func (this *Engine) Draw(useTriangles bool, indexStart, indexCount int) {
	var gltype gl.Enum
	if useTriangles {
		gltype = gl.TRIANGLES
	} else {
		gltype = gl.LINES
	}
	gl.DrawElements(gltype, indexCount, gl.UNSIGNED_SHORT, indexStart*2)

	if err := gl.GetError(); err != 0 {
		log.Printf("Draw gl error: %v \r\n", err)
	}
}

func (this *Engine) CreateShaderProgram(vertexCode, fragmentCode, defines string) gl.Program {

	if defines != "" {
		defines = defines + "\n"
	}

	vertexCode_str := defines + vertexCode
	fragmentCode_str := defines + fragmentCode

	shaderProgram, err := glutil.CreateProgram(vertexCode_str, fragmentCode_str)

	if err != nil {
		log.Print(err)
		return gl.Program{}
	}

	return shaderProgram

}

func (this *Engine) GetUniforms(shaderProgram gl.Program, uniformsNames []string) []gl.Uniform {
	results := make([]gl.Uniform, 0)

	for index := 0; index < len(uniformsNames); index++ {
		uniform := gl.GetUniformLocation(shaderProgram, uniformsNames[index])
		results = append(results, uniform)
	}

	return results
}

func (this *Engine) GetAttributes(shaderProgram gl.Program, attributesNames []string) []gl.Attrib {

	results := make([]gl.Attrib, 0)

	for index := 0; index < len(attributesNames); index++ {
		attri := gl.GetAttribLocation(shaderProgram, attributesNames[index])
		results = append(results, attri)
	}

	return results
}

func (this *Engine) EnableEffect(effect IEffect) {
	if effect == nil || effect.GetAttributesCount() == 0 || reflect.DeepEqual(this._currentEffect, effect) {
		return
	}
	this._buffersCache._cachedVertexBuffer = nil

	// Use program
	gl.UseProgram(effect.GetProgram())

	for index := 0; index < effect.GetAttributesCount(); index++ {
		// Attributes
		order := effect.GetAttribute(index)

		if order.Valid() {
			gl.EnableVertexAttribArray(effect.GetAttribute(index))
		}
	}

	this._currentEffect = effect
}

func (this *Engine) SetMatrix(uniform gl.Uniform, m *math32.Matrix4) {
	if !uniform.Valid() {
		log.Println("SetMatrix uniform.Valid Failed")
		return
	}
	gl.UniformMatrix4fv(uniform, m.ToArray32())
}

func (this *Engine) SetVector2(uniform gl.Uniform, v *math32.Vector2) {
	if !uniform.Valid() {
		log.Println("SetVector2 uniform.Valid Failed")
		return
	}
	gl.Uniform2f(uniform, v.X, v.Y)
}

func (this *Engine) SetVector3(uniform gl.Uniform, v *math32.Vector3) {
	if !uniform.Valid() {
		log.Println("SetVector3 uniform.Valid Failed")
		return
	}
	gl.Uniform3f(uniform, v.X, v.Y, v.Z)
}

func (this *Engine) SetFloat2(uniform gl.Uniform, x, y float32) {
	if !uniform.Valid() {
		log.Println("SetFloat2 uniform.Valid Failed")
		return
	}
	gl.Uniform2f(uniform, x, y)
}

func (this *Engine) SetFloat3(uniform gl.Uniform, x, y, z float32) {
	if !uniform.Valid() {
		log.Println("SetFloat3 uniform.Valid Failed")
		return
	}
	gl.Uniform3f(uniform, x, y, z)
}

func (this *Engine) SetBool(uniform gl.Uniform, b bool) {
	if !uniform.Valid() {
		log.Println("SetBool uniform.Valid Failed")
		return
	}
	bval := 0
	if b {
		bval = 1
	}
	gl.Uniform1i(uniform, bval)
}

func (this *Engine) SetFloat4(uniform gl.Uniform, x, y, z, w float32) {
	if !uniform.Valid() {
		log.Println("SetFloat4 uniform.Valid Failed")
		return
	}
	gl.Uniform4f(uniform, x, y, z, w)
}

func (this *Engine) SetColor3(uniform gl.Uniform, v *math32.Color3) {
	if !uniform.Valid() {
		log.Println("SetColor3 uniform.Valid Failed")
		return
	}
	gl.Uniform3f(uniform, v.R, v.G, v.B)
}

func (this *Engine) SetColor4(uniform gl.Uniform, v *math32.Color4) {
	if !uniform.Valid() {
		log.Println("SetColor4 uniform.Valid Failed")
		return
	}
	gl.Uniform4f(uniform, v.R, v.G, v.B, v.A)
}

// States

func (this *Engine) SetState(culling bool) {
	// Culling
	if this._currentState.Culling != culling {
		if culling {
			var culltype gl.Enum
			//Fix Gl.Front
			if this.CullBackFaces == true {
				culltype = gl.BACK
			} else {
				culltype = gl.FRONT
			}
			gl.CullFace(culltype)
			gl.Enable(gl.CULL_FACE)
		} else {
			gl.Disable(gl.CULL_FACE)
		}

		this._currentState.Culling = culling
	}
}
func (this *Engine) SetDepthBuffer(enable bool) {
	if enable {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}
func (this *Engine) SetDepthWrite(enable bool) {
	gl.DepthMask(enable)
}
func (this *Engine) SetColorWrite(enable bool) {
	gl.ColorMask(enable, enable, enable, enable)
}

func (this *Engine) SetAlphaMode(mode int) {
	switch mode {
	case core.ALPHA_DISABLE:
		this.SetDepthWrite(true)
		gl.Disable(gl.BLEND)
		break
	case core.ALPHA_COMBINE:
		this.SetDepthWrite(false)
		gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ZERO, gl.ONE)
		gl.Enable(gl.BLEND)
		break
	case core.ALPHA_ADD:
		this.SetDepthWrite(false)
		gl.BlendFuncSeparate(gl.ONE, gl.ONE, gl.ZERO, gl.ONE)
		gl.Enable(gl.BLEND)
		break
	}
}
func (this *Engine) SetAlphaTesting(enable bool) {
	this._alphaTest = enable
}
func (this *Engine) GetAlphaTesting() bool {
	return this._alphaTest
}

func (this *Engine) WipeCaches() {
	this._activeTexturesCache = make([]*gl.GLTextureBuffer, 0)
	this._currentEffect = nil
	this._currentState = &gl.GLCullState{
		Culling: false,
	}
	this._buffersCache = &EngineBuffersCache{}
}

func (this *Engine) GetExponantOfTwo(value int, max int) int {
	count := 1

	count = count * 2
	for count < value {
		count = count * 2
	}

	if count > max {
		count = max
	}

	return count
}

func (this *Engine) GetScaled(img *image.RGBA, newWidth int, newHeight int) *image.RGBA {
	m := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
	img, _ = m.(*image.RGBA)
	return img
}

func (this *Engine) CreateTexture(url string, noMipmap bool, invertY int, scene *Scene) *gl.GLTextureBuffer {

	texture := gl.NewGLTextureBuffer()
	texture.Tex = gl.CreateTexture()
	texture.IsReady = false
	texture.Url = url
	texture.NoMipmap = noMipmap
	texture.References = 1

	onload := func(img *image.RGBA) {
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()
		canvas_width := this.GetExponantOfTwo(width, int(this._caps.MaxTextureSize))
		canvas_height := this.GetExponantOfTwo(height, int(this._caps.MaxTextureSize))
		pixelData := img.Pix

		isPot := (width == canvas_width && height == canvas_height)
		if !isPot {

			img = this.GetScaled(img, canvas_width, canvas_height)
			pixelData = img.Pix
		}

		gl.BindTexture(gl.TEXTURE_2D, texture.Tex)

		gl.TexImage2D(gl.TEXTURE_2D, 0, canvas_width, canvas_height, gl.RGBA, gl.UNSIGNED_BYTE, pixelData)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		if noMipmap == true {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		} else {
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
			gl.GenerateMipmap(gl.TEXTURE_2D)
		}

		this._activeTexturesCache = make([]*gl.GLTextureBuffer, 0)
		texture.BaseWidth = (int)(width)
		texture.BaseHeight = (int)(height)
		texture.Width = (int)(canvas_width)
		texture.Height = (int)(canvas_height)
		texture.IsReady = true
		scene.RemovePendingData(url)

	}
	onfailed := func(err error) {
		scene.RemovePendingData(url)
	}
	scene.AddPendingData(url)
	tools.LoadImage(url, onload, onfailed)

	this._loadedTexturesCache = append(this._loadedTexturesCache, texture)

	return texture

}

func (this *Engine) CreateDynamicTexture(size int, generateMipMaps bool) *gl.GLTextureBuffer {

	texture := gl.NewGLTextureBuffer()
	texture.Tex = gl.CreateTexture()

	width := this.GetExponantOfTwo(size, (int)(this._caps.MaxTextureSize))
	height := this.GetExponantOfTwo(size, (int)(this._caps.MaxTextureSize))

	gl.BindTexture(gl.TEXTURE_2D, texture.Tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	if !generateMipMaps {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	}

	this._activeTexturesCache = make([]*gl.GLTextureBuffer, 0)

	texture.BaseWidth = (int)(width)
	texture.BaseHeight = (int)(height)
	texture.Width = (int)(width)
	texture.Height = (int)(height)

	texture.IsReady = false
	texture.GenerateMipMaps = generateMipMaps
	texture.References = 1

	this._loadedTexturesCache = append(this._loadedTexturesCache, texture)

	return texture
}

func (this *Engine) UpdateDynamicTexture(texture *gl.GLTextureBuffer, img *image.RGBA) {
	gl.BindTexture(gl.TEXTURE_2D, texture.Tex)
	pixelData := img.Pix

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D, 0, width, height, gl.RGBA, gl.UNSIGNED_BYTE, pixelData)
	if texture.GenerateMipMaps {
		gl.GenerateMipmap(gl.TEXTURE_2D)
	}

	this._activeTexturesCache = make([]*gl.GLTextureBuffer, 0)
	texture.IsReady = true
}
func (this *Engine) UpdateVideoTexture(texture *gl.GLTextureBuffer, img *image.RGBA) {

	this.UpdateDynamicTexture(texture, img)
}

func (this *Engine) CreateRenderTargetTexture(size int, generateMipMaps bool) *gl.GLTextureBuffer {
	texture := gl.NewGLTextureBuffer()
	texture.Tex = gl.CreateTexture()

	var minFilter int
	if generateMipMaps {
		minFilter = gl.LINEAR_MIPMAP_NEAREST
	} else {
		minFilter = gl.LINEAR
	}

	gl.BindTexture(gl.TEXTURE_2D, texture.Tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, minFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, size, size, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	// Create the depth buffer
	depthBuffer := gl.CreateRenderbuffer()
	gl.BindRenderbuffer(gl.RENDERBUFFER, depthBuffer)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT16, size, size)

	// Create the framebuffer
	framebuffer := gl.CreateFramebuffer()
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture.Tex, 0)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, depthBuffer)

	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		log.Printf("framebuffer create failed: %v", status)
		return texture
	}

	texture.FrameBuf = framebuffer
	texture.DepthBuf = depthBuffer
	texture.Width = size
	texture.Height = size
	texture.IsReady = true
	texture.GenerateMipMaps = generateMipMaps
	texture.References = 1
	this._activeTexturesCache = make([]*gl.GLTextureBuffer, 0)

	this._loadedTexturesCache = append(this._loadedTexturesCache, texture)

	return texture
}

func cascadeLoad(scene *Scene, rootUrl string, extensions []string, index int, loadedImages []*image.RGBA, onfinish func([]*image.RGBA)) {

	url := rootUrl + extensions[index]
	scene.AddPendingData(url)

	tools.LoadImage(url, func(img *image.RGBA) {
		loadedImages = append(loadedImages, img)

		scene.RemovePendingData(url)

		if index != len(extensions)-1 {
			cascadeLoad(scene, rootUrl, extensions, index+1, loadedImages, onfinish)
		} else {
			onfinish(loadedImages)
		}

	}, func(error) {
		scene.RemovePendingData(url)
	})
}

func (this *Engine) CreateCubeTexture(rootUrl string, scene *Scene, extensions []string) *gl.GLTextureBuffer {
	if extensions == nil {
		extensions = []string{"_px.jpg", "_py.jpg", "_pz.jpg", "_nx.jpg", "_ny.jpg", "_nz.jpg"}
	}

	texture := gl.NewGLTextureBuffer()
	texture.Tex = gl.CreateTexture()

	texture.IsCube = true
	texture.Url = rootUrl
	texture.References = 1

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture.Tex)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	onfinish := func(imgs []*image.RGBA) {

		width := imgs[0].Bounds().Dx()
		height := width
		canvas_width := this.GetExponantOfTwo(width, int(this._caps.MaxTextureSize))
		canvas_height := canvas_width

		faces := []gl.Enum{
			gl.TEXTURE_CUBE_MAP_POSITIVE_X, gl.TEXTURE_CUBE_MAP_POSITIVE_Y, gl.TEXTURE_CUBE_MAP_POSITIVE_Z,
			gl.TEXTURE_CUBE_MAP_NEGATIVE_X, gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, gl.TEXTURE_CUBE_MAP_NEGATIVE_Z,
		}

		gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture.Tex)
		gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
		for index := 0; index < len(faces); index++ {

			img := imgs[index]
			pixelData := img.Pix
			isPot := (width == canvas_width && height == canvas_height)
			if !isPot {

				img = this.GetScaled(img, canvas_width, canvas_height)
				pixelData = img.Pix
			}

			gl.TexImage2D(faces[index], 0, canvas_width, canvas_height, gl.RGBA, gl.UNSIGNED_BYTE, pixelData)

		}

		gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

		gl.BindTexture(gl.TEXTURE_CUBE_MAP, gl.Texture{})

		this._activeTexturesCache = make([]*gl.GLTextureBuffer, 0)

		texture.Width = width
		texture.Height = height
		texture.IsReady = true

	}
	loadedImages := make([]*image.RGBA, 0)
	cascadeLoad(scene, rootUrl, extensions, 0, loadedImages, onfinish)

	return texture
}

func (this *Engine) ReleaseTexture(texture *gl.GLTextureBuffer) {

	if texture.FrameBuf.Valid() {
		gl.DeleteFramebuffer(texture.FrameBuf)
	}

	if texture.DepthBuf.Valid() {
		gl.DeleteRenderbuffer(texture.DepthBuf)
	}

	gl.DeleteTexture(texture.Tex)

	// Unbind channels
	for channel := 0; channel < this._caps.MaxTexturesImageUnits; channel++ {
		val := gl.TEXTURE0 + channel
		gl.ActiveTexture((gl.Enum)(val))
		gl.BindTexture(gl.TEXTURE_2D, gl.Texture{})
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, gl.Texture{})
		this._activeTexturesCache[channel] = nil
	}

}

func (this *Engine) BindSamplers(effect IEffect) {
	gl.UseProgram(effect.GetProgram())
	samplers := effect.GetSamplers()
	for index, sampler := range samplers {
		uniform := effect.GetUniform(sampler)
		gl.Uniform1i(uniform, index)
	}
	this._currentEffect = nil
}
func (this *Engine) SetTexture(channel int, texture *gl.GLTextureBuffer) {
	if texture == nil || !texture.IsReady {
		if this._activeTexturesCache[channel] != nil {
			val := gl.TEXTURE0 + channel
			gl.ActiveTexture((gl.Enum)(val))
			gl.BindTexture(gl.TEXTURE_2D, gl.Texture{})
			gl.BindTexture(gl.TEXTURE_CUBE_MAP, gl.Texture{})
			this._activeTexturesCache[channel] = nil
		}
		return
	}

	//更新
	if texture.UpdateFunc != nil {
		ret := texture.UpdateFunc()
		if !ret {
			this._activeTexturesCache[channel] = nil
		}
	}

	if len(this._activeTexturesCache) > 0 && reflect.DeepEqual(this._activeTexturesCache[channel], texture) {
		return
	}

	val := gl.TEXTURE0 + channel
	gl.ActiveTexture((gl.Enum)(val))

	if texture.IsCube {
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture.Tex)

		if texture.CacheCoordinatesMode != texture.CoordinatesMode {
			texture.CacheCoordinatesMode = texture.CoordinatesMode

			var texval int
			if texture.CoordinatesMode == gl.CUBIC_MODE {
				texval = gl.REPEAT
			} else {
				texval = gl.CLAMP_TO_EDGE
			}

			gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, texval)
			gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, texval)
		}
	} else {
		gl.BindTexture(gl.TEXTURE_2D, texture.Tex)

		if texture.CachedWrapU != texture.WrapU {
			texture.CachedWrapU = texture.WrapU

			switch texture.WrapU {
			case gl.WRAP_ADDRESSMODE:
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
				break
			case gl.CLAMP_ADDRESSMODE:
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
				break
			case gl.MIRROR_ADDRESSMODE:
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
				break
			}
		}

		if texture.CachedWrapV != texture.WrapV {
			texture.CachedWrapV = texture.WrapV
			switch texture.WrapV {
			case gl.WRAP_ADDRESSMODE:
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
				break
			case gl.CLAMP_ADDRESSMODE:
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
				break
			case gl.MIRROR_ADDRESSMODE:
				gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)
				break
			}
		}

	}
}

// Dispose
func (this *Engine) Dispose() {
	// Release scenes

	for _, scene := range this.Scenes {
		scene.Dispose()
	}
	this.Scenes = make([]*Scene, 0)

	for _, effect := range this.CompiledEffects {
		gl.DeleteProgram(effect.GetProgram())
	}

}
