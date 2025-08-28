#ifndef BAR_PATH_PHYS_DATA_CPU
#define BAR_PATH_PHYS_DATA_CPU

#include <math.h>
#include <stdint.h>
#include "../../glue/glue.h"

#ifdef __cplusplus
extern "C" {
#endif

	int64_t calcBarPathPhysData(
		int64_t timeLen,
		double_t* time,
		posVec2_t* pos,
		velVec2_t* vel,
		accVec2_t* acc,
		jerkVec2_t* jerk,
		workVec2_t* work,
		impulseVec2_t* impulse,
		forceVec2_t* force,
		barPathCalcConf_t *bpOpts,
		physDataConf_t *pdOpts
	);

#ifdef __cplusplus
}
#endif

#endif
