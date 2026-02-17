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

#include <unistd.h>
#include "../../clib/glue.h"
#include "../../clib/errors.h"

namespace BarPathTracker {

static void display_frame(const AVFrame *frame) {
    int x, y;
    uint8_t *p0, *p;

	usleep(1000);
    // if (frame->pts != AV_NOPTS_VALUE) {
    //     if (last_pts != AV_NOPTS_VALUE) {
    //         /* sleep roughly the right amount of time;
    //          * usleep is in microseconds, just like AV_TIME_BASE. */
    //         delay = av_rescale_q(frame->pts - last_pts,
    //                              time_base, AV_TIME_BASE_Q);
    //         if (delay > 0 && delay < 1000000)
    //             usleep(delay);
    //     }
    //     last_pts = frame->pts;
    // }

    /* Trivial ASCII grayscale display. */
    p0 = frame->data[0];
    puts("\033c");
    for (y = 0; y < frame->height; y++) {
        p = p0;
        for (x = 0; x < frame->width; x++)
            putchar(" .-+#"[*(p++) / 52]);
        putchar('\n');
        p0 += frame->linesize[0];
    }
	printf("%d %d\n", frame->format, AV_PIX_FMT_NV12);
    fflush(stdout);
}

// The goal of this class is to implement the following cmd using the FFmpeg API:
//
// ./_deps/ffmpeg/bin/ffmpeg -i <video> -vf "scale=w=600:h=800,format=nv12" -
//
// That is: it takes a video source, decodes it on the CPU, scales it on the CPU,
// and makes sure the output format is nv12.
//
// Cmd that is helpful when debugging FFmpeg:
// ./_deps/ffmpeg/bin/ffmpeg -i <video> -vf "scale=w=600:h=800,format=nv12" -frames:v 1 output.png
struct SwDecode {
public:
	const char *file = NULL;

private:
	const char *filterDesc = "scale=78:24,transpose=cclock";

	AVFormatContext *inputCtx = NULL;
	int videoStreamIdx = 0;
	const AVCodec *decoder = NULL;
	AVCodecContext *decoderCtx = NULL;

	AVFilterContext *bufferSinkCtx;
	AVFilterContext *bufferSrcCtx;
	AVFilterGraph *filterGraph;

	AVPacket *packet;
	AVFrame *frame;
	AVFrame *filtFrame;

private:
	enum BarPathTrackerErrCode_t openInputFile() {
		int ret = 0;
	
		if (avformat_open_input(&this->inputCtx, this->file, NULL, NULL) < 0) {
			return CouldNotOpenVideoFileErr;
		}
		if (avformat_find_stream_info(this->inputCtx, NULL) < 0) {
			return CouldNotFindInputStreamInfoErr;
		}
		if ((ret = av_find_best_stream(
			this->inputCtx, AVMEDIA_TYPE_VIDEO, -1, -1, &this->decoder, 0
		)) < 0) {
			return CouldNotFindVideoStreamErr;
		}
		this->videoStreamIdx = ret;
		
		if (!(this->decoderCtx = avcodec_alloc_context3(this->decoder))) {
			return CouldNotAllocDecoderCtxErr;
		}
		avcodec_parameters_to_context(
			this->decoderCtx,
			this->inputCtx->streams[this->videoStreamIdx]->codecpar
		);
		if (avcodec_open2(this->decoderCtx, this->decoder, NULL)) {
			return CouldNotOpenCodecForStreamErr;
		}
		return NoBarPathTrackerErr;
	}

	enum BarPathTrackerErrCode_t initFilters() {
		enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;
		char args[512];
		const AVFilter *bufferSrc  = avfilter_get_by_name("buffer");
		const AVFilter *bufferSink = avfilter_get_by_name("buffersink");
		AVFilterInOut *outputs = avfilter_inout_alloc();
		AVFilterInOut *inputs  = avfilter_inout_alloc();
		AVRational timeBase = this->inputCtx->streams[this->videoStreamIdx]->time_base;
	
		this->filterGraph = avfilter_graph_alloc();
		if (!outputs || !inputs || !this->filterGraph) {
			err = CouldNotAllocFilterGraphErr;
			goto end;
		}
	
		// buffer video source: decoded frames from the decoder will be inserted here
		snprintf(
			args, sizeof(args),
			"video_size=%dx%d:pix_fmt=%d:time_base=%d/%d:pixel_aspect=%d/%d",
			this->decoderCtx->width, this->decoderCtx->height,
			this->decoderCtx->pix_fmt,
			timeBase.num, timeBase.den,
			this->decoderCtx->sample_aspect_ratio.num,
			this->decoderCtx->sample_aspect_ratio.den
		);
	
		if (avfilter_graph_create_filter(
			&this->bufferSrcCtx, bufferSrc, "in", args, NULL, this->filterGraph
		) < 0) {
			err = CouldNotCreateBufferSourceErr;
			goto end;
		}
	
		// buffer video sink: to terminate the filter chain
		if (!(this->bufferSinkCtx = avfilter_graph_alloc_filter(
			this->filterGraph, bufferSink, "out"
		))) {
			err = CouldNotCreateBufferSinkErr;
			goto end;
		}
	
		if (av_opt_set(
			this->bufferSinkCtx, "pixel_formats", "nv12", AV_OPT_SEARCH_CHILDREN
		)) {
			err = CouldNotSetSinkPixFmtErr;
			goto end;
		}
	
		if (avfilter_init_dict(this->bufferSinkCtx, NULL)) {
			err = CouldNotInitBufferSinkErr;
			goto end;
		}
	
		// Set the endpoints for the filter graph. The filter_graph will
		// be linked to the graph described by filters_descr.
	
		// The buffer source output must be connected to the input pad of
		// the first filter described by filters_descr; since the first
		// filter input label is not specified, it is set to "in" by
		// default.
		outputs->name		= av_strdup("in");
		outputs->filter_ctx	= this->bufferSrcCtx;
		outputs->pad_idx	= 0;
		outputs->next		= NULL;
	
		// The buffer sink input must be connected to the output pad of
		// the last filter described by filters_descr; since the last
		// filter output label is not specified, it is set to "out" by
		// default.
		inputs->name		= av_strdup("out");
		inputs->filter_ctx	= this->bufferSinkCtx;
		inputs->pad_idx		= 0;
		inputs->next		= NULL;
	
		if (avfilter_graph_parse_ptr(
			this->filterGraph, filterDesc, &inputs, &outputs, NULL
		) < 0) {
			err = CouldNotParseFilterErr;
			goto end;
		}
	
		if (avfilter_graph_config(this->filterGraph, NULL) < 0) {
			err = CouldNotConfigureFilterErr;
			goto end;
		}
	
	end:
		avfilter_inout_free(&inputs);
		avfilter_inout_free(&outputs);
		return err;
	}

	enum BarPathTrackerErrCode_t decodeAll() {
		enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;
		while (err == NoBarPathTrackerErr) {
			if (av_read_frame(this->inputCtx, this->packet) < 0) {
				err = CouldNotReadFrameErr;
				break;
			}
			if (this->videoStreamIdx == this->packet->stream_index) {
				if (avcodec_send_packet(this->decoderCtx, this->packet) < 0) {
					err = CouldNotSendPacketErr;
					break;
				}
				err = this->decodeOne();
			}
			av_packet_unref(this->packet);
		}
		return err;
	}

	enum BarPathTrackerErrCode_t decodeOne() {
		int ret = 0;
		enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;
		while (err == NoBarPathTrackerErr) {
			ret = avcodec_receive_frame(this->decoderCtx, this->frame);
			if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) break;
			if (ret < 0) return CouldNotReceiveFrameErr;
	
			// push the decoded frame into the filtergraph
			this->frame->pts = this->frame->best_effort_timestamp;
			if (av_buffersrc_add_frame_flags(
				this->bufferSrcCtx, this->frame, AV_BUFFERSRC_FLAG_KEEP_REF
			) < 0) {
				err = CouldNotAddFrameToFilterGraphErr;
				break;
			}
	
			// pull filtered frames from the filtergraph
			while (err == NoBarPathTrackerErr) {
				ret = av_buffersink_get_frame(this->bufferSinkCtx, this->filtFrame);
				if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) break;
				if (ret < 0) return CouldNotGetFrameFromFilterGraphErr;
				display_frame(this->filtFrame);
				av_frame_unref(this->filtFrame);
			}
			av_frame_unref(frame);
		}
		if (ret == AVERROR_EOF) {
			// signal EOF to the filtergraph
			if (av_buffersrc_add_frame_flags(this->bufferSinkCtx, NULL, 0) < 0) {
				av_log(NULL, AV_LOG_ERROR, "Error while closing the filtergraph\n");
				return CouldNotReceiveFrameErr;
			}
	
			// pull remaining frames from the filtergraph
			while (true) {
				ret = av_buffersink_get_frame(this->bufferSinkCtx, this->filtFrame);
				if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) break;
				if (ret < 0) {
					return CouldNotGetFrameFromFilterGraphErr;
				}
				display_frame(this->filtFrame);
				av_frame_unref(this->filtFrame);
			}
		}
		return err;
	}

public:
	SwDecode(const char *file): file(file) {}

	enum BarPathTrackerErrCode_t decode() {
		enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;
		if (!(this->frame = av_frame_alloc())) {
			err = CouldNotAllocFrameErr;
			goto done;
		}
		if (!(this->filtFrame = av_frame_alloc())) {
			err = CouldNotAllocFrameErr;
			goto done;
		}
		if (!(this->packet = av_packet_alloc())) {
			err = CouldNotAllocPacketErr;
			goto done;
		}
		if ((err = this->openInputFile()) != NoBarPathTrackerErr) goto done;
		if ((err = this->initFilters()) != NoBarPathTrackerErr) goto done;
		if ((err = this->decodeAll()) != NoBarPathTrackerErr) goto done;
	
	done:
		avfilter_graph_free(&this->filterGraph);
		avcodec_free_context(&this->decoderCtx);
		avformat_close_input(&this->inputCtx);
		av_frame_free(&this->frame);
		av_frame_free(&this->filtFrame);
		av_packet_free(&this->packet);
		return err;
	}
};

};
