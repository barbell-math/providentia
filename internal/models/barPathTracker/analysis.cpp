extern "C" {
	#include <libavutil/hwcontext.h>
}

#include <limits>
#include <cassert>
#include <cstdlib>
#include <stdint.h>
#include <unistd.h>
#include <iostream>
#include <functional>

#include "cpu.h"
#include "../../clib/glue.h"
#include "../../clib/errors.h"

namespace BarPathTracker {

class Pixel {
public:
	int x;
	int y;

	Pixel(): x(-1), y(-1) {}
	constexpr Pixel(int x, int y): x(x), y(y) {}

	static const uint8_t black = std::numeric_limits<uint8_t>::min();
	static const uint8_t white = std::numeric_limits<uint8_t>::max();

	Pixel operator+(const Pixel other) const {
		return Pixel(this->x+other.x, this->y+other.y);
	}
	Pixel operator-(const Pixel other) const {
		return Pixel(this->x-other.x, this->y-other.y);
	}

	friend bool operator==(const Pixel& a, const Pixel& b) {
		return a.x == b.x && a.y == b.y;
	}

	friend std::ostream& operator<<(std::ostream& os, Pixel p) {
		os << "Pixel{ X: " << p.x << " Y: " << p.y << "}";
		return os;
	}
};

class Frame {
public:
	uint8_t *frame = NULL;
	size_t width = 0;
	size_t height = 0;

public:
	Frame() {}
	~Frame() { if (this->frame!=NULL) free(this->frame); }

	void init(size_t width, size_t height) {
		if (this->frame == NULL) {
			this->frame = (uint8_t*)calloc(width*height, sizeof(uint8_t));
			if (this->frame == NULL) {
				// TODO - return oom err?? or throw??
			}
			this->width = width;
			this->height = height;
		} else {
			this->fill(Pixel::black);
		}
		assert(this->width == width && "Width of frames changed?");
		assert(this->height == height && "Height of frames changed?");
	}

	bool pixIsValid(const Pixel &pix) const {
		return (pix.x >=0 && pix.x<this->width) && (pix.y>=0 && pix.y<this->height);
	}

	void fill(uint8_t val) {
		for (int i=0; i<this->width*this->height; i++) this->frame[i]=val;
	}

	uint8_t& operator[](size_t x, size_t y) const {
		assert(x>=0 && x < this->width && "Invalid X dimension");
		assert(y>=0 && y < this->height && "Invalid Y dimension");
		return this->frame[y*this->width+x];
	}
	uint8_t& operator[](const Pixel &pixel) const {
		assert(pixel.x>=0 && pixel.x < this->width && "Invalid X dimension");
		assert(pixel.y>=0 && pixel.y < this->height && "Invalid Y dimension");
		return this->frame[pixel.y*this->width+pixel.x];
	}
};

class MeanAdaptiveThresholding {
public:
	const int neighborhoodSize = 11;
	const int thresholdOffset = 2;
	std::function<enum BarPathTrackerErrCode_t(const Frame &frame)> callback;
private:
	Frame frame;

private:
	double getNeighborhoodAvg(
		const AVFrame *iterFrame,
		size_t centerX,
		size_t centerY
	) {
		double cntr = 0;
		double total = 0;
		uint8_t *data = iterFrame->data[0];
		for (
			size_t y=std::max((size_t)0, centerY-this->neighborhoodSize/2);
			y<std::min(centerY+this->neighborhoodSize/2+1, (size_t)iterFrame->height);
			y++
		) {
			for (
				size_t x=std::max((size_t)0, centerX-this->neighborhoodSize/2);
				x<std::min(centerX+this->neighborhoodSize/2+1, (size_t)iterFrame->width);
				x++
			) {
				// Note: line size will not always be the same as frame->width
				// line size is aligned to word boundaries
				total += (double)(data[y*iterFrame->linesize[0] + x]);
				cntr++;
			}
		}
		if (cntr == 0) return 0;
		return total/cntr;
	}

public:
	MeanAdaptiveThresholding(
		const int neighborhoodSize,
		const int thresholdOffset,
		std::function<enum BarPathTrackerErrCode_t(const Frame &frame)> callback
	):
		neighborhoodSize(neighborhoodSize),
		thresholdOffset(thresholdOffset),
		callback(callback) {}

	enum BarPathTrackerErrCode_t onCallback(const AVFrame *iterFrame) {
		this->frame.init(iterFrame->width, iterFrame->height);
		uint8_t *data = iterFrame->data[0];
		for (size_t y=0; y<iterFrame->height; y++) {
			for (size_t x=0; x<iterFrame->width; x++) {
				double threshold = this->getNeighborhoodAvg(iterFrame, x, y);
				uint8_t *pix = &data[y*iterFrame->linesize[0] + x];
				uint8_t *newPix = &this->frame[x, y];
				if (*pix <= threshold-this->thresholdOffset) {
					*newPix = Pixel::white;
				} else {
					*newPix = Pixel::black;
				}
			}
		}
		// this->displayFrame(frame);
		// this->displayModifiedFrame();
		// goSaveImage(this->frame.frame, this->frame.width, this->frame.height);
		// return DecoderDoesNotSupportVulkanErr;
		return this->callback(this->frame);
	}

	// void displayFrame(const AVFrame *frame) {
	//     int x, y;
	//     uint8_t *p0, *p;
	// 
	// 	// usleep(1000);
	//     // if (frame->pts != AV_NOPTS_VALUE) {
	//     //     if (last_pts != AV_NOPTS_VALUE) {
	//     //         /* sleep roughly the right amount of time;
	//     //          * usleep is in microseconds, just like AV_TIME_BASE. */
	//     //         delay = av_rescale_q(frame->pts - last_pts,
	//     //                              time_base, AV_TIME_BASE_Q);
	//     //         if (delay > 0 && delay < 1000000)
	//     //             usleep(delay);
	//     //     }
	//     //     last_pts = frame->pts;
	//     // }
	// 
	//     /* Trivial ASCII grayscale display. */
	//     p0 = frame->data[0];
	//     puts("\033c");
	//     for (y = 0; y < frame->height; y++) {
	//         p = p0;
	//         for (x = 0; x < frame->width; x++)
	//             putchar(" .-+#"[*(p++) / 52]);
	//         putchar('\n');
	//         p0 += frame->linesize[0];
	//     }
	// 	printf("%d %d\n", frame->format, AV_PIX_FMT_NV12);
	//     fflush(stdout);
	// }

	// void displayModifiedFrame() {
	//     int x, y;
	//     uint8_t *p0, *p;
	// 
	// 	// usleep(1000);
	//     // if (frame->pts != AV_NOPTS_VALUE) {
	//     //     if (last_pts != AV_NOPTS_VALUE) {
	//     //         /* sleep roughly the right amount of time;
	//     //          * usleep is in microseconds, just like AV_TIME_BASE. */
	//     //         delay = av_rescale_q(frame->pts - last_pts,
	//     //                              time_base, AV_TIME_BASE_Q);
	//     //         if (delay > 0 && delay < 1000000)
	//     //             usleep(delay);
	//     //     }
	//     //     last_pts = frame->pts;
	//     // }
	// 
	//     /* Trivial ASCII grayscale display. */
	//     p0 = this->modifiedFrame;
	//     puts("\033c");
	//     for (y = 0; y < this->height; y++) {
	//         p = p0;
	//         for (x = 0; x < this->width; x++)
	//             putchar(" .-+#"[*(p++) / 52]);
	//         putchar('\n');
	//         p0 += this->width;
	//     }
	//     fflush(stdout);
	// }
};

class MooreNeighborTracing {
public:
	size_t minContourLen = 1000;

private:
	Frame frame;
	std::vector<Pixel> border;
	std::vector<Pixel> contours;

	static constexpr uint8_t unconsidered = Pixel::black;
	static constexpr uint8_t considered = Pixel::white;

	enum Direction { N=0, NE=1, E=2, SE=3, S=4, SW=5, W=6, NW=7 };
	static constexpr Pixel cwMooreNeighborhood[] = {
		Pixel(0, 1),  Pixel(1, 1),   Pixel(1, 0),  Pixel(1,-1),
		Pixel(0, -1), Pixel(-1, -1), Pixel(-1, 0), Pixel(-1, 1),
	};

	inline enum Direction rotCW(enum Direction dir) {
		return (enum Direction)((dir+1)%8);
	}

	enum Direction getDirection(const Pixel &from, const Pixel &to) {
		size_t idx = 0;
		const Pixel dir = from-to;
		for (auto& iterDir : MooreNeighborTracing::cwMooreNeighborhood) {
			if (dir == iterDir) {
				return (enum Direction)idx;
			}
			idx++;
		}
		assert(false && "Unreachable: invalid direction");
	}

	void followBorder(
		const Frame &iterFrame,
		const Pixel &startBlack,
		const Pixel &startWhite
	) {
		this->border.emplace_back(startWhite);
		this->frame[startWhite] = MooreNeighborTracing::considered;

		Pixel curPix = startWhite;
		Pixel prevPix = startBlack;
		while (true) {
			bool foundNextPix = false;
			enum Direction origDir = this->getDirection(prevPix, curPix);
			enum Direction movedDir = this->rotCW(origDir);
			for (size_t i=0; i<std::size(MooreNeighborTracing::cwMooreNeighborhood); i++) {
				Pixel iterPix = curPix + MooreNeighborTracing::cwMooreNeighborhood[(int)movedDir];
				if (!iterFrame.pixIsValid(iterPix)) continue;
				if (iterFrame[iterPix] == Pixel::white) {
					this->border.emplace_back(iterPix);
					this->frame[iterPix] = MooreNeighborTracing::considered;
					prevPix = curPix;
					curPix = iterPix;
					foundNextPix = true;
					break;
				}
				movedDir = this->rotCW(movedDir);
			}
			if (!foundNextPix) break;
			if (curPix == startWhite) break;
		}

		if (this->border.size() > this->minContourLen) {
			this->contours.insert(
				this->contours.end(), this->border.begin(), this->border.end()
			);
		}
		this->border.clear();
	}

public:
	enum BarPathTrackerErrCode_t onCallback(const Frame &iterFrame) {
		this->contours.clear();
		this->frame.init(iterFrame.width, iterFrame.height);

		Pixel curPix;
		Pixel prevPix;
		for(size_t y=0; y<iterFrame.height; y++) {
			for (size_t x=1; x<iterFrame.width; x++) {
				prevPix = Pixel(x-1, y);
				curPix = Pixel(x, y);
				if (this->frame[curPix] == MooreNeighborTracing::considered) continue;
				if (iterFrame[prevPix] != Pixel::black) continue;
				if (iterFrame[curPix] != Pixel::white) continue;
				this->followBorder(iterFrame, prevPix, curPix);
			}
		}

		goSaveImage(this->frame.frame, this->frame.width, this->frame.height);
		return CouldNotAllocPacketErr;
	}
};

};
