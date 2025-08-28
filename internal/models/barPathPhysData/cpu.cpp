#include <cmath>
#include "cpu.h"
#include "cgoStructs.h"

extern "C" int64_t calcBarPathPhysData(
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
) {
	double_t h=time[1]-time[0];

	for (int i=1; i<timeLen; i++) {
		double_t iterH=time[i]-time[i-1];
		if (iterH<0) {
			return TimeSeriesNotIncreasingErr.v;
		}
		if (std::fabs(iterH-h) > pdOpts->TimeDeltaEps) {
			return TimeSeriesNotMonotonicErr.v;
		}
	}

	// For an explanation of the formulas refer to here:
	// http://code.barbellmath.net/barbell-math/providentia/wiki/Numerical-Difference-Methods
	if (bpOpts->ApproxErr.v==SecondOrder.v) {
		for (int i=2; i<timeLen-2; i++) {
			vel[i].X=(-pos[i-1].X+pos[i+1].X)/(2*h);
			vel[i].Y=(-pos[i-1].Y+pos[i+1].Y)/(2*h);

			acc[i].X=(pos[i-1].X-2*pos[i].X+pos[i+1].X)/(h*h);
			acc[i].Y=(pos[i-1].Y-2*pos[i].Y+pos[i+1].Y)/(h*h);

			jerk[i].X=(
				-pos[i-2].X+2*pos[i-1].X-2*pos[i+1].X+pos[i+2].X
			)/(
				2*h*h*h
			);
			jerk[i].Y=(
				-pos[i-2].Y+2*pos[i-1].Y-2*pos[i+1].Y+pos[i+2].Y
			)/(
				2*h*h*h
			);
		}

		// Smear edges to the ends of the results rather than computing forward and
		// backward difference formulas. Running those calculations would provide
		// little benefit while significantly increasing complexity and maintenance
		for (int i=0; i<2 && i<timeLen; i++) {
			vel[i]=vel[2];
			acc[i]=acc[2];
			jerk[i]=jerk[2];
		}
		for (int i=timeLen-2; i<timeLen; i++) {
			vel[i]=vel[timeLen-3];
			acc[i]=acc[timeLen-3];
			jerk[i]=jerk[timeLen-3];
		}
	} else if (bpOpts->ApproxErr.v==FourthOrder.v) {
		for (int i=3; i<timeLen-3; i++) {
			vel[i].X=(pos[i-2].X-8*pos[i-1].X+8*pos[i+1].X-pos[i+2].X)/(12*h);
			vel[i].Y=(pos[i-2].Y-8*pos[i-1].Y+8*pos[i+1].Y-pos[i+2].Y)/(12*h);

			acc[i].X=(
				-pos[i-2].X+16*pos[i-1].X-30*pos[i].X+16*pos[i+1].X-pos[i+2].X
			)/(
				12*h*h
			);
			acc[i].Y=(
				-pos[i-2].Y+16*pos[i-1].Y-30*pos[i].Y+16*pos[i+1].Y-pos[i+2].Y
			)/(
				12*h*h
			);

			jerk[i].X=(
				pos[i-3].X-8*pos[i-2].X+13*pos[i-1].X
				-13*pos[i+1].X+8*pos[i+2].X-pos[i+3].X
			)/(
				8*h*h*h
			);
			jerk[i].Y=(
				pos[i-3].Y-8*pos[i-2].Y+13*pos[i-1].Y
				-13*pos[i+1].Y+8*pos[i+2].Y-pos[i+3].Y
			)/(
				8*h*h*h
			);
		}

		// Smear edges to the ends of the results rather than computing forward and
		// backward difference formulas. Running those calculations would provide
		// little benefit while significantly increasing complexity and maintenance
		for (int i=0; i<3 && i<timeLen; i++) {
			vel[i]=vel[3];
			acc[i]=acc[3];
			jerk[i]=jerk[3];
		}
		for (int i=timeLen-3; i<timeLen; i++) {
			vel[i]=vel[timeLen-4];
			acc[i]=acc[timeLen-4];
			jerk[i]=jerk[timeLen-4];
		}
	} else {
		return InvalidApproximationErrErr.v;
	}

	// TODO - calc work, impulse, force

	return 0;
}
