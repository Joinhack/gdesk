#include <dispatch/dispatch.h>
#include <CoreFoundation/CoreFoundation.h>
#include <CoreGraphics/CoreGraphics.h>
#include <IOSurface/IOSurface.h>


CGDisplayStreamRef DisplayStreamCreateWithDispatchQueue(void *cap, CGDirectDisplayID display, uint outputWidth, uint outputHeight, CFDictionaryRef properties, dispatch_queue_t queue);

void dict_release(CFDictionaryRef p);

CFDictionaryRef dict_create(double throttle, int8_t queue_length, int cursor, int letterbox);

void q_release(dispatch_queue_t q);

void CaptureStop(void *p);

void CaptureComplete(void *p, IOSurfaceRef);