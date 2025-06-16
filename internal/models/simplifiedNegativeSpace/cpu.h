#ifndef SIMPLIFIED_NEGATIVE_SPACE_CPU
#define SIMPLIFIED_NEGATIVE_SPACE_CPU

#include <math.h>
#include <stdint.h>
#include "cgoStructs.h"

#ifdef __cplusplus
extern "C" {
#endif
	
	void calcModelStates(
		int64_t clientID,
		int32_t modelID,
		trainingLog_t *data,
		int64_t dataLen,
		int64_t startCalcsIdx,
		modelState_t *outValues,
		int64_t outValuesLen,
		opts_t *opts
	);

#ifdef __cplusplus
}
#endif

#endif
