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

type Item struct {
	Name     string
	Price    int64
	Discount int64
}

const logicalOperationsGrl = `
	rule HondaRule1 "When Honda Price is greater than 2000, 30% discount is applicable" {
		When
			Item.Discount == 0 && Item.Name == "Honda" && Item.Price > 2000
		Then
			Log("Item Price:" + Item.Price + " is > 2000. 30% Discount will be applicable");
			Item.Discount = 30;
	}

	rule HondaRule2 "When Honda Price is less than 2000, 0% discount is applicable" {
		When
			Item.Discount == 0 && Item.Name == "Honda" && Item.Price < 2000
		Then
			Log("Item Price:" + Item.Price + " is < 2000. 30% Discount will be applicable");
			Item.Discount = 0;
			Retract("HondaRule2");
	}

	rule HeroOrSuzikiRule1 "When Item Name is Hero or Suziki, 40% discount is applicable" {
		When
		Item.Discount == 0 && (Item.Name == "Hero" || Item.Name == "Suziki")
	Then
		Log("Item Name is " + Item.Name + ". 40% Discount will be applicable");
		Item.Discount = 40;
	}
`
/* Setting Retract in HondaRule2
	Refer https://github.com/hyperjumptech/grule-rule-engine/blob/master/docs/FAQ_en.md
	for more details
*/

func TestItems(t *testing.T) {
	testData := []*struct {
		Item         *Item
		WantDiscount int64
	}{
		{
			Item: &Item{
				Name:  "Honda",
				Price: 2300,
			},
			WantDiscount: 30,
		},
		{
			Item: &Item{
				Name:  "Honda",
				Price: 1500,
			},
			WantDiscount: 0,
		},
		{
			Item: &Item{
				Name:  "Hero",
				Price: 1500,
			},
			WantDiscount: 40,
		},
		{
			Item: &Item{
				Name:  "Suziki",
				Price: 1500,
			},
			WantDiscount: 40,
		},
	}

	// Prepare knowledgebase library and load it with our rule.
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	byteArr := pkg.NewBytesResource([]byte(logicalOperationsGrl))
	err := rb.BuildRuleFromResource("ItemTutorial", "0.0.1", byteArr)
	assert.NoError(t, err)

	engine := &engine.GruleEngine{
		MaxCycle: 10,
	}
	knowledgeBase := lib.NewKnowledgeBaseInstance("ItemTutorial", "0.0.1")

	for i, td := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dataCtx := ast.NewDataContext()
			err := dataCtx.Add("Item", td.Item)
			assert.NoError(t, err)

			err = engine.Execute(dataCtx, knowledgeBase)
			assert.NoError(t, err)
			assert.Equal(t, td.WantDiscount, td.Item.Discount)
		})
	}
}
