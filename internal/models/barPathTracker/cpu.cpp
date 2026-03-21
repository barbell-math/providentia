#include <iostream>
#include "cpu.h"
#include "../../clib/glue.h"
#include "./hwDecode.cpp"
#include "./swDecode.cpp"

namespace BarPathTracker {

extern "C" enum BarPathTrackerErrCode_t CalcBarPathTrackerData() {
	enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;

	std::cout << "HELLO from tracker!" << std::endl;

	const char* file = "/home/jack/Documents/research/vid/lifter3SquatTest.mp4";
	// const char* file = "/home/jack/Documents/research/vid/lifter1DeadliftTestShort.mp4";

	SwSuzukiAbeFindContours contourFinder;
	SwMeanAdaptiveThresholding thresholding(11, 2, [&contourFinder](const SwFrame *frame) {
		return contourFinder.onCallback(frame);
	});
	SwDecode swDecode(file, [&thresholding](const AVFrame *frame) {
		return thresholding.onCallback(frame);
	});
	err = swDecode.decode();

	// HwDecode hwDecode(file);
	// err = hwDecode.decode();

	std::cout << "Made it to the end " << err << std::endl;
	return err;
}

};
