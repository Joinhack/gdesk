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

int NV12ToI420(
    char* src_y,
    int src_stride_y,
    char *src_uv,
    int src_stride_uv,
    char *dst_y,
    int dst_stride_y,
    char* dst_u,
    int dst_stride_u,
    char* dst_v,
    int dst_stride_v,
    int width,
    int height
);

int i420_to_rgb(int width, int height, void *src, void *to) {
    vpx_image_t *img = vpx_img_wrap(NULL, VPX_IMG_FMT_I420, width, height, 16, NULL);
    int u = img->planes[1] - img->planes[0];
    int v = img->planes[2] - img->planes[0];
    int src_stride_uv = img->stride[1];
    int src_stride_y = img->stride[0];
    char *src_y = ((char*)src);
    char *src_u = ((char*)src) + u;
    char *src_v = ((char*)src) + v;
    vpx_img_free(img);
    I420ToRAW(src_y, src_stride_y, src_u, src_stride_uv, src_v, src_stride_uv, (char*)to, width*3, width, height);
    return 0;
}

int nv12_to_i420(void *src_y, int src_stride_y, void *src_uv, int src_stride_uv, int width, int height, void *dst) {
    vpx_image_t *img = vpx_img_wrap(NULL, VPX_IMG_FMT_I420, width, height, 16, NULL);
    int u = img->planes[1] - img->planes[0];
    int v = img->planes[2] - img->planes[0];
    int dst_stride_uv = img->stride[1];
    int dst_stride_y = img->stride[0];
    char *dst_y = ((char*)dst);
    char *dst_u = ((char*)dst) + u;
    char *dst_v = ((char*)dst) + v;
    vpx_img_free(img);
    NV12ToI420((char*)src_y, src_stride_y, (char*)src_uv, src_stride_uv, dst_y, dst_stride_y, dst_u, dst_stride_uv, dst_v, dst_stride_uv , width, height);
    return 0;
}

