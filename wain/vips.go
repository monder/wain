package wain

/*
#cgo pkg-config: vips
#include "vips.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

func VipsInit() {
	C.vips_wain_init()
}

func VipsResize(data []byte, r ResizeOptions) ([]byte, error) {

	length := C.size_t(0)
	var ptr unsafe.Pointer

	err := C.vips_wain_resize(
		unsafe.Pointer(&data[0]), C.size_t(len(data)),
		&ptr, &length,
		C.int(r.Width), C.int(r.Height),
	)
	if err != 0 {
		s := C.GoString(C.vips_error_buffer())
		C.vips_error_clear()
		return nil, errors.New(s)
	}

	buf := C.GoBytes(ptr, C.int(length))
	C.g_free(C.gpointer(ptr))

	return buf, nil
}
