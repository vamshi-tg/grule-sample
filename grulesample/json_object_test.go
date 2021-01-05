package grulesample

import (
	"encoding/json"
	"testing"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/stretchr/testify/assert"
)

/*
	type TPayloadMap = map[string]interface{} will not work with the latest version of grule engine
*/
type TPayloadMap = map[string]int

type Result struct {
	Val float64
}
type HasuraData struct {
	Payload TPayloadMap
	Result  float64
}

const arithmeticGrlTwo = `
	rule Addition "Should perform addition" {
		When
			Result.Val == 0 && AO["operation"] == 23
		Then
			Result.Val = AO["operand_a"] + AO["operand_b"];
	}`

func createMapFromJSON(jsonPayload string) TPayloadMap {
	var jsonByt []byte = []byte(jsonPayload)

	var hasuraPayload TPayloadMap
	if err := json.Unmarshal(jsonByt, &hasuraPayload); err != nil {
		panic(err)
	}
	return hasuraPayload
}

func TestJSONObject(t *testing.T) {
	jsonPayload := `{"operation": 23,"operand_a":20, "operand_b":30}`
	payloadMap := createMapFromJSON(jsonPayload)
	result := new(Result)

	// Prepare knowledgebase library and load it with our rule.
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	byteArr := pkg.NewBytesResource([]byte(arithmeticGrlTwo))
	err := rb.BuildRuleFromResource("ArithmeticOperationTutorial", "0.0.1", byteArr)
	assert.NoError(t, err)

	engine := &engine.GruleEngine{
		MaxCycle: 10,
	}
	knowledgeBase := lib.NewKnowledgeBaseInstance("ArithmeticOperationTutorial", "0.0.1")

	dataCtx := ast.NewDataContext()
	err = dataCtx.Add("AO", payloadMap)
	assert.NoError(t, err)
	err = dataCtx.Add("Result", result)

	err = engine.Execute(dataCtx, knowledgeBase)
	assert.NoError(t, err)
	assert.Equal(t, float64(50), result.Val)
}
