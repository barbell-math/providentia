#include <cmath>
#include <vector>
#include <iostream>
#include <algorithm>
#include <bits/stdc++.h>
#include "cpu.h"
#include "../../clib/glue.h"

namespace BarPathCalc {

template <typename T>
struct PointInTime {
	double Time;
	T Value;
};

struct AbsVec2YOps : Math::Vec2 {
	friend bool operator>(const AbsVec2YOps l, const AbsVec2YOps r) {
		return fabs(l.Y) > fabs(r.Y);
	}
	friend bool operator<(const AbsVec2YOps l, const AbsVec2YOps r) {
		return fabs(l.Y) < fabs(r.Y);
	}
	friend bool operator>=(const AbsVec2YOps l, const AbsVec2YOps r) {
		return fabs(l.Y) >= fabs(r.Y);
	}
	friend bool operator<=(const AbsVec2YOps l, const AbsVec2YOps r) {
		return fabs(l.Y) <= fabs(r.Y);
	}
	friend bool operator!=(const AbsVec2YOps l, const AbsVec2YOps r) {
		return fabs(l.Y) != fabs(r.Y);
	}
	friend bool operator==(const AbsVec2YOps l, const AbsVec2YOps r) {
		return fabs(l.Y) == fabs(r.Y);
	}
};

struct Vec2MagOps : Math::Vec2 {
	friend bool operator>(const Vec2MagOps l, const Vec2MagOps r) {
		return Math::Mag(l) > Math::Mag(r);
	}
	friend bool operator<(const Vec2MagOps l, const Vec2MagOps r) {
		return Math::Mag(l) < Math::Mag(r);
	}
	friend bool operator>=(const Vec2MagOps l, const Vec2MagOps r) {
		return Math::Mag(l) >= Math::Mag(r);
	}
	friend bool operator<=(const Vec2MagOps l, const Vec2MagOps r) {
		return Math::Mag(l) <= Math::Mag(r);
	}
	friend bool operator!=(const Vec2MagOps l, const Vec2MagOps r) {
		return Math::Mag(l) != Math::Mag(r);
	}
	friend bool operator==(const Vec2MagOps l, const Vec2MagOps r) {
		return Math::Mag(l) == Math::Mag(r);
	}
};

enum BarPathCalcErrCode_t validateSuppliedData(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
) {
	double_t h=data->time[1]-data->time[0];

	for (int i=1; i<data->timeLen; i++) {
		double_t iterH=data->time[i]-data->time[i-1];
		if (iterH<0) {
			return TimeSeriesNotIncreasingErr;
		}
		if (fabs(iterH-h) > opts->TimeDeltaEps) {
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

	switch (opts->ApproxErr) {
	case SecondOrder:
		Math::CalcFirstThreeDerivatives<Math::SecondOrderApprox>(
			Slice<Math::Vec2>((Math::Vec2*)data->pos, data->timeLen),
			Slice<Math::Vec2>((Math::Vec2*)data->vel, data->timeLen),
			Slice<Math::Vec2>((Math::Vec2*)data->acc, data->timeLen),
			Slice<Math::Vec2>((Math::Vec2*)data->jerk, data->timeLen),
			h
		);
		break;
	case FourthOrder:
		Math::CalcFirstThreeDerivatives<Math::FourthOrderApprox>(
			Slice<Math::Vec2>((Math::Vec2*)data->pos, data->timeLen),
			Slice<Math::Vec2>((Math::Vec2*)data->vel, data->timeLen),
			Slice<Math::Vec2>((Math::Vec2*)data->acc, data->timeLen),
			Slice<Math::Vec2>((Math::Vec2*)data->jerk, data->timeLen),
			h
		);
		break;
	default:
		return InvalidApproximationErrErr;
	}

	return NoErr;
}

enum BarPathCalcErrCode_t runSmoother(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
) {
	Math::Vec2 _tmps[3]={};
	double _weights[5]={
		opts->SmootherWeight1,
		opts->SmootherWeight2,
		opts->SmootherWeight3,
		opts->SmootherWeight4,
		opts->SmootherWeight5,
	};
	Math::CenteredRollingWeightedAvg(
		Slice<Math::Vec2>((Math::Vec2*)data->vel, data->timeLen),
		FixedSlice<double, 5>(_weights),
		FixedRing<Math::Vec2, 3>(_tmps)
	);
	Math::CenteredRollingWeightedAvg(
		Slice<Math::Vec2>((Math::Vec2*)data->acc, data->timeLen),
		FixedSlice<double, 5>(_weights),
		FixedRing<Math::Vec2, 3>(_tmps)
	);
	Math::CenteredRollingWeightedAvg(
		Slice<Math::Vec2>((Math::Vec2*)data->jerk, data->timeLen),
		FixedSlice<double, 5>(_weights),
		FixedRing<Math::Vec2, 3>(_tmps)
	);

	return NoErr;
}

enum BarPathCalcErrCode_t calcHigherOrderData(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
) {
	Slice<Math::Vec2> vel((Math::Vec2*)data->vel, data->timeLen);
	Slice<Math::Vec2> acc((Math::Vec2*)data->acc, data->timeLen);
	Slice<Math::Vec2> force((Math::Vec2*)data->force, data->timeLen);
	Slice<Math::Vec2> impulse((Math::Vec2*)data->impulse, data->timeLen);

	for (int i=0; i<data->timeLen; i++) {
		// For an explanation of some of these formulas:
		// http://code.barbellmath.net/barbell-math/providentia/wiki/Bar-Path-Calcs
		force[i]=acc[i]*data->mass;
		impulse[i]=vel[i]*data->mass;
		data->power[i]=vel[i].Dot(force[i]);
		data->work[i]=(data->mass/2)*(vel[i].Dot(vel[i]));
	}

	return NoErr;
}

// TODO - see if the near zero filter can be removed.
// How to find "resting position" of bar?
enum BarPathCalcErrCode_t calcRepSplits(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
) {
	Slice<size_t> repCenters(data->reps);
	size_t numMaxes = Math::NLargestMaximums(
		Slice<AbsVec2YOps>((AbsVec2YOps*)data->pos, data->timeLen),
		repCenters,
		(AbsVec2YOps)Math::Vec2{ .X=0, .Y=0, },	// zeros because abs
		opts->NoiseFilter
	);
	std::sort(repCenters.begin(), repCenters.end());

	Slice<Math::Vec2> vel((Math::Vec2*)data->vel, data->timeLen);
	Slice<Math::Vec2> pos((Math::Vec2*)data->pos, data->timeLen);
	for (size_t i=0; i<numMaxes; i++) {
		size_t repCenter=repCenters[i];
		data->repSplit[i].StartIdx = 0;
		data->repSplit[i].EndIdx = data->timeLen;

		for (int j=repCenter+2; j<data->timeLen; j++) {
			if (
				std::signbit(vel[j].Y)!=std::signbit(vel[j-1].Y) &&
				fabs(pos[i].Y)<opts->NearZeroFilter
			) {
				data->repSplit[i].EndIdx=j;
				break;
			}
		}

		for (int j=repCenter-2; j>=0; j--) {
			if (
				std::signbit(vel[j].Y)!=std::signbit(vel[j+1].Y) &&
				fabs(pos[i].Y)<opts->NearZeroFilter
			) {
				data->repSplit[i].StartIdx=j;
				break;
			}
		}
	}

	return NoErr;
}

template <typename T, typename U>
void setRepMinMaxVal(
	Slice<double> time,
	Slice<T> data,
	split_t repSplit,
	PointInTime<U> *minPointInTime,
	PointInTime<U> *maxPointInTime,
	U (*valueTransform)(const T& t)=[](const T& t){ return t; }
) {
	Slice<T> subSlice(data(repSplit.StartIdx, repSplit.EndIdx));
	typename Slice<T>::Iterator min=std::min_element(subSlice.begin(), subSlice.end());
	typename Slice<T>::Iterator max=std::max_element(subSlice.begin(), subSlice.end());
	size_t minIdx=(min-subSlice.begin())+(size_t)repSplit.StartIdx;
	size_t maxIdx=(max-subSlice.begin())+(size_t)repSplit.StartIdx;

	minPointInTime->Time=time[minIdx];
	minPointInTime->Value=valueTransform(*min);
	maxPointInTime->Time=time[maxIdx];
	maxPointInTime->Value=valueTransform(*max);
}

template <typename T, typename U>
void setRepAvgVal(
	Slice<T> data,
	split_t repSplit,
	U *avgVal,
	U (*valueTransform)(const T& t)=[](const T& t){ return t; }
) {
	*avgVal=U{};
	Slice<T> subSlice(data(repSplit.StartIdx, repSplit.EndIdx));
	if (subSlice.Len()==0) { return; }
	U tot=std::accumulate(subSlice.begin(), subSlice.end(), *avgVal);
	*avgVal=valueTransform(tot)/subSlice.Len();
}

enum BarPathCalcErrCode_t calcRepStats(
	barPathData_t* data,
	barPathCalcHyperparams_t* opts
) {
	Slice<double> time((double*)data->time, data->timeLen);
	Slice<Vec2MagOps> vel((Vec2MagOps*)data->vel, data->timeLen);
	Slice<Vec2MagOps> acc((Vec2MagOps*)data->acc, data->timeLen);
	Slice<Vec2MagOps> force((Vec2MagOps*)data->force, data->timeLen);
	Slice<Vec2MagOps> impulse((Vec2MagOps*)data->impulse, data->timeLen);
	Slice<double> power((double*)data->power, data->timeLen);
	Slice<double> work((double*)data->work, data->timeLen);
	auto transform=[](const Vec2MagOps& t){ return Math::Mag((Math::Vec2)t); };

	for (size_t i=0; i<(size_t)data->reps; i++) {
		split_t repSplit=data->repSplit[i];
		if (repSplit.EndIdx-repSplit.StartIdx==0) {
			continue;
		}

		setRepMinMaxVal<Vec2MagOps, double>(
			time, vel, repSplit,
			(PointInTime<double>*)&data->minVel[i],
			(PointInTime<double>*)&data->maxVel[i],
			transform
		);
		setRepMinMaxVal<Vec2MagOps, double>(
			time, acc, repSplit,
			(PointInTime<double>*)&data->minAcc[i],
			(PointInTime<double>*)&data->maxAcc[i],
			transform
		);
		setRepMinMaxVal<Vec2MagOps, double>(
			time, force, repSplit,
			(PointInTime<double>*)&data->minForce[i],
			(PointInTime<double>*)&data->maxForce[i],
			transform
		);
		setRepMinMaxVal<Vec2MagOps, double>(
			time, impulse, repSplit,
			(PointInTime<double>*)&data->minImpulse[i],
			(PointInTime<double>*)&data->maxImpulse[i],
			transform
		);
		setRepMinMaxVal<double, double>(
			time, power, repSplit,
			(PointInTime<double>*)&data->minPower[i],
			(PointInTime<double>*)&data->maxPower[i]
		);
		setRepMinMaxVal<double, double>(
			time, work, repSplit,
			(PointInTime<double>*)&data->minWork[i],
			(PointInTime<double>*)&data->maxWork[i]
		);

		setRepAvgVal(power, repSplit, data->avgPower);
		setRepAvgVal(work, repSplit, data->avgWork);
	}
	return NoErr;
}


extern "C" enum BarPathCalcErrCode_t CalcBarPathPhysData(
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

};
