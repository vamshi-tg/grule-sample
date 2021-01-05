package grulesample

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

type MyPoGo struct {
	InVal  string
	OutVal string
}

func (p *MyPoGo) GetStringLength(str string) int {
	return len(str)
}

func (p *MyPoGo) AppendString(aString, subString string) string {
	return fmt.Sprintf("%s%s", aString, subString)
}

const functionsGrl = `
	rule Rule1 {
		when
			Pogo.GetStringLength(Pogo.InVal) < 4
		then
			Pogo.OutVal = Pogo.AppendString(Pogo.InVal, "Grooling");
			Complete();
	}

	rule Rule2 {
		when
			Pogo.GetStringLength(Pogo.InVal) > 4
		then
			Pogo.OutVal = Pogo.InVal;
			Complete();
	}
`

func TestFunctions(t *testing.T) {
	testData := []*struct {
		Pogo       *MyPoGo
		WantOutVal string
	}{
		{
			Pogo: &MyPoGo{
				InVal: "Go",
			},
			WantOutVal: "GoGrooling",
		},
		{
			Pogo: &MyPoGo{
				InVal: "Google",
			},
			WantOutVal: "Google",
		},
	}
	// Prepare knowledgebase library and load it with our rule.
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	byteArr := pkg.NewBytesResource([]byte(functionsGrl))
	err := rb.BuildRuleFromResource("FunctionOperationsTutorial", "0.0.1", byteArr)
	assert.NoError(t, err)

	engine := &engine.GruleEngine{
		MaxCycle: 10,
	}
	knowledgeBase := lib.NewKnowledgeBaseInstance("FunctionOperationsTutorial", "0.0.1")

	for i, td := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dataCtx := ast.NewDataContext()
			err := dataCtx.Add("Pogo", td.Pogo)
			assert.NoError(t, err)

			err = engine.Execute(dataCtx, knowledgeBase)
			assert.NoError(t, err)
			assert.Equal(t, td.WantOutVal, td.Pogo.OutVal)
		})
	}
}
