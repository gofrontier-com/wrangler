package alerts

import (
	"github.com/gofrontier-com/go-utils/output"
	"github.com/gofrontier-com/wrangler/pkg/core"
	"github.com/gofrontier-com/wrangler/pkg/serializers"
	"github.com/rs/zerolog/log"
)

type CliBudgetAlertService struct {
}

func (h CliBudgetAlertService) HandleViolations(violations []core.BudgetRuleViolation) error {
	log.Trace().
		Int("violation_count", len(violations)).
		Msg("Processing violations")
	rows := [][]string{
		{"Resource ID", "Rule name", "Condition", "Date", "Budget amount", "Actual amount"},
	}
	for _, v := range violations {
		rows = append(rows, []string{
			v.ResourceId,
			v.Name,
			v.Description,
			v.Date,
			core.FormatCurrency(v.BudgetAmount, v.Currency),
			core.FormatCurrency(v.ActualAmount, v.Currency),
		})
	}

	output.PrintlnWarn(
		serializers.SerializeTable(rows, serializers.TableOptions{
			FirstRowIsHeader: true,
			HasBorder:        true,
			HeaderAlignment:  serializers.AlignLeft,
			Alignment:        serializers.AlignLeft,
		}),
	)

	return nil
}

func init() {
	core.RegisterService("alert_provider.cli", CliBudgetAlertService{})
}
