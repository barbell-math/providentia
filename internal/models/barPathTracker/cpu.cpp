extern "C" {
	#include <time.h>
	#include <libavcodec/avcodec.h>
	#include <libavformat/avformat.h>
	#include <libavutil/mem.h>
	#include <libavutil/pixdesc.h>
	#include <libavutil/hwcontext.h>
	#include <libavutil/opt.h>
	#include <libavutil/avassert.h>
	#include <libavutil/imgutils.h>
}

#include <iostream>
#include "../../clib/glue.h"


namespace BarPathTracker {

extern "C" enum BarPathTrackerErrCode_t CalcBarPathTrackerData() {
	std::cout << "HELLO" << std::endl;

	const char* device = "GPU";
    enum AVHWDeviceType type;
    type = av_hwdevice_find_type_by_name(device);
    if (type == AV_HWDEVICE_TYPE_NONE) {
        fprintf(stderr, "Device type %s is not supported.\n", device);
        fprintf(stderr, "Available device types:");
        while((type = av_hwdevice_iterate_types(type)) != AV_HWDEVICE_TYPE_NONE)
            fprintf(stderr, " %s", av_hwdevice_get_type_name(type));
        fprintf(stderr, "\n");
    }

	return NoBarPathTrackerErr;
}

};
