#include "capture_darwin.h"

void q_release(dispatch_queue_t q) {
	dispatch_release(q);
}

CFDictionaryRef dict_create(double throttle, int8_t queue_length, int cursor, int letterbox) {
	CFNumberRef throttle_ref = CFNumberCreate(NULL, kCFNumberFloat64Type, &throttle);
	CFNumberRef q_len_ref = CFNumberCreate(NULL, kCFNumberSInt8Type, &queue_length);
	CFStringRef keys[4] = {
		kCGDisplayStreamShowCursor,
		kCGDisplayStreamPreserveAspectRatio,
		kCGDisplayStreamMinimumFrameTime,
		kCGDisplayStreamQueueDepth,
	};
	void* vals[4] = {
		(void*)(cursor == 1?kCFBooleanTrue:kCFBooleanFalse),
		(void*)(letterbox == 1?kCFBooleanTrue:kCFBooleanFalse),
		(void*)throttle_ref,
		(void*)q_len_ref,
	};
	CFRelease(throttle_ref);
	CFRelease(q_len_ref);
	return (void*)CFDictionaryCreate(NULL, 
			(const void**)&keys, 
			(const void**)&vals, 
			4, 
			&kCFTypeDictionaryKeyCallBacks, 
			&kCFTypeDictionaryValueCallBacks);
}

void dict_release(CFDictionaryRef p) {
	CFRelease(p);
}

CGDisplayStreamRef DisplayStreamCreateWithDispatchQueue(void *cap, CGDirectDisplayID display, uint outputWidth, uint outputHeight, CFDictionaryRef properties, dispatch_queue_t queue) {
	return CGDisplayStreamCreateWithDispatchQueue(display, outputWidth, outputHeight, '420v', properties, queue, ^(CGDisplayStreamFrameStatus status, uint64_t displayTime, IOSurfaceRef frameSurface, CGDisplayStreamUpdateRef updateRef) {
		if (kCGDisplayStreamFrameStatusStopped == status) {
			CaptureStop(cap);
		} else if (kCGDisplayStreamFrameStatusFrameComplete == status) {
			CaptureComplete(cap, frameSurface);
		}
	});
}
