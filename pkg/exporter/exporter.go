/*
Copyright © 2021 Loïc Yavercovski

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
package exporter

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var labelNames = []string {
	"computerName",
}


func getData() string {
	client := &http.Client{}
	var login string = viper.GetString("username")
	var token string = viper.GetString("password")
	var url string = viper.GetString("address")
	req, err := http.NewRequest("GET", url+"/computer/api/json", nil)
	if err != nil {
		log.Println(err)
		return "unknown"
	}
	if !viper.GetBool("disable-authentication") {
		req.SetBasicAuth(login,token)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return "unknown"
		}
		return string(body)
	}
	log.Printf("getData: failed to get jenkins computer data => %s\n", resp.Status)
	return "unknown"
}

func promWatchJenkinsNodes(registry prometheus.Registry) {
	maintenanceState := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jenkins_node_maintenance_status",
		Help: "The maintenance status of a jenkins computer node 0:ONLINE 1:MAINTENANCE 2:OFFLINE",
	}, labelNames)
	busyState := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jenkins_node_busy_status",
		Help: "The busy status of a jenkins computer node 0:IDLE 1:BUSY",
	}, labelNames)
	failGetData := promauto.NewCounter(prometheus.CounterOpts{
		Name: "jenkins_node_exporter_failure",
		Help: "The number of failure to get/parse api data since startup",
	})
	registry.MustRegister(maintenanceState)
	registry.MustRegister(busyState)
	registry.MustRegister(failGetData)

	var listOfKnownComputer = make(map[string]bool)

	go func() {
		for {
			apiJsonData := getData()
			var computerList computerList
			if err := json.Unmarshal([]byte(apiJsonData), &computerList); err != nil {
				log.Printf("promWatchJenkinsNodes: error parsing JSON (%s) => retrying in 5s\n", err.Error())
				failGetData.Inc()
				time.Sleep(5 * time.Second)
				continue
			}
			for _, computer := range computerList.Computers {
				_, found := listOfKnownComputer[computer.DisplayName]
				if !found {
					listOfKnownComputer[computer.DisplayName] = true
				}

				maintenanceState.WithLabelValues(computer.GetLabelValues()...).Set(computer.GetMaintenanceStatus())
				busyState.WithLabelValues(computer.GetLabelValues()...).Set(computer.GetBusyStatus())
			}
			L:
			for k,_ := range listOfKnownComputer {

				for _, computer := range computerList.Computers {
					if computer.DisplayName == k  {
						continue L
					}
				}
				log.Printf("computer %v was removed from master, removing metric..", k)
				maintenanceState.DeleteLabelValues(k)
				busyState.DeleteLabelValues(k)
				delete(listOfKnownComputer, k)
			}


			time.Sleep(10 * time.Second)
		}
	}()
}

func Entrypoint() {
	r := prometheus.NewRegistry()
	promWatchJenkinsNodes(*r)
	handler := promhttp.HandlerFor(r,promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)
	addr := fmt.Sprintf(":%d",viper.GetInt("port"))
	http.ListenAndServe(addr, nil)
}
