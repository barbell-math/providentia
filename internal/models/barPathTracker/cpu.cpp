#include <iostream>
#include "../../clib/glue.h"
#include "./hwDecode.cpp"
#include "./swDecode.cpp"

namespace BarPathTracker {

extern "C" enum BarPathTrackerErrCode_t CalcBarPathTrackerData() {
	enum BarPathTrackerErrCode_t err = NoBarPathTrackerErr;

	std::cout << "HELLO from tracker!" << std::endl;

	const char* file = "/home/jack/Downloads/VID_20250930_165654085.mp4";

	SwDecode swDecode(file);
	err = swDecode.decode();

	// HwDecode hwDecode(file);
	// err = hwDecode.decode();

	std::cout << "Made it to the end " << err << std::endl;
	return err;
}

};
