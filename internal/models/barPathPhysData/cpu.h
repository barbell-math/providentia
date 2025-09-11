#ifndef BAR_PATH_PHYS_DATA_CPU
#define BAR_PATH_PHYS_DATA_CPU

#include "../../clib/glue.h"

#ifdef __cplusplus
extern "C" {
#endif

	enum BarPathCalcErrCode_t calcBarPathPhysData(
		barPathData_t* data,
		barPathCalcConf_t* opts
	);

#ifdef __cplusplus
}
#endif

#endif
