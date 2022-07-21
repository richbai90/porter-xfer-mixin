package testdata
import _ "embed"

//go:embed build-input.yaml
var BuildInput string

//go:embed invalid-input.yaml
var InvalidInput string

//go:embed step-input.yaml
var StepInput string