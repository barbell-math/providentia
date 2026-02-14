extern "C" {
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
#include "../../clib/errors.h"


static AVPixelFormat hwPixFmt;

namespace BarPathTracker {

struct HwDecode {
	const char *file = NULL;

	int videoStream = 0;
	enum AVHWDeviceType type = AV_HWDEVICE_TYPE_NONE;
	AVPacket *packet = NULL;
	AVFormatContext *input_ctx = NULL;
	const AVCodec *decoder = NULL;
	AVCodecContext *decoder_ctx = NULL;
	AVStream *video = NULL;

	// static enum AVPixelFormat hwPixFmt;
	// TODO - these are static...
	AVBufferRef *hw_device_ctx = NULL;

private:
enum BarPathTrackerErrCode_t checkVAAPIDeviceAvailable() {
	this->type = av_hwdevice_find_type_by_name("vaapi");
	if (type == AV_HWDEVICE_TYPE_NONE) {
		return VAAPINotSupportedErr;

		// Helpfull for debugging - TODO make logging work with Go somehow?
		// fprintf(stderr, "Available device types:");
		// while((type = av_hwdevice_iterate_types(type)) != AV_HWDEVICE_TYPE_NONE)
		// 	fprintf(stderr, " %s", av_hwdevice_get_type_name(type));
		// fprintf(stderr, "\n");
	}
	return NoBarPathTrackerErr;
}

enum BarPathTrackerErrCode_t selectVAAPIDevice() {
	for (int i=0;; i++) {
		const AVCodecHWConfig *config = avcodec_get_hw_config(this->decoder, i);
		if (!config) {
			// fprintf(stderr, "Decoder %s does not support device type %s.\n",
			// 		decoder->name, av_hwdevice_get_type_name(type));
			return DecoderDoesNotSupportVAAPIErr;
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


enum BarPathTrackerErrCode_t hwDecoderInit(
	AVCodecContext *ctx,
	const enum AVHWDeviceType type
) {
	int avErr = 0;
	enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;
	// TODO - allow device to be specified here - https://stackoverflow.com/a/68742185
	if ((avErr = av_hwdevice_ctx_create(&hw_device_ctx, type, NULL, NULL, 0)) < 0) {
	  // fprintf(stderr, "Failed to create specified HW device.\n");
	  return CouldNotCreateHwDeviceErr;
	}
	ctx->hw_device_ctx = av_buffer_ref(hw_device_ctx);
	return err;
}

enum BarPathTrackerErrCode_t filterInit() {}

enum BarPathTrackerErrCode_t run() {
	int ret = 0;
	while (ret >= 0) {
		if ((ret = av_read_frame(this->input_ctx, this->packet)) < 0) break;
		if (this->videoStream == this->packet->stream_index) {
			ret = decodeFrame(this->decoder_ctx, this->packet);
		}
		av_packet_unref(this->packet);
	}
	ret = decodeFrame(this->decoder_ctx, NULL);	// Flush the decoder
	return NoBarPathTrackerErr;
}

// TODO
// enum BarPathTrackerErrCode_t decodeFrame(AVCodecContext *avctx, AVPacket *packet) {
int decodeFrame(AVCodecContext *avctx, AVPacket *packet) {
	AVFrame *frame = NULL;
	AVFrame *swFrame = NULL;
	AVFrame *tmpFrame = NULL;
	uint8_t *buffer = NULL;
	int size;
	int ret = 0;

	ret = avcodec_send_packet(avctx, packet);
	if (ret < 0) {
		fprintf(stderr, "Error during decoding\n");
		return ret;
		// return NoBarPathTrackerErr; // TODO - proper error
	}

	while (1) {
		if (!(frame = av_frame_alloc()) || !(swFrame = av_frame_alloc())) {
			fprintf(stderr, "Can not alloc frame\n");
			ret = AVERROR(ENOMEM);
			// TODO - set error
			goto fail;
		}

		ret = avcodec_receive_frame(avctx, frame);
		if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
			av_frame_free(&frame);
			av_frame_free(&swFrame);
			return 0;
			// return NoBarPathTrackerErr; - this is the correct return
		} else if (ret < 0) {
			fprintf(stderr, "Error while decoding\n");
			// TODO - set error
			goto fail;
		}

		if (frame->format == hwPixFmt) {
			/* retrieve data from GPU to CPU */
			if ((ret = av_hwframe_transfer_data(swFrame, frame, 0)) < 0) {
				fprintf(stderr, "Error transferring the data to system memory\n");
				// TODO - set error
				goto fail;
			}
			tmpFrame = swFrame;
		} else {
			tmpFrame = frame;
		}

		size = av_image_get_buffer_size(
			(enum AVPixelFormat)tmpFrame->format,
			tmpFrame->width, tmpFrame->height, 1
		);
		buffer = (uint8_t*)av_malloc(size);
		if (!buffer) {
			fprintf(stderr, "Can not alloc buffer\n");
			ret = AVERROR(ENOMEM);
			// TODO - set error
			goto fail;
		}
		ret = av_image_copy_to_buffer(
			buffer, size,
			(const uint8_t * const *)tmpFrame->data,
			(const int *)tmpFrame->linesize,
			(enum AVPixelFormat)tmpFrame->format,
			tmpFrame->width, tmpFrame->height, 1
		);
		if (ret < 0) {
			fprintf(stderr, "Can not copy image to buffer\n");
			// TODO - set error
			goto fail;
		}
		// Check if the frame is a planar YUV 4:2:0, 12bpp
		// That is the format of the provided .mp4 file
		// RGB formats will definitely not give a gray image
		// Other YUV image may do so, but untested, so give a warning
		if (tmpFrame->format != AV_PIX_FMT_YUV420P) {
			printf("%d %d\n", tmpFrame->format, AV_PIX_FMT_NV12);
			printf("Warning: the generated file may not be a grayscale image, but could e.g. be just the R component if the video format is RGB");
		}

		// if ((ret = fwrite(buffer, 1, size, output_file)) < 0) {
		// 	fprintf(stderr, "Failed to dump raw data.\n");
		// 	// TODO - set error
		// 	goto fail;
		// }

	fail:
		av_frame_free(&frame);
		av_frame_free(&swFrame);
		av_freep(&buffer);
		if (ret < 0)
			return ret;
	}
}

public:
HwDecode(const char *file): file(file) {}

enum BarPathTrackerErrCode_t decode() {
	int ret = 0;
	enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;
	if ((err = this->checkVAAPIDeviceAvailable()) != NoBarPathTrackerErr) {
		goto done;
	}
	if ((this->packet = av_packet_alloc()) == NULL) {
		err = CouldNotAllocateAVPacketErr;
		goto done;
	}
	if (avformat_open_input(&this->input_ctx, this->file, NULL, NULL) != 0) {
		// fprintf(stderr, "Cannot open input file '%s'\n", argv[2]);
		err = CouldNotOpenVideoFileErr;
		goto done;
	}
	if (avformat_find_stream_info(this->input_ctx, NULL) < 0) {
		// fprintf(stderr, "Cannot find input stream information.\n");
		err = CouldNotFindInputStreamInfoErr;
		goto done;
	}
	if ((this->videoStream = av_find_best_stream(
		this->input_ctx, AVMEDIA_TYPE_VIDEO, -1, -1, &this->decoder, 0
	)) < 0 ) {
		// fprintf(stderr, "Cannot find a video stream in the input file\n");
		err = CouldNotFindVideoStreamErr;
		goto done;
	}
	if ((err = this->selectVAAPIDevice()) != NoBarPathTrackerErr) {
		goto done;
	}
	if (!(this->decoder_ctx = avcodec_alloc_context3(this->decoder))) {
		throw Err::OOM{ .Desc="Could not alloc avcodec" };
	}
	this->video = this->input_ctx->streams[this->videoStream];
	if (avcodec_parameters_to_context(
		this->decoder_ctx, this->video->codecpar
	) < 0) {
		avcodec_free_context(&decoder_ctx);
		err = AVCodecParametersToCtxtErr;
		goto done;
	}
	this->decoder_ctx->get_format = getHwFormat;
	if ((err = hwDecoderInit(decoder_ctx, type)) < 0) {
		goto done;
	}
	if ((ret = avcodec_open2(decoder_ctx, decoder, NULL)) < 0) {
		// fprintf(stderr, "Failed to open codec for stream #%u\n", video_stream);
		err = CouldNotOpenCodecForStreamErr;
		goto done;
	}
	if ((err = run()) != 0) {
		goto done;
	}

done:
	// if (this->output_file) fclose(this->output_file);
	av_packet_free(&this->packet);
	avcodec_free_context(&this->decoder_ctx);
	avformat_close_input(&this->input_ctx);
	av_buffer_unref(&this->hw_device_ctx);
	return err;
}

};


extern "C" enum BarPathTrackerErrCode_t CalcBarPathTrackerData() {
	std::cout << "HELLO from tracker!" << std::endl;

	const char* file = "/home/jack/Downloads/VID_20250930_165654085.mp4";
	HwDecode hwDecode(file);
	hwDecode.decode();


	std::cout << "Made it to the end" << std::endl;
	return NoBarPathTrackerErr;
}

};
