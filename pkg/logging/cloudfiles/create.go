package cloudfiles

import (
	"io"

	"github.com/fastly/cli/pkg/common"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/fastly"
)

// CreateCommand calls the Fastly API to create Cloudfiles logging endpoints.
type CreateCommand struct {
	common.Base
	manifest manifest.Data

	// required
	EndpointName string // Can't shaddow common.Base method Name().
	Token        string
	Version      int
	User         string
	AccessKey    string
	BucketName   string

	// optional
	Path              common.OptionalString
	Region            common.OptionalString
	PublicKey         common.OptionalString
	Period            common.OptionalUint
	GzipLevel         common.OptionalUint
	MessageType       common.OptionalString
	TimestampFormat   common.OptionalString
	Format            common.OptionalString
	FormatVersion     common.OptionalUint
	ResponseCondition common.OptionalString
	Placement         common.OptionalString
}

// NewCreateCommand returns a usable command registered under the parent.
func NewCreateCommand(parent common.Registerer, globals *config.Data) *CreateCommand {
	var c CreateCommand

	c.Globals = globals
	c.manifest.File.Read(manifest.Filename)
	c.CmdClause = parent.Command("create", "Create a Cloudfiles logging endpoint on a Fastly service version").Alias("add")

	c.CmdClause.Flag("name", "The name of the Cloudfiles logging object. Used as a primary key for API access").Short('n').Required().StringVar(&c.EndpointName)
	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("version", "Number of service version").Required().IntVar(&c.Version)

	c.CmdClause.Flag("user", "The username for your Cloudfile account").Required().StringVar(&c.User)
	c.CmdClause.Flag("access-key", "Your Cloudfile account access key").Required().StringVar(&c.AccessKey)
	c.CmdClause.Flag("bucket", "The name of your Cloudfiles container").Required().StringVar(&c.BucketName)

	c.CmdClause.Flag("path", "The path to upload logs to").Action(c.Path.Set).StringVar(&c.Path.Value)
	c.CmdClause.Flag("region", "The region to stream logs to. One of: DFW-Dallas, ORD-Chicago, IAD-Northern Virginia, LON-London, SYD-Sydney, HKG-Hong Kong").Action(c.Region.Set).StringVar(&c.Region.Value)
	c.CmdClause.Flag("placement", "Where in the generated VCL the logging call should be placed, overriding any format_version default. Can be none or waf_debug").Action(c.Placement.Set).StringVar(&c.Placement.Value)
	c.CmdClause.Flag("period", "How frequently log files are finalized so they can be available for reading (in seconds, default 3600)").Action(c.Period.Set).UintVar(&c.Period.Value)
	c.CmdClause.Flag("gzip-level", "What level of GZIP encoding to have when dumping logs (default 0, no compression)").Action(c.GzipLevel.Set).UintVar(&c.GzipLevel.Value)
	c.CmdClause.Flag("format", "Apache style log formatting").Action(c.Format.Set).StringVar(&c.Format.Value)
	c.CmdClause.Flag("format-version", "The version of the custom logging format used for the configured endpoint. Can be either 2 (default) or 1").Action(c.FormatVersion.Set).UintVar(&c.FormatVersion.Value)
	c.CmdClause.Flag("response-condition", "The name of an existing condition in the configured endpoint, or leave blank to always execute").Action(c.ResponseCondition.Set).StringVar(&c.ResponseCondition.Value)
	c.CmdClause.Flag("message-type", "How the message should be formatted. One of: classic (default), loggly, logplex or blank").Action(c.MessageType.Set).StringVar(&c.MessageType.Value)
	c.CmdClause.Flag("timestamp-format", `strftime specified timestamp formatting (default "%Y-%m-%dT%H:%M:%S.000")`).Action(c.TimestampFormat.Set).StringVar(&c.TimestampFormat.Value)
	c.CmdClause.Flag("public-key", "A PGP public key that Fastly will use to encrypt your log files before writing them to disk").Action(c.PublicKey.Set).StringVar(&c.PublicKey.Value)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used by the API client library.
func (c *CreateCommand) createInput() (*fastly.CreateCloudfilesInput, error) {
	var input fastly.CreateCloudfilesInput

	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	input.Service = serviceID
	input.Version = c.Version
	input.Name = fastly.String(c.EndpointName)
	input.User = fastly.String(c.User)
	input.AccessKey = fastly.String(c.AccessKey)
	input.BucketName = fastly.String(c.BucketName)

	if c.Path.Valid {
		input.Path = fastly.String(c.Path.Value)
	}

	if c.Region.Valid {
		input.Region = fastly.String(c.Region.Value)
	}

	if c.Placement.Valid {
		input.Placement = fastly.String(c.Placement.Value)
	}

	if c.Period.Valid {
		input.Period = fastly.Uint(c.Period.Value)
	}

	if c.GzipLevel.Valid {
		input.GzipLevel = fastly.Uint(c.GzipLevel.Value)
	}

	if c.Format.Valid {
		input.Format = fastly.String(c.Format.Value)
	}

	if c.FormatVersion.Valid {
		input.FormatVersion = fastly.Uint(c.FormatVersion.Value)
	}

	if c.ResponseCondition.Valid {
		input.ResponseCondition = fastly.String(c.ResponseCondition.Value)
	}

	if c.MessageType.Valid {
		input.MessageType = fastly.String(c.MessageType.Value)
	}

	if c.TimestampFormat.Valid {
		input.TimestampFormat = fastly.String(c.TimestampFormat.Value)
	}

	if c.PublicKey.Valid {
		input.PublicKey = fastly.String(c.PublicKey.Value)
	}

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *CreateCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	d, err := c.Globals.Client.CreateCloudfiles(input)
	if err != nil {
		return err
	}

	text.Success(out, "Created Cloudfiles logging endpoint %s (service %s version %d)", d.Name, d.ServiceID, d.Version)
	return nil
}
