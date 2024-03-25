package core

import "time"

type CliParameters struct {
	ConfigFile   string
	OutputFmt    string
	DisableStdin bool
}

type BudgetAlertService interface {
	HandleViolations(violations []BudgetRuleViolation) error
}

type BudgetRuleViolation struct {
	ResourceId   string
	Name         string
	Description  string
	Date         string
	BudgetAmount float64
	ActualAmount float64
	Currency     Currency
}

type Currency string
type Period string

const (
	Monthly Period   = "monthly"
	Daily   Period   = "daily"
	GBP     Currency = "GBP"
	USD     Currency = "USD"
	EUR     Currency = "EUR"
)

var currencies []Currency = []Currency{GBP, USD, EUR}

var currencySymbols = map[Currency]string{
	GBP: "£",
	USD: "$",
	EUR: "€",
}

type CostDataProvider interface {
	GetData(params CliParameters) ([]CostRecord, error)
}

type CostRecord struct {
	ResourceId string
	Timestamp  time.Time
	Period     Period
	Value      float64
	Currency   Currency
	Baseline   float64
	Category   string
}
