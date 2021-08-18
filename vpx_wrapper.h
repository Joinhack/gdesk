#include <stdlib.h>
#include "vpx/vpx_image.h"
#include "vpx/vp8cx.h"
#include "vpx/vpx_encoder.h"


int i420_to_rgb(int width, int height, void *src, void *to);

int nv12_to_i420(void *src_y, int src_stride_y, void *src_uv, int src_stride_uv, int width, int height, void *dst);

int frame_data(vpx_codec_ctx_t* ctx, vpx_codec_iter_t* iter, int *key, char **data, int *len, int *pts);

int encode(vpx_codec_ctx_t *ctx, void *img, int w, int h, int64_t pts);