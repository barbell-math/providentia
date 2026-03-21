#ifndef BAR_PATH_PHYS_DATA_GPU
#define BAR_PATH_PHYS_DATA_GPU

#include "../../clib/glue.h"

#ifdef __cplusplus
extern "C" {
#endif
	extern void goSaveImage(unsigned char *data, int width, int height);

	enum BarPathTrackerErrCode_t CalcBarPathTrackerData(
		// barPathData_t* data,
		// barPathCalcHyperparams_t* opts
	);

#ifdef __cplusplus
}
#endif

#endif
