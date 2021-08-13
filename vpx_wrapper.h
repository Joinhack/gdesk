#include "vpx/vpx_image.h"

int i420_to_rgb(int width, int height, void *src, void *to);

int nv12_to_i420(void *src_y, int src_stride_y, void *src_uv, int src_stride_uv, int width, int height, void *dst);

