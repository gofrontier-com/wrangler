package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type BudgetRuleType string

const (
	PercentageRule BudgetRuleType = "percentage"
	FixedRule      BudgetRuleType = "fixed"
	OverrunRule    BudgetRuleType = "overrun"
)

type Budget struct {
	ResourceId    string        `mapstructure:"resource_id"`
	Category      *string       `mapstructure:"category,omitempty"`
	MonthlyAmount *float64      `mapstructure:"monthly_amount,omitempty"`
	DailyAmount   *float64      `mapstructure:"daily_amount,omitempty"`
	Currency      Currency      `mapstructure:"currency"`
	Rules         *[]BudgetRule `mapstructure:"rules,omitempty"`
}

func (b *Budget) HasAmount() bool {
	return b.MonthlyAmount != nil || b.DailyAmount != nil
}

func (b *Budget) GetAmount(period Period) *float64 {
	if period == Monthly {
		return b.MonthlyAmount
	} else if period == Daily {
		return b.DailyAmount
	} else {
		return nil
	}
}

type BudgetRule struct {
	Name       string         `mapstructure:"name,omitempty"`
	Type       BudgetRuleType `mapstructure:"type"`
	Value      float64        `mapstructure:"value"`
	Period     Period         `mapstructure:"period,omitempty"`
	Currency   Currency       `mapstructure:"currency,omitempty"`
	Categories []string       `mapstructure:"categories,omitempty"`
}

func (r *BudgetRule) Evaluate(budget Budget, record CostRecord) bool {
	log.Trace().
		Fields(map[string]interface{}{
			"rule": r,
		}).Msg("")
	// if rule is not scoped to a specific charge period then apply rule to any charge
	period := r.Period
	if period == "" {
		period = record.Period
	}
	if period == "" {
		log.Warn().
			Str("rule", r.Name).
			Msg("Unable to determine charge period, rule will be skipped")
		return false
	} else if period != record.Period {
		log.Warn().
			Str("rule", r.Name).
			Str("rule_period", string(period)).
			Str("charge_period", string(record.Period)).
			Msg("Unable to determine charge period, rule will be skipped")
		return false
	}

	amountPtr := budget.GetAmount(period)
	if r.Type != FixedRule && amountPtr == nil {
		log.Warn().
			Str("rule", r.Name).
			Str("rule_period", string(period)).
			Msgf("No %s budget found, rule will be skipped", period)
		return false
	}

	if r.Type == FixedRule {
		return record.Value >= r.Value
	} else if r.Type == PercentageRule {
		amount := *budget.GetAmount(period)
		return record.Value >= (r.Value/100)*amount
	} else if r.Type == OverrunRule {
		amount := *budget.GetAmount(period)
		return (record.Value - amount) >= r.Value
	} else {
		return false
	}
}

func (r *BudgetRule) GetDescription() string {
	switch r.Type {
	case PercentageRule:
		return fmt.Sprintf("actual amount >= %.2f%% of budget", r.Value)
	case FixedRule:
		return fmt.Sprintf("actual amount >= %s", FormatCurrency(r.Value, r.Currency))
	case OverrunRule:
		return fmt.Sprintf("actual amount >= %s overrun", FormatCurrency(r.Value, r.Currency))
	default:
		return ""
	}
}

func CheckBudgets(data []CostRecord, config Config) ([]BudgetRuleViolation, error) {
	violations := []BudgetRuleViolation{}
	for _, record := range data {
		log.Debug().
			Fields(map[string]interface{}{
				"record": record,
			}).
			Msg("")
		budgetComparer := func(b Budget) bool {
			return strings.EqualFold(b.ResourceId, record.ResourceId) &&
				(b.Currency == "" || b.Currency == record.Currency)
		}
		budget, budgetExists := FindInSlice[Budget](*config.Budgets, budgetComparer)
		if !budgetExists && record.Baseline > 0 {
			log.Trace().
				Str("resource_id", record.ResourceId).
				Msgf("No %s budget defined for resource", record.Period)
			budget = Budget{
				ResourceId: record.ResourceId,
				Currency:   record.Currency,
			}
			if record.Period == Monthly {
				budget.MonthlyAmount = &record.Baseline
			} else if record.Period == Daily {
				budget.DailyAmount = &record.Baseline
			}
			log.Trace().
				Fields(map[string]interface{}{
					"monthly_amount": budget.GetAmount(Monthly),
					"daily_amount":   budget.GetAmount(Daily),
				}).Msg("Budget estimated from baseline")
		}
		if budget.HasAmount() {
			if budget.Rules != nil {
				rules := *budget.Rules
				log.Trace().
					Str("resource_id", record.ResourceId).
					Int("rule_count", len(rules)).
					Msg("Evaluating local rules")
				localViolations, err := getViolations(rules, budget, record)
				if err != nil {
					return nil, fmt.Errorf("error evaluating local rules %+v", err)
				}
				violations = append(violations, localViolations...)
			} else {
				log.Trace().Msg("No local rules found.")
			}

			if config.Rules != nil {
				rules := *config.Rules
				log.Trace().
					Str("resource_id", record.ResourceId).
					Int("rule_count", len(rules)).
					Msg("Evaluating global rules")
				globalViolations, err := getViolations(rules, budget, record)
				if err != nil {
					return nil, err
				}
				violations = append(violations, globalViolations...)
			} else {
				log.Trace().Msg("No global rules found.")
			}
		}
	}

	return violations, nil
}

func getViolations(rules []BudgetRule, budget Budget, record CostRecord) ([]BudgetRuleViolation, error) {
	violations := []BudgetRuleViolation{}
	for _, rule := range rules {
		shouldEvaluate := true
		if len(rule.Categories) > 0 {
			_, shouldEvaluate = FindInSlice[string](rule.Categories, func(category string) bool {
				return strings.EqualFold(record.Category, category)
			})
			if shouldEvaluate {
				log.Trace().
					Str("resource_id", record.ResourceId).
					Str("category_match", record.Category).
					Msg("Rule matched by category")
			}
		}
		shouldEvaluate = shouldEvaluate && (rule.Currency == "" || rule.Currency == record.Currency)
		if shouldEvaluate && rule.Evaluate(budget, record) {
			violations = append(violations, BudgetRuleViolation{
				ResourceId:   record.ResourceId,
				Name:         rule.Name,
				Description:  rule.GetDescription(),
				Date:         record.Timestamp.Format(time.DateOnly),
				BudgetAmount: *budget.GetAmount(record.Period),
				ActualAmount: record.Value,
				Currency:     record.Currency,
			})
		}
	}
	return violations, nil
}
