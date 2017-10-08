package main

import (
	"time"
	"log"
	"fmt"
	//"strconv"

	// Golearn
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/linear_models"
	"github.com/sjwhitworth/golearn/evaluation"

	// Goml
	b "github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/linear"

	"github.com/levabd/go-atifraud-ml/lib/go/services"
	"github.com/uniplaces/carbon"
	"strconv"
	"io/ioutil"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
	defer timeTrack(time.Now(), "main")

	startTime := time.Now()
	startNanosecond := startTime.Nanosecond()
	println(startTime.Minute(), startTime.Second(), startNanosecond)

	_, floatUAClasses, fullFeatures := services.PrepareData(carbon.Now().SubMonths(2).Unix(), carbon.Now().Unix())

	println(len(floatUAClasses), len(fullFeatures))

	fmt.Println(len(fullFeatures[0]))
	fmt.Println(fullFeatures[0])

	_, userAgentFloatCodes := services.LoadFittedUserAgentCodes()
	// valuesFeaturesOrder := services.LoadFittedValuesFeaturesOrder()
	// tryGolearn(fullFeatures, valuesFeaturesOrder, floatUAClasses, userAgentFloatCodes)
	tryGoml(fullFeatures, floatUAClasses, userAgentFloatCodes, 0.1)

	end := time.Now()
	println(end.Minute(), end.Second(), end.Nanosecond(), end.Nanosecond() - startNanosecond)
}

func tryGoml(fullFeatures [][]float64, floatUAClasses []float64, userAgentFloatCodes map[string]float64, decisionBoundary float64) {

	classNumbers := len(userAgentFloatCodes)

	fullSampleSize := len(fullFeatures)
	percent70 := int(float64(fullSampleSize) * 0.7)
	xTrain	:= fullFeatures[0:percent70]
	xTest	:= fullFeatures[percent70 + 1:fullSampleSize - 1]
	yTrain	:= floatUAClasses[0:percent70]
	yTest	:= floatUAClasses[percent70 + 1:fullSampleSize - 1]

	/*cm := services.ConfusionMatrix{}
	for _, y := range yTest {
		if y == 1.0 {
			cm.Positive++
		}
		if y == 0.0 {
			cm.Negative++
		}
	}*/

	model := linear.NewSoftmax(b.BatchGA, 0.00001, 0, classNumbers,1000, xTrain, yTrain)
	model.Output = ioutil.Discard

	err := model.Learn()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(yTest[0])
	fmt.Println(model.Predict(xTest[0]))
	/*for i := range xTest {
		prediction, err := model.Predict(xTest[i])
		if err != nil {
			fmt.Println(err)
		}
		y := int(yTest[i])
		positive := prediction[0] >= decisionBoundary

		if y == 1 && positive {
			cm.TruePositive++
		}
		if y == 1 && !positive {
			cm.FalseNegative++
		}
		if y == 0 && positive {
			cm.FalsePositive++
		}
		if y == 0 && !positive {
			cm.TrueNegative++
		}
	}

	// Calculate Evaluation Metrics
	cm.Recall = float64(cm.TruePositive) / float64(cm.Positive)
	cm.Precision = float64(cm.TruePositive) / (float64(cm.TruePositive) + float64(cm.FalsePositive))
	cm.Accuracy = float64(float64(cm.TruePositive)+float64(cm.TrueNegative)) / float64(float64(cm.Positive)+float64(cm.Negative))*/
}

//noinspection GoUnusedFunction,GoUnusedParameter
func tryGolearn(fullFeatures [][]float64, valuesFeaturesOrder map[string]int, floatUAClasses []float64, userAgentFloatCodes map[string]float64) {
	featuresSize := len(fullFeatures[0])
	attrs := make([]base.Attribute, featuresSize + 1)
	// Features
	for key, value := range services.OrderPairsFeaturesOrder {
		attrs[value] = base.NewFloatAttribute(key)
	}
	for key, value := range valuesFeaturesOrder {
		attrs[value+len(services.OrderPairsFeaturesOrder)] = base.NewFloatAttribute(key)
	}
	// User Agents
	attrs[featuresSize] = base.NewFloatAttribute("UserAgent")
	// Insert a standard class
	// attrs[featuresSize].GetSysValFromString("NaN")

	instance := base.NewDenseInstances()

	// Add the attributes
	specs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		specs[i] = instance.AddAttribute(a)
	}

	err := instance.AddClassAttribute(attrs[len(attrs)-1])
	if err != nil {
		panic(err)
	}

	// Allocate space
	instance.Extend(len(fullFeatures))
	// Write the data
	for row, rowFeatures := range fullFeatures {
		for col, value := range rowFeatures {
			instance.Set(specs[col], row, specs[col].GetAttribute().GetSysValFromString(strconv.FormatFloat(value, 'f', 1, 32)))
		}
		instance.Set(specs[featuresSize], row, specs[featuresSize].GetAttribute().GetSysValFromString(strconv.FormatFloat(floatUAClasses[row], 'f', 1, 32)))
	}

	now := time.Now()
	println(now.Minute(), now.Second(), now.Nanosecond())
	fmt.Println("Features collected into instance")

	//Initialises a new classifier
	//cls, err := linear_models.NewLogisticRegression("l2", 100, 1e-6)
	cls, err := linear_models.NewLogisticRegression("l2", 1, 1e-6)
	if err != nil {
		panic(err)
	}

	//Do a training-test split
	trainData, testData := base.InstancesTrainTestSplit(instance, 0.70)

	now = time.Now()
	println(now.Minute(), now.Second(), now.Nanosecond())
	fmt.Println("Train Test sets created. Fitting Started.")
	cls.Fit(trainData)

	//Calculates the Euclidean distance and returns the most popular label
	predictions, err := cls.Predict(testData)
	if err != nil {
		panic(err)
	}

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}

	fmt.Println(evaluation.GetSummary(confusionMat))
}