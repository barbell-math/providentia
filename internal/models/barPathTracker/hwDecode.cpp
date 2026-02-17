extern "C" {
	#include <libavcodec/avcodec.h>
	#include <libavformat/avformat.h>
	#include <libavutil/mem.h>
	#include <libavutil/pixdesc.h>
	#include <libavutil/hwcontext.h>
	#include <libavutil/opt.h>
	#include <libavutil/avassert.h>
	#include <libavutil/imgutils.h>
	#include <libavfilter/buffersink.h>
	#include <libavfilter/buffersrc.h>
}

#include <iostream>
#include "../../clib/glue.h"
#include "../../clib/errors.h"

// Note - it is ok to make these static because when this code is called from
// go it will always be in it's own thread.
static AVPixelFormat hwPixFmt;
static AVBufferRef *hwDeviceCtx = NULL;
// static AVBufferRef *filterHwDeviceCtx = NULL;

namespace BarPathTracker {

// ./_deps/ffmpeg/bin/ffmpeg -init_hw_device vulkan=vk:0 -filter_hw_device vk -hwaccel vulkan -hwaccel_output_format vulkan -i ~/Downloads/VID_20250930_165654085.mp4 -vf "scale_vulkan=w=600:h=800:format=nv12,hwdownload,format=nv12" -frames:v 1 output.png
struct HwDecode {
	const char *filterDesc = "scale_vulkan=w=800:h=600:format=nv12,hwdownload";
	const char *file = NULL;

	int videoStream = 0;
	enum AVHWDeviceType type = AV_HWDEVICE_TYPE_NONE;
	AVFilterGraph *filterGraph = NULL;

	const AVCodec *decoder = NULL;
	AVCodecContext *decoderCtx = NULL;
	AVFormatContext *inputCtx = NULL;
	AVFilterContext *bufferSrcCtx = NULL;
	AVFilterContext *bufferSinkCtx = NULL;

	AVStream *video = NULL;
	AVPacket *packet = NULL;

	AVFrame *frame = NULL;
	AVFrame *swFrame = NULL;
	// AVFrame *filtFrame = NULL;

private:

enum BarPathTrackerErrCode_t selectVulkanDevice() {
	for (int i=0;; i++) {
		const AVCodecHWConfig *config = avcodec_get_hw_config(this->decoder, i);
		if (!config) {
			return DecoderDoesNotSupportVulkanErr;
		}
		if (
			config->methods & AV_CODEC_HW_CONFIG_METHOD_HW_DEVICE_CTX &&
			config->device_type == this->type
		) {
			hwPixFmt = config->pix_fmt;
			break;
		}
	}
	return NoBarPathTrackerErr;
}

static enum AVPixelFormat getHwFormat(
	AVCodecContext *ctx,
	const enum AVPixelFormat *pixFmts
) {
	const enum AVPixelFormat *p;
	for (p = pixFmts; *p != -1; p++) {
		if (*p == hwPixFmt) return *p;
	}

	// TODO - make logging work
	// fprintf(stderr, "Failed to get HW surface format.\n");
	return AV_PIX_FMT_NONE;
}

enum BarPathTrackerErrCode_t initFilters() {
	enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;
	int ret = 0;
	char args[512];
	const AVFilter *bufferSrc  = avfilter_get_by_name("buffer");
	const AVFilter *bufferSink = avfilter_get_by_name("buffersink");
	AVFilterInOut *outputs = avfilter_inout_alloc();
	AVFilterInOut *inputs  = avfilter_inout_alloc();
	AVRational timeBase = this->inputCtx->streams[this->videoStream]->time_base;

	this->filterGraph = avfilter_graph_alloc();
	if (!outputs || !inputs || !this->filterGraph) {
		err = CouldNotAllocFilterGraphErr;
		goto end;
	}

	// buffer video source: the decoded frames from the decoder will be inserted here
	snprintf(
		args, sizeof(args),
		"video_size=%dx%d:pix_fmt=%d:time_base=%d/%d:pixel_aspect=%d/%d",
		this->decoderCtx->width,
		this->decoderCtx->height,
		hwPixFmt,
		// this->decoderCtx->pix_fmt,
		timeBase.num, timeBase.den,
		this->decoderCtx->sample_aspect_ratio.num,
		this->decoderCtx->sample_aspect_ratio.den
	);
	printf("%s\n", args);

	// this->bufferSrcCtx = avfilter_graph_alloc_filter(
	// 	this->filterGraph, bufferSrc, "in"
	// );
	// if (!this->bufferSrcCtx) {
	// 	err = CouldNotCreateBufferSourceErr;
	// 	goto end;
	// }
	// // this->bufferSrcCtx->hw_device_ctx = av_buffer_ref(hwDeviceCtx);
	// // printf("%p\n", this->bufferSrcCtx->hw_device_ctx);
	// if ((ret = avfilter_init_str(this->bufferSrcCtx, args)) <0 ) {
	// 	std::cout << "HERE?" << std::endl;
	// 	err = CouldNotCreateBufferSourceErr;
	// 	goto end;
	// }

	if ((ret = avfilter_graph_create_filter(
		&this->bufferSrcCtx, bufferSrc, "in", args, NULL, this->filterGraph
	)) < 0) {
		err = CouldNotCreateBufferSourceErr;
		goto end;
	}

	// // buffer video source: to start the filter chain
	// this->bufferSrcCtx = avfilter_graph_alloc_filter(
	// 	this->filterGraph, bufferSrc, "in"
	// );
	// if (!this->bufferSrcCtx) {
	// 	err = CouldNotCreateBufferSourceErr;
	// 	goto end;
	// }

	// if ((ret = av_opt_set(
	// 	this->bufferSrcCtx, "pixel_formats", "nv12", AV_OPT_SEARCH_CHILDREN
	// )) < 0) {
	// 	printf("Here...\n");
	// 	err = CouldNotSetOutputPixelFmtErr;//TODO - correct error
	// 	goto end;
	// }

	// if ((ret = avfilter_init_str(this->bufferSrcCtx, args)) < 0) {
	// 	err = CouldNotCreateBufferSourceErr;
	// 	goto end;
	// }

	// buffer video sink: to terminate the filter chain
	this->bufferSinkCtx = avfilter_graph_alloc_filter(
		this->filterGraph, bufferSink, "out"
	);
	if (!this->bufferSinkCtx) {
		err = CouldNotCreateBufferSinkErr;
		goto end;
	}

	if ((ret = av_opt_set(
		this->bufferSinkCtx, "pixel_formats", "nv12", AV_OPT_SEARCH_CHILDREN
	)) < 0) {
		err = CouldNotSetSinkPixFmtErr;
		goto end;
	}

	if ((ret = avfilter_init_dict(this->bufferSinkCtx, NULL)) < 0) {
		err = CouldNotInitBufferSinkErr;
		goto end;
	}

	 // Set the endpoints for the filter graph. The filter_graph will
	 // be linked to the graph described by filterDesc.

	// The buffer source output must be connected to the input pad of
	// the first filter described by filterDesc; since the first
	// filter input label is not specified, it is set to "in" by
	// default.
	outputs->name		= av_strdup("in");
	outputs->filter_ctx	= this->bufferSrcCtx;
	outputs->pad_idx	= 0;
	outputs->next		= NULL;

	// The buffer sink input must be connected to the output pad of
	// the last filter described by filterDesc; since the last
	// filter output label is not specified, it is set to "out" by
	// default.
	inputs->name		= av_strdup("out");
	inputs->filter_ctx 	= this->bufferSinkCtx;
	inputs->pad_idx		= 0;
	inputs->next		= NULL;

	if ((ret = avfilter_graph_parse_ptr(
		this->filterGraph, HwDecode::filterDesc, &inputs, &outputs, NULL
	)) < 0) {
		err = CouldNotParseFilterErr;
		goto end;
	}

	if ((ret = avfilter_graph_config(this->filterGraph, NULL)) < 0) {
		err = CouldNotConfigureFilterErr;
		goto end;
	}

end:
	avfilter_inout_free(&inputs);
	avfilter_inout_free(&outputs);
	return err;
}

enum BarPathTrackerErrCode_t hwDecoderInit() {
	int ret = 0;
	enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;
	// TODO - allow device to be specified here - https://stackoverflow.com/a/68742185
	if ((ret = av_hwdevice_ctx_create(
		&hwDeviceCtx, this->type, NULL, NULL, 0
	)) < 0) {
	  return CouldNotCreateHwDeviceErr;
	}
	this->decoderCtx->hw_device_ctx = av_buffer_ref(hwDeviceCtx);
	return err;
}

enum BarPathTrackerErrCode_t run() {
	int ret = 0;
	enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;

	if ((this->packet = av_packet_alloc()) == NULL) return CouldNotAllocPacketErr;
	if ((this->frame = av_frame_alloc()) == NULL) return CouldNotAllocFrameErr;
	if ((this->swFrame = av_frame_alloc()) == NULL) return CouldNotAllocFrameErr;

	while (err == NoBarPathTrackerErr) {
		if ((ret = av_read_frame(this->inputCtx, this->packet)) < 0) {
			err = CouldNotReadFrameErr;
			break;
		}
		if (this->videoStream == this->packet->stream_index) {
			if ((ret =avcodec_send_packet(this->decoderCtx, this->packet)) < 0) {
				return CouldNotSendPacketErr;
			}
			err = decodeFrame();
		}
		av_packet_unref(this->packet);
	}

	avcodec_send_packet(this->decoderCtx, NULL); // Flush the decoder
	av_frame_free(&this->frame);
	av_frame_free(&this->swFrame);
	return err;
}

enum BarPathTrackerErrCode_t decodeFrame() {
	int ret = 0;

	while (true) {
		ret = avcodec_receive_frame(this->decoderCtx, this->frame);
		if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
			return NoBarPathTrackerErr;
		} else if (ret < 0) {
			return CouldNotReceiveFrameErr;
		}

		// TODO - push frame through filter graph
		// frame->pts = frame->best_effort_timestamp;

		// /* push the decoded frame into the filtergraph */
		// if (av_buffersrc_add_frame_flags(buffersrc_ctx, frame, AV_BUFFERSRC_FLAG_KEEP_REF) < 0) {
		// 	av_log(NULL, AV_LOG_ERROR, "Error while feeding the filtergraph\n");
		// 	break;
		// }

		// /* pull filtered frames from the filtergraph */
		// while (1) {
		// 	ret = av_buffersink_get_frame(buffersink_ctx, filt_frame);
		// 	if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF)
		// 		break;
		// 	if (ret < 0)
		// 		goto end;
		// 	display_frame(filt_frame, buffersink_ctx->inputs[0]->time_base);
		//	processFrame();
		// 	av_frame_unref(filt_frame);
		// }

		processFrame();
	}
}

enum BarPathTrackerErrCode_t processFrame() {
	int ret = 0;
	AVFrame *tmpFrame = NULL;
	if (this->frame->format == hwPixFmt) {
		printf("HERE?\n");
		if ((ret = av_hwframe_transfer_data(this->swFrame, this->frame, 0)) < 0) {
			return CouldNotTransferDataFromGPUToCPUErr;
		}
		tmpFrame = this->swFrame;
	} else {
		tmpFrame = this->frame;
	}

	if (tmpFrame->format != AV_PIX_FMT_YUV420P) {
		printf("%d %d\n", tmpFrame->format, AV_PIX_FMT_NV12);
		printf("Warning: the generated file may not be a grayscale image, but could e.g. be just the R component if the video format is RGB\n");
	}
	printf("%u\n", *tmpFrame->data[0]);

	return NoBarPathTrackerErr;
}

public:

HwDecode(const char *file): file(file) {}

enum BarPathTrackerErrCode_t decode() {
	int ret = 0;
	enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;
	if ((this->type = av_hwdevice_find_type_by_name("vulkan")) == AV_HWDEVICE_TYPE_NONE) {
		return VulkanNotSupportedErr;
	}
	if (avformat_open_input(&this->inputCtx, this->file, NULL, NULL) != 0) {
		err = CouldNotOpenVideoFileErr;
		goto done;
	}
	if (avformat_find_stream_info(this->inputCtx, NULL) < 0) {
		err = CouldNotFindInputStreamInfoErr;
		goto done;
	}
	if ((this->videoStream = av_find_best_stream(
		this->inputCtx, AVMEDIA_TYPE_VIDEO, -1, -1, &this->decoder, 0
	)) < 0 ) {
		err = CouldNotFindVideoStreamErr;
		goto done;
	}
	if ((err = this->selectVulkanDevice()) != NoBarPathTrackerErr) goto done;
	if (!(this->decoderCtx = avcodec_alloc_context3(this->decoder))) {
		throw Err::OOM{ .Desc="Could not alloc avcodec" };	// TODO - remove?
	}
	this->video = this->inputCtx->streams[this->videoStream];
	if (avcodec_parameters_to_context(
		this->decoderCtx, this->video->codecpar
	) < 0) {
		err = AVCodecParametersToCtxtErr;
		goto done;
	}
	this->decoderCtx->get_format = getHwFormat;
	if ((err = hwDecoderInit()) < 0) {
		goto done;
	}
	if ((ret = avcodec_open2(this->decoderCtx, this->decoder, NULL)) < 0) {
		err = CouldNotOpenCodecForStreamErr;
		goto done;
	}
	if ((err = initFilters()) != NoBarPathTrackerErr) goto done;
	if ((err = run()) != 0) goto done;

done:
	avfilter_graph_free(&this->filterGraph);
	avcodec_free_context(&this->decoderCtx);
	avformat_close_input(&this->inputCtx);
	av_packet_free(&this->packet);
	av_buffer_unref(&hwDeviceCtx);
	// av_buffer_unref(&filterHwDeviceCtx);
	return err;
}

};

};
