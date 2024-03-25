package wrangler

import (
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/gofrontier-com/go-utils/output"
	_ "github.com/gofrontier-com/wrangler/pkg/alerts"
	"github.com/gofrontier-com/wrangler/pkg/core"
	_ "github.com/gofrontier-com/wrangler/pkg/data"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	configFile   string
	outputFmt    string
	disableStdin bool
	verbose      bool
)

func printHeader(title string, version string, tagline string) {
	output.Println()
	output.Println(strings.Repeat("=", 80))
	appFigure := figure.NewFigure(title, "slant", true)
	appFigure.Print()
	output.PrintfInfo("Version: %s", version)
	if len(tagline) > 0 {
		tagFigure := figure.NewFigure(tagline, "pepper", true)
		tagFigure.Print()
	}
	output.Println(strings.Repeat("-", 80))
}

func getDataProviders() []core.CostDataProvider {
	return core.GetServices[core.CostDataProvider]()
}

func getAlertProviders() []core.BudgetAlertService {
	return core.GetServices[core.BudgetAlertService]()
}

func NewRootCmd(version string, commit string, date string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "wrangler",
		Short:   "Wrangler is a command line tool for cost management.",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			core.ConfigureLogger(verbose)
			// TODO: validate input
			printHeader(cmd.Name(), version, "Keeping costs under control")

			log.Info().Msg("Reading configuration...")
			config, err := core.LoadConfig(configFile)
			if err != nil {
				log.Fatal().
					Err(err).
					Str("file", configFile).
					Msg("Unable to load config.")
			}

			cliParams := core.CliParameters{
				ConfigFile:   configFile,
				OutputFmt:    outputFmt,
				DisableStdin: disableStdin,
			}

			records := []core.CostRecord{}
			for _, svc := range getDataProviders() {
				data, err := svc.GetData(cliParams)
				if err != nil {
					log.Fatal().
						Err(err).
						Str("provider", reflect.TypeOf(svc).Name()).
						Msg("Provider failed with error")
				}

				records = append(records, data...)
			}

			log.Info().Msgf("Evaluating %d record(s)...", len(records))
			violations, err := core.CheckBudgets(records, config)
			if err != nil {
				return err
			}

			violationCount := len(violations)
			if violationCount > 0 {
				log.Info().Msgf("%d violation(s) found", violationCount)

				for _, svc := range getAlertProviders() {
					err := svc.HandleViolations(violations)
					if err != nil {
						log.Fatal().
							Err(err).
							Str("provider", reflect.TypeOf(svc).Name()).
							Msg("Provider failed with error")
					}
				}

				log.Fatal().Msg("Failed with violations")
			}

			// TODO: save config

			log.Info().Msg("Check complete.")
			return nil
		},
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if strings.Contains(version, ".") {
		cmd.SetVersionTemplate("{{printf \"v%s\" .Version}}\n")
	} else {
		cmd.SetVersionTemplate("{{printf \"%s\" .Version}}\n")
	}
	cmd.Flags().StringVarP(&configFile, "config", "c", path.Join(wd, ".wrangler.yaml"), "configuration file")
	cmd.Flags().BoolVar(&disableStdin, "no-stdin", false, "disable reading data from stdin")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
	return cmd
}
