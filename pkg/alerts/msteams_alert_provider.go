package alerts

import (
	"github.com/gofrontier-com/wrangler/pkg/core"
	"github.com/rs/zerolog/log"
)

type MSTeamsBudgetAlertService struct {
}

func (h MSTeamsBudgetAlertService) HandleViolations(violations []core.BudgetRuleViolation) error {
	log.Trace().
		Int("violation_count", len(violations)).
		Msg("Processing violations")
	return nil
}

func init() {
	core.RegisterService("alert_provider.msteams", MSTeamsBudgetAlertService{})
}
