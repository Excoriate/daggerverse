package terragrunt

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type TerragruntValidator interface {
	workDirHasHCLFiles(workdir string) error
	hclFileIsNotEmpty(hclFile string) error
	workDirHasTerragruntFiles(workdir string) error
}

type TerragruntValidatorImpl struct {
}

func NewValidator() TerragruntValidator {
	return &TerragruntValidatorImpl{}
}

func (v *TerragruntValidatorImpl) workDirHasHCLFiles(workdir string) error {
	var hasHCLFiles bool

	err := filepath.Walk(workdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".hcl") {
			hasHCLFiles = true
		}
		return nil
	})

	if err != nil {
		return err
	}
	if !hasHCLFiles {
		return errors.New("no .hcl files found in the directory")
	}
	return nil
}

func (v *TerragruntValidatorImpl) hclFileIsNotEmpty(hclFile string) error {
	if !strings.HasSuffix(hclFile, "terragrunt.hcl") {
		return errors.New("file is not terragrunt.hcl")
	}

	fileInfo, err := os.Stat(hclFile)
	if err != nil {
		return err
	}
	if fileInfo.Size() == 0 {
		return errors.New("terragrunt.hcl file is empty")
	}
	return nil
}

func (v *TerragruntValidatorImpl) workDirHasTerragruntFiles(workdir string) error {
	terragruntFilePath := filepath.Join(workdir, "terragrunt.hcl")
	if _, err := os.Stat(terragruntFilePath); os.IsNotExist(err) {
		return errors.New("terragrunt.hcl file not found in the directory")
	}
	return nil
}
