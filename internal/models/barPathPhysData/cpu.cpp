#include <cmath>
#include "cpu.h"
#include "../../glue/glue.h"

extern "C" int64_t calcBarPathPhysData(
	barPathData_t* data,
	barPathCalcConf_t* bpOpts,
	physDataConf_t* pdOpts
) {
	double_t h=data->time[1]-data->time[0];

	for (int i=1; i<data->timeLen; i++) {
		double_t iterH=data->time[i]-data->time[i-1];
		if (iterH<0) {
			return TimeSeriesNotIncreasingErr;
		}
		if (std::fabs(iterH-h) > pdOpts->TimeDeltaEps) {
			return TimeSeriesNotMonotonicErr;
		}
	}

	// For an explanation of the formulas refer to here:
	// http://code.barbellmath.net/barbell-math/providentia/wiki/Numerical-Difference-Methods
	switch (bpOpts->ApproxErr) {
	case SecondOrder:
		for (int i=2; i<data->timeLen-2; i++) {
			data->vel[i].X=(-data->pos[i-1].X+data->pos[i+1].X)/(2*h);
			data->vel[i].Y=(-data->pos[i-1].Y+data->pos[i+1].Y)/(2*h);

			data->acc[i].X=(
				data->pos[i-1].X-2*data->pos[i].X+data->pos[i+1].X
			)/(h*h);
			data->acc[i].Y=(
				data->pos[i-1].Y-2*data->pos[i].Y+data->pos[i+1].Y
			)/(h*h);

			data->jerk[i].X=(
				-data->pos[i-2].X
				+2*data->pos[i-1].X
				-2*data->pos[i+1].X
				+data->pos[i+2].X
			)/(
				2*h*h*h
			);
			data->jerk[i].Y=(
				-data->pos[i-2].Y
				+2*data->pos[i-1].Y
				-2*data->pos[i+1].Y
				+data->pos[i+2].Y
			)/(
				2*h*h*h
			);
		}

		// Smear edges to the ends of the results rather than computing forward and
		// backward difference formulas. Running those calculations would provide
		// little benefit while significantly increasing complexity and maintenance
		for (int i=0; i<2 && i<data->timeLen; i++) {
			data->vel[i]=data->vel[2];
			data->acc[i]=data->acc[2];
			data->jerk[i]=data->jerk[2];
		}
		for (int i=data->timeLen-2; i<data->timeLen; i++) {
			data->vel[i]=data->vel[data->timeLen-3];
			data->acc[i]=data->acc[data->timeLen-3];
			data->jerk[i]=data->jerk[data->timeLen-3];
		}

		break;
	case FourthOrder:
		for (int i=3; i<data->timeLen-3; i++) {
			data->vel[i].X=(
				data->pos[i-2].X
				-8*data->pos[i-1].X
				+8*data->pos[i+1].X
				-data->pos[i+2].X
			)/(12*h);
			data->vel[i].Y=(
				data->pos[i-2].Y
				-8*data->pos[i-1].Y
				+8*data->pos[i+1].Y
				-data->pos[i+2].Y
			)/(12*h);

			data->acc[i].X=(
				-data->pos[i-2].X
				+16*data->pos[i-1].X
				-30*data->pos[i].X
				+16*data->pos[i+1].X
				-data->pos[i+2].X
			)/(
				12*h*h
			);
			data->acc[i].Y=(
				-data->pos[i-2].Y
				+16*data->pos[i-1].Y
				-30*data->pos[i].Y
				+16*data->pos[i+1].Y
				-data->pos[i+2].Y
			)/(
				12*h*h
			);

			data->jerk[i].X=(
				data->pos[i-3].X
				-8*data->pos[i-2].X
				+13*data->pos[i-1].X
				-13*data->pos[i+1].X
				+8*data->pos[i+2].X
				-data->pos[i+3].X
			)/(
				8*h*h*h
			);
			data->jerk[i].Y=(
				data->pos[i-3].Y
				-8*data->pos[i-2].Y
				+13*data->pos[i-1].Y
				-13*data->pos[i+1].Y
				+8*data->pos[i+2].Y
				-data->pos[i+3].Y
			)/(
				8*h*h*h
			);
		}

		// Smear edges to the ends of the results rather than computing forward and
		// backward difference formulas. Running those calculations would provide
		// little benefit while significantly increasing complexity and maintenance
		for (int i=0; i<3 && i<data->timeLen; i++) {
			data->vel[i]=data->vel[3];
			data->acc[i]=data->acc[3];
			data->jerk[i]=data->jerk[3];
		}
		for (int i=data->timeLen-3; i<data->timeLen; i++) {
			data->vel[i]=data->vel[data->timeLen-4];
			data->acc[i]=data->acc[data->timeLen-4];
			data->jerk[i]=data->jerk[data->timeLen-4];
		}

		break;
	default:
		return InvalidApproximationErrErr;
	}

	float_t wTot=(
		bpOpts->SmootherWeight1+
		bpOpts->SmootherWeight2+
		bpOpts->SmootherWeight3+
		bpOpts->SmootherWeight4+
		bpOpts->SmootherWeight5
	);
	for (int i=2; wTot>0 && i<data->timeLen-2; i++) {
		data->vel[i].X=(
			data->vel[i-2].X*bpOpts->SmootherWeight1+
			data->vel[i-1].X*bpOpts->SmootherWeight2+
			data->vel[i].X*bpOpts->SmootherWeight3+
			data->vel[i+1].X*bpOpts->SmootherWeight4+
			data->vel[i+2].X*bpOpts->SmootherWeight5
		)/(wTot);
		data->vel[i].Y=(
			data->vel[i-2].Y*bpOpts->SmootherWeight1+
			data->vel[i-1].Y*bpOpts->SmootherWeight2+
			data->vel[i].Y*bpOpts->SmootherWeight3+
			data->vel[i+1].Y*bpOpts->SmootherWeight4+
			data->vel[i+2].Y*bpOpts->SmootherWeight5
		)/(wTot);

		data->acc[i].X=(
			data->acc[i-2].X*bpOpts->SmootherWeight1+
			data->acc[i-1].X*bpOpts->SmootherWeight2+
			data->acc[i].X*bpOpts->SmootherWeight3+
			data->acc[i+1].X*bpOpts->SmootherWeight4+
			data->acc[i+2].X*bpOpts->SmootherWeight5
		)/(wTot);
		data->acc[i].Y=(
			data->acc[i-2].Y*bpOpts->SmootherWeight1+
			data->acc[i-1].Y*bpOpts->SmootherWeight2+
			data->acc[i].Y*bpOpts->SmootherWeight3+
			data->acc[i+1].Y*bpOpts->SmootherWeight4+
			data->acc[i+2].Y*bpOpts->SmootherWeight5
		)/(wTot);

		data->jerk[i].X=(
			data->jerk[i-2].X*bpOpts->SmootherWeight1+
			data->jerk[i-1].X*bpOpts->SmootherWeight2+
			data->jerk[i].X*bpOpts->SmootherWeight3+
			data->jerk[i+1].X*bpOpts->SmootherWeight4+
			data->jerk[i+2].X*bpOpts->SmootherWeight5
		)/(wTot);
		data->jerk[i].Y=(
			data->jerk[i-2].Y*bpOpts->SmootherWeight1+
			data->jerk[i-1].Y*bpOpts->SmootherWeight2+
			data->jerk[i].Y*bpOpts->SmootherWeight3+
			data->jerk[i+1].Y*bpOpts->SmootherWeight4+
			data->jerk[i+2].Y*bpOpts->SmootherWeight5
		)/(wTot);
	}

	// For an explanation of some of these formulas:
	// http://code.barbellmath.net/barbell-math/providentia/wiki/Bar-Path-Calcs
	for (int i=0; i<data->timeLen; i++) {
		data->force[i].X=data->Mass*data->acc[i].X;
		data->force[i].Y=data->Mass*data->acc[i].Y;
		data->power[i]=data->force[i].X*data->vel[i].X + data->force[i].Y*data->vel[i].Y;
		data->impulse[i].X = data->Mass*data->vel[i].X;
		data->impulse[i].Y = data->Mass*data->vel[i].Y;
		data->work[i]=(data->Mass/2) * (data->vel[i].X*data->vel[i].X + data->vel[i].Y*data->vel[i].Y);
	}

	return NoErr;
}
