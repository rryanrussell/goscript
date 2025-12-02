package app

type ContextId string

var (
	_2d            ContextId = "2d"
	Bitmaprenderer ContextId = "bitmaprenderer"
	Webgl          ContextId = "webgl"
	Webgl2         ContextId = "webgl2"
)

type ExtensionName string

var (
	EXT_blend_minmax_                  ExtensionName = "EXT_blend_minmax"
	EXT_color_buffer_float             ExtensionName = "EXT_color_buffer_float"
	EXT_color_buffer_half_float        ExtensionName = "EXT_color_buffer_half_float"
	EXT_float_blend                    ExtensionName = "EXT_float_blend"
	EXT_texture_filter_anisotropic     ExtensionName = "EXT_texture_filter_anisotropic"
	EXT_frag_depth                     ExtensionName = "EXT_frag_depth"
	EXT_shader_texture_lod             ExtensionName = "EXT_shader_texture_lod"
	EXT_sRGB                           ExtensionName = "EXT_sRGB"
	KHR_parallel_shader_compile        ExtensionName = "KHR_parallel_shader_compile"
	OES_vertex_array_object            ExtensionName = "OES_vertex_array_object"
	OVR_multiview2                     ExtensionName = "OVR_multiview2"
	WEBGL_color_buffer_float           ExtensionName = "WEBGL_color_buffer_float"
	WEBGL_compressed_texture_astc      ExtensionName = "WEBGL_compressed_texture_astc"
	WEBGL_compressed_texture_etc       ExtensionName = "WEBGL_compressed_texture_etc"
	WEBGL_compressed_texture_etc1      ExtensionName = "WEBGL_compressed_texture_etc1"
	WEBGL_compressed_texture_pvrtc     ExtensionName = "WEBGL_compressed_texture_pvrtc"
	WEBGL_compressed_texture_s3tc_srgb ExtensionName = "WEBGL_compressed_texture_s3tc_srgb"
	WEBGL_debug_shaders                ExtensionName = "WEBGL_debug_shaders"
	WEBGL_draw_buffers                 ExtensionName = "WEBGL_draw_buffers"
	WEBGL_lose_context                 ExtensionName = "WEBGL_lose_context"
	WEBGL_depth_texture                ExtensionName = "WEBGL_depth_texture"
	WEBGL_debug_renderer_info          ExtensionName = "WEBGL_debug_renderer_info"
	WEBGL_compressed_texture_s3tc      ExtensionName = "WEBGL_compressed_texture_s3tc"
	OES_texture_half_float_linear      ExtensionName = "OES_texture_half_float_linear"
	OES_texture_half_float             ExtensionName = "OES_texture_half_float"
	OES_texture_float_linear           ExtensionName = "OES_texture_float_linear"
	OES_texture_float                  ExtensionName = "OES_texture_float"
	OES_standard_derivatives           ExtensionName = "OES_standard_derivatives"
	OES_element_index_uint             ExtensionName = "OES_element_index_uint"
	ANGLE_instanced_arrays             ExtensionName = "ANGLE_instanced_arrays"
)

type EventInterface string

var (
	AnimationEvent_                EventInterface = "AnimationEvent"
	AnimationPlaybackEvent         EventInterface = "AnimationPlaybackEvent"
	AudioProcessingEvent           EventInterface = "AudioProcessingEvent"
	BeforeUnloadEvent              EventInterface = "BeforeUnloadEvent"
	BlobEvent                      EventInterface = "BlobEvent"
	ClipboardEvent                 EventInterface = "ClipboardEvent"
	CloseEvent                     EventInterface = "CloseEvent"
	CompositionEvent               EventInterface = "CompositionEvent"
	CustomEvent                    EventInterface = "CustomEvent"
	DeviceMotionEvent              EventInterface = "DeviceMotionEvent"
	DeviceOrientationEvent         EventInterface = "DeviceOrientationEvent"
	DragEvent                      EventInterface = "DragEvent"
	ErrorEvent                     EventInterface = "ErrorEvent"
	FocusEvent                     EventInterface = "FocusEvent"
	FontFaceSetLoadEvent           EventInterface = "FontFaceSetLoadEvent"
	FormDataEvent                  EventInterface = "FormDataEvent"
	GamepadEvent                   EventInterface = "GamepadEvent"
	HashChangeEvent                EventInterface = "HashChangeEvent"
	IDBVersionChangeEvent          EventInterface = "IDBVersionChangeEvent"
	InputEvent                     EventInterface = "InputEvent"
	KeyboardEvent                  EventInterface = "KeyboardEvent"
	MediaEncryptedEvent            EventInterface = "MediaEncryptedEvent"
	MediaKeyMessageEvent           EventInterface = "MediaKeyMessageEvent"
	MediaQueryListEvent            EventInterface = "MediaQueryListEvent"
	MediaRecorderErrorEvent        EventInterface = "MediaRecorderErrorEvent"
	MediaStreamTrackEvent          EventInterface = "MediaStreamTrackEvent"
	MessageEvent                   EventInterface = "MessageEvent"
	MouseEvent                     EventInterface = "MouseEvent"
	MouseEvents                    EventInterface = "MouseEvents"
	MutationEvent                  EventInterface = "MutationEvent"
	MutationEvents                 EventInterface = "MutationEvents"
	OfflineAudioCompletionEvent    EventInterface = "OfflineAudioCompletionEvent"
	PageTransitionEvent            EventInterface = "PageTransitionEvent"
	PaymentMethodChangeEvent       EventInterface = "PaymentMethodChangeEvent"
	PaymentRequestUpdateEvent      EventInterface = "PaymentRequestUpdateEvent"
	PointerEvent                   EventInterface = "PointerEvent"
	PopStateEvent                  EventInterface = "PopStateEvent"
	ProgressEvent                  EventInterface = "ProgressEvent"
	PromiseRejectionEvent          EventInterface = "PromiseRejectionEvent"
	RTCDTMFToneChangeEvent         EventInterface = "RTCDTMFToneChangeEvent"
	RTCDataChannelEvent            EventInterface = "RTCDataChannelEvent"
	RTCPeerConnectionIceErrorEvent EventInterface = "RTCPeerConnectionIceErrorEvent"
	RTCPeerConnectionIceEvent      EventInterface = "RTCPeerConnectionIceEvent"
	RTCTrackEvent                  EventInterface = "RTCTrackEvent"
	SecurityPolicyViolationEvent   EventInterface = "SecurityPolicyViolationEvent"
	SpeechSynthesisErrorEvent      EventInterface = "SpeechSynthesisErrorEvent"
	SpeechSynthesisEvent           EventInterface = "SpeechSynthesisEvent"
	StorageEvent                   EventInterface = "StorageEvent"
	SubmitEvent                    EventInterface = "SubmitEvent"
	TouchEvent                     EventInterface = "TouchEvent"
	TrackEvent                     EventInterface = "TrackEvent"
	TransitionEvent                EventInterface = "TransitionEvent"
	UIEvent                        EventInterface = "UIEvent"
	UIEvents                       EventInterface = "UIEvents"
	WebGLContextEvent              EventInterface = "WebGLContextEvent"
	WheelEvent                     EventInterface = "WheelEvent"
)

type HTMLCanvasElement struct {
	HTMLElement
	// rename:height
	Height float64
	// rename:width
	Width float64
}

// rename:captureStream
func (HTMLCanvasElement) CaptureStream(frameRequestRate float64) MediaStream { panic("extern") }

// rename:getContext
func (HTMLCanvasElement) GetContext(contextId ContextId, options CanvasRenderingContext2DSettings) CanvasRenderingContext2D {
	panic("extern")
}

// rename:getContext
func (HTMLCanvasElement) GetContext0(contextId ContextId, options ImageBitmapRenderingContextSettings) ImageBitmapRenderingContext {
	panic("extern")
}

// rename:getContext
func (HTMLCanvasElement) GetContext1(contextId ContextId, options WebGLContextAttributes) WebGLRenderingContext {
	panic("extern")
}

// rename:getContext
func (HTMLCanvasElement) GetContext2(contextId string, options any) RenderingContext { panic("extern") }

// rename:toBlob
func (HTMLCanvasElement) ToBlob(callback BlobCallback, type0 string, quality any) { panic("extern") }

// rename:toDataURL
func (HTMLCanvasElement) ToDataURL(type0 string, quality any) string { panic("extern") }

// rename:addEventListener
func (HTMLCanvasElement) AddEventListener(type0 any, listener any, options any) { panic("extern") }

// rename:addEventListener
func (HTMLCanvasElement) AddEventListener0(type0 string, listener EventListenerOrEventListenerObject, options any) {
	panic("extern")
}

// rename:removeEventListener
func (HTMLCanvasElement) RemoveEventListener(type0 any, listener any, options any) { panic("extern") }

// rename:removeEventListener
func (HTMLCanvasElement) RemoveEventListener0(type0 string, listener EventListenerOrEventListenerObject, options any) {
	panic("extern")
}

type WebGLRenderingContext struct {
	WebGLRenderingContextBase
	WebGLRenderingContextOverloads
}

type WebGLRenderingContextBase struct {
	// rename:canvas
	Canvas HTMLCanvasElement
	// rename:drawingBufferHeight
	DrawingBufferHeight GLsizei
	// rename:drawingBufferWidth
	DrawingBufferWidth                           GLsizei
	ACTIVE_ATTRIBUTES                            GLenum
	ACTIVE_TEXTURE                               GLenum
	ACTIVE_UNIFORMS                              GLenum
	ALIASED_LINE_WIDTH_RANGE                     GLenum
	ALIASED_POINT_SIZE_RANGE                     GLenum
	ALPHA                                        GLenum
	ALPHA_BITS                                   GLenum
	ALWAYS                                       GLenum
	ARRAY_BUFFER                                 GLenum
	ARRAY_BUFFER_BINDING                         GLenum
	ATTACHED_SHADERS                             GLenum
	BACK                                         GLenum
	BLEND                                        GLenum
	BLEND_COLOR                                  GLenum
	BLEND_DST_ALPHA                              GLenum
	BLEND_DST_RGB                                GLenum
	BLEND_EQUATION                               GLenum
	BLEND_EQUATION_ALPHA                         GLenum
	BLEND_EQUATION_RGB                           GLenum
	BLEND_SRC_ALPHA                              GLenum
	BLEND_SRC_RGB                                GLenum
	BLUE_BITS                                    GLenum
	BOOL                                         GLenum
	BOOL_VEC2                                    GLenum
	BOOL_VEC3                                    GLenum
	BOOL_VEC4                                    GLenum
	BROWSER_DEFAULT_WEBGL                        GLenum
	BUFFER_SIZE                                  GLenum
	BUFFER_USAGE                                 GLenum
	BYTE                                         GLenum
	CCW                                          GLenum
	CLAMP_TO_EDGE                                GLenum
	COLOR_ATTACHMENT0                            GLenum
	COLOR_BUFFER_BIT                             GLenum
	COLOR_CLEAR_VALUE                            GLenum
	COLOR_WRITEMASK                              GLenum
	COMPILE_STATUS                               GLenum
	COMPRESSED_TEXTURE_FORMATS                   GLenum
	CONSTANT_ALPHA                               GLenum
	CONSTANT_COLOR                               GLenum
	CONTEXT_LOST_WEBGL                           GLenum
	CULL_FACE                                    GLenum
	CULL_FACE_MODE                               GLenum
	CURRENT_PROGRAM                              GLenum
	CURRENT_VERTEX_ATTRIB                        GLenum
	CW                                           GLenum
	DECR                                         GLenum
	DECR_WRAP                                    GLenum
	DELETE_STATUS                                GLenum
	DEPTH_ATTACHMENT                             GLenum
	DEPTH_BITS                                   GLenum
	DEPTH_BUFFER_BIT                             GLenum
	DEPTH_CLEAR_VALUE                            GLenum
	DEPTH_COMPONENT                              GLenum
	DEPTH_COMPONENT16                            GLenum
	DEPTH_FUNC                                   GLenum
	DEPTH_RANGE                                  GLenum
	DEPTH_STENCIL                                GLenum
	DEPTH_STENCIL_ATTACHMENT                     GLenum
	DEPTH_TEST                                   GLenum
	DEPTH_WRITEMASK                              GLenum
	DITHER                                       GLenum
	DONT_CARE                                    GLenum
	DST_ALPHA                                    GLenum
	DST_COLOR                                    GLenum
	DYNAMIC_DRAW                                 GLenum
	ELEMENT_ARRAY_BUFFER                         GLenum
	ELEMENT_ARRAY_BUFFER_BINDING                 GLenum
	EQUAL                                        GLenum
	FASTEST                                      GLenum
	FLOAT                                        GLenum
	FLOAT_MAT2                                   GLenum
	FLOAT_MAT3                                   GLenum
	FLOAT_MAT4                                   GLenum
	FLOAT_VEC2                                   GLenum
	FLOAT_VEC3                                   GLenum
	FLOAT_VEC4                                   GLenum
	FRAGMENT_SHADER                              GLenum
	FRAMEBUFFER                                  GLenum
	FRAMEBUFFER_ATTACHMENT_OBJECT_NAME           GLenum
	FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE           GLenum
	FRAMEBUFFER_ATTACHMENT_TEXTURE_CUBE_MAP_FACE GLenum
	FRAMEBUFFER_ATTACHMENT_TEXTURE_LEVEL         GLenum
	FRAMEBUFFER_BINDING                          GLenum
	FRAMEBUFFER_COMPLETE                         GLenum
	FRAMEBUFFER_INCOMPLETE_ATTACHMENT            GLenum
	FRAMEBUFFER_INCOMPLETE_DIMENSIONS            GLenum
	FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT    GLenum
	FRAMEBUFFER_UNSUPPORTED                      GLenum
	FRONT                                        GLenum
	FRONT_AND_BACK                               GLenum
	FRONT_FACE                                   GLenum
	FUNC_ADD                                     GLenum
	FUNC_REVERSE_SUBTRACT                        GLenum
	FUNC_SUBTRACT                                GLenum
	GENERATE_MIPMAP_HINT                         GLenum
	GEQUAL                                       GLenum
	GREATER                                      GLenum
	GREEN_BITS                                   GLenum
	HIGH_FLOAT                                   GLenum
	HIGH_INT                                     GLenum
	IMPLEMENTATION_COLOR_READ_FORMAT             GLenum
	IMPLEMENTATION_COLOR_READ_TYPE               GLenum
	INCR                                         GLenum
	INCR_WRAP                                    GLenum
	INT                                          GLenum
	INT_VEC2                                     GLenum
	INT_VEC3                                     GLenum
	INT_VEC4                                     GLenum
	INVALID_ENUM                                 GLenum
	INVALID_FRAMEBUFFER_OPERATION                GLenum
	INVALID_OPERATION                            GLenum
	INVALID_VALUE                                GLenum
	INVERT                                       GLenum
	KEEP                                         GLenum
	LEQUAL                                       GLenum
	LESS                                         GLenum
	LINEAR                                       GLenum
	LINEAR_MIPMAP_LINEAR                         GLenum
	LINEAR_MIPMAP_NEAREST                        GLenum
	LINES                                        GLenum
	LINE_LOOP                                    GLenum
	LINE_STRIP                                   GLenum
	LINE_WIDTH                                   GLenum
	LINK_STATUS                                  GLenum
	LOW_FLOAT                                    GLenum
	LOW_INT                                      GLenum
	LUMINANCE                                    GLenum
	LUMINANCE_ALPHA                              GLenum
	MAX_COMBINED_TEXTURE_IMAGE_UNITS             GLenum
	MAX_CUBE_MAP_TEXTURE_SIZE                    GLenum
	MAX_FRAGMENT_UNIFORM_VECTORS                 GLenum
	MAX_RENDERBUFFER_SIZE                        GLenum
	MAX_TEXTURE_IMAGE_UNITS                      GLenum
	MAX_TEXTURE_SIZE                             GLenum
	MAX_VARYING_VECTORS                          GLenum
	MAX_VERTEX_ATTRIBS                           GLenum
	MAX_VERTEX_TEXTURE_IMAGE_UNITS               GLenum
	MAX_VERTEX_UNIFORM_VECTORS                   GLenum
	MAX_VIEWPORT_DIMS                            GLenum
	MEDIUM_FLOAT                                 GLenum
	MEDIUM_INT                                   GLenum
	MIRRORED_REPEAT                              GLenum
	NEAREST                                      GLenum
	NEAREST_MIPMAP_LINEAR                        GLenum
	NEAREST_MIPMAP_NEAREST                       GLenum
	NEVER                                        GLenum
	NICEST                                       GLenum
	NONE                                         GLenum
	NOTEQUAL                                     GLenum
	NO_ERROR                                     GLenum
	ONE                                          GLenum
	ONE_MINUS_CONSTANT_ALPHA                     GLenum
	ONE_MINUS_CONSTANT_COLOR                     GLenum
	ONE_MINUS_DST_ALPHA                          GLenum
	ONE_MINUS_DST_COLOR                          GLenum
	ONE_MINUS_SRC_ALPHA                          GLenum
	ONE_MINUS_SRC_COLOR                          GLenum
	OUT_OF_MEMORY                                GLenum
	PACK_ALIGNMENT                               GLenum
	POINTS                                       GLenum
	POLYGON_OFFSET_FACTOR                        GLenum
	POLYGON_OFFSET_FILL                          GLenum
	POLYGON_OFFSET_UNITS                         GLenum
	RED_BITS                                     GLenum
	RENDERBUFFER                                 GLenum
	RENDERBUFFER_ALPHA_SIZE                      GLenum
	RENDERBUFFER_BINDING                         GLenum
	RENDERBUFFER_BLUE_SIZE                       GLenum
	RENDERBUFFER_DEPTH_SIZE                      GLenum
	RENDERBUFFER_GREEN_SIZE                      GLenum
	RENDERBUFFER_HEIGHT                          GLenum
	RENDERBUFFER_INTERNAL_FORMAT                 GLenum
	RENDERBUFFER_RED_SIZE                        GLenum
	RENDERBUFFER_STENCIL_SIZE                    GLenum
	RENDERBUFFER_WIDTH                           GLenum
	RENDERER                                     GLenum
	REPEAT                                       GLenum
	REPLACE                                      GLenum
	RGB                                          GLenum
	RGB565                                       GLenum
	RGB5_A1                                      GLenum
	RGBA                                         GLenum
	RGBA4                                        GLenum
	SAMPLER_2D                                   GLenum
	SAMPLER_CUBE                                 GLenum
	SAMPLES                                      GLenum
	SAMPLE_ALPHA_TO_COVERAGE                     GLenum
	SAMPLE_BUFFERS                               GLenum
	SAMPLE_COVERAGE                              GLenum
	SAMPLE_COVERAGE_INVERT                       GLenum
	SAMPLE_COVERAGE_VALUE                        GLenum
	SCISSOR_BOX                                  GLenum
	SCISSOR_TEST                                 GLenum
	SHADER_TYPE                                  GLenum
	SHADING_LANGUAGE_VERSION                     GLenum
	SHORT                                        GLenum
	SRC_ALPHA                                    GLenum
	SRC_ALPHA_SATURATE                           GLenum
	SRC_COLOR                                    GLenum
	STATIC_DRAW                                  GLenum
	STENCIL_ATTACHMENT                           GLenum
	STENCIL_BACK_FAIL                            GLenum
	STENCIL_BACK_FUNC                            GLenum
	STENCIL_BACK_PASS_DEPTH_FAIL                 GLenum
	STENCIL_BACK_PASS_DEPTH_PASS                 GLenum
	STENCIL_BACK_REF                             GLenum
	STENCIL_BACK_VALUE_MASK                      GLenum
	STENCIL_BACK_WRITEMASK                       GLenum
	STENCIL_BITS                                 GLenum
	STENCIL_BUFFER_BIT                           GLenum
	STENCIL_CLEAR_VALUE                          GLenum
	STENCIL_FAIL                                 GLenum
	STENCIL_FUNC                                 GLenum
	STENCIL_INDEX8                               GLenum
	STENCIL_PASS_DEPTH_FAIL                      GLenum
	STENCIL_PASS_DEPTH_PASS                      GLenum
	STENCIL_REF                                  GLenum
	STENCIL_TEST                                 GLenum
	STENCIL_VALUE_MASK                           GLenum
	STENCIL_WRITEMASK                            GLenum
	STREAM_DRAW                                  GLenum
	SUBPIXEL_BITS                                GLenum
	TEXTURE                                      GLenum
	TEXTURE0                                     GLenum
	TEXTURE1                                     GLenum
	TEXTURE10                                    GLenum
	TEXTURE11                                    GLenum
	TEXTURE12                                    GLenum
	TEXTURE13                                    GLenum
	TEXTURE14                                    GLenum
	TEXTURE15                                    GLenum
	TEXTURE16                                    GLenum
	TEXTURE17                                    GLenum
	TEXTURE18                                    GLenum
	TEXTURE19                                    GLenum
	TEXTURE2                                     GLenum
	TEXTURE20                                    GLenum
	TEXTURE21                                    GLenum
	TEXTURE22                                    GLenum
	TEXTURE23                                    GLenum
	TEXTURE24                                    GLenum
	TEXTURE25                                    GLenum
	TEXTURE26                                    GLenum
	TEXTURE27                                    GLenum
	TEXTURE28                                    GLenum
	TEXTURE29                                    GLenum
	TEXTURE3                                     GLenum
	TEXTURE30                                    GLenum
	TEXTURE31                                    GLenum
	TEXTURE4                                     GLenum
	TEXTURE5                                     GLenum
	TEXTURE6                                     GLenum
	TEXTURE7                                     GLenum
	TEXTURE8                                     GLenum
	TEXTURE9                                     GLenum
	TEXTURE_2D                                   GLenum
	TEXTURE_BINDING_2D                           GLenum
	TEXTURE_BINDING_CUBE_MAP                     GLenum
	TEXTURE_CUBE_MAP                             GLenum
	TEXTURE_CUBE_MAP_NEGATIVE_X                  GLenum
	TEXTURE_CUBE_MAP_NEGATIVE_Y                  GLenum
	TEXTURE_CUBE_MAP_NEGATIVE_Z                  GLenum
	TEXTURE_CUBE_MAP_POSITIVE_X                  GLenum
	TEXTURE_CUBE_MAP_POSITIVE_Y                  GLenum
	TEXTURE_CUBE_MAP_POSITIVE_Z                  GLenum
	TEXTURE_MAG_FILTER                           GLenum
	TEXTURE_MIN_FILTER                           GLenum
	TEXTURE_WRAP_S                               GLenum
	TEXTURE_WRAP_T                               GLenum
	TRIANGLES                                    GLenum
	TRIANGLE_FAN                                 GLenum
	TRIANGLE_STRIP                               GLenum
	UNPACK_ALIGNMENT                             GLenum
	UNPACK_COLORSPACE_CONVERSION_WEBGL           GLenum
	UNPACK_FLIP_Y_WEBGL                          GLenum
	UNPACK_PREMULTIPLY_ALPHA_WEBGL               GLenum
	UNSIGNED_BYTE                                GLenum
	UNSIGNED_INT                                 GLenum
	UNSIGNED_SHORT                               GLenum
	UNSIGNED_SHORT_4_4_4_4                       GLenum
	UNSIGNED_SHORT_5_5_5_1                       GLenum
	UNSIGNED_SHORT_5_6_5                         GLenum
	VALIDATE_STATUS                              GLenum
	VENDOR                                       GLenum
	VERSION                                      GLenum
	VERTEX_ATTRIB_ARRAY_BUFFER_BINDING           GLenum
	VERTEX_ATTRIB_ARRAY_ENABLED                  GLenum
	VERTEX_ATTRIB_ARRAY_NORMALIZED               GLenum
	VERTEX_ATTRIB_ARRAY_POINTER                  GLenum
	VERTEX_ATTRIB_ARRAY_SIZE                     GLenum
	VERTEX_ATTRIB_ARRAY_STRIDE                   GLenum
	VERTEX_ATTRIB_ARRAY_TYPE                     GLenum
	VERTEX_SHADER                                GLenum
	VIEWPORT                                     GLenum
	ZERO                                         GLenum
}

// rename:activeTexture
func (WebGLRenderingContextBase) ActiveTexture(texture GLenum) { panic("extern") }

// rename:attachShader
func (WebGLRenderingContextBase) AttachShader(program WebGLProgram, shader WebGLShader) {
	panic("extern")
}

// rename:bindAttribLocation
func (WebGLRenderingContextBase) BindAttribLocation(program WebGLProgram, index GLuint, name string) {
	panic("extern")
}

// rename:bindBuffer
func (WebGLRenderingContextBase) BindBuffer(target GLenum, buffer WebGLBuffer) { panic("extern") }

// rename:bindFramebuffer
func (WebGLRenderingContextBase) BindFramebuffer(target GLenum, framebuffer WebGLFramebuffer) {
	panic("extern")
}

// rename:bindRenderbuffer
func (WebGLRenderingContextBase) BindRenderbuffer(target GLenum, renderbuffer WebGLRenderbuffer) {
	panic("extern")
}

// rename:bindTexture
func (WebGLRenderingContextBase) BindTexture(target GLenum, texture WebGLTexture) { panic("extern") }

// rename:blendColor
func (WebGLRenderingContextBase) BlendColor(red GLclampf, green GLclampf, blue GLclampf, alpha GLclampf) {
	panic("extern")
}

// rename:blendEquation
func (WebGLRenderingContextBase) BlendEquation(mode GLenum) { panic("extern") }

// rename:blendEquationSeparate
func (WebGLRenderingContextBase) BlendEquationSeparate(modeRGB GLenum, modeAlpha GLenum) {
	panic("extern")
}

// rename:blendFunc
func (WebGLRenderingContextBase) BlendFunc(sfactor GLenum, dfactor GLenum) { panic("extern") }

// rename:blendFuncSeparate
func (WebGLRenderingContextBase) BlendFuncSeparate(srcRGB GLenum, dstRGB GLenum, srcAlpha GLenum, dstAlpha GLenum) {
	panic("extern")
}

// rename:checkFramebufferStatus
func (WebGLRenderingContextBase) CheckFramebufferStatus(target GLenum) GLenum { panic("extern") }

// rename:clear
func (WebGLRenderingContextBase) Clear(mask GLbitfield) { panic("extern") }

// rename:clearColor
func (WebGLRenderingContextBase) ClearColor(red GLclampf, green GLclampf, blue GLclampf, alpha GLclampf) {
	panic("extern")
}

// rename:clearDepth
func (WebGLRenderingContextBase) ClearDepth(depth GLclampf) { panic("extern") }

// rename:clearStencil
func (WebGLRenderingContextBase) ClearStencil(s GLint) { panic("extern") }

// rename:colorMask
func (WebGLRenderingContextBase) ColorMask(red GLboolean, green GLboolean, blue GLboolean, alpha GLboolean) {
	panic("extern")
}

// rename:compileShader
func (WebGLRenderingContextBase) CompileShader(shader WebGLShader) { panic("extern") }

// rename:copyTexImage2D
func (WebGLRenderingContextBase) CopyTexImage2D(target GLenum, level GLint, internalformat GLenum, x GLint, y GLint, width GLsizei, height GLsizei, border GLint) {
	panic("extern")
}

// rename:copyTexSubImage2D
func (WebGLRenderingContextBase) CopyTexSubImage2D(target GLenum, level GLint, xoffset GLint, yoffset GLint, x GLint, y GLint, width GLsizei, height GLsizei) {
	panic("extern")
}

// rename:createBuffer
func (WebGLRenderingContextBase) CreateBuffer() WebGLBuffer { panic("extern") }

// rename:createFramebuffer
func (WebGLRenderingContextBase) CreateFramebuffer() WebGLFramebuffer { panic("extern") }

// rename:createProgram
func (WebGLRenderingContextBase) CreateProgram() WebGLProgram { panic("extern") }

// rename:createRenderbuffer
func (WebGLRenderingContextBase) CreateRenderbuffer() WebGLRenderbuffer { panic("extern") }

// rename:createShader
func (WebGLRenderingContextBase) CreateShader(type0 GLenum) WebGLShader { panic("extern") }

// rename:createTexture
func (WebGLRenderingContextBase) CreateTexture() WebGLTexture { panic("extern") }

// rename:cullFace
func (WebGLRenderingContextBase) CullFace(mode GLenum) { panic("extern") }

// rename:deleteBuffer
func (WebGLRenderingContextBase) DeleteBuffer(buffer WebGLBuffer) { panic("extern") }

// rename:deleteFramebuffer
func (WebGLRenderingContextBase) DeleteFramebuffer(framebuffer WebGLFramebuffer) { panic("extern") }

// rename:deleteProgram
func (WebGLRenderingContextBase) DeleteProgram(program WebGLProgram) { panic("extern") }

// rename:deleteRenderbuffer
func (WebGLRenderingContextBase) DeleteRenderbuffer(renderbuffer WebGLRenderbuffer) { panic("extern") }

// rename:deleteShader
func (WebGLRenderingContextBase) DeleteShader(shader WebGLShader) { panic("extern") }

// rename:deleteTexture
func (WebGLRenderingContextBase) DeleteTexture(texture WebGLTexture) { panic("extern") }

// rename:depthFunc
func (WebGLRenderingContextBase) DepthFunc(func0 GLenum) { panic("extern") }

// rename:depthMask
func (WebGLRenderingContextBase) DepthMask(flag GLboolean) { panic("extern") }

// rename:depthRange
func (WebGLRenderingContextBase) DepthRange(zNear GLclampf, zFar GLclampf) { panic("extern") }

// rename:detachShader
func (WebGLRenderingContextBase) DetachShader(program WebGLProgram, shader WebGLShader) {
	panic("extern")
}

// rename:disable
func (WebGLRenderingContextBase) Disable(cap GLenum) { panic("extern") }

// rename:disableVertexAttribArray
func (WebGLRenderingContextBase) DisableVertexAttribArray(index GLuint) { panic("extern") }

// rename:drawArrays
func (WebGLRenderingContextBase) DrawArrays(mode GLenum, first GLint, count GLsizei) { panic("extern") }

// rename:drawElements
func (WebGLRenderingContextBase) DrawElements(mode GLenum, count GLsizei, type0 GLenum, offset GLintptr) {
	panic("extern")
}

// rename:enable
func (WebGLRenderingContextBase) Enable(cap GLenum) { panic("extern") }

// rename:enableVertexAttribArray
func (WebGLRenderingContextBase) EnableVertexAttribArray(index GLuint) { panic("extern") }

// rename:finish
func (WebGLRenderingContextBase) Finish() { panic("extern") }

// rename:flush
func (WebGLRenderingContextBase) Flush() { panic("extern") }

// rename:framebufferRenderbuffer
func (WebGLRenderingContextBase) FramebufferRenderbuffer(target GLenum, attachment GLenum, renderbuffertarget GLenum, renderbuffer WebGLRenderbuffer) {
	panic("extern")
}

// rename:framebufferTexture2D
func (WebGLRenderingContextBase) FramebufferTexture2D(target GLenum, attachment GLenum, textarget GLenum, texture WebGLTexture, level GLint) {
	panic("extern")
}

// rename:frontFace
func (WebGLRenderingContextBase) FrontFace(mode GLenum) { panic("extern") }

// rename:generateMipmap
func (WebGLRenderingContextBase) GenerateMipmap(target GLenum) { panic("extern") }

// rename:getActiveAttrib
func (WebGLRenderingContextBase) GetActiveAttrib(program WebGLProgram, index GLuint) WebGLActiveInfo {
	panic("extern")
}

// rename:getActiveUniform
func (WebGLRenderingContextBase) GetActiveUniform(program WebGLProgram, index GLuint) WebGLActiveInfo {
	panic("extern")
}

// rename:getAttachedShaders
func (WebGLRenderingContextBase) GetAttachedShaders(program WebGLProgram) []WebGLShader {
	panic("extern")
}

// rename:getAttribLocation
func (WebGLRenderingContextBase) GetAttribLocation(program WebGLProgram, name string) GLint {
	panic("extern")
}

// rename:getBufferParameter
func (WebGLRenderingContextBase) GetBufferParameter(target GLenum, pname GLenum) any { panic("extern") }

// rename:getContextAttributes
func (WebGLRenderingContextBase) GetContextAttributes() WebGLContextAttributes { panic("extern") }

// rename:getError
func (WebGLRenderingContextBase) GetError() GLenum { panic("extern") }

// rename:getExtension
func (WebGLRenderingContextBase) GetExtension(extensionName ExtensionName) EXT_blend_minmax {
	panic("extern")
}

// rename:getExtension
func (WebGLRenderingContextBase) GetExtension0(name string) any { panic("extern") }

// rename:getFramebufferAttachmentParameter
func (WebGLRenderingContextBase) GetFramebufferAttachmentParameter(target GLenum, attachment GLenum, pname GLenum) any {
	panic("extern")
}

// rename:getParameter
func (WebGLRenderingContextBase) GetParameter(pname GLenum) any { panic("extern") }

// rename:getProgramInfoLog
func (WebGLRenderingContextBase) GetProgramInfoLog(program WebGLProgram) string { panic("extern") }

// rename:getProgramParameter
func (WebGLRenderingContextBase) GetProgramParameter(program WebGLProgram, pname GLenum) any {
	panic("extern")
}

// rename:getRenderbufferParameter
func (WebGLRenderingContextBase) GetRenderbufferParameter(target GLenum, pname GLenum) any {
	panic("extern")
}

// rename:getShaderInfoLog
func (WebGLRenderingContextBase) GetShaderInfoLog(shader WebGLShader) string { panic("extern") }

// rename:getShaderParameter
func (WebGLRenderingContextBase) GetShaderParameter(shader WebGLShader, pname GLenum) any {
	panic("extern")
}

// rename:getShaderPrecisionFormat
func (WebGLRenderingContextBase) GetShaderPrecisionFormat(shadertype GLenum, precisiontype GLenum) WebGLShaderPrecisionFormat {
	panic("extern")
}

// rename:getShaderSource
func (WebGLRenderingContextBase) GetShaderSource(shader WebGLShader) string { panic("extern") }

// rename:getSupportedExtensions
func (WebGLRenderingContextBase) GetSupportedExtensions() []string { panic("extern") }

// rename:getTexParameter
func (WebGLRenderingContextBase) GetTexParameter(target GLenum, pname GLenum) any { panic("extern") }

// rename:getUniform
func (WebGLRenderingContextBase) GetUniform(program WebGLProgram, location WebGLUniformLocation) any {
	panic("extern")
}

// rename:getUniformLocation
func (WebGLRenderingContextBase) GetUniformLocation(program WebGLProgram, name string) WebGLUniformLocation {
	panic("extern")
}

// rename:getVertexAttrib
func (WebGLRenderingContextBase) GetVertexAttrib(index GLuint, pname GLenum) any { panic("extern") }

// rename:getVertexAttribOffset
func (WebGLRenderingContextBase) GetVertexAttribOffset(index GLuint, pname GLenum) GLintptr {
	panic("extern")
}

// rename:hint
func (WebGLRenderingContextBase) Hint(target GLenum, mode GLenum) { panic("extern") }

// rename:isBuffer
func (WebGLRenderingContextBase) IsBuffer(buffer WebGLBuffer) GLboolean { panic("extern") }

// rename:isContextLost
func (WebGLRenderingContextBase) IsContextLost() bool { panic("extern") }

// rename:isEnabled
func (WebGLRenderingContextBase) IsEnabled(cap GLenum) GLboolean { panic("extern") }

// rename:isFramebuffer
func (WebGLRenderingContextBase) IsFramebuffer(framebuffer WebGLFramebuffer) GLboolean {
	panic("extern")
}

// rename:isProgram
func (WebGLRenderingContextBase) IsProgram(program WebGLProgram) GLboolean { panic("extern") }

// rename:isRenderbuffer
func (WebGLRenderingContextBase) IsRenderbuffer(renderbuffer WebGLRenderbuffer) GLboolean {
	panic("extern")
}

// rename:isShader
func (WebGLRenderingContextBase) IsShader(shader WebGLShader) GLboolean { panic("extern") }

// rename:isTexture
func (WebGLRenderingContextBase) IsTexture(texture WebGLTexture) GLboolean { panic("extern") }

// rename:lineWidth
func (WebGLRenderingContextBase) LineWidth(width GLfloat) { panic("extern") }

// rename:linkProgram
func (WebGLRenderingContextBase) LinkProgram(program WebGLProgram) { panic("extern") }

// rename:pixelStorei
func (WebGLRenderingContextBase) PixelStorei(pname GLenum, param any) { panic("extern") }

// rename:polygonOffset
func (WebGLRenderingContextBase) PolygonOffset(factor GLfloat, units GLfloat) { panic("extern") }

// rename:renderbufferStorage
func (WebGLRenderingContextBase) RenderbufferStorage(target GLenum, internalformat GLenum, width GLsizei, height GLsizei) {
	panic("extern")
}

// rename:sampleCoverage
func (WebGLRenderingContextBase) SampleCoverage(value GLclampf, invert GLboolean) { panic("extern") }

// rename:scissor
func (WebGLRenderingContextBase) Scissor(x GLint, y GLint, width GLsizei, height GLsizei) {
	panic("extern")
}

// rename:shaderSource
func (WebGLRenderingContextBase) ShaderSource(shader WebGLShader, source string) { panic("extern") }

// rename:stencilFunc
func (WebGLRenderingContextBase) StencilFunc(func0 GLenum, ref GLint, mask GLuint) { panic("extern") }

// rename:stencilFuncSeparate
func (WebGLRenderingContextBase) StencilFuncSeparate(face GLenum, func0 GLenum, ref GLint, mask GLuint) {
	panic("extern")
}

// rename:stencilMask
func (WebGLRenderingContextBase) StencilMask(mask GLuint) { panic("extern") }

// rename:stencilMaskSeparate
func (WebGLRenderingContextBase) StencilMaskSeparate(face GLenum, mask GLuint) { panic("extern") }

// rename:stencilOp
func (WebGLRenderingContextBase) StencilOp(fail GLenum, zfail GLenum, zpass GLenum) { panic("extern") }

// rename:stencilOpSeparate
func (WebGLRenderingContextBase) StencilOpSeparate(face GLenum, fail GLenum, zfail GLenum, zpass GLenum) {
	panic("extern")
}

// rename:texParameterf
func (WebGLRenderingContextBase) TexParameterf(target GLenum, pname GLenum, param GLfloat) {
	panic("extern")
}

// rename:texParameteri
func (WebGLRenderingContextBase) TexParameteri(target GLenum, pname GLenum, param GLint) {
	panic("extern")
}

// rename:uniform1f
func (WebGLRenderingContextBase) Uniform1f(location WebGLUniformLocation, x GLfloat) { panic("extern") }

// rename:uniform1i
func (WebGLRenderingContextBase) Uniform1i(location WebGLUniformLocation, x GLint) { panic("extern") }

// rename:uniform2f
func (WebGLRenderingContextBase) Uniform2f(location WebGLUniformLocation, x GLfloat, y GLfloat) {
	panic("extern")
}

// rename:uniform2i
func (WebGLRenderingContextBase) Uniform2i(location WebGLUniformLocation, x GLint, y GLint) {
	panic("extern")
}

// rename:uniform3f
func (WebGLRenderingContextBase) Uniform3f(location WebGLUniformLocation, x GLfloat, y GLfloat, z GLfloat) {
	panic("extern")
}

// rename:uniform3i
func (WebGLRenderingContextBase) Uniform3i(location WebGLUniformLocation, x GLint, y GLint, z GLint) {
	panic("extern")
}

// rename:uniform4f
func (WebGLRenderingContextBase) Uniform4f(location WebGLUniformLocation, x GLfloat, y GLfloat, z GLfloat, w GLfloat) {
	panic("extern")
}

// rename:uniform4i
func (WebGLRenderingContextBase) Uniform4i(location WebGLUniformLocation, x GLint, y GLint, z GLint, w GLint) {
	panic("extern")
}

// rename:useProgram
func (WebGLRenderingContextBase) UseProgram(program WebGLProgram) { panic("extern") }

// rename:validateProgram
func (WebGLRenderingContextBase) ValidateProgram(program WebGLProgram) { panic("extern") }

// rename:vertexAttrib1f
func (WebGLRenderingContextBase) VertexAttrib1f(index GLuint, x GLfloat) { panic("extern") }

// rename:vertexAttrib1fv
func (WebGLRenderingContextBase) VertexAttrib1fv(index GLuint, values []float32) { panic("extern") }

// rename:vertexAttrib2f
func (WebGLRenderingContextBase) VertexAttrib2f(index GLuint, x GLfloat, y GLfloat) { panic("extern") }

// rename:vertexAttrib2fv
func (WebGLRenderingContextBase) VertexAttrib2fv(index GLuint, values []float32) { panic("extern") }

// rename:vertexAttrib3f
func (WebGLRenderingContextBase) VertexAttrib3f(index GLuint, x GLfloat, y GLfloat, z GLfloat) {
	panic("extern")
}

// rename:vertexAttrib3fv
func (WebGLRenderingContextBase) VertexAttrib3fv(index GLuint, values []float32) { panic("extern") }

// rename:vertexAttrib4f
func (WebGLRenderingContextBase) VertexAttrib4f(index GLuint, x GLfloat, y GLfloat, z GLfloat, w GLfloat) {
	panic("extern")
}

// rename:vertexAttrib4fv
func (WebGLRenderingContextBase) VertexAttrib4fv(index GLuint, values []float32) { panic("extern") }

// rename:vertexAttribPointer
func (WebGLRenderingContextBase) VertexAttribPointer(index GLuint, size GLint, type0 GLenum, normalized GLboolean, stride GLsizei, offset GLintptr) {
	panic("extern")
}

// rename:viewport
func (WebGLRenderingContextBase) Viewport(x GLint, y GLint, width GLsizei, height GLsizei) {
	panic("extern")
}

type WebGLRenderingContextOverloads struct {
}

// rename:bufferData
func (WebGLRenderingContextOverloads) BufferData(target GLenum, size GLsizeiptr, usage GLenum) {
	panic("extern")
}

// rename:bufferData
func (WebGLRenderingContextOverloads) BufferData0(target GLenum, data BufferSource, usage GLenum) {
	panic("extern")
}

// rename:bufferSubData
func (WebGLRenderingContextOverloads) BufferSubData(target GLenum, offset GLintptr, data BufferSource) {
	panic("extern")
}

// rename:compressedTexImage2D
func (WebGLRenderingContextOverloads) CompressedTexImage2D(target GLenum, level GLint, internalformat GLenum, width GLsizei, height GLsizei, border GLint, data ArrayBufferView) {
	panic("extern")
}

// rename:compressedTexSubImage2D
func (WebGLRenderingContextOverloads) CompressedTexSubImage2D(target GLenum, level GLint, xoffset GLint, yoffset GLint, width GLsizei, height GLsizei, format GLenum, data ArrayBufferView) {
	panic("extern")
}

// rename:readPixels
func (WebGLRenderingContextOverloads) ReadPixels(x GLint, y GLint, width GLsizei, height GLsizei, format GLenum, type0 GLenum, pixels ArrayBufferView) {
	panic("extern")
}

// rename:texImage2D
func (WebGLRenderingContextOverloads) TexImage2D(target GLenum, level GLint, internalformat GLint, width GLsizei, height GLsizei, border GLint, format GLenum, type0 GLenum, pixels ArrayBufferView) {
	panic("extern")
}

// rename:texImage2D
func (WebGLRenderingContextOverloads) TexImage2D0(target GLenum, level GLint, internalformat GLint, format GLenum, type0 GLenum, source TexImageSource) {
	panic("extern")
}

// rename:texSubImage2D
func (WebGLRenderingContextOverloads) TexSubImage2D(target GLenum, level GLint, xoffset GLint, yoffset GLint, width GLsizei, height GLsizei, format GLenum, type0 GLenum, pixels ArrayBufferView) {
	panic("extern")
}

// rename:texSubImage2D
func (WebGLRenderingContextOverloads) TexSubImage2D0(target GLenum, level GLint, xoffset GLint, yoffset GLint, format GLenum, type0 GLenum, source TexImageSource) {
	panic("extern")
}

// rename:uniform1fv
func (WebGLRenderingContextOverloads) Uniform1fv(location WebGLUniformLocation, v []float32) {
	panic("extern")
}

// rename:uniform1iv
func (WebGLRenderingContextOverloads) Uniform1iv(location WebGLUniformLocation, v []int32) {
	panic("extern")
}

// rename:uniform2fv
func (WebGLRenderingContextOverloads) Uniform2fv(location WebGLUniformLocation, v []float32) {
	panic("extern")
}

// rename:uniform2iv
func (WebGLRenderingContextOverloads) Uniform2iv(location WebGLUniformLocation, v []int32) {
	panic("extern")
}

// rename:uniform3fv
func (WebGLRenderingContextOverloads) Uniform3fv(location WebGLUniformLocation, v []float32) {
	panic("extern")
}

// rename:uniform3iv
func (WebGLRenderingContextOverloads) Uniform3iv(location WebGLUniformLocation, v []int32) {
	panic("extern")
}

// rename:uniform4fv
func (WebGLRenderingContextOverloads) Uniform4fv(location WebGLUniformLocation, v []float32) {
	panic("extern")
}

// rename:uniform4iv
func (WebGLRenderingContextOverloads) Uniform4iv(location WebGLUniformLocation, v []int32) {
	panic("extern")
}

// rename:uniformMatrix2fv
func (WebGLRenderingContextOverloads) UniformMatrix2fv(location WebGLUniformLocation, transpose GLboolean, value []float32) {
	panic("extern")
}

// rename:uniformMatrix3fv
func (WebGLRenderingContextOverloads) UniformMatrix3fv(location WebGLUniformLocation, transpose GLboolean, value []float32) {
	panic("extern")
}

// rename:uniformMatrix4fv
func (WebGLRenderingContextOverloads) UniformMatrix4fv(location WebGLUniformLocation, transpose GLboolean, value []float32) {
	panic("extern")
}

// rename:document
type DocumentType struct {
	Node
	DocumentAndElementEventHandlers
	DocumentOrShadowRoot
	FontFaceSource
	GlobalEventHandlers
	NonElementParentNode
	ParentNode
	XPathEvaluatorBase
	URL string
	// rename:alinkColor
	AlinkColor string
	// rename:all
	All HTMLAllCollection
	// rename:anchors
	Anchors any
	// rename:applets
	Applets HTMLCollection
	// rename:bgColor
	BgColor string
	// rename:body
	Body HTMLElement
	// rename:characterSet
	CharacterSet string
	// rename:charset
	Charset string
	// rename:compatMode
	CompatMode string
	// rename:contentType
	ContentType string
	// rename:cookie
	Cookie string
	// rename:currentScript
	CurrentScript HTMLOrSVGScriptElement
	// rename:defaultView
	DefaultView any
	// rename:designMode
	DesignMode string
	// rename:dir
	Dir string
	// rename:doctype
	Doctype struct{}
	// rename:documentElement
	DocumentElement HTMLElement
	// rename:documentURI
	DocumentURI string
	// rename:domain
	Domain string
	// rename:embeds
	Embeds any
	// rename:fgColor
	FgColor string
	// rename:forms
	Forms any
	// rename:fullscreen
	Fullscreen bool
	// rename:fullscreenEnabled
	FullscreenEnabled bool
	// rename:head
	Head HTMLHeadElement
	// rename:hidden
	Hidden bool
	// rename:images
	Images any
	// rename:implementation
	Implementation DOMImplementation
	// rename:inputEncoding
	InputEncoding string
	// rename:lastModified
	LastModified string
	// rename:linkColor
	LinkColor string
	// rename:links
	Links any
	// rename:onfullscreenchange
	Onfullscreenchange any
	// rename:onfullscreenerror
	Onfullscreenerror any
	// rename:onpointerlockchange
	Onpointerlockchange any
	// rename:onpointerlockerror
	Onpointerlockerror any
	// rename:onreadystatechange
	Onreadystatechange any
	// rename:onvisibilitychange
	Onvisibilitychange any
	// rename:ownerDocument
	OwnerDocument Null
	// rename:pictureInPictureEnabled
	PictureInPictureEnabled bool
	// rename:plugins
	Plugins any
	// rename:readyState
	ReadyState DocumentReadyState
	// rename:referrer
	Referrer string
	// rename:rootElement
	RootElement SVGSVGElement
	// rename:scripts
	Scripts any
	// rename:scrollingElement
	ScrollingElement Element
	// rename:timeline
	Timeline DocumentTimeline
	// rename:title
	Title string
	// rename:visibilityState
	VisibilityState DocumentVisibilityState
	// rename:vlinkColor
	VlinkColor string
}

// rename:adoptNode
func (DocumentType) AdoptNode(node Node) Node { panic("extern") }

// rename:captureEvents
func (DocumentType) CaptureEvents() { panic("extern") }

// rename:caretRangeFromPoint
func (DocumentType) CaretRangeFromPoint(x float64, y float64) Range { panic("extern") }

// rename:clear
func (DocumentType) Clear() { panic("extern") }

// rename:close
func (DocumentType) Close() { panic("extern") }

// rename:createAttribute
func (DocumentType) CreateAttribute(localName string) Attr { panic("extern") }

// rename:createAttributeNS
func (DocumentType) CreateAttributeNS(namespace string, qualifiedName string) Attr { panic("extern") }

// rename:createCDATASection
func (DocumentType) CreateCDATASection(data string) CDATASection { panic("extern") }

// rename:createComment
func (DocumentType) CreateComment(data string) Comment { panic("extern") }

// rename:createDocumentFragment
func (DocumentType) CreateDocumentFragment() DocumentFragment { panic("extern") }

// rename:createElement
func (DocumentType) CreateElement(tagName any, options ElementCreationOptions) any { panic("extern") }

// rename:createElement
func (DocumentType) CreateElement0(tagName string, options ElementCreationOptions) HTMLElement {
	panic("extern")
}

// rename:createElementNS
func (DocumentType) CreateElementNS(namespaceURI string, qualifiedName string) HTMLElement {
	panic("extern")
}

// rename:createElementNS
func (DocumentType) CreateElementNS0(namespaceURI string, qualifiedName any) any { panic("extern") }

// rename:createElementNS
func (DocumentType) CreateElementNS1(namespaceURI string, qualifiedName string, options ElementCreationOptions) Element {
	panic("extern")
}

// rename:createElementNS
func (DocumentType) CreateElementNS2(namespace string, qualifiedName string, options any) Element {
	panic("extern")
}

// rename:createEvent
func (DocumentType) CreateEvent(eventInterface EventInterface) AnimationEvent { panic("extern") }

// rename:createEvent
func (DocumentType) CreateEvent0(eventInterface string) Event { panic("extern") }

// rename:createNodeIterator
func (DocumentType) CreateNodeIterator(root Node, whatToShow float64, filter NodeFilter) NodeIterator {
	panic("extern")
}

// rename:createProcessingInstruction
func (DocumentType) CreateProcessingInstruction(target string, data string) ProcessingInstruction {
	panic("extern")
}

// rename:createRange
func (DocumentType) CreateRange() Range { panic("extern") }

// rename:createTextNode
func (DocumentType) CreateTextNode(data string) Text { panic("extern") }

// rename:createTreeWalker
func (DocumentType) CreateTreeWalker(root Node, whatToShow float64, filter NodeFilter) TreeWalker {
	panic("extern")
}

// rename:execCommand
func (DocumentType) ExecCommand(commandId string, showUI bool, value string) bool { panic("extern") }

// rename:exitFullscreen
func (DocumentType) ExitFullscreen() any { panic("extern") }

// rename:exitPictureInPicture
func (DocumentType) ExitPictureInPicture() any { panic("extern") }

// rename:exitPointerLock
func (DocumentType) ExitPointerLock() { panic("extern") }

// rename:getElementById
func (DocumentType) GetElementById(elementId string) HTMLElement { panic("extern") }

// rename:getElementsByClassName
func (DocumentType) GetElementsByClassName(classNames string) any { panic("extern") }

// rename:getElementsByName
func (DocumentType) GetElementsByName(elementName string) any { panic("extern") }

// rename:getElementsByTagName
func (DocumentType) GetElementsByTagName(qualifiedName any) any { panic("extern") }

// rename:getElementsByTagName
func (DocumentType) GetElementsByTagName0(qualifiedName string) any { panic("extern") }

// rename:getElementsByTagNameNS
func (DocumentType) GetElementsByTagNameNS(namespaceURI string, localName string) any {
	panic("extern")
}

// rename:getSelection
func (DocumentType) GetSelection() Selection { panic("extern") }

// rename:hasFocus
func (DocumentType) HasFocus() bool { panic("extern") }

// rename:hasStorageAccess
func (DocumentType) HasStorageAccess() any { panic("extern") }

// rename:importNode
func (DocumentType) ImportNode(node Node, deep bool) Node { panic("extern") }

// rename:open
func (DocumentType) Open(unused1 string, unused2 string) DocumentType { panic("extern") }

// rename:open
func (DocumentType) Open0(url any, name string, features string) WindowProxy { panic("extern") }

// rename:queryCommandEnabled
func (DocumentType) QueryCommandEnabled(commandId string) bool { panic("extern") }

// rename:queryCommandIndeterm
func (DocumentType) QueryCommandIndeterm(commandId string) bool { panic("extern") }

// rename:queryCommandState
func (DocumentType) QueryCommandState(commandId string) bool { panic("extern") }

// rename:queryCommandSupported
func (DocumentType) QueryCommandSupported(commandId string) bool { panic("extern") }

// rename:queryCommandValue
func (DocumentType) QueryCommandValue(commandId string) string { panic("extern") }

// rename:releaseEvents
func (DocumentType) ReleaseEvents() { panic("extern") }

// rename:requestStorageAccess
func (DocumentType) RequestStorageAccess() any { panic("extern") }

// rename:write
func (DocumentType) Write(text []string) { panic("extern") }

// rename:writeln
func (DocumentType) Writeln(text []string) { panic("extern") }

// rename:addEventListener
func (DocumentType) AddEventListener(type0 any, listener any, options any) { panic("extern") }

// rename:addEventListener
func (DocumentType) AddEventListener0(type0 string, listener EventListenerOrEventListenerObject, options any) {
	panic("extern")
}

// rename:removeEventListener
func (DocumentType) RemoveEventListener(type0 any, listener any, options any) { panic("extern") }

// rename:removeEventListener
func (DocumentType) RemoveEventListener0(type0 string, listener EventListenerOrEventListenerObject, options any) {
	panic("extern")
}

type GLboolean = bool
type GLenum = float64
type GLfloat = float64
type GLint = float64
type GLsizei = float64
type GLuint = float64
type MediaStream struct{}
type CanvasRenderingContext2D struct{}
type WebGLProgram struct{}
type WebGLBuffer struct{}
type GLbitfield = float64
type FontFaceSource struct{}
type Comment struct{}
type NodeFilter struct{}
type CDATASection struct{}
type Node struct{}
type WebGLShader struct{}
type WebGLRenderbuffer struct{}
type ArrayBufferView struct{}
type HTMLAllCollection struct{}
type ElementCreationOptions struct{}
type AnimationEvent struct{}
type WebGLActiveInfo struct{}
type TexImageSource struct{}
type XPathEvaluatorBase struct{}
type WebGLTexture struct{}
type BufferSource struct{}
type CanvasRenderingContext2DSettings struct{}
type EventListenerOrEventListenerObject struct{}
type WebGLUniformLocation struct{}
type GlobalEventHandlers struct{}
type Range struct{}
type RenderingContext struct{}
type ImageBitmapRenderingContextSettings struct{}
type GLclampf = float64
type GLintptr = float64
type EXT_blend_minmax struct{}
type DocumentVisibilityState struct{}
type DocumentFragment struct{}
type HTMLElement struct{}
type NonElementParentNode struct{}
type HTMLCollection struct{}
type SVGSVGElement struct{}
type Event struct{}
type TreeWalker struct{}
type WebGLShaderPrecisionFormat struct{}
type WebGLContextAttributes struct{}
type HTMLHeadElement struct{}
type DocumentTimeline struct{}
type HTMLOrSVGScriptElement struct{}
type Attr struct{}
type GLsizeiptr struct{}
type DocumentAndElementEventHandlers struct{}
type ImageBitmapRenderingContext struct{}
type ParentNode struct{}
type DOMImplementation struct{}
type DocumentReadyState struct{}
type Element struct{}
type BlobCallback struct{}
type NodeIterator struct{}
type ProcessingInstruction struct{}
type Text struct{}
type Selection struct{}
type DocumentOrShadowRoot struct{}
type Null struct{}
type WindowProxy struct{}
type WebGLFramebuffer struct{}

// rename:document
var Document DocumentType = DocumentType{}
