#include <cmath>
#include <vector>
#include <iostream>
#include <algorithm>
#include <bits/stdc++.h>
#include "cpu.h"
#include "../../clib/glue.h"
#include "../../clib/common.h"

enum BarPathCalcErrCode_t validateSuppliedData(
	barPathData_t* data,
	barPathCalcConf_t* opts
) {
	double_t h=data->time[1]-data->time[0];

	for (int i=1; i<data->timeLen; i++) {
		double_t iterH=data->time[i]-data->time[i-1];
		if (iterH<0) {
			return TimeSeriesNotIncreasingErr;
		}
		if (std::fabs(iterH-h) > opts->TimeDeltaEps) {
			return TimeSeriesNotMonotonicErr;
		}
	}

	return NoErr;
}

// For an explanation of the formulas refer to here:
// http://code.barbellmath.net/barbell-math/providentia/wiki/Numerical-Difference-Methods
enum BarPathCalcErrCode_t calcDerivatives(
	barPathData_t* data,
	barPathCalcConf_t* opts
) {
	double_t h=data->time[1]-data->time[0];

	switch (opts->ApproxErr) {
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

		return NoErr;
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

		return NoErr;
	default:
		return InvalidApproximationErrErr;
	}
}

enum BarPathCalcErrCode_t runSmoother(
	barPathData_t* data,
	barPathCalcConf_t* opts
) {
	float_t wTot=(
		opts->SmootherWeight1+
		opts->SmootherWeight2+
		opts->SmootherWeight3+
		opts->SmootherWeight4+
		opts->SmootherWeight5
	);
	for (int i=2; wTot>0 && i<data->timeLen-2; i++) {
		data->vel[i].X=(
			data->vel[i-2].X*opts->SmootherWeight1+
			data->vel[i-1].X*opts->SmootherWeight2+
			data->vel[i].X*opts->SmootherWeight3+
			data->vel[i+1].X*opts->SmootherWeight4+
			data->vel[i+2].X*opts->SmootherWeight5
		)/(wTot);
		data->vel[i].Y=(
			data->vel[i-2].Y*opts->SmootherWeight1+
			data->vel[i-1].Y*opts->SmootherWeight2+
			data->vel[i].Y*opts->SmootherWeight3+
			data->vel[i+1].Y*opts->SmootherWeight4+
			data->vel[i+2].Y*opts->SmootherWeight5
		)/(wTot);

		data->acc[i].X=(
			data->acc[i-2].X*opts->SmootherWeight1+
			data->acc[i-1].X*opts->SmootherWeight2+
			data->acc[i].X*opts->SmootherWeight3+
			data->acc[i+1].X*opts->SmootherWeight4+
			data->acc[i+2].X*opts->SmootherWeight5
		)/(wTot);
		data->acc[i].Y=(
			data->acc[i-2].Y*opts->SmootherWeight1+
			data->acc[i-1].Y*opts->SmootherWeight2+
			data->acc[i].Y*opts->SmootherWeight3+
			data->acc[i+1].Y*opts->SmootherWeight4+
			data->acc[i+2].Y*opts->SmootherWeight5
		)/(wTot);

		data->jerk[i].X=(
			data->jerk[i-2].X*opts->SmootherWeight1+
			data->jerk[i-1].X*opts->SmootherWeight2+
			data->jerk[i].X*opts->SmootherWeight3+
			data->jerk[i+1].X*opts->SmootherWeight4+
			data->jerk[i+2].X*opts->SmootherWeight5
		)/(wTot);
		data->jerk[i].Y=(
			data->jerk[i-2].Y*opts->SmootherWeight1+
			data->jerk[i-1].Y*opts->SmootherWeight2+
			data->jerk[i].Y*opts->SmootherWeight3+
			data->jerk[i+1].Y*opts->SmootherWeight4+
			data->jerk[i+2].Y*opts->SmootherWeight5
		)/(wTot);
	}

	return NoErr;
}


// For an explanation of some of these formulas:
// http://code.barbellmath.net/barbell-math/providentia/wiki/Bar-Path-Calcs
enum BarPathCalcErrCode_t calcHigherOrderData(
	barPathData_t* data,
	barPathCalcConf_t* opts
) {
	for (int i=0; i<data->timeLen; i++) {
		data->force[i].X=data->mass*data->acc[i].X;
		data->force[i].Y=data->mass*data->acc[i].Y;
		data->power[i]=(
			data->force[i].X*data->vel[i].X + data->force[i].Y*data->vel[i].Y
		);
		data->impulse[i].X = data->mass*data->vel[i].X;
		data->impulse[i].Y = data->mass*data->vel[i].Y;
		data->work[i]=(data->mass/2) * (
			data->vel[i].X*data->vel[i].X + data->vel[i].Y*data->vel[i].Y
		);
	}

	return NoErr;
}

enum BarPathCalcErrCode_t calcRepSplits(
	barPathData_t* data,
	barPathCalcConf_t* opts
) {
	int centersAdded=0;
	std::vector<TimestampedVal> repCenters(data->reps);
	for (int i=1; i<data->timeLen-1; i++) {
		if (std::signbit(data->vel[i].Y)!=std::signbit(data->vel[i-1].Y)) {
			TimestampedVal localMin=TimestampedVal{
				.idx = i,
				.time = data->time[i],
				.value=-std::fabs(data->pos[i].Y),
			};
			if (centersAdded<data->reps) {
				repCenters[centersAdded]=localMin;
				std::make_heap(
					repCenters.begin(), repCenters.end(),
					TimestampedVal::sortByValue
				);
				centersAdded++;
			} else if (localMin.value<repCenters[0].value) {
				std::pop_heap(
					repCenters.begin(), repCenters.end(),
					TimestampedVal::sortByValue
				);
				repCenters[repCenters.size()-1]=localMin;
				std::push_heap(
					repCenters.begin(), repCenters.end(),
					TimestampedVal::sortByValue
				);
			}
		}
	}
	std::sort(repCenters.begin(), repCenters.end(), TimestampedVal::sortByTime);

	for (int i=0; i<repCenters.size(); i++) {
		TimestampedVal& iterRep=repCenters[i];

		for (int j=iterRep.idx+2; j<data->timeLen; j++) {
			if (
				std::signbit(data->vel[j].Y)!=std::signbit(data->vel[j-1].Y) &&
				std::fabs(data->pos[i].Y)<opts->NearZeroFilter
			) {
				data->repSplit[i].EndIdx=j;
				break;
			}
		}

		for (int j=iterRep.idx-2; j>=0; j--) {
			if (
				std::signbit(data->vel[j].Y)!=std::signbit(data->vel[j+1].Y) &&
				std::fabs(data->pos[i].Y)<opts->NearZeroFilter
			) {
				data->repSplit[i].StartIdx=j;
				break;
			}
		}
	}

	// for (int i=0; i<data->reps; i++){
	// 	std::cout << "(" 
	// 		<< data->repSplit[i].StartIdx << "[" << data->time[data->repSplit[i].StartIdx] << "], " 
	// 		<< data->repSplit[i].EndIdx << "[" << data->time[data->repSplit[i].EndIdx] << "], " 
	// 	<< ") ";
	// }
	// std::cout << std::endl;

	return NoErr;
}

enum BarPathCalcErrCode_t calcRepStats(
	barPathData_t* data,
	barPathCalcConf_t* opts
) {
	for (int i=0; i<data->reps; i++) {
		data->minVel[i].Value=std::numeric_limits<double_t>::infinity();
		data->maxVel[i].Value=-std::numeric_limits<double_t>::infinity();
		data->minAcc[i].Value=std::numeric_limits<double_t>::infinity();
		data->maxAcc[i].Value=-std::numeric_limits<double_t>::infinity();
		data->minForce[i].Value=std::numeric_limits<double_t>::infinity();
		data->maxForce[i].Value=-std::numeric_limits<double_t>::infinity();
		data->minImpulse[i].Value=std::numeric_limits<double_t>::infinity();
		data->maxImpulse[i].Value=-std::numeric_limits<double_t>::infinity();
		data->avgWork[i]=0;
		data->minWork[i].Value=std::numeric_limits<double_t>::infinity();
		data->maxWork[i].Value=-std::numeric_limits<double_t>::infinity();
		data->avgPower[i]=-1;
		data->minPower[i].Value=std::numeric_limits<double_t>::infinity();
		data->maxPower[i].Value=-std::numeric_limits<double_t>::infinity();
		int avgCntr=0;

		for (
			int j=data->repSplit[i].StartIdx;
			j<data->repSplit[i].EndIdx && j<data->timeLen;
			j++
		) {
			float velMag=Vec2::mag(data->vel[j].X, data->vel[j].Y);
			if (data->minVel[i].Value>velMag) {
				data->minVel[i]=velPointInTime_t{
					.Time=data->time[j], .Value=velMag,
				};
			}
			if (data->maxVel[i].Value<velMag) {
				data->maxVel[i]=velPointInTime_t{
					.Time=data->time[j], .Value=velMag,
				};
			}
			
			float accMag=Vec2::mag(data->acc[j].X, data->acc[j].Y);
			if (data->minAcc[i].Value>accMag) {
				data->minAcc[i]=accPointInTime_t{
					.Time=data->time[j], .Value=accMag,
				};
			}
			if (data->maxAcc[i].Value<accMag) {
				data->maxAcc[i]=accPointInTime_t{
					.Time=data->time[j], .Value=accMag,
				};
			}
			
			float forceMag=Vec2::mag(data->force[j].X, data->force[j].Y);
			if (data->minForce[i].Value>forceMag) {
				data->minForce[i]=newtonPointInTime_t{
					.Time=data->time[j], .Value=forceMag,
				};
			}
			if (data->maxForce[i].Value<forceMag) {
				data->maxForce[i]=newtonPointInTime_t{
					.Time=data->time[j], .Value=forceMag,
				};
			}

			float impulseMag=Vec2::mag(data->impulse[j].X, data->impulse[j].Y);
			if (data->minImpulse[i].Value>impulseMag) {
				data->minImpulse[i]=newtonSecPointInTime_t{
					.Time=data->time[j], .Value=impulseMag,
				};
			}
			if (data->maxImpulse[i].Value<impulseMag) {
				data->maxImpulse[i]=newtonSecPointInTime_t{
					.Time=data->time[j], .Value=impulseMag,
				};
			}

			data->avgPower[i]+=data->power[j];
			if (data->minPower[i].Value>data->power[j]) {
				data->minPower[i]=wattPointInTime_t{
					.Time=data->time[j], .Value=data->power[j],
				};
			}
			if (data->maxPower[i].Value<data->power[j]) {
				data->maxPower[i]=wattPointInTime_t{
					.Time=data->time[j], .Value=data->power[j],
				};
			}

			data->avgWork[i]+=data->work[j];
			if (data->minWork[i].Value>data->work[j]) {
				data->minWork[i]=joulePointInTime_t{
					.Time=data->time[j], .Value=data->work[j],
				};
			}
			if (data->maxWork[i].Value<data->work[j]) {
				data->maxWork[i]=joulePointInTime_t{
					.Time=data->time[j], .Value=data->work[j],
				};
			}

			avgCntr++;
		}

		if (avgCntr>0) {
			data->avgPower[i]/=avgCntr;
			data->avgWork[i]/=avgCntr;
		} else {
			data->avgPower[i]=0;
			data->avgWork[i]=0;
		}
	}
	return NoErr;
}

extern "C" enum BarPathCalcErrCode_t calcBarPathPhysData(
	barPathData_t* data,
	barPathCalcConf_t* opts
) {
	BarPathCalcErrCode_t err = validateSuppliedData(data, opts);
	if (err!=NoErr) {
		return  err;
	}

	err = calcDerivatives(data, opts);
	if (err!=NoErr) {
		return  err;
	}

	err = runSmoother(data, opts);
	if (err!=NoErr) {
		return  err;
	}

	err = calcHigherOrderData(data, opts);
	if (err!=NoErr) {
		return  err;
	}

	err = calcRepSplits(data, opts);
	if (err!=NoErr) {
		return  err;
	}

	err = calcRepStats(data, opts);
	if (err!=NoErr) {
		return err;
	}

	return NoErr;
}
