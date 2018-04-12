package statusservice

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
	//"database/sql"
	//"errors"
	//"fmt"
	log "github.com/Sirupsen/logrus"
	crv1 "github.com/crunchydata/postgres-operator/apis/cr/v1"
	"github.com/crunchydata/postgres-operator/apiserver"
	msgs "github.com/crunchydata/postgres-operator/apiservermsgs"
	"github.com/crunchydata/postgres-operator/kubeapi"
	//"github.com/crunchydata/postgres-operator/util"
	//_ "github.com/lib/pq"
	//"github.com/spf13/viper"
	//"k8s.io/api/core/v1"
	//meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/util/validation"
	//"strconv"
	//"strings"
	//"time"
)

func StatusCluster(name, selector string) msgs.StatusResponse {
	var err error

	response := msgs.StatusResponse{}
	response.Results = make([]msgs.StatusDetail, 0)
	response.Status = msgs.Status{Code: msgs.Ok, Msg: ""}

	log.Debug("selector is " + selector)
	if selector == "" && name == "all" {
		log.Debug("selector is empty and name is all")
	} else {
		if selector == "" {
			selector = "name=" + name
		}
	}

	//get a list of matching clusters
	clusterList := crv1.PgclusterList{}
	err = kubeapi.GetpgclustersBySelector(apiserver.RESTClient, &clusterList, selector, apiserver.Namespace)

	if err != nil {
		response.Status.Code = msgs.Error
		response.Status.Msg = err.Error()
		return response
	}

	//loop thru each cluster

	log.Debug("clusters found len is %d\n", len(clusterList.Items))

	for _, c := range clusterList.Items {
		result := msgs.StatusDetail{}
		result.Name = c.Name
		result.Working = true

		response.Results = append(response.Results, result)
	}

	return response
}
