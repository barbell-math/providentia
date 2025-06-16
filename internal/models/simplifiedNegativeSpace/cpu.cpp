#include "cpu.h"
#include <cmath>
#include <cstdint>
#include <limits>
#include <iostream>
#include "Eigen/Dense"
#include "cgoStructs.h"

typedef struct modelData {
	Eigen::Array<double_t, 7, 1> curCoef;
	Eigen::Matrix<double_t, 7, 1> curResponses;
	Eigen::Matrix<double_t, 7, 7> curDesignMatrix;
	trainingLog_t* curDataPoint;
	trainingLog_t* iterDataPoint;
	modelState_t* curModelState;
	double_t curSmallestDelta;
} modelData_t;

void addDataPoint(modelData_t* md, opts_t* opts);

extern "C" void calcModelStates(
	int64_t clientID,
	int32_t modelID,
	trainingLog_t* historicalData,
	int64_t historicalDataLen,
	trainingLog_t* needsCalc,
	int64_t needsCalcLen,
	modelState_t* outValues,
	int64_t outValuesLen,
	opts_t* opts
) {
	for (int64_t i=0; i<needsCalcLen; i++) {
		outValues[i]=modelState_t{
			.ClientID = clientID,
			.TrainingLogID=needsCalc[i].ID,
			.ModelID =  modelID,
			.V1=0,
			.V2=0,
			.V3=0,
			.V4=0,
			.V5=0,
			.V6=0,
			.V7=0,
			.V8=opts->Alpha,
			.V9=opts->Beta,
			.V10=opts->Gamma,
			.TimeFrame=0,
			.Mse=std::numeric_limits<double_t>::infinity(),
			.PredWeight=0,
		};

		modelData_t md=modelData_t{
			.curCoef=Eigen::Matrix<double_t, 7, 1>::Constant(0),
			.curResponses=Eigen::Matrix<double_t, 7, 1>::Constant(0),
			.curDesignMatrix=Eigen::Matrix<double, 7, 7>::Constant(0),
			.curDataPoint=&needsCalc[i],
			.iterDataPoint=nullptr,
			.curModelState=&outValues[i],
			.curSmallestDelta=std::numeric_limits<double_t>::infinity(),
		};

		uint64_t iterCntr=0;
		for (int64_t j=i-1; j>=0 && iterCntr<opts->MaxIters; j--) {
			if (needsCalc[i].ExerciseID!=needsCalc[j].ExerciseID) {
				continue;
			}
			iterCntr++;
			md.iterDataPoint=&needsCalc[j];
			addDataPoint(&md, opts);
		}
		for (int64_t j=historicalDataLen-1; j>=0 && iterCntr<opts->MaxIters; j--) {
			if (needsCalc[i].ExerciseID!=needsCalc[j].ExerciseID) {
				continue;
			}
			iterCntr++;
			md.iterDataPoint=&needsCalc[j];
			addDataPoint(&md, opts);
		}
	}
}

void addDataPoint(modelData_t* md, opts_t* opts) {
	double_t reps=(double_t)md->iterDataPoint->Reps;
	double_t sets=(double_t)md->iterDataPoint->Sets;
	double_t interSessionCntr=(double_t)md->iterDataPoint->InterSessionCntr;
	double_t interWorkoutCntr=(double_t)md->iterDataPoint->InterWorkoutCntr;

	double_t effortTerm=pow(md->iterDataPoint->Effort, opts->Alpha);
	double_t repsTerm=pow(reps-1, opts->Beta);
	double_t setsTerm=pow(sets-1, opts->Gamma);
	double_t repsAndEffortTerm=repsTerm*effortTerm;
	double_t setsAndEffortTerm=setsTerm*effortTerm;
	double_t setsAndRepsTerm=setsTerm*repsTerm;
	double_t setsAndRepsAndEffortTerm=effortTerm*setsAndRepsTerm;

	Eigen::Array<double_t, 7, 1> baseTerms {
		1,
		interSessionCntr,
		interWorkoutCntr,
		effortTerm,
		setsAndRepsAndEffortTerm,
		repsAndEffortTerm,
		setsAndEffortTerm,
	};

	Eigen::Array<double_t, 7, 7> designMatrixAdditions=Eigen::Array<double_t, 7, 7>::Constant(0);
	designMatrixAdditions.row(0)=baseTerms;
	designMatrixAdditions.row(1)=baseTerms*interSessionCntr;
	designMatrixAdditions.row(2)=baseTerms*interWorkoutCntr;
	designMatrixAdditions.row(3)=baseTerms*effortTerm;
	designMatrixAdditions.row(4)=baseTerms*setsAndRepsAndEffortTerm;
	designMatrixAdditions.row(5)=baseTerms*repsAndEffortTerm;
	designMatrixAdditions.row(6)=baseTerms*setsAndEffortTerm;

	Eigen::Array<double_t, 7, 1> responseAdditions=baseTerms;
	responseAdditions*=md->iterDataPoint->Weight;

	md->curDesignMatrix+=designMatrixAdditions.matrix();
	md->curResponses+=responseAdditions.matrix();

	md->curCoef=(md->curDesignMatrix*md->curDesignMatrix).
		ldlt().
		solve(md->curDesignMatrix*md->curResponses).
		array().max(0);

	double_t predictedWeight=(md->curCoef*baseTerms).sum();
	double_t delta=predictedWeight-md->curDataPoint->Weight;
	if (delta < md->curSmallestDelta) {
		md->curSmallestDelta=delta;
		md->curModelState->Mse=delta*delta;
		md->curModelState->TimeFrame=md->iterDataPoint->DaysSince;
		md->curModelState->PredWeight=predictedWeight;

		md->curModelState->V1=md->curCoef(0,0);
		md->curModelState->V2=md->curCoef(1,0);
		md->curModelState->V3=md->curCoef(2,0);
		md->curModelState->V4=md->curCoef(3,0);
		md->curModelState->V5=md->curCoef(4,0);
		md->curModelState->V6=md->curCoef(5,0);
		md->curModelState->V7=md->curCoef(6,0);
	}
}
