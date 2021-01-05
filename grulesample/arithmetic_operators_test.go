package grulesample

import (
	"strconv"
	"testing"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/stretchr/testify/assert"
)

type ArithmeticOperation struct {
	Operation string
	OperandA  float32
	OperandB  float32
	Result    float32
}

const arithmeticGrl = `
	rule Addition "Should perform addition" {
		When
			AO.Result == 0 && AO.Operation == "addition"
		Then
			AO.Result = AO.OperandA + AO.OperandB;
	}

	rule Subtraction "Should perform subtraction" {
		When
			AO.Result == 0 && AO.Operation == "subtraction"
		Then
			AO.Result = AO.OperandA - AO.OperandB;
	}

	rule Multiplication "Should perform multiplication" {
		When
			AO.Result == 0 && AO.Operation == "multiplication"
		Then
			AO.Result = AO.OperandA * AO.OperandB;
	}

	rule Division "Should perform division" {
		When
			AO.Result == 0 && AO.Operation == "division"
		Then
			AO.Result = AO.OperandA / AO.OperandB;
	}

	rule Modulus "Should perform modulo division" {
		When
			AO.Result == 0 && AO.Operation == "modulus"
		Then
			AO.Result = 8 % 3;
	}
`

func TestArithmeticOperations(t *testing.T) {
	testData := []*struct {
		AoObj      *ArithmeticOperation
		WantResult float32
	}{
		{
			AoObj: &ArithmeticOperation{
				Operation: "addition",
				OperandA:  2.5,
				OperandB:  3,
			},
			WantResult: 5.5,
		},
		{
			AoObj: &ArithmeticOperation{
				Operation: "multiplication",
				OperandA:  4,
				OperandB:  2,
			},
			WantResult: 8,
		},
		{
			AoObj: &ArithmeticOperation{
				Operation: "subtraction",
				OperandA:  30,
				OperandB:  45,
			},
			WantResult: -15,
		},
		{
			AoObj: &ArithmeticOperation{
				Operation: "division",
				OperandA:  44,
				OperandB:  5,
			},
			WantResult: 8.8,
		},
		{
			AoObj: &ArithmeticOperation{
				Operation: "modulus",
			},
			WantResult: 2,
		},
	}

	// Prepare knowledgebase library and load it with our rule.
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	byteArr := pkg.NewBytesResource([]byte(arithmeticGrl))
	err := rb.BuildRuleFromResource("ArithmeticOperationTutorial", "0.0.1", byteArr)
	assert.NoError(t, err)

	engine := &engine.GruleEngine{
		MaxCycle: 10,
	}
	knowledgeBase := lib.NewKnowledgeBaseInstance("ArithmeticOperationTutorial", "0.0.1")

	for i, td := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dataCtx := ast.NewDataContext()
			err := dataCtx.Add("AO", td.AoObj)
			assert.NoError(t, err)

			err = engine.Execute(dataCtx, knowledgeBase)
			assert.NoError(t, err)
			assert.Equal(t, td.WantResult, td.AoObj.Result)
		})
	}
}
