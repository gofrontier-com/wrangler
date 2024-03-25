package core

import "testing"

func float64Ptr(f float64) *float64 {
	return &f
}

func TestCheckBudgets_ReturnsViolationWhenMonthlyFixedRuleIsExceeded(t *testing.T) {
	config := Config{
		Budgets: &[]Budget{
			{
				ResourceId:    "app-1",
				MonthlyAmount: float64Ptr(200),
			},
			{
				ResourceId:    "app-2",
				MonthlyAmount: float64Ptr(50),
				Rules: &[]BudgetRule{
					{
						Name:   "test-rule",
						Type:   FixedRule,
						Period: Monthly,
						Value:  20,
					},
				},
			},
		},
	}
	data := []CostRecord{
		{
			ResourceId: "app-2",
			Value:      100,
			Currency:   GBP,
			Period:     Monthly,
		},
	}
	actual, err := CheckBudgets(data, config)
	expected := []BudgetRuleViolation{
		{
			ResourceId:   "app-2",
			Description:  "actual amount >= 20.00",
			BudgetAmount: 50,
			ActualAmount: 100,
		},
	}
	if err != nil {
		t.Errorf("CheckBudgets(data, config) failed with error %t", err)
	}
	if len(actual) != len(expected) {
		t.Errorf("CheckBudgets(data, config) = %d; want %d", len(actual), len(expected))
	}

	violation := actual[0]
	if violation.ResourceId != expected[0].ResourceId {
		t.Errorf("Violation.ResourceId = %s; want %s", violation.ResourceId, expected[0].ResourceId)
	}
	if violation.Description != expected[0].Description {
		t.Errorf("Violation.Description = %s; want %s", violation.Description, expected[0].Description)
	}
	if violation.BudgetAmount != expected[0].BudgetAmount {
		t.Errorf("Violation.BudgetAmount = %f; want %f", violation.BudgetAmount, expected[0].BudgetAmount)
	}
	if violation.ActualAmount != expected[0].ActualAmount {
		t.Errorf("Violation.ActualAmount = %f; want %f", violation.ActualAmount, expected[0].ActualAmount)
	}
}

func TestCheckBudgets_ReturnsViolationWhenDailyFixedRuleIsExceeded(t *testing.T) {
	config := Config{
		Budgets: &[]Budget{
			{
				ResourceId:  "app-1",
				DailyAmount: float64Ptr(15),
			},
			{
				ResourceId:  "app-2",
				DailyAmount: float64Ptr(20),
				Rules: &[]BudgetRule{
					{
						Name:   "test-rule",
						Type:   FixedRule,
						Period: Daily,
						Value:  10,
					},
				},
			},
		},
	}
	data := []CostRecord{
		{
			ResourceId: "app-2",
			Value:      100,
			Currency:   GBP,
			Period:     Daily,
		},
	}
	actual, err := CheckBudgets(data, config)
	expected := []BudgetRuleViolation{
		{
			ResourceId:   "app-2",
			Description:  "actual amount >= 10.00",
			BudgetAmount: 20,
			ActualAmount: 100,
		},
	}
	if err != nil {
		t.Errorf("CheckBudgets(data, config) failed with error %t", err)
	}
	if len(actual) != len(expected) {
		t.Errorf("CheckBudgets(data, config) = %d; want %d", len(actual), len(expected))
	}

	violation := actual[0]
	if violation.ResourceId != expected[0].ResourceId {
		t.Errorf("Violation.ResourceId = %s; want %s", violation.ResourceId, expected[0].ResourceId)
	}
	if violation.Description != expected[0].Description {
		t.Errorf("Violation.Description = %s; want %s", violation.Description, expected[0].Description)
	}
	if violation.BudgetAmount != expected[0].BudgetAmount {
		t.Errorf("Violation.BudgetAmount = %f; want %f", violation.BudgetAmount, expected[0].BudgetAmount)
	}
	if violation.ActualAmount != expected[0].ActualAmount {
		t.Errorf("Violation.ActualAmount = %f; want %f", violation.ActualAmount, expected[0].ActualAmount)
	}
}

func TestCheckBudgets_ReturnsNoViolationsWhenNoBudgets(t *testing.T) {
	config := Config{
		Budgets: &[]Budget{},
	}
	data := []CostRecord{
		{
			ResourceId: "test-svc",
			Value:      100,
		},
	}
	actual, err := CheckBudgets(data, config)
	expected := []BudgetRuleViolation{}
	if err != nil {
		t.Errorf("CheckBudgets(data, config) failed with error %t", err)
	}
	if len(actual) != len(expected) {
		t.Errorf("CheckBudgets(data, config) = %d; want %d", len(actual), len(expected))
	}
}

func TestCheckBudgets_ReturnsNoViolationsWhenNoBudgetsMatchResource(t *testing.T) {
	config := Config{
		Budgets: &[]Budget{
			{
				ResourceId:    "app-1",
				MonthlyAmount: float64Ptr(100),
			},
			{
				ResourceId:    "app-2",
				MonthlyAmount: float64Ptr(50),
			},
		},
	}
	data := []CostRecord{
		{
			ResourceId: "app-3",
			Value:      100,
		},
	}
	actual, err := CheckBudgets(data, config)
	expected := []BudgetRuleViolation{}
	if err != nil {
		t.Errorf("CheckBudgets(data, config) failed with error %t", err)
	}
	if len(actual) != len(expected) {
		t.Errorf("CheckBudgets(data, config) = %d; want %d", len(actual), len(expected))
	}
}

func TestCheckBudgets_ReturnsNoViolationsWhenNoBudgetsMatchCurrency(t *testing.T) {
	config := Config{
		Budgets: &[]Budget{
			{
				ResourceId:    "app-1",
				MonthlyAmount: float64Ptr(200),
			},
			{
				ResourceId:    "app-2",
				MonthlyAmount: float64Ptr(50),
				Currency:      EUR,
				Rules: &[]BudgetRule{
					{
						Name:   "test-rule",
						Type:   FixedRule,
						Value:  10,
						Period: Monthly,
					},
				},
			},
		},
	}
	data := []CostRecord{
		{
			ResourceId: "app-2",
			Value:      100,
			Currency:   GBP,
		},
	}
	actual, err := CheckBudgets(data, config)
	expected := []BudgetRuleViolation{}
	if err != nil {
		t.Errorf("CheckBudgets(data, config) failed with error %t", err)
	}
	if len(actual) != len(expected) {
		t.Errorf("CheckBudgets(data, config) = %d; want %d", len(actual), len(expected))
	}
}

/** Budget **/

func TestBudget_GetAmount_ReturnsDailyAmount(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
		DailyAmount:   float64Ptr(5.00),
	}
	actual := *budget.GetAmount(Daily)
	expected := 5.00
	if actual != expected {
		t.Errorf("GetAmount(Daily) = %f; want %f", actual, expected)
	}
}

func TestBudget_GetAmount_ReturnsMonthlyAmount(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
		DailyAmount:   float64Ptr(5.00),
	}
	actual := *budget.GetAmount(Monthly)
	expected := 100.00
	if actual != expected {
		t.Errorf("GetAmount(Daily) = %f; want %f", actual, expected)
	}
}

func TestBudget_HasAmount_ReturnsTrueWhenMonthlyBudgetIsSet(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
	}
	actual := budget.HasAmount()
	expected := true
	if actual != expected {
		t.Errorf("HasAmount() = %t; want %t", actual, expected)
	}
}

func TestBudget_HasAmount_ReturnsTrueWhenDailyBudgetIsSet(t *testing.T) {
	budget := Budget{
		DailyAmount: float64Ptr(10.00),
	}
	actual := budget.HasAmount()
	expected := true
	if actual != expected {
		t.Errorf("HasAmount() = %t; want %t", actual, expected)
	}
}

func TestBudget_HasAmount_ReturnsFalseWhenNoAmountIsSet(t *testing.T) {
	budget := Budget{}
	actual := budget.HasAmount()
	expected := false
	if actual != expected {
		t.Errorf("HasAmount() = %t; want %t", actual, expected)
	}
}

/** Fixed budget rules **/

func TestFixedBudgetRule_Evaluate_ReturnsTrueWhenCostValueExceedsRuleValue(t *testing.T) {
	budget := Budget{
		DailyAmount: float64Ptr(0),
	}
	rule := BudgetRule{
		Name:  "test-rule",
		Type:  FixedRule,
		Value: 10.00,
	}
	record := CostRecord{
		Value:  11.00,
		Period: Daily,
	}
	actual := rule.Evaluate(budget, record)
	expected := true
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestFixedBudgetRule_Evaluate_ReturnsFalseWhenNoDailyAmountSet(t *testing.T) {
	budget := Budget{}
	rule := BudgetRule{
		Name: "test-rule",
		Type: FixedRule,
	}
	record := CostRecord{}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestFixedBudgetRule_Evaluate_ReturnsFalseWhenNoMonthlyAmountSet(t *testing.T) {
	budget := Budget{}
	rule := BudgetRule{
		Name: "test-rule",
		Type: FixedRule,
	}
	record := CostRecord{}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestFixedBudgetRule_Evaluate_ReturnsFalseWhenCostValueIsEqualToRuleValue(t *testing.T) {
	budget := Budget{}
	rule := BudgetRule{
		Name:  "test-rule",
		Type:  FixedRule,
		Value: 10.00,
	}
	record := CostRecord{
		Value: 10.00,
	}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestFixedBudgetRule_Evaluate_ReturnsFalseWhenCostValueIsLessThanRuleValue(t *testing.T) {
	budget := Budget{}
	rule := BudgetRule{
		Name:  "test-rule",
		Type:  FixedRule,
		Value: 10.00,
	}
	record := CostRecord{
		Value: 9.00,
	}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

/** Percentage budget rules **/

func TestPercentageBudgetRule_Evaluate_ReturnsTrueWhenCostValueExceedsPercentageOfDailyBudgetUnder100Percent(t *testing.T) {
	budget := Budget{
		DailyAmount: float64Ptr(10.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   PercentageRule,
		Period: Daily,
		Value:  50,
	}
	record := CostRecord{
		Value:  7.00,
		Period: Daily,
	}
	actual := rule.Evaluate(budget, record)
	expected := true
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestPercentageBudgetRule_Evaluate_ReturnsTrueWhenCostValueExceedsPercentageOfMonthlyBudgetUnder100Percent(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   PercentageRule,
		Period: Monthly,
		Value:  50,
	}
	record := CostRecord{
		Value:  60.00,
		Period: Monthly,
	}
	actual := rule.Evaluate(budget, record)
	expected := true
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestPercentageBudgetRule_Evaluate_ReturnsTrueWhenCostValueExceedsPercentageOfDailyBudgetOf100Percent(t *testing.T) {
	budget := Budget{
		DailyAmount: float64Ptr(10.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   PercentageRule,
		Period: Daily,
		Value:  100,
	}
	record := CostRecord{
		Value:  11.00,
		Period: Daily,
	}
	actual := rule.Evaluate(budget, record)
	expected := true
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestPercentageBudgetRule_Evaluate_ReturnsTrueWhenCostValueExceedsPercentageOfMonthlyBudgetOf100Percent(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   PercentageRule,
		Period: Monthly,
		Value:  100,
	}
	record := CostRecord{
		Value:  101.00,
		Period: Monthly,
	}
	actual := rule.Evaluate(budget, record)
	expected := true
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestPercentageBudgetRule_Evaluate_ReturnsTrueWhenCostValueExceedsPercentageOfDailyBudgetOver100Percent(t *testing.T) {
	budget := Budget{
		DailyAmount: float64Ptr(10.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   PercentageRule,
		Period: Daily,
		Value:  135,
	}
	record := CostRecord{
		Value:  15.00,
		Period: Daily,
	}
	actual := rule.Evaluate(budget, record)
	expected := true
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestPercentageBudgetRule_Evaluate_ReturnsTrueWhenCostValueExceedsPercentageOfMonthlyBudgetOver100Percent(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   PercentageRule,
		Period: Monthly,
		Value:  135,
	}
	record := CostRecord{
		Value:  150.00,
		Period: Monthly,
	}
	actual := rule.Evaluate(budget, record)
	expected := true
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestPercentageBudgetRule_Evaluate_ReturnsFalseWhenNoMonthlyAmountSet(t *testing.T) {
	budget := Budget{}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   PercentageRule,
		Period: Monthly,
	}
	record := CostRecord{}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestPercentageBudgetRule_Evaluate_ReturnsFalseWhenNoDailyAmountSet(t *testing.T) {
	budget := Budget{}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   PercentageRule,
		Period: Daily,
	}
	record := CostRecord{}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestPercentageBudgetRule_Evaluate_ReturnsFalseForMonthlyBudgetWhenDailyCost(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
	}
	rule := BudgetRule{
		Name:  "test-rule",
		Type:  PercentageRule,
		Value: 135,
	}
	record := CostRecord{
		Value:  10,
		Period: Daily,
	}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

/** Overrun budget rules **/

func TestOverrunBudgetRule_Evaluate_ReturnsTrueWhenCostValueExceedsDailyBudgetAmountPlusRuleValue(t *testing.T) {
	budget := Budget{
		DailyAmount: float64Ptr(10.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   OverrunRule,
		Period: Daily,
		Value:  5,
	}
	record := CostRecord{
		Value:  17.00,
		Period: Daily,
	}
	actual := rule.Evaluate(budget, record)
	expected := true
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestOverrunBudgetRule_Evaluate_ReturnsTrueWhenCostValueExceedsMonthlyBudgetAmountPlusRuleValue(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   OverrunRule,
		Period: Monthly,
		Value:  20,
	}
	record := CostRecord{
		Value:  125.00,
		Period: Monthly,
	}
	actual := rule.Evaluate(budget, record)
	expected := true
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestOverrunBudgetRule_Evaluate_ReturnsFalseWhenCostValueDoesNotExceedMonthlyBudget(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   OverrunRule,
		Period: Monthly,
		Value:  20,
	}
	record := CostRecord{
		Value: 80.00,
	}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestOverrunBudgetRule_Evaluate_ReturnsFalseWhenCostValueDoesNotExceedDailyBudget(t *testing.T) {
	budget := Budget{
		DailyAmount: float64Ptr(10.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   OverrunRule,
		Period: Daily,
		Value:  10,
	}
	record := CostRecord{
		Value: 8.00,
	}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestOverrunBudgetRule_Evaluate_ReturnsFalseWhenCostValueExceedsMonthlyBudgetButNotOverrun(t *testing.T) {
	budget := Budget{
		MonthlyAmount: float64Ptr(100.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   OverrunRule,
		Period: Monthly,
		Value:  20,
	}
	record := CostRecord{
		Value: 110.00,
	}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestOverrunBudgetRule_Evaluate_ReturnsFalseWhenCostValueExceedsDailyBudgetButNotOverrun(t *testing.T) {
	budget := Budget{
		DailyAmount: float64Ptr(10.00),
	}
	rule := BudgetRule{
		Name:   "test-rule",
		Type:   OverrunRule,
		Period: Daily,
		Value:  5,
	}
	record := CostRecord{
		Value: 13.00,
	}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}

func TestOverrunBudgetRule_Evaluate_ReturnsFalseWhenNoValidChargePeriod(t *testing.T) {
	budget := Budget{}
	rule := BudgetRule{
		Name: "test-rule",
		Type: OverrunRule,
	}
	record := CostRecord{}
	actual := rule.Evaluate(budget, record)
	expected := false
	if actual != expected {
		t.Errorf("Evaluate(budget, record) = %t; want %t", actual, expected)
	}
}
