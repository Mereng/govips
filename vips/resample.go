package vips

// #include "resample.h"
import "C"
import "unsafe"

// Kernel represents VipsKernel type
type Kernel int

// Kernel enum
const (
	KernelAuto     Kernel = -1
	KernelNearest  Kernel = C.VIPS_KERNEL_NEAREST
	KernelLinear   Kernel = C.VIPS_KERNEL_LINEAR
	KernelCubic    Kernel = C.VIPS_KERNEL_CUBIC
	KernelLanczos2 Kernel = C.VIPS_KERNEL_LANCZOS2
	KernelLanczos3 Kernel = C.VIPS_KERNEL_LANCZOS3
	KernelMitchell Kernel = C.VIPS_KERNEL_MITCHELL
)

// Interpalate represents VipsInterpolate type
type Interpolate int

const (
	// Bicubic interpolation value.
	Bicubic Interpolate = iota
	// Bilinear interpolation value.
	Bilinear
	// Nohalo interpolation value.
	Nohalo
	// Nearest neighbour interpolation value.
	Nearest
)

var interpolations = map[Interpolate]string{
	Bicubic:  "bicubic",
	Bilinear: "bilinear",
	Nohalo:   "nohalo",
	Nearest:  "nearest",
}

func (i Interpolate) String() string {
	return interpolations[i]
}

// https://libvips.github.io/libvips/API/current/libvips-resample.html#vips-resize
func vipsResize(in *C.VipsImage, scale float64, kernel Kernel) (*C.VipsImage, error) {
	incOpCounter("resize")
	var out *C.VipsImage

	// libvips recommends Lanczos3 as the default kernel
	if kernel == KernelAuto {
		kernel = KernelLanczos3
	}

	if err := C.resize_image(in, &out, C.double(scale), C.double(-1), C.int(kernel)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-resample.html#vips-resize
func vipsResizeWithVScale(in *C.VipsImage, scale, vscale float64, kernel Kernel) (*C.VipsImage, error) {
	incOpCounter("resize")
	var out *C.VipsImage

	if err := C.resize_image(in, &out, C.double(scale), C.gdouble(vscale), C.int(kernel)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

func vipsThumbnail(in *C.VipsImage, width, height int, crop Interesting) (*C.VipsImage, error) {
	incOpCounter("thumbnail")
	var out *C.VipsImage

	if err := C.thumbnail_image(in, &out, C.int(width), C.int(height), C.int(crop)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-resample.html#vips-mapim
func vipsMapim(in *C.VipsImage, index *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("mapim")
	var out *C.VipsImage

	if err := C.mapim(in, &out, index); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-histogram.html#vips-maplut
func vipsMaplut(in *C.VipsImage, lut *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("maplut")
	var out *C.VipsImage

	if err := C.maplut(in, &out, lut); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

func vipsAffine(in *C.VipsImage, a, b, c, d float64, interpolator Interpolate) (*C.VipsImage, error) {
	incOpCounter("affine")
	var out *C.VipsImage

	cstring := C.CString(interpolator.String())
	i := C.vips_interpolate_new(cstring)

	defer C.free(unsafe.Pointer(cstring))

	if err := C.affine_image(in, &out, C.double(a), C.double(b), C.double(c), C.double(d), i); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}
