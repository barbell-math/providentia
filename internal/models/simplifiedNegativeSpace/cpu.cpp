#include "cpu.h"
#include <cmath>
#include <cstdint>
#include <limits>
#include <iostream>
#include "Eigen/Dense"
#include "cgoStructs.h"

// TODO - REFACTOR THIS SHIT
// namespace SimplifiedNegativeSpace {
// 	class Optimizer {
// 	private:
// 		Eigen::Array<double_t, 7, 1> curCoef;
// 		Eigen::Matrix<double_t, 7, 1> curResponses;
// 		Eigen::Matrix<double_t, 7, 7> curDesignMatrix;
// 		trainingLog_t* optimizingFor;
// 		trainingLog_t* iterDataPoint;
// 		modelState_t* curModelState;
// 		double_t curSmallestDelta;
// 	};
// }

typedef struct optimizerData {
	Eigen::Array<double_t, 7, 1> curCoef;
	Eigen::Matrix<double_t, 7, 1> curResponses;
	Eigen::Matrix<double_t, 7, 7> curDesignMatrix;
	trainingLog_t* optimizingFor;
	trainingLog_t* iterDataPoint;
	modelState_t* curModelState;
	double_t curSmallestDelta;
} optimizerData_t;

void addDataPoint(optimizerData_t* optimizerData, opts_t* opts);
double_t makePrediction(optimizerData_t* optimizerData, opts* opts);
void updateModel(optimizerData_t* optimizerData, double_t prediction);

extern "C" void calcModelStates(
	int64_t clientID,
	int32_t modelID,
	trainingLog_t* data,
	int64_t dataLen,
	int64_t startCalcsIdx,
	modelState_t* outValues,
	int64_t outValuesLen,
	opts_t* opts
) {
	for (int64_t i=startCalcsIdx; i<dataLen; i++) {
		outValues[i]=modelState_t{
			.ClientID = clientID,
			.TrainingLogID=data[i].ID,
			.ModelID =  modelID,
			.V1=0, .V2=0, .V3=0, .V4=0, .V5=0, .V6=0, .V7=0,
			.V8=opts->Alpha, .V9=opts->Beta, .V10=opts->Gamma,
			.TimeFrame=0,
			.Mse=std::numeric_limits<double_t>::infinity(),
			.PredWeight=0,
		};

		optimizerData_t optimizerData=optimizerData_t{
			.curCoef=Eigen::Matrix<double_t, 7, 1>::Constant(0),
			.curResponses=Eigen::Matrix<double_t, 7, 1>::Constant(0),
			.curDesignMatrix=Eigen::Matrix<double_t, 7, 7>::Constant(0),
			.optimizingFor=&data[i],
			.iterDataPoint=nullptr,
			.curModelState=&outValues[i],
			.curSmallestDelta=std::numeric_limits<double_t>::infinity(),
		};

		uint64_t iterCntr=0;
		for (int64_t j=i-1; j>=0 && iterCntr<opts->MaxIters; j--) {
			if (data[i].ExerciseID!=data[j].ExerciseID) {
				continue;
			}
			iterCntr++;
			optimizerData.iterDataPoint=&data[j];
			addDataPoint(&optimizerData, opts);
			double_t prediction=makePrediction(&optimizerData, opts);
			updateModel(&optimizerData, prediction);
		}
	}
}

void addDataPoint(optimizerData_t* optimizerData, opts_t* opts) {
	double_t reps=(double_t)optimizerData->iterDataPoint->Reps;
	double_t sets=(double_t)optimizerData->iterDataPoint->Sets;
	double_t interSessionCntr=(double_t)(
		optimizerData->iterDataPoint->InterSessionCntr
	);
	double_t interWorkoutCntr=(double_t)(
		optimizerData->iterDataPoint->InterWorkoutCntr
	);

	double_t effortTerm=pow(optimizerData->iterDataPoint->Effort, opts->Alpha);
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
	responseAdditions*=optimizerData->iterDataPoint->Weight;

	optimizerData->curDesignMatrix+=designMatrixAdditions.matrix();
	optimizerData->curResponses+=responseAdditions.matrix();

	optimizerData->curCoef=(
		optimizerData->curDesignMatrix*optimizerData->curDesignMatrix
	).
		ldlt().
		solve(optimizerData->curDesignMatrix*optimizerData->curResponses).
		array().max(0);
}

double_t makePrediction(optimizerData_t* optimizerData, opts_t* opts) {
	double_t reps=(double_t)optimizerData->optimizingFor->Reps;
	double_t sets=(double_t)optimizerData->optimizingFor->Sets;
	double_t interSessionCntr=(double_t)(
		optimizerData->optimizingFor->InterSessionCntr
	);
	double_t interWorkoutCntr=(double_t)(
		optimizerData->optimizingFor->InterWorkoutCntr
	);

	double_t effortTerm=pow(optimizerData->optimizingFor->Effort, opts->Alpha);
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
	return (optimizerData->curCoef*baseTerms).sum();
}

void updateModel(optimizerData_t* optimizerData, double_t prediction) {
	double_t delta=prediction-optimizerData->optimizingFor->Weight;
	if (delta < optimizerData->curSmallestDelta && prediction>0) {
		optimizerData->curSmallestDelta=delta;
		optimizerData->curModelState->Mse=delta*delta;
		optimizerData->curModelState->TimeFrame=(
			optimizerData->iterDataPoint->DaysSince
			- optimizerData->optimizingFor->DaysSince
		);
		optimizerData->curModelState->PredWeight=prediction;

		optimizerData->curModelState->V1=optimizerData->curCoef(0,0);
		optimizerData->curModelState->V2=optimizerData->curCoef(1,0);
		optimizerData->curModelState->V3=optimizerData->curCoef(2,0);
		optimizerData->curModelState->V4=optimizerData->curCoef(3,0);
		optimizerData->curModelState->V5=optimizerData->curCoef(4,0);
		optimizerData->curModelState->V6=optimizerData->curCoef(5,0);
		optimizerData->curModelState->V7=optimizerData->curCoef(6,0);
	}
}
