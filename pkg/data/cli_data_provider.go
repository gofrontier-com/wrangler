package data

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/gofrontier-com/wrangler/pkg/core"
	"github.com/rs/zerolog/log"
)

type CliCostDataProvider struct {
}

func (p CliCostDataProvider) GetData(params core.CliParameters) ([]core.CostRecord, error) {
	if params.DisableStdin {
		return nil, nil
	}

	reader := csv.NewReader(os.Stdin)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	log.Info().Msg("Reading data from stdin...")
	records := []core.CostRecord{}
	rowNumber := 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Unable to read CSV row")
			continue
		}

		log.Debug().Msgf("Raw data: %s", row)
		record, err := costRecordFromCsv(row)
		if err != nil {
			return records, fmt.Errorf("parsing failed: %+v", err)
		}

		log.Debug().Fields(map[string]interface{}{
			"record": record,
		}).Msgf("Row %d parsed successfully", rowNumber)
		records = append(records, record)
	}
	return records, nil
}

func costRecordFromCsv(data []string) (core.CostRecord, error) {
	if len(data) < 4 {
		return core.CostRecord{}, errors.New("insufficient number of fields")
	}

	log.Debug().Msgf("%d fields detected", len(data))
	log.Debug().
		Str("value", data[0]).
		Str("target_type", reflect.TypeOf(time.Now()).Name()).
		Msg("Parsing timestamp field")

	var ts time.Time
	unixTimestamp, convErr := strconv.ParseInt(data[0], 10, 64)
	if convErr != nil {
		log.Debug().
			Err(convErr).
			Msg("Value cannot be parsed as Unix timestamp")
		unix, err := time.Parse(time.RFC3339, data[0])
		if err != nil {
			return core.CostRecord{}, err
		}

		ts = unix
	} else {
		ts = time.Unix(unixTimestamp, 0)
	}

	log.Debug().
		Str("value", data[3]).
		Str("target_type", reflect.TypeOf(float64(0)).Name()).
		Msg("Parsing value field")
	value, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		return core.CostRecord{}, err
	}

	var curr core.Currency
	var cat string
	if len(data) > 5 {
		log.Debug().
			Str("value", data[3]).
			Str("target_type", reflect.TypeOf(core.Currency("")).Name()).
			Msg("Parsing value field")
		curr = core.Currency(data[4])
		cat = data[5]
	} else if len(data) > 4 {
		curr = core.Currency(data[4])
	}

	record := core.CostRecord{
		ResourceId: data[1],
		Timestamp:  ts,
		Period:     core.Period(data[2]),
		Value:      value,
		Baseline:   -1,
		Currency:   curr,
		Category:   cat,
	}

	return record, nil
}

func init() {
	core.RegisterService("data_provider.cli", CliCostDataProvider{})
}
