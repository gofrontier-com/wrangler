package serializers

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
)

type TableAlignment int

const (
	AlignLeft   TableAlignment = tablewriter.ALIGN_LEFT
	AlignRight  TableAlignment = tablewriter.ALIGN_RIGHT
	AlignCenter TableAlignment = tablewriter.ALIGN_CENTER
)

type TableOptions struct {
	FirstRowIsHeader bool
	HasBorder        bool
	HeaderAlignment  TableAlignment
	Alignment        TableAlignment
}

func SerializeTable(data [][]string, options TableOptions) string {
	buffer := new(bytes.Buffer)
	table := tablewriter.NewWriter(buffer)

	startIndex := 0
	if options.FirstRowIsHeader {
		table.SetHeader(data[0]) // use first row as header
		startIndex = 1
	}

	table.AppendBulk(data[startIndex:])

	// Set table formatting options
	table.SetBorder(options.HasBorder)
	table.SetHeaderAlignment(int(options.HeaderAlignment))
	table.SetAlignment(int(options.Alignment))
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(false)
	table.Render()

	return buffer.String()
}
