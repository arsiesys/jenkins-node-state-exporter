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
	"github.com/spf13/viper"
	"strings"
)

type computerList struct {
	Computers []computer `json:"computer"`
}

type label struct {
	Name string `json:"name"`
}

type computer struct {
	DisplayName        string  `json:"displayName"`
	Idle               bool    `json:"idle"`
	Offline            bool    `json:"offline"`
	OfflineCauseReason string  `json:"offlineCauseReason"`
	TemporarilyOffline bool    `json:"temporarilyOffline"`
	AssignedLabels     []label `json:"assignedLabels"`
}

func (c *computer) GetCustomTagFromAssignedLabels() string {
	// If no role is assigned, defaulting to "worker"
	prefix := viper.GetString("labelrole")
	role := "worker"
	for _, label := range c.AssignedLabels {
		labelName := label.Name
		if strings.HasPrefix(labelName, prefix) {
			extractedRole := labelName[len(prefix):]
			if len(extractedRole) > 0 {
				role = extractedRole
			}
		}
	}
	return role
}

func (c *computer) GetLabelValues() []string {
	return []string{
		c.DisplayName,
		c.GetCustomTagFromAssignedLabels(),
	}
}

func (c *computer) GetLabelValuesString() string {
	return strings.Join(c.GetLabelValues(), "_")
}

func (c *computer) GetBusyStatus() float64 {
	// 0: Idle
	// 1: Busy
	if c.Idle {
		return 0
	}
	return 1
}

func (c *computer) GetMaintenanceStatus() float64 {
	// O: Online
	// 1: Maintenance
	// 2: Offline
	if c.Offline {
		if c.TemporarilyOffline {
			return 1
		} else {
			return 2
		}
	}
	return 0
}
