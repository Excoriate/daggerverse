package main

import "fmt"

type IacTerragruntCMDError struct {
	ErrWrapped error
	Message    string
}

const iacTerragruntCMDErrorBaselineMsg = "IacTerragrunt command error"

func (e *IacTerragruntCMDError) Error() string {
	if e.ErrWrapped == nil {
		return fmt.Sprintf("%s: %s", iacTerragruntCMDErrorBaselineMsg, e.Message)
	}

	return fmt.Sprintf("%s: %s: %s", iacTerragruntCMDErrorBaselineMsg, e.Message, e.ErrWrapped.Error())
}
