// Copyright © 2017 Aqua Security Software Ltd. <info@aquasec.com>
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

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// listCmd gets  list of available manifesto data for this image
var listCmd = &cobra.Command{
	Use:   "list [IMAGE]",
	Short: "List currently stored metadata for the container image",
	Long:  `Display a list of the metadata stored for the specified container image.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		name := args[0]
		repoName, imageName := repoAndTaggedNames(name)
		metadataImageName := imageNameForManifest(repoName)

		raw, err := dockerGetData(metadataImageName)
		if err != nil {
			fmt.Printf("No manifesto data stored for image '%s'\n", imageName)
			os.Exit(1)
		}
		var mml MetadataManifestList
		json.Unmarshal(raw, &mml)

		found := false
		for _, v := range mml.Tags {
			// TODO: for now this is checking against the image name including the tag, but since this
			// can be moved we should really be finding the SHA for the tag and using that as the key in
			// the manifesto data.
			if v.Tag == imageName {
				fmt.Printf("Metadata types stored for image '%s':\n", imageName)
				found = true
				for _, m := range v.MetadataManifest {
					fmt.Printf("    %s\n", m.Type)
				}
			}
		}

		if !found {
			fmt.Printf("No metadata stored for image '%s'\n", imageName)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
