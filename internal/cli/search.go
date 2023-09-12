package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/noalino/boursorama-finance-go/internal/lib"
	"github.com/noalino/boursorama-finance-go/internal/utils"
)

type searchFlags struct {
	pretty  bool
	verbose bool
}

func (cli *Cli) RegisterSearchAction() {
	search := cli.NewSubCommand("search", "Search a financial asset\n")
	search.LongDescription(`Search a financial asset by name or ISIN and return the following information:
Symbol, Name, Category, Last price

Usage: quotes search [name | ISIN]`)

	// Flags
	flags := &searchFlags{
		pretty:  false,
		verbose: false,
	}
	search.BoolFlag("pretty", "Display output in a table.", &flags.pretty)
	search.BoolFlag("verbose", "Log more info.", &flags.verbose)

	// Action
	search.Action(func() error {
		otherArgs := search.OtherArgs()
		if len(otherArgs) == 0 {
			return errors.New("too few arguments, please refer to the documentation by using `quotes search -help`")
		}

		query := lib.SearchQuery{Value: otherArgs[0]}

		utils.PrintfOrVoid(flags.verbose, "Searching for '%s'...\n", query.Value)
		assets, err := lib.Search(query)
		if err != nil {
			return err
		}

		if len(assets) == 0 {
			fmt.Println("No result found.")
			return nil
		}

		utils.PrintlnOrVoid(flags.verbose, "Results found:")

		if flags.pretty {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Symbol", "Name", "Market", "Last price"})
			table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
			table.SetCenterSeparator("|")
			table.SetRowLine(true)

			lines := [][]string{}
			for _, asset := range assets {
				line := []string{asset.Symbol, asset.Name, asset.Market, asset.LastPrice}
				lines = append(lines, line)
			}

			table.AppendBulk(lines)
			table.Render()
		} else {
			fmt.Println("symbol,name,market,last price")
			for _, asset := range assets {
				fmt.Printf("%s,%s,%s,%s\n", asset.Symbol, asset.Name, asset.Market, asset.LastPrice)
			}
		}

		return nil
	})
}
