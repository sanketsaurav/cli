package dictionaryitem

import (
	"fmt"
	"io"

	"github.com/fastly/cli/pkg/common"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/fastly"
)

// ListCommand calls the Fastly API to list dictionary items.
type ListCommand struct {
	common.Base
	manifest manifest.Data

	// required
	DictionaryID string
}

// NewListCommand returns a usable command registered under the parent.
func NewListCommand(parent common.Registerer, globals *config.Data) *ListCommand {
	var c ListCommand
	c.Globals = globals
	c.manifest.File.Read(manifest.Filename)

	c.CmdClause = parent.Command("list", "List items in an dictionary")

	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("dictionary-id", "The ID of the dictionary").Required().StringVar(&c.DictionaryID)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *ListCommand) createInput() (*fastly.ListDictionaryItemsInput, error) {
	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	input := fastly.ListDictionaryItemsInput{
		Service:    serviceID,
		Dictionary: c.DictionaryID,
	}

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *ListCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	items, err := c.Globals.Client.ListDictionaryItems(input)
	if err != nil {
		return err
	}

	if !c.Globals.Verbose() {
		tw := text.NewTable(out)
		tw.AddHeader("SERVICE", "DICTIONARY ID", "KEY", "VALUE")
		for _, i := range items {
			tw.AddLine(i.ServiceID, i.DictionaryID, i.ItemKey, i.ItemValue)
		}
		tw.Print()
		return nil
	}

	fmt.Fprintf(out, "Service ID: %s\n", input.Service)
	fmt.Fprintf(out, "Dictionary ID: %s\n", input.Dictionary)
	for i, item := range items {
		fmt.Fprintf(out, "\tItem %d/%d\n", i+1, len(items))
		fmt.Fprintf(out, "\t\tKey: %s\n", item.ItemKey)
		fmt.Fprintf(out, "\t\tValue: %s\n", item.ItemValue)

		if item.CreatedAt != nil {
			fmt.Fprintf(out, "\t\tCreated (UTC): %s\n", item.CreatedAt.UTC().Format(common.TimeFormat))
		}
		if item.UpdatedAt != nil {
			fmt.Fprintf(out, "\t\tLast edited (UTC): %s\n", item.UpdatedAt.UTC().Format(common.TimeFormat))
		}
		if item.DeletedAt != nil {
			fmt.Fprintf(out, "\t\tDeleted (UTC): %s\n", item.DeletedAt.UTC().Format(common.TimeFormat))
		}
	}
	fmt.Fprintln(out)

	return nil
}
