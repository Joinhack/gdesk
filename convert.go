package gdesk


/*
#cgo pkg-config: vpx
#cgo LDFLAGS: -lyuv
#include "vpx_wrapper.h"
*/
import "C"

import (
	"unsafe"
	"reflect"
)

func i420ToRgb(w,h int, src, dest []byte) {
	s := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&src)).Data)
	d := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&dest)).Data)
	C.i420_to_rgb(C.int(w), C.int(h), s, d)
}