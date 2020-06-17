package dictionary

import (
	"io"

	"github.com/fastly/cli/pkg/common"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/fastly"
)

// UpdateCommand calls the Fastly API to update dictionaries.
type UpdateCommand struct {
	common.Base
	manifest manifest.Data

	// required
	DictionaryName string // Can't shaddow common.Base method Name().
	Version        int

	// optional
	NewName   common.OptionalString
	WriteOnly common.OptionalBool
}

// NewUpdateCommand returns a usable command registered under the parent.
func NewUpdateCommand(parent common.Registerer, globals *config.Data) *UpdateCommand {
	var c UpdateCommand
	c.Globals = globals

	c.CmdClause = parent.Command("update", "Update a dictionary on a Fastly service version")

	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("version", "Number of service version").Required().IntVar(&c.Version)
	c.CmdClause.Flag("name", "Dictionary name").Short('n').Required().StringVar(&c.DictionaryName)
	c.CmdClause.Flag("new-name", "New dictionary name").Action(c.NewName.Set).StringVar(&c.NewName.Value)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used
// by the API client library.
func (c *UpdateCommand) createInput() (*fastly.UpdateDictionaryInput, error) {
	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	dictionary, err := c.Globals.Client.GetDictionary(&fastly.GetDictionaryInput{
		Service: serviceID,
		Version: c.Version,
		Name:    c.DictionaryName,
	})
	if err != nil {
		return nil, err
	}

	input := fastly.UpdateDictionaryInput{
		Service: dictionary.ServiceID,
		Version: dictionary.Version,
		Name:    dictionary.Name,
	}

	if c.NewName.Valid {
		input.NewName = c.NewName.Value
	}

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *UpdateCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	d, err := c.Globals.Client.UpdateDictionary(input)
	if err != nil {
		return err
	}

	text.Success(out, "Updated dictionary %s (service %s version %d)", d.Name, d.ServiceID, d.Version)
	return nil
}
