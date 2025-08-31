#ifndef BAR_PATH_PHYS_DATA_CPU
#define BAR_PATH_PHYS_DATA_CPU

#include <math.h>
#include <stdint.h>
#include "../../glue/glue.h"

#ifdef __cplusplus
extern "C" {
#endif

	int64_t calcBarPathPhysData(
		barPathData_t* data,
		barPathCalcConf_t* bpOpts,
		physDataConf_t* pdOpts
	);

#ifdef __cplusplus
}
#endif

#endif
