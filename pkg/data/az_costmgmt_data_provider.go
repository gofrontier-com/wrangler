package data

import (
	"github.com/gofrontier-com/wrangler/pkg/core"
)

type AzCostMgmtDataProvider struct {
}

func (p AzCostMgmtDataProvider) GetData(params core.CliParameters) ([]core.CostRecord, error) {
	return []core.CostRecord{}, nil
}

func init() {
	core.RegisterService("data_provider.azcostmgmt", AzCostMgmtDataProvider{})
}
