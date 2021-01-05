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

/*
	def get_bid_strategy_value():
			# If SFDC currency is not equal to DV360 currency, then convert SFDC CPM rate to DV360 Currency
			goal_type = computed_inputs.goal_type
			prod_type = com.product.product_type
			cpm = computed_inputs.cpm
			bid = cpm
			if goal_type == 'Conversion':
					bid = 1.5 * cpm
			if goal_type == 'Click' or goal_type == 'Video':
					bid = cpm
			if cpm == 0:
					if prod_type in ['Display', 'Display - CPM']:
							bid = 5
					elif prod_type in ['Video', 'Video - CPM']:
							bid = 7
					elif prod_type in ['Native', 'Native - CPM']:
							bid = 5
			if user_inputs.client_restrictions['Include Site List(Channel ID)']:
					bid = bid*1.1
			if user_inputs.client_restrictions['Include Geography']:
					bid = bid*1.25
			if user_inputs.client_restrictions['Include KCT List']:
					bid = bid*1.15
			return bid
*/

type ComputedInputs struct {
	GoalType string
	Bid      float64
	CPM      float64
}

type Product struct {
	ProductType string
}

type UserInputs struct {
	ClientRestrictions map[string]bool
}

var defaultUserInputs = &UserInputs{
	ClientRestrictions: map[string]bool{
		"Include Site List(Channel ID)": false,
		"Include Geography":             false,
		"Include KCT List":              false,
	},
}

const (
	Display    = "Display"
	DisplayCpm = "Display - CPM"
	Video      = "Video"
	VideoCpm   = "Video - CPM"
	Native     = "Native"
	NativeCpm  = "Native - CPM"
)

const (
	BidStrategyRule1 = "BidStrategyRule1"
	BidStrategyRule2 = "BidStrategyRule2"
	BidStrategyRule3 = "BidStrategyRule3"
	BidStrategyRule4 = "BidStrategyRule4"
	BidStrategyRule5 = "BidStrategyRule5"
	BidStrategyRule6 = "BidStrategyRule6"
	BidStrategyRule7 = "BidStrategyRule7"
	BidStrategyRule8 = "BidStrategyRule8"
)

const BidStrategyRulesGrl = `
	rule BidStrategyRule1 "Conversion" salience -1 {
		When
			ComputedInputs.GoalType == "Conversion"
		Then
			ComputedInputs.Bid = 1.5 * ComputedInputs.CPM;
			Retract("BidStrategyRule1");
	}

	rule BidStrategyRule2 "Click/Video" salience -2 {
		When
			ComputedInputs.GoalType == "Click" || ComputedInputs.GoalType == "Video"
		Then
			ComputedInputs.Bid = ComputedInputs.CPM;
			Retract("BidStrategyRule2");
	}

	rule BidStrategyRule3 "When CPM is zero, and Product Type Display" salience -3 {
		When
			ComputedInputs.CPM == 0 &&
			(Product.ProductType == "Display" || Product.ProductType == "Display - CPM")
		Then
			ComputedInputs.Bid = 5;
			Retract("BidStrategyRule3");
	}

	rule BidStrategyRule4 "When CPM is zero, and Product Type Video" salience -4 {
		When
			ComputedInputs.CPM == 0 &&
			(Product.ProductType == "Video" || Product.ProductType == "Video - CPM")
		Then
			ComputedInputs.Bid = 6;
			Retract("BidStrategyRule4");
	}

	rule BidStrategyRule5 "When CPM is zero, and Product Type Native" salience -5 {
		When
			ComputedInputs.CPM == 0 &&
			(Product.ProductType == "Native" || Product.ProductType == "Native - CPM")
		Then
			ComputedInputs.Bid = 7;
			Retract("BidStrategyRule5");
	}

	rule BidStrategyRule6 "When 'Include Site List(Channel ID)' in UserInputs is true" salience -6 {
		When
			UserInputs.ClientRestrictions["Include Site List(Channel ID)"]
		Then
			ComputedInputs.Bid = ComputedInputs.Bid * 2;
			Retract("BidStrategyRule6");
	}

	rule BidStrategyRule7 "When 'Include Geography' in UserInputs is true" salience -7 {
		When
			UserInputs.ClientRestrictions["Include Geography"]
		Then
			ComputedInputs.Bid = ComputedInputs.Bid * 3;
			Retract("BidStrategyRule7");
	}

	rule BidStrategyRule8 "When 'Include KCT List' in UserInputs is true" salience -8 {
		When
			UserInputs.ClientRestrictions["Include KCT List"]
		Then
			ComputedInputs.Bid = ComputedInputs.Bid * 4;
			Retract("BidStrategyRule8");
	}
`

/* Setting Retract
Refer https://github.com/hyperjumptech/grule-rule-engine/blob/master/docs/FAQ_en.md
for more details
*/

func TestMapOperations(t *testing.T) {
	testData := []*struct {
		MatchingRules []string
		CI            *ComputedInputs
		Product       *Product
		UserInputs    *UserInputs
		WantBid       float64
	}{
		{
			MatchingRules: []string{BidStrategyRule1},
			CI: &ComputedInputs{
				GoalType: "Conversion",
				CPM:      10,
			},
			UserInputs: defaultUserInputs,
			Product: &Product{
				ProductType: Display,
			},
			WantBid: 15,
		},
		{
			MatchingRules: []string{BidStrategyRule1, BidStrategyRule3},
			CI: &ComputedInputs{
				GoalType: "Conversion",
				CPM:      0,
			},
			UserInputs: defaultUserInputs,
			Product: &Product{
				ProductType: Display,
			},
			WantBid: 5,
		},
		{
			MatchingRules: []string{BidStrategyRule2},
			CI: &ComputedInputs{
				GoalType: "Click",
				CPM:      12,
			},
			UserInputs: defaultUserInputs,
			Product: &Product{
				ProductType: "Meta",
			},
			WantBid: 12,
		},
		{
			MatchingRules: []string{BidStrategyRule2},
			CI: &ComputedInputs{
				GoalType: "Video",
				CPM:      14,
			},
			UserInputs: defaultUserInputs,
			Product: &Product{
				ProductType: DisplayCpm,
			},
			WantBid: 14,
		},
		{
			MatchingRules: []string{BidStrategyRule2, BidStrategyRule3},
			CI: &ComputedInputs{
				GoalType: "Click",
				CPM:      0,
			},
			UserInputs: defaultUserInputs,
			Product: &Product{
				ProductType: Display,
			},
			WantBid: 5,
		},
		{
			MatchingRules: []string{BidStrategyRule2, BidStrategyRule3},
			CI: &ComputedInputs{
				GoalType: "Video",
				CPM:      0,
			},
			UserInputs: defaultUserInputs,
			Product: &Product{
				ProductType: DisplayCpm,
			},
			WantBid: 5,
		},
		{
			MatchingRules: []string{BidStrategyRule2, BidStrategyRule3, BidStrategyRule6},
			CI: &ComputedInputs{
				GoalType: "Video",
				CPM:      0,
			},
			Product: &Product{
				ProductType: DisplayCpm,
			},
			UserInputs: &UserInputs{
				ClientRestrictions: map[string]bool{
					"Include Site List(Channel ID)": true,
					"Include Geography":             false,
					"Include KCT List":              false,
				},
			},
			WantBid: 10,
		},
		{
			MatchingRules: []string{BidStrategyRule2, BidStrategyRule3, BidStrategyRule6, BidStrategyRule7, BidStrategyRule8},
			CI: &ComputedInputs{
				GoalType: "Video",
				CPM:      0,
			},
			Product: &Product{
				ProductType: DisplayCpm,
			},
			UserInputs: &UserInputs{
				ClientRestrictions: map[string]bool{
					"Include Site List(Channel ID)": true,
					"Include Geography":             true,
					"Include KCT List":              true,
				},
			},
			WantBid: 120,
		},
	}

	// Prepare knowledgebase library and load it with our rule.
	lib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(lib)
	byteArr := pkg.NewBytesResource([]byte(BidStrategyRulesGrl))
	err := rb.BuildRuleFromResource("Tutorial", "0.0.1", byteArr)
	assert.NoError(t, err)

	engine := &engine.GruleEngine{
		MaxCycle: 10,
	}

	for i, td := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			knowledgeBase := lib.NewKnowledgeBaseInstance("Tutorial", "0.0.1")

			// Add required data to the DataContext
			dataCtx := ast.NewDataContext()
			err := dataCtx.Add("ComputedInputs", td.CI)
			assert.NoError(t, err)
			err = dataCtx.Add("Product", td.Product)
			assert.NoError(t, err)
			err = dataCtx.Add("UserInputs", td.UserInputs)
			assert.NoError(t, err)

			rules, _ := engine.FetchMatchingRules(dataCtx, knowledgeBase)
			fmt.Printf("\nNo. of matching rules: %d\n", len(rules))
			var gotRules []string
			for _, rule := range rules {
				gotRules = append(gotRules, rule.RuleName)
			}
			fmt.Println(gotRules)
			assert.Equal(t, td.MatchingRules, gotRules)

			err = engine.Execute(dataCtx, knowledgeBase)
			assert.NoError(t, err)
			assert.Equal(t, td.WantBid, td.CI.Bid)
		})
	}
}
