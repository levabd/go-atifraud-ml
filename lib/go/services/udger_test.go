package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
	"testing"
	"fmt"
	"os"
)

func init() {
	helpers.LoadEnv()
}

func TestUdgerInitiation(t *testing.T) {
	_assert := assert.New(t)

	u, err := GetUdgerInstance()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(-1)
	}

	_assert.Equal(607, len(u.Browsers), "main_table be 607 in len")
	_assert.Equal(154, len(u.OS), "value_table be 154 in len")
	_assert.Equal(8, len(u.Devices), "ordered_table be 7 in len")
}
