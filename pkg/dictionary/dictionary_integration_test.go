package dictionary_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/fastly/cli/pkg/app"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/mock"
	"github.com/fastly/cli/pkg/testutil"
	"github.com/fastly/cli/pkg/update"
	"github.com/fastly/go-fastly/fastly"
)

func TestDictionaryCreate(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:      []string{"dictionary", "create", "--service-id", "123", "--version", "1"},
			wantError: "error parsing arguments: required flag --name not provided",
		},
		{
			args:       []string{"dictionary", "create", "--service-id", "123", "--version", "1", "--name", "dictionary"},
			api:        mock.API{CreateDictionaryFn: createDictionaryOK},
			wantOutput: "Created dictionary dictionary (service 123 version 1)",
		},
		{
			args:      []string{"dictionary", "create", "--service-id", "123", "--version", "1", "--name", "dictionary"},
			api:       mock.API{CreateDictionaryFn: createDictionaryError},
			wantError: errTest.Error(),
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertStringContains(t, out.String(), testcase.wantOutput)
		})
	}
}

func TestDictionaryList(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:       []string{"dictionary", "list", "--service-id", "123", "--version", "1"},
			api:        mock.API{ListDictionariesFn: listDictionariesOK},
			wantOutput: listDictionarysShortOutput,
		},
		{
			args:       []string{"dictionary", "list", "--service-id", "123", "--version", "1", "--verbose"},
			api:        mock.API{ListDictionariesFn: listDictionariesOK},
			wantOutput: listDictionarysVerboseOutput,
		},
		{
			args:       []string{"dictionary", "list", "--service-id", "123", "--version", "1", "-v"},
			api:        mock.API{ListDictionariesFn: listDictionariesOK},
			wantOutput: listDictionarysVerboseOutput,
		},
		{
			args:       []string{"dictionary", "--verbose", "list", "--service-id", "123", "--version", "1"},
			api:        mock.API{ListDictionariesFn: listDictionariesOK},
			wantOutput: listDictionarysVerboseOutput,
		},
		{
			args:      []string{"dictionary", "list", "--service-id", "123", "--version", "1"},
			api:       mock.API{ListDictionariesFn: listDictionariesError},
			wantError: errTest.Error(),
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertString(t, testcase.wantOutput, out.String())
		})
	}
}

func TestDictionaryDescribe(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:      []string{"dictionary", "describe", "--service-id", "123", "--version", "1"},
			wantError: "error parsing arguments: required flag --name not provided",
		},
		{
			args:      []string{"dictionary", "describe", "--service-id", "123", "--version", "1", "--name", "dictionary"},
			api:       mock.API{GetDictionaryFn: getDictionaryError},
			wantError: errTest.Error(),
		},
		{
			args:       []string{"dictionary", "describe", "--service-id", "123", "--version", "1", "--name", "dictionary"},
			api:        mock.API{GetDictionaryFn: getDictionaryOK},
			wantOutput: describeDictionaryOutput,
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertString(t, testcase.wantOutput, out.String())
		})
	}
}

func TestDictionaryUpdate(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:      []string{"dictionary", "update", "--service-id", "123", "--version", "1", "--new-name", "new-dictionary"},
			wantError: "error parsing arguments: required flag --name not provided",
		},
		{
			args: []string{"dictionary", "update", "--service-id", "123", "--version", "1", "--name", "dictionary", "--new-name", "new-dictionary"},
			api: mock.API{
				GetDictionaryFn:    getDictionaryError,
				UpdateDictionaryFn: updateDictionaryOK,
			},
			wantError: errTest.Error(),
		},
		{
			args: []string{"dictionary", "update", "--service-id", "123", "--version", "1", "--name", "dictionary", "--new-name", "new-dictionary"},
			api: mock.API{
				GetDictionaryFn:    getDictionaryOK,
				UpdateDictionaryFn: updateDictionaryError,
			},
			wantError: errTest.Error(),
		},
		{
			args: []string{"dictionary", "update", "--service-id", "123", "--version", "1", "--name", "dictionary", "--new-name", "new-dictionary"},
			api: mock.API{
				GetDictionaryFn:    getDictionaryOK,
				UpdateDictionaryFn: updateDictionaryOK,
			},
			wantOutput: "Updated dictionary new-dictionary (service 123 version 1)",
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertStringContains(t, out.String(), testcase.wantOutput)
		})
	}
}

func TestDictionaryDelete(t *testing.T) {
	for _, testcase := range []struct {
		args       []string
		api        mock.API
		wantError  string
		wantOutput string
	}{
		{
			args:      []string{"dictionary", "delete", "--service-id", "123", "--version", "1"},
			wantError: "error parsing arguments: required flag --name not provided",
		},
		{
			args:      []string{"dictionary", "delete", "--service-id", "123", "--version", "1", "--name", "dictionary"},
			api:       mock.API{DeleteDictionaryFn: deleteDictionaryError},
			wantError: errTest.Error(),
		},
		{
			args:       []string{"dictionary", "delete", "--service-id", "123", "--version", "1", "--name", "dictionary"},
			api:        mock.API{DeleteDictionaryFn: deleteDictionaryOK},
			wantOutput: "Deleted dictionary dictionary (service 123 version 1)",
		},
	} {
		t.Run(strings.Join(testcase.args, " "), func(t *testing.T) {
			var (
				args                           = testcase.args
				env                            = config.Environment{}
				file                           = config.File{}
				appConfigFile                  = "/dev/null"
				clientFactory                  = mock.APIClient(testcase.api)
				httpClient                     = http.DefaultClient
				versioner     update.Versioner = nil
				in            io.Reader        = nil
				out           bytes.Buffer
			)
			err := app.Run(args, env, file, appConfigFile, clientFactory, httpClient, versioner, in, &out)
			testutil.AssertErrorContains(t, err, testcase.wantError)
			testutil.AssertStringContains(t, out.String(), testcase.wantOutput)
		})
	}
}

var errTest = errors.New("fixture error")

func createDictionaryOK(i *fastly.CreateDictionaryInput) (*fastly.Dictionary, error) {
	s := fastly.Dictionary{
		ServiceID: i.Service,
		Version:   i.Version,
		Name:      i.Name,
	}

	return &s, nil
}

func createDictionaryError(i *fastly.CreateDictionaryInput) (*fastly.Dictionary, error) {
	return nil, errTest
}

func listDictionariesOK(i *fastly.ListDictionariesInput) ([]*fastly.Dictionary, error) {
	return []*fastly.Dictionary{
		{
			ServiceID: i.Service,
			Version:   i.Version,
			ID:        "1",
			Name:      "dictionary1",
			WriteOnly: false,
		},
		{
			ServiceID: i.Service,
			Version:   i.Version,
			ID:        "2",
			Name:      "dictionary2",
			WriteOnly: true,
		},
	}, nil
}

func listDictionariesError(i *fastly.ListDictionariesInput) ([]*fastly.Dictionary, error) {
	return nil, errTest
}

var listDictionarysShortOutput = strings.TrimSpace(`
SERVICE  VERSION  ID  NAME         WRITE ONLY
123      1        1   dictionary1  false
123      1        2   dictionary2  true
`) + "\n"

var listDictionarysVerboseOutput = strings.TrimSpace(`
Fastly API token not provided
Fastly API endpoint: https://api.fastly.com
Service ID: 123
Version: 1
	Dictionary 1/2
		ID: 1
		Name: dictionary1
		Write only: false
	Dictionary 2/2
		ID: 2
		Name: dictionary2
		Write only: true
`) + "\n\n"

func getDictionaryOK(i *fastly.GetDictionaryInput) (*fastly.Dictionary, error) {
	return &fastly.Dictionary{
		ServiceID: i.Service,
		Version:   i.Version,
		ID:        "1",
		Name:      "dictionary",
		WriteOnly: false,
	}, nil
}

func getDictionaryError(i *fastly.GetDictionaryInput) (*fastly.Dictionary, error) {
	return nil, errTest
}

var describeDictionaryOutput = strings.TrimSpace(`
Service ID: 123
Version: 1
ID: 1
Name: dictionary
Write only: false
`) + "\n"

func updateDictionaryOK(i *fastly.UpdateDictionaryInput) (*fastly.Dictionary, error) {
	return &fastly.Dictionary{
		ServiceID: i.Service,
		Version:   i.Version,
		ID:        "1",
		Name:      "new-dictionary",
		WriteOnly: false,
	}, nil
}

func updateDictionaryError(i *fastly.UpdateDictionaryInput) (*fastly.Dictionary, error) {
	return nil, errTest
}

func deleteDictionaryOK(i *fastly.DeleteDictionaryInput) error {
	return nil
}

func deleteDictionaryError(i *fastly.DeleteDictionaryInput) error {
	return errTest
}
