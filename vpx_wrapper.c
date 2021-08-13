#include <stdio.h>
#include "vpx_wrapper.h"

int I420ToRAW(
    char* src_y,
    int src_stride_y,
    char *src_u,
    int src_stride_u,
    char *src_v,
    int src_stride_v,
    char* dst_rgba,
    int dst_stride_rgba,
    int width,
    int height
);

int i420_to_rgb(int width, int height, void *src, void *to) {
    vpx_image_t img = {};
    vpx_img_wrap(&img, VPX_IMG_FMT_I420, width, height, 16, NULL);
    int u = img.planes[1] - img.planes[0];
    int v = img.planes[2] - img.planes[0];
    int src_stride_uv = img.stride[1];
    int src_stride_y = img.stride[0];
    char *src_y = ((char*)src);
    char *src_u = ((char*)src) + u;
    char *src_v = ((char*)src) + v;
    I420ToRAW(src_y, src_stride_y, src_u, src_stride_uv, src_v, src_stride_uv, (char*)to, width*3, width, height);
    return 0;
}

