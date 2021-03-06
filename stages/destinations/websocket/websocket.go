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
package websocket

import (
	"bytes"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/streamsets/datacollector-edge/api"
	"github.com/streamsets/datacollector-edge/container/common"
	"github.com/streamsets/datacollector-edge/stages/lib/datagenerator"
	"github.com/streamsets/datacollector-edge/stages/stagelibrary"
	"log"
	"net/http"
)

const (
	LIBRARY    = "streamsets-datacollector-basic-lib"
	STAGE_NAME = "com_streamsets_pipeline_stage_destination_websocket_WebSocketDTarget"
)

type WebSocketClientDestination struct {
	*common.BaseStage
	Conf WebSocketTargetConfig `ConfigDefBean:"conf"`
}

type WebSocketTargetConfig struct {
	ResourceUrl               string                                  `ConfigDef:"type=STRING,required=true"`
	Headers                   map[string]string                       `ConfigDef:"type=MAP,required=true"`
	DataFormat                string                                  `ConfigDef:"type=STRING,required=true"`
	DataGeneratorFormatConfig datagenerator.DataGeneratorFormatConfig `ConfigDefBean:"dataGeneratorFormatConfig"`
}

func init() {
	stagelibrary.SetCreator(LIBRARY, STAGE_NAME, func() api.Stage {
		return &WebSocketClientDestination{BaseStage: &common.BaseStage{}}
	})
}

func (w *WebSocketClientDestination) Init(stageContext api.StageContext) error {
	if err := w.BaseStage.Init(stageContext); err != nil {
		return err
	}
	log.Println("[DEBUG] WebSocketClientDestination Init method")
	return w.Conf.DataGeneratorFormatConfig.Init(w.Conf.DataFormat)
}

func (w *WebSocketClientDestination) Write(batch api.Batch) error {
	log.Println("[DEBUG] WebSocketClientDestination write method = " + w.Conf.ResourceUrl)
	recordWriterFactory := w.Conf.DataGeneratorFormatConfig.RecordWriterFactory
	if recordWriterFactory == nil {
		return errors.New("recordWriterFactory is null")
	}

	var requestHeader = http.Header{}
	if w.Conf.Headers != nil {
		for key, value := range w.Conf.Headers {
			requestHeader.Set(key, value)
		}
	}

	c, _, err := websocket.DefaultDialer.Dial(w.Conf.ResourceUrl, requestHeader)
	if err != nil {
		return err
	}

	for _, record := range batch.GetRecords() {
		recordBuffer := bytes.NewBuffer([]byte{})
		recordWriter, err := recordWriterFactory.CreateWriter(w.GetStageContext(), recordBuffer)
		if err != nil {
			return err
		}
		err = recordWriter.WriteRecord(record)
		if err != nil {
			return err
		}
		recordWriter.Flush()
		recordWriter.Close()

		err = c.WriteMessage(websocket.TextMessage, recordBuffer.Bytes())
		if err != nil {
			log.Println("[ERROR] write:", err)
			w.GetStageContext().ToError(err, record)
		}
	}

	defer c.Close()
	return nil
}
