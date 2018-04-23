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
	"flag"
	log "github.com/Sirupsen/logrus"
	msgs "github.com/crunchydata/postgres-operator/apiservermsgs"
	"github.com/crunchydata/postgres-operator/kubeapi"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

var kubeconfig = flag.String("kubeconfig", "./config", "path to kubeconfig")

const OUTPUTPATH = "/pgo-status/status"
const SLEEP = 30

var Clientset *kubernetes.Clientset
var Namespace string

func main() {

	debugFlag := os.Getenv("CRUNCHY_DEBUG")
	if debugFlag == "true" {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug flag set to true")
	} else {
		log.Info("debug flag set to false")
	}
	Namespace = os.Getenv("NAMESPACE")
	if Namespace == "" {
		log.Error("namespace env var not set")
		os.Exit(2)
	} else {
		log.Info("namespace set to " + Namespace)
	}

	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	log.Infoln("pgo-status starts")

	go collectStatus()

	for i := 0; i < 3; time.Sleep(SLEEP * time.Second) {
		log.Info("sleeping in main loop")
	}
}

func collectStatus() {
	log.Info("collecting status")

	for {
		time.Sleep(SLEEP * time.Second)
		log.Info("sleeping in collect status loop")

		agg := msgs.StatusDetail{}
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
	pods, err := kubeapi.GetPods(Clientset, "name=postgres-operator", Namespace)
	if err != nil {
		log.Error(err)
		return "error"
	}

	for _, p := range pods.Items {
		if string(p.Status.Phase) == "Running" {
			log.Info("found a postgres-operator pod running phase")
			log.Info("start time is " + p.Status.StartTime.String())
			return p.Status.StartTime.String()
		}
	}
	log.Error("a running postgres-operator pod is not found")

	//postgres-operator pod status.startTime
	//t := time.Now()
	//return t.Format(time.RFC3339)
	return "error"
}

func getNumBackups() int {
	//count the number of Jobs with pgbackup=true and completionTime not nil
	return 1
}

func getNumClaims() int {
	//count number of PVCs with pgremove=true
	return 2
}

func getNumDatabases() int {
	//count number of Deployments with pg-cluster
	return 3
}

func getVolumeCap() string {
	//sum all PVCs storage capacity
	return "11G"
}

//this might move to the df command instead
func getLowCaps() []string {
	caps := make([]string, 0)
	caps = append(caps, "db1")
	return caps
}

func getDBTags() map[string]string {
	//count all pods with pg-cluster, sum by image tag value
	tags := make(map[string]string)
	tags["one"] = "uno"
	return tags
}

func getNotReady() []string {
	//show all pods with pg-cluster that have status.Phase of not Running
	agg := make([]string, 0)
	agg = append(agg, "db1")
	return agg
}
