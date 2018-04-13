package main

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
	log "github.com/Sirupsen/logrus"
	msgs "github.com/crunchydata/postgres-operator/apiservermsgs"
	"io/ioutil"
	"os"
	"time"
)

const OUTPUTPATH = "/pgo-status/status"

func main() {

	debugFlag := os.Getenv("CRUNCHY_DEBUG")
	if debugFlag == "true" {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug flag set to true")
	} else {
		log.Info("debug flag set to false")
	}

	log.Infoln("pgo-status starts")

	go collectStatus()

	for i := 0; i < 3; time.Sleep(3 * time.Second) {
		log.Info("sleeping in main loop")
	}
}

func collectStatus() {
	log.Info("collecting status")

	for {
		time.Sleep(3 * time.Second)
		log.Info("sleeping in collect status loop")

		agg := msgs.StatusAgg{}
		agg.OperatorStartTime = getOperatorStart()
		agg.NumBackups = getNumBackups()
		agg.NumClaims = getNumClaims()
		agg.NumDatabases = getNumDatabases()
		agg.VolumeCap = getVolumeCap()
		agg.LowCaps = getLowCaps()
		agg.DbTags = getDBTags()
		agg.NotReady = getNotReady()

		b, err := json.MarshalIndent(agg, "", " ")
		if err != nil {
			log.Error(err)
		} else {
			writeOutput(b)
		}
	}
}

func writeOutput(b []byte) {
	err := ioutil.WriteFile(OUTPUTPATH, b, 0644)
	if err != nil {
		log.Error(err)
	}
}

func getOperatorStart() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}

func getNumBackups() int {
	return 1
}

func getNumClaims() int {
	return 2
}

func getNumDatabases() int {
	return 3
}

func getVolumeCap() string {
	return "11G"
}

func getLowCaps() []string {
	caps := make([]string, 0)
	caps = append(caps, "db1")
	return caps
}

func getDBTags() map[string]string {
	tags := make(map[string]string)
	tags["one"] = "uno"
	return tags
}

func getNotReady() []string {
	agg := make([]string, 0)
	agg = append(agg, "db1")
	return agg
}
