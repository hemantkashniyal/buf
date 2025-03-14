// Copyright 2020-2021 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/bufbuild/buf/private/pkg/app"
	"github.com/bufbuild/buf/private/pkg/app/appcmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/multierr"
)

const (
	programName = "spdx-go-data"

	pkgFlagName = "package"

	licenseListURL = "https://raw.githubusercontent.com/spdx/license-list-data/v3.12/json/licenses.json"
)

var (
	httpClient = &http.Client{}

	bannedLowercaseIDs = []string{
		"custom",
		"none",
	}
)

func main() {
	appcmd.Main(context.Background(), newCommand())
}

func newCommand() *appcmd.Command {
	flags := newFlags()
	return &appcmd.Command{
		Use:  programName,
		Args: cobra.NoArgs,
		Run: func(ctx context.Context, container app.Container) error {
			return run(ctx, container, flags)
		},
		BindFlags: flags.Bind,
	}
}

type flags struct {
	Pkg string
}

func newFlags() *flags {
	return &flags{}
}

func (f *flags) Bind(flagSet *pflag.FlagSet) {
	flagSet.StringVar(
		&f.Pkg,
		pkgFlagName,
		"",
		"The name of the generated package.",
	)
	_ = cobra.MarkFlagRequired(flagSet, pkgFlagName)
}

func run(ctx context.Context, container app.Container, flags *flags) error {
	licenseInfos, err := getLicenseInfos(ctx)
	if err != nil {
		return err
	}
	data, err := getGolangFileData(ctx, flags.Pkg, licenseInfos)
	if err != nil {
		return err
	}
	_, err = container.Stdout().Write(data)
	return err
}

// sorted by lowercase of ID
// bans "custom", "none"
func getLicenseInfos(ctx context.Context) (_ []*licenseInfo, retErr error) {
	request, err := http.NewRequestWithContext(ctx, "GET", licenseListURL, nil)
	if err != nil {
		return nil, err
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		retErr = multierr.Append(retErr, response.Body.Close())
	}()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected HTTP status code %d to be %d", response.StatusCode, http.StatusOK)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	licenseInfoList := &licenseInfoList{}
	if err := json.Unmarshal(data, licenseInfoList); err != nil {
		return nil, err
	}
	lowercaseIDMap := make(map[string]struct{})
	for _, licenseInfo := range licenseInfoList.LicenseInfos {
		lowercaseID := strings.ToLower(licenseInfo.ID)
		if _, ok := lowercaseIDMap[lowercaseID]; ok {
			return nil, fmt.Errorf("duplicate lowercase ID: %q", lowercaseID)
		}
		lowercaseIDMap[lowercaseID] = struct{}{}
	}
	for _, bannedLowercaseID := range bannedLowercaseIDs {
		if _, ok := lowercaseIDMap[bannedLowercaseID]; ok {
			return nil, fmt.Errorf("banned lowercase ID: %q", bannedLowercaseID)
		}
	}
	sort.Slice(
		licenseInfoList.LicenseInfos,
		func(i int, j int) bool {
			return strings.ToLower(licenseInfoList.LicenseInfos[i].ID) <
				strings.ToLower(licenseInfoList.LicenseInfos[j].ID)
		},
	)
	return licenseInfoList.LicenseInfos, nil
}

func getGolangFileData(
	ctx context.Context,
	pkg string,
	licenseInfos []*licenseInfo,
) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	_, _ = buffer.WriteString(`// Code generated by ` + programName + `. DO NOT EDIT.

package ` + pkg + `

import "strings"

// LicenseInfo is SPDX license information.
//
// See https://spdx.org/licenses.
type LicenseInfo interface {
	// The SPDX identifier.
	//
	// Case-sensitive.
	ID() string
}

// GetLicenseInfo gets the LicenseInfo for the id.
//
// The parameter id is not case-sensitive.
func GetLicenseInfo(id string) (LicenseInfo, bool) {
	licenseInfo, ok := lowercaseIDToLicenseInfo[strings.ToLower(id)]
	return licenseInfo, ok
}

var lowercaseIDToLicenseInfo = map[string]*licenseInfo{
`)

	for _, licenseInfo := range licenseInfos {
		_, _ = buffer.WriteString(`"` + strings.ToLower(licenseInfo.ID) + `": &licenseInfo{
id: "` + licenseInfo.ID + `",
},
`)
	}

	_, _ = buffer.WriteString(`}

type licenseInfo struct {
	id string
}

func (l *licenseInfo) ID() string {
	return l.id
}`)
	return format.Source(buffer.Bytes())
}

type licenseInfoList struct {
	LicenseInfos []*licenseInfo `json:"licenses,omitempty" yaml:"licenses,omitempty"`
}

type licenseInfo struct {
	ID string `json:"licenseId,omitempty" yaml:"licenseId,omitempty"`
}
