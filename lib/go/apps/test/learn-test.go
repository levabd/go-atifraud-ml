package main

import (
	"fmt"
	//"strconv"

	// Golearn
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/linear_models"
)


func main() {
	// Load data
	X, _ := base.ParseCSVToInstances("exams.csv", true)
	Y, _ := base.ParseCSVToInstances("exams.csv", true)

	// Setup the problem
	//lr, _ := linear_models.NewLogisticRegression("l2", 1.0, 1e-6)
	lr := linear_models.NewLinearRegression()

	lr.Fit(X)

	Z, _ := lr.Predict(Y)
	fmt.Println(Z.RowString(0))
	fmt.Println(Z.RowString(1))
}
