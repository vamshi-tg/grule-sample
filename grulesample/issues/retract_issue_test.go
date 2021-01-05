package issues

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/stretchr/testify/assert"
)

type ComputedInputs struct {
	Name string
	Bid  float64
	Rate float64
}

const Grl = `
	rule SampleRule "Hero/Honda" {
		When
			ComputedInputs.Name == "Hero" || ComputedInputs.Name == "Honda"
		Then
			Log("SampleRule");
			ComputedInputs.Bid = ComputedInputs.Rate;
			Retract("SampleRule");
	}
`

func TestItems(t *testing.T) {
	testData := []*struct {
		CI      *ComputedInputs
		WantBid float64
	}{
		{
			CI: &ComputedInputs{
				Name: "Hero",
				Rate: 10,
			},
			WantBid: 10,
		},
		{
			CI: &ComputedInputs{
				Name: "Honda",
				Rate: 10,
			},
			WantBid: 10,
		},
	}

	// Prepare knowledgebase library and load it with our rule.
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	byteArr := pkg.NewBytesResource([]byte(Grl))
	err := rb.BuildRuleFromResource("Tutorial", "0.0.1", byteArr)
	assert.NoError(t, err)

	engine := &engine.GruleEngine{
		MaxCycle: 10,
	}
	knowledgeBase := lib.NewKnowledgeBaseInstance("Tutorial", "0.0.1")

	for i, td := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dataCtx := ast.NewDataContext()
			err := dataCtx.Add("ComputedInputs", td.CI)
			assert.NoError(t, err)

			rules, _ := engine.FetchMatchingRules(dataCtx, knowledgeBase)
			fmt.Printf("\nNo. of matching rules: %d\n", len(rules))

			err = engine.Execute(dataCtx, knowledgeBase)
			assert.NoError(t, err)
			assert.Equal(t, td.WantBid, td.CI.Bid)
		})
	}
}
