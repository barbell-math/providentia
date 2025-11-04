#include <cmath>
#include <vector>
#include <iostream>
#include <algorithm>
#include <bits/stdc++.h>
#include "cpu.h"
#include "../../clib/glue.h"
#include "../../clib/common.h"
#include "../../clib/slice.h"

enum BarPathCalcErrCode_t validateSuppliedData(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
) {
	double_t h=data->time[1]-data->time[0];

	for (int i=1; i<data->timeLen; i++) {
		double_t iterH=data->time[i]-data->time[i-1];
		if (iterH<0) {
			std::cout << i << ": " << data->time[i] << ", " << data->time[i-1] << std::endl;
			return TimeSeriesNotIncreasingErr;
		}
		if (fabs(iterH-h) > opts->TimeDeltaEps) {
			std::cout << fabs(iterH-h) << std::endl;
			return TimeSeriesNotMonotonicErr;
		}
	}

	return NoErr;
}

enum BarPathCalcErrCode_t calcDerivatives(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
) {
	double_t h=data->time[1]-data->time[0];
	int smearLen = 0;

	switch (opts->ApproxErr) {
	case SecondOrder:
		smearLen = 2;
		Math::calcFirstThreeDerivatives<5>(
			Slice<Vec2>(data->pos, data->timeLen),
			Slice<Vec2>(data->vel, data->timeLen),
			Slice<Vec2>(data->acc, data->timeLen),
			Slice<Vec2>(data->jerk, data->timeLen),
			h
		);
		break;
	case FourthOrder:
		smearLen = 3;
		Math::calcFirstThreeDerivatives<7>(
			Slice<Vec2>(data->pos, data->timeLen),
			Slice<Vec2>(data->vel, data->timeLen),
			Slice<Vec2>(data->acc, data->timeLen),
			Slice<Vec2>(data->jerk, data->timeLen),
			h
		);
		break;
	default:
		return InvalidApproximationErrErr;
	}

	// Smear edges to the ends of the results rather than computing forward and
	// backward difference formulas. Running those calculations would provide
	// little benefit while significantly increasing complexity and maintenance
	for (int i=0; i<smearLen && i<data->timeLen; i++) {
		data->vel[i]=data->vel[smearLen];
		data->acc[i]=data->acc[smearLen];
		data->jerk[i]=data->jerk[smearLen];
	}
	for (int i=data->timeLen-smearLen; i<data->timeLen; i++) {
		data->vel[i]=data->vel[data->timeLen-smearLen-1];
		data->acc[i]=data->acc[data->timeLen-smearLen-1];
		data->jerk[i]=data->jerk[data->timeLen-smearLen-1];
	}

	return NoErr;
}

void smootherFunc(
	Array<Vec2, 5> window,
	double_t wTot,
	barPathCalcHyperparams_t* opts
) {
	window[2].X=(
		window[0].X*opts->SmootherWeight1+
		window[1].X*opts->SmootherWeight2+
		window[2].X*opts->SmootherWeight3+
		window[3].X*opts->SmootherWeight4+
		window[4].X*opts->SmootherWeight5
	)/(wTot);
	window[2].Y=(
		window[0].Y*opts->SmootherWeight1+
		window[1].Y*opts->SmootherWeight2+
		window[2].Y*opts->SmootherWeight3+
		window[3].Y*opts->SmootherWeight4+
		window[4].Y*opts->SmootherWeight5
	)/(wTot);
}

enum BarPathCalcErrCode_t runSmoother(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
) {
	float_t wTot=(
		opts->SmootherWeight1+
		opts->SmootherWeight2+
		opts->SmootherWeight3+
		opts->SmootherWeight4+
		opts->SmootherWeight5
	);
	Slice<Vec2> vel = Slice<Vec2>(data->vel, data->timeLen);
	Slice<Vec2> acc = Slice<Vec2>(data->acc, data->timeLen);
	Slice<Vec2> jerk = Slice<Vec2>(data->jerk, data->timeLen);
	for (int i=0; wTot>0 && i<data->timeLen-4; i++) {
		smootherFunc(Array<Vec2, 5>(vel,i), wTot, opts);
		smootherFunc(Array<Vec2, 5>(acc,i), wTot, opts);
		smootherFunc(Array<Vec2, 5>(jerk,i), wTot, opts);
	}

	return NoErr;
}

// For an explanation of some of these formulas:
// http://code.barbellmath.net/barbell-math/providentia/wiki/Bar-Path-Calcs
enum BarPathCalcErrCode_t calcHigherOrderData(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
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
	barPathCalcHyperparams_t* opts
) {
	int centersAdded=0;
	std::vector<TimestampedVal> repCenters(data->reps);
	for (int i=1; i<data->timeLen-1; i++) {
		if (std::signbit(data->vel[i].Y)!=std::signbit(data->vel[i-1].Y)) {
			TimestampedVal localMin=TimestampedVal{
				.Idx = i,
				.Time = data->time[i],
				.Value=-fabs(data->pos[i].Y),
			};
			if (centersAdded<data->reps) {
				repCenters[centersAdded]=localMin;
				std::make_heap(
					repCenters.begin(), repCenters.end(),
					TimestampedVal::sortByValue
				);
				centersAdded++;
			} else if (localMin.Value<repCenters[0].Value) {
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

	for (size_t i=0; i<repCenters.size(); i++) {
		TimestampedVal& iterRep=repCenters[i];

		for (int j=iterRep.Idx+2; j<data->timeLen; j++) {
			if (
				std::signbit(data->vel[j].Y)!=std::signbit(data->vel[j-1].Y) &&
				fabs(data->pos[i].Y)<opts->NearZeroFilter
			) {
				data->repSplit[i].EndIdx=j;
				break;
			}
		}

		for (int j=iterRep.Idx-2; j>=0; j--) {
			if (
				std::signbit(data->vel[j].Y)!=std::signbit(data->vel[j+1].Y) &&
				fabs(data->pos[i].Y)<opts->NearZeroFilter
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

void setPointInTimeMinMax(
	PointInTime *curMin,
	PointInTime *curMax,
	PointInTime val
) {
	if (curMin->Value>val.Value) {
		*curMin=val;
	}
	if (curMax->Value<val.Value) {
		*curMax=val;
	}
}

enum BarPathCalcErrCode_t calcRepStats(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
) {
	for (int i=0; i<data->reps; i++) {
		data->minVel[i].Value=INFINITY;
		data->maxVel[i].Value=-INFINITY;
		data->minAcc[i].Value=INFINITY;
		data->maxAcc[i].Value=-INFINITY;
		data->minForce[i].Value=INFINITY;
		data->maxForce[i].Value=-INFINITY;
		data->minImpulse[i].Value=INFINITY;
		data->maxImpulse[i].Value=-INFINITY;
		data->avgWork[i]=0;
		data->minWork[i].Value=INFINITY;
		data->maxWork[i].Value=-INFINITY;
		data->avgPower[i]=-1;	// todo - wtf???
		data->minPower[i].Value=INFINITY;
		data->maxPower[i].Value=-INFINITY;
		int avgCntr=0;

		for (
			int j=data->repSplit[i].StartIdx;
			j<data->repSplit[i].EndIdx && j<data->timeLen;
			j++
		) {
			setPointInTimeMinMax(
				(PointInTime*)(&data->minVel[i]),
				(PointInTime*)(&data->maxVel[i]),
				PointInTime{
					.Time=data->time[j], .Value=Math::mag(*(Vec2*)(&data->vel[j]))
				}
			);
			setPointInTimeMinMax(
				(PointInTime*)(&data->minAcc[i]),
				(PointInTime*)(&data->maxAcc[i]),
				PointInTime{
					.Time=data->time[j], .Value=Math::mag(*(Vec2*)(&data->acc[j]))
				}
			);
			setPointInTimeMinMax(
				(PointInTime*)(&data->minForce[i]),
				(PointInTime*)(&data->maxForce[i]),
				PointInTime{
					.Time=data->time[j], .Value=Math::mag(*(Vec2*)(&data->force[j]))
				}
			);
			setPointInTimeMinMax(
				(PointInTime*)(&data->minImpulse[i]),
				(PointInTime*)(&data->maxImpulse[i]),
				PointInTime{
					.Time=data->time[j], .Value=Math::mag(*(Vec2*)(&data->impulse[j]))
				}
			);
			setPointInTimeMinMax(
				(PointInTime*)(&data->minPower[i]),
				(PointInTime*)(&data->maxPower[i]),
				PointInTime{
					.Time=data->time[j], .Value=data->power[j],
				}
			);
			setPointInTimeMinMax(
				(PointInTime*)(&data->minWork[i]),
				(PointInTime*)(&data->maxWork[i]),
				PointInTime{
					.Time=data->time[j], .Value=data->work[j],
				}
			);
			

			data->avgPower[i]+=data->power[j];
			data->avgWork[i]+=data->work[j];
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
	barPathCalcHyperparams_t* opts
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
