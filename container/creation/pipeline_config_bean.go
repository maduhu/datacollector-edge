/*
 * Copyright 2017 StreamSets Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package creation

import (
	"github.com/streamsets/datacollector-edge/container/common"
)

type PipelineConfigBean struct {
	Version              string
	ExecutionMode        string
	DeliveryGuarantee    string
	ShouldRetry          bool
	RetryAttempts        float64
	MemoryLimit          string
	MemoryLimitExceeded  string
	NotifyOnStates       []interface{}
	EmailIDs             []interface{}
	Constants            map[string]interface{}
	BadRecordsHandling   string
	StatsAggregatorStage string
	RateLimit            float64
	MaxRunners           float64
}

func NewPipelineConfigBean(pipelineConfig common.PipelineConfiguration) PipelineConfigBean {
	pipelineConfigBean := PipelineConfigBean{}

	for _, config := range pipelineConfig.Configuration {
		if config.Value == nil {
			continue
		}
		switch config.Name {
		case "executionMode":
			pipelineConfigBean.ExecutionMode = config.Value.(string)
		case "deliveryGuarantee":
			pipelineConfigBean.DeliveryGuarantee = config.Value.(string)
		case "shouldRetry":
			pipelineConfigBean.ShouldRetry = config.Value.(bool)
		case "retryAttempts":
			pipelineConfigBean.RetryAttempts = config.Value.(float64)
		case "memoryLimit":
			pipelineConfigBean.MemoryLimit = config.Value.(string)
		case "memoryLimitExceeded":
			pipelineConfigBean.MemoryLimitExceeded = config.Value.(string)
		case "notifyOnStates":
			pipelineConfigBean.NotifyOnStates = config.Value.([]interface{})
		case "emailIDs":
			pipelineConfigBean.EmailIDs = config.Value.([]interface{})
		case "constants":
			constants := config.Value.([]interface{})
			pipelineConfigBean.Constants = make(map[string]interface{})
			for _, constant := range constants {
				constantMap := constant.(map[string]interface{})
				key := constantMap["key"]
				pipelineConfigBean.Constants[key.(string)] = constantMap["value"]
			}
		case "badRecordsHandling":
			pipelineConfigBean.BadRecordsHandling = config.Value.(string)
		case "statsAggregatorStage":
			pipelineConfigBean.StatsAggregatorStage = config.Value.(string)
		case "rateLimit":
			pipelineConfigBean.RateLimit = config.Value.(float64)
		case "maxRunners":
			pipelineConfigBean.MaxRunners = config.Value.(float64)
		}
	}

	return pipelineConfigBean
}

func GetDefaultPipelineConfigs() []common.Config {
	pipelineConfigs := []common.Config{
		common.Config{Name: "executionMode", Value: "STANDALONE"},
		common.Config{Name: "deliveryGuarantee", Value: "AT_LEAST_ONCE"},
		common.Config{Name: "shouldRetry", Value: true},
		common.Config{Name: "retryAttempts", Value: -1},
		common.Config{Name: "memoryLimit", Value: "${jvm:maxMemoryMB() * 0.65}"},
		common.Config{Name: "memoryLimitExceeded", Value: "STOP_PIPELINE"},
		common.Config{Name: "notifyOnStates", Value: []string{"RUN_ERROR", "STOPPED", "FINISHED"}},
		common.Config{Name: "emailIDs", Value: []string{}},
		common.Config{Name: "constants", Value: []string{}},
		common.Config{Name: "badRecordsHandling", Value: "streamsets-datacollector-basic-lib::com_streamsets_pipeline_stage_destination_devnull_ToErrorNullDTarget::1"},
		common.Config{Name: "clusterSlaveMemory", Value: 1024},
		common.Config{Name: "clusterSlaveJavaOpts", Value: "-XX:+UseConcMarkSweepGC -XX:+UseParNewGC -Dlog4j.debug"},
		common.Config{Name: "clusterLauncherEnv", Value: []string{}},
		common.Config{Name: "mesosDispatcherURL", Value: nil},
		common.Config{Name: "hdfsS3ConfDir", Value: nil},
		common.Config{Name: "rateLimit", Value: 0},
		common.Config{Name: "maxRunners", Value: 0},
		common.Config{Name: "webhookConfigs", Value: []interface{}{}},
		common.Config{Name: "statsAggregatorStage", Value: "streamsets-datacollector-basic-lib::com_streamsets_pipeline_stage_destination_devnull_StatsDpmDirectlyDTarget::1"},
	}

	return pipelineConfigs
}

func GetTrashErrorStageInstance() common.StageConfiguration {
	return common.StageConfiguration{
		InstanceName:  "Discard_ErrorStage",
		Library:       "streamsets-datacollector-basic-lib",
		StageName:     "com_streamsets_pipeline_stage_destination_devnull_ToErrorNullDTarget",
		StageVersion:  "1",
		Configuration: []common.Config{},
		UiInfo:        map[string]interface{}{},
		InputLanes:    []string{},
		OutputLanes:   []string{},
		EventLanes:    []string{},
	}
}

func GetDefaultStatsAggregatorStageInstance() common.StageConfiguration {
	return common.StageConfiguration{
		InstanceName:  "WritetoDPMdirectly_StatsAggregatorStage",
		Library:       "streamsets-datacollector-basic-lib",
		StageName:     "com_streamsets_pipeline_stage_destination_devnull_StatsDpmDirectlyDTarget",
		StageVersion:  "1",
		Configuration: []common.Config{},
		UiInfo:        map[string]interface{}{},
		InputLanes:    []string{},
		OutputLanes:   []string{},
		EventLanes:    []string{},
	}
}
