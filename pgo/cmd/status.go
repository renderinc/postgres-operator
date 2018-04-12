package cmd

/*
 Copyright 2017-2018 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	msgs "github.com/crunchydata/postgres-operator/apiservermsgs"
	"github.com/spf13/cobra"
	"net/http"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "status on Clusters",
	Long: `status displays status of Clusters
				For example:

				pgo status
				pgo status --selector=env=research
				.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("status called")
		if Selector == "" && len(args) == 0 {
			fmt.Println(`You must specify the name of the clusters to test.`)
		} else {
			showStatus(args)
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVarP(&Selector, "selector", "s", "", "The selector to use for cluster filtering ")
	statusCmd.Flags().StringVarP(&OutputFormat, "output", "o", "", "The output format, json is currently supported")
}

func showStatus(args []string) {

	log.Debugf("showStatus called %v\n", args)

	log.Debug("selector is " + Selector)
	if len(args) == 0 && Selector != "" {
		args = make([]string, 1)
		args[0] = "all"
	}

	for _, arg := range args {
		url := APIServerURL + "/status/" + arg + "?selector=" + Selector
		log.Debug(url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal("NewRequest: ", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(BasicAuthUsername, BasicAuthPassword)

		resp, err := httpclient.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return
		}
		log.Debugf("%v\n", resp)
		StatusCheck(resp)

		defer resp.Body.Close()

		var response msgs.StatusResponse

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Printf("%v\n", resp.Body)
			log.Error(err)
			log.Println(err)
			return
		}

		if OutputFormat == "json" {
			b, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Println(string(b))
			return
		}

		if len(response.Results) == 0 {
			fmt.Println("nothing found")
			return
		}

		for _, result := range response.Results {
			fmt.Println("")
			fmt.Printf("cluster : %s ", result.Name)
			if result.Working {
				fmt.Printf("%s\n", GREEN("working"))
			} else {
				fmt.Printf("%s\n", RED("NOT working"))
			}
		}

	}
}
