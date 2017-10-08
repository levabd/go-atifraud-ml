package services

import "fmt"

// ConfusionMatrix describes a confusion matrix
type ConfusionMatrix struct {
	Positive      int
	Negative      int
	TruePositive  int
	TrueNegative  int
	FalsePositive int
	FalseNegative int
	Recall        float64
	Precision     float64
	Accuracy      float64
}

func (cm ConfusionMatrix) String() string {
	return fmt.Sprintf("\tPositives: %d\n\tNegatives: %d\n\tTrue Positives: %d\n\tTrue Negatives: %d\n\tFalse Positives: %d\n\tFalse Negatives: %d\n\n\tRecall: %.2f\n\tPrecision: %.2f\n\tAccuracy: %.2f\n",
		cm.Positive, cm.Negative, cm.TruePositive, cm.TrueNegative, cm.FalsePositive, cm.FalseNegative, cm.Recall, cm.Precision, cm.Accuracy)
}