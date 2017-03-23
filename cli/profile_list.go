// Copyright 2017 CoreOS Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/coreos/torcx/pkg/torcx"
)

var (
	cmdProfileList = &cobra.Command{
		Use:   "list",
		Short: "list available profiles",
		RunE:  runProfileList,
	}
)

func init() {
	cmdProfile.AddCommand(cmdProfileList)
}

func runProfileList(cmd *cobra.Command, args []string) error {
	var err error

	commonCfg, err := fillCommonRuntime()
	if err != nil {
		return errors.Wrap(err, "common configuration failed")
	}
	profileCfg, err := fillProfileRuntime(commonCfg)
	if err != nil {
		return errors.Wrap(err, "profile configuration failed")
	}

	profileDirs := []string{
		filepath.Join(torcx.VENDOR_DIR, "profiles.d"),
		filepath.Join(commonCfg.ConfDir, "profiles.d"),
	}
	localProfiles, err := torcx.ListProfiles(profileDirs)
	if err != nil {
		return errors.Wrap(err, "profiles listing failed")
	}
	profNames := make([]string, 0, len(localProfiles))
	for k := range localProfiles {
		profNames = append(profNames, k)
	}

	var curName, curPath *string
	if profileCfg.CurrentProfileName != "" {
		curName = &profileCfg.CurrentProfileName
	}
	if profileCfg.CurrentProfilePath != "" {
		curPath = &profileCfg.CurrentProfilePath
	}

	profListOut := ProfileList{
		Kind: TorcxProfileListV0,
		Value: profileList{
			CurrentProfileName: curName,
			CurrentProfilePath: curPath,
			NextProfileName:    profileCfg.NextProfile,
			Profiles:           profNames,
		},
	}

	jsonOut := json.NewEncoder(os.Stdout)
	jsonOut.SetIndent("", "  ")
	err = jsonOut.Encode(profListOut)

	return err
}