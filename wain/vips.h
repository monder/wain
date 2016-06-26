#include <stdlib.h>
#include <string.h>
#include <vips/vips.h>

int vips_wain_init() {
    if (vips_init("wain-vips")) vips_error_exit("unable to start VIPS");
    vips_concurrency_set(1);
	vips_cache_set_max_mem(100 * 1048576); // 100Mb
	vips_cache_set_max(500);
    
    return 0;
}

static double calculate_shrink(VipsImage *im, gboolean crop, int width, int height) {
	VipsDirection direction;
	double horizontal = (double)im->Xsize / (width + 0.1);
	double vertical = (double)im->Ysize / (height + 0.1);

    if (height == 0) {
        return horizontal;
    }

	if (crop) {
		if (horizontal < vertical)
			direction = VIPS_DIRECTION_HORIZONTAL;
		else
			direction = VIPS_DIRECTION_VERTICAL;
	} else {
		if (horizontal < vertical)
			direction = VIPS_DIRECTION_VERTICAL;
		else
			direction = VIPS_DIRECTION_HORIZONTAL;
	}

	return (direction == VIPS_DIRECTION_HORIZONTAL ? horizontal : vertical);  
}

int vips_wain_resize(void *src, size_t src_len, void **dst, size_t *dst_len, int width, int height) {

    VipsImage *image = vips_image_new_from_buffer(src, src_len, NULL, NULL);
    VipsImage *tmp = NULL;

    if (height == 0) {
        double shrink = calculate_shrink(image, FALSE, width, height);
        vips_resize(image, &tmp, 1.0 / shrink, NULL);
        g_object_unref(image);
        image = tmp;
    } else {
        VipsImage *background = NULL;
        VipsImage *foreground = NULL;
        // Create background
        {
            // Shrink
            double shrink = calculate_shrink(image, TRUE, width, height);
            vips_resize(image, &tmp, 1.0 / shrink, NULL);
            // Crop
            int left = (tmp->Xsize - width) / 2;
		    int top = (tmp->Ysize - height) / 2;
		    vips_extract_area(tmp, &background, left, top, width, height, NULL);
            g_object_unref(tmp);
            // Blur
            vips_gaussblur(background, &tmp, 30, NULL);
            g_object_unref(background);
        }
        // Resize image and add it to bg
        {
            // Resize
            double shrink = calculate_shrink(image, FALSE, width, height);
            vips_resize(image, &foreground, 1.0 / shrink, NULL);
            // Combine
            int left = (width - foreground->Xsize) / 2;
		    int top = (height - foreground->Ysize) / 2;
            vips_insert(tmp, foreground, &background, left, top, NULL);
            g_object_unref(tmp);
            g_object_unref(foreground);
        }
        g_object_unref(image);
        image = background;
    }

    int ret = vips_jpegsave_buffer(
        image,
        dst, dst_len,
        "strip", TRUE,
        "Q", 96,
        "optimize_coding", TRUE,
        "interlace", TRUE,
        "trellis_quant", TRUE,
        "optimize_scans", TRUE,
        "no-subsample", FALSE,
        NULL
    );
    g_object_unref(image);
    return ret;
}


