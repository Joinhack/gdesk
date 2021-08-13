package gdesk

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework IOSurface
#include "capture_darwin.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

type Display struct {
	inner C.CGDirectDisplayID
}

func GetPrimaryDisplay() *Display {
	return &Display{inner: C.CGMainDisplayID()}
}

func (d *Display) String() string {
	return fmt.Sprintf("%d", d.inner)
}

func (d *Display) Width() uint {
	return uint(C.CGDisplayPixelsWide(d.inner))
}

func (d *Display) Height() uint {
	return uint(C.CGDisplayPixelsHigh(d.inner))
}

type Frame struct {
	inner  C.IOSurfaceRef
	data   []byte
	inited int32
}

func NewFrame(surface C.IOSurfaceRef) *Frame {
	frame := &Frame{
		inner: surface,
	}
	C.CFRetain(C.CFTypeRef(surface))
	C.IOSurfaceIncrementUseCount(surface)
	C.IOSurfaceLock(surface, C.kIOSurfaceLockReadOnly, nil)
	runtime.SetFinalizer(frame, (*Frame).Release)
	frame.inited = 1
	return frame
}

func (f *Frame) convert(w, h int) {
	surface := f.inner
	plane0_addr := C.IOSurfaceGetBaseAddressOfPlane(surface, 0)
	plane0_l := int(C.IOSurfaceGetBytesPerRowOfPlane(surface, 0))

	plane1_addr := C.IOSurfaceGetBaseAddressOfPlane(surface, 1)
	plane1_l := int(C.IOSurfaceGetBytesPerRowOfPlane(surface, 1))

	plane0_header := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(plane0_addr)),
		Cap:  plane0_l,
		Len:  plane0_l,
	}

	plane1_header := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(plane1_addr)),
		Cap:  plane1_l,
		Len:  plane1_l,
	}

	plane0 := *(*[]byte)(unsafe.Pointer(&plane0_header))
	plane1 := *(*[]byte)(unsafe.Pointer(&plane1_header))
	
	f.data = nv12_to_i420(plane0, plane1, w, h)
}

func (f *Frame) Release() {
	if atomic.LoadInt32(&f.inited) == 0 {
		return
	}
	for {
		if atomic.CompareAndSwapInt32(&f.inited, f.inited, 0) {
			break
		}
	}
	C.IOSurfaceUnlock(f.inner, C.kIOSurfaceLockReadOnly, nil)
	C.IOSurfaceDecrementUseCount(f.inner)
	C.CFRelease(C.CFTypeRef(f.inner))
}

type Capturer struct {
	ref     unsafe.Pointer
	queue   C.dispatch_queue_t
	dict    C.CFDictionaryRef
	stream  C.CGDisplayStreamRef
	stopped int32
	frame   unsafe.Pointer
	dW, dH 	int
}

//export CaptureStop
func CaptureStop(c unsafe.Pointer) {
	var cap = (*Capturer)(c)
	for {
		if atomic.CompareAndSwapInt32(&cap.stopped, cap.stopped, 1) {
			break
		}
	}
}

//export CaptureComplete
func CaptureComplete(c unsafe.Pointer, surf C.IOSurfaceRef) {
	var cap = (*Capturer)(c)
	frame := unsafe.Pointer(NewFrame(surf))
	for {
		if atomic.CompareAndSwapPointer(&cap.frame, cap.frame, frame) {
			break
		}
	}
}

func NewCapturer() *Capturer {
	return &Capturer{
		stopped: 1,
	}
}

func (cap *Capturer) GetFrame() *Frame {
	f := (*Frame)(atomic.LoadPointer(&cap.frame))
	if f != nil {
		f.convert(cap.dW, cap.dH)
	}
	return f
}

func (cap *Capturer) Start(display *Display) error {
	cap.queue = C.dispatch_queue_create(C.CString("capture queue"), C.dispatch_queue_attr_t(nil))
	cap.dict = C.dict_create(0.0, 8, 0, 1)
	w := C.uint(display.Width())
	h := C.uint(display.Height())

	cap.stream = C.DisplayStreamCreateWithDispatchQueue(unsafe.Pointer(cap), display.inner, w, h, cap.dict, cap.queue)
	rs := C.CGDisplayStreamStart(cap.stream)
	if rs != C.kCGErrorSuccess {
		return errors.New(fmt.Sprintf("error start stream, code %d", int(rs)))
	}
	for {
		if atomic.CompareAndSwapInt32(&cap.stopped, cap.stopped, 0) {
			break
		}
	}
	cap.dW = int(w)
	cap.dH = int(h)
	return nil
}

func (cap *Capturer) Stop() {
	C.CGDisplayStreamStop(cap.stream)
	for {
		if atomic.LoadInt32(&cap.stopped) == 1 {
			break
		}
	}
	time.Sleep(30 * time.Millisecond)
	C.q_release(cap.queue)
	C.dict_release(cap.dict)
	return
}
