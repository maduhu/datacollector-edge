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
package el

import (
	"context"
	"errors"
	"github.com/streamsets/datacollector-edge/api"
	"github.com/streamsets/datacollector-edge/api/fieldtype"
	"testing"
)

type MockRecord struct {
}

func (r *MockRecord) GetHeader() api.Header {
	return nil
}

func (r *MockRecord) Get(fieldPath ...string) (*api.Field, error) {
	if len(fieldPath) > 0 {
		switch fieldPath[0] {
		case "/a":
			return &api.Field{
				Type:  fieldtype.MAP,
				Value: map[string]string{"b": "value"},
			}, nil
		case "/a/b":
			return &api.Field{
				Type:  fieldtype.MAP,
				Value: "Test Value",
			}, nil
		case "/inValid":
			return &api.Field{}, errors.New("invalid fieldPath '/inValid'")
		default:
			return &api.Field{}, nil
		}
	}
	return &api.Field{}, nil
}

func (r *MockRecord) Set(field *api.Field) *api.Field {
	return &api.Field{}
}

func (r *MockRecord) GetFieldPaths() map[string]bool {
	return map[string]bool{}
}

func (r *MockRecord) Clone() api.Record {
	return r
}

func (r *MockRecord) SetField(fieldPath string, field *api.Field) (*api.Field, error) {
	return nil, errors.New("unsupported operation")
}

func (r *MockRecord) Delete(fieldPath string) (*api.Field, error) {
	return nil, errors.New("unsupported operation")
}

func TestRecordEL(test *testing.T) {
	evaluationTests := []EvaluationTest{
		{
			Name:       "Test function record:type",
			Expression: "${record:type('/a')}",
			Expected:   "MAP",
		},
		{
			Name:       "Test function record:type - Error 1",
			Expression: "${record:type()}",
			Expected:   "The function 'record:type' requires 1 arguments but was passed 0",
			ErrorCase:  true,
		},
		{
			Name:       "Test function record:type - Error 2",
			Expression: "${record:type('/inValid')}",
			Expected:   "invalid fieldPath '/inValid'",
			ErrorCase:  true,
		},
		{
			Name:       "Test function record:value",
			Expression: "${record:value('/a/b')}",
			Expected:   "Test Value",
		},
		{
			Name:       "Test function record:value - Error 1",
			Expression: "${record:value()}",
			Expected:   "The function 'record:value' requires 1 arguments but was passed 0",
			ErrorCase:  true,
		},
		{
			Name:       "Test function record:value - Error 2",
			Expression: "${record:value('/inValid')}",
			Expected:   "invalid fieldPath '/inValid'",
			ErrorCase:  true,
		},
		{
			Name:       "Test function record:valueOrDefault",
			Expression: "${record:valueOrDefault('notValid', 'test default value')}",
			Expected:   "test default value",
		},
		{
			Name:       "Test function record:valueOrDefault",
			Expression: "${record:valueOrDefault('/a/b', 'test default value')}",
			Expected:   "Test Value",
		},
		{
			Name:       "Test function record:valueOrDefault - Error 1",
			Expression: "${record:valueOrDefault()}",
			Expected:   "The function 'record:valueOrDefault' requires 2 arguments but was passed 0",
			ErrorCase:  true,
		},
		{
			Name:       "Test function record:valueOrDefault - Error 2",
			Expression: "${record:valueOrDefault('/inValid', 'inValid')}",
			Expected:   "invalid fieldPath '/inValid'",
			ErrorCase:  true,
		},

		{
			Name:       "Test function record:exists",
			Expression: "${record:exists('/a/b')}",
			Expected:   true,
		},
		{
			Name:       "Test function record:exists",
			Expression: "${record:exists('/a/b/c')}",
			Expected:   false,
		},
		{
			Name:       "Test function record:exists",
			Expression: "${record:exists('/inValid')}",
			Expected:   false,
		},
		{
			Name:       "Test function record:exists - Error 1",
			Expression: "${record:exists()}",
			Expected:   "The function 'record:exists' requires 1 arguments but was passed 0",
			ErrorCase:  true,
		},
	}

	record := &MockRecord{}
	recordContext := context.WithValue(context.Background(), RECORD_CONTEXT_VAR, record)
	RunEvaluationTests(evaluationTests, []Definitions{&RecordEL{Context: recordContext}}, test)
}

func TestRecordEL_withOutContext(test *testing.T) {
	evaluationTests := []EvaluationTest{
		{
			Name:       "Test function record:type",
			Expression: "${record:type('/a')}",
			Expected:   "record context is not set",
			ErrorCase:  true,
		},
		{
			Name:       "Test function record:value",
			Expression: "${record:value('/a')}",
			Expected:   "record context is not set",
			ErrorCase:  true,
		},
		{
			Name:       "Test function record:valueOrDefault",
			Expression: "${record:valueOrDefault('/a', 'defaultValue')}",
			Expected:   "record context is not set",
			ErrorCase:  true,
		},
		{
			Name:       "Test function record:exists",
			Expression: "${record:exists('/a')}",
			Expected:   "record context is not set",
			ErrorCase:  true,
		},
	}
	RunEvaluationTests(evaluationTests, []Definitions{&RecordEL{}}, test)
}
