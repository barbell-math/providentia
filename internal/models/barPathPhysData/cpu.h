#ifndef BAR_PATH_PHYS_DATA_CPU
#define BAR_PATH_PHYS_DATA_CPU

#include "../../clib/glue.h"

#ifdef __cplusplus
extern "C" {
#endif

	enum BarPathCalcErrCode_t CalcBarPathPhysData(
		barPathData_t* data,
		barPathCalcHyperparams_t* opts
	);

#ifdef __cplusplus
}
#endif

#endif
