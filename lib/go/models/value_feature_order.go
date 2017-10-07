package models


type ValueFeatureOrder struct {
	ID				uint
	FeatureName		string
	Order			int
}

func (ValueFeatureOrder) TableName() string {
	return "value_features_order"
}