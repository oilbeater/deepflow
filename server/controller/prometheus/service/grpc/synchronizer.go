/**
 * Copyright (c) 2023 Yunshan Networks
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

package grpc

import (
	"github.com/op/go-logging"
	"golang.org/x/net/context"

	"github.com/deepflowio/deepflow/message/trident"
	"github.com/deepflowio/deepflow/server/controller/prometheus"
	"github.com/deepflowio/deepflow/server/controller/trisolaris/services/grpc/statsd"
)

var log = logging.MustGetLogger("prometheus.synchronizer.grpc")

type SynchronizerEvent struct{}

func NewSynchronizerEvent() *SynchronizerEvent {
	return &SynchronizerEvent{}
}

func (e *SynchronizerEvent) GetLabelIDs(ctx context.Context, in *trident.PrometheusLabelRequest) (*trident.PrometheusLabelResponse, error) {
	if len(in.GetRequestLabels()) != 0 || len(in.GetRequestTargets()) != 0 {
		log.Debugf("PrometheusLabelRequest: %+v", in)
	}
	synchronizer := prometheus.NewSynchronizer()
	resp, err := synchronizer.Sync(in)
	statsd.GetPrometheusLabelIDsCounterSingleton().Fill(synchronizer.GetStatsdCounter())
	if len(in.GetRequestLabels()) != 0 || len(in.GetRequestTargets()) != 0 {
		log.Debugf("PrometheusLabelResponse: %+v", resp)
	}
	log.Infof("counter detail: %+v", synchronizer.GetStatsdCounter())
	if err != nil {
		log.Errorf("encode str error: %+v", err)
		return &trident.PrometheusLabelResponse{}, nil
	}
	return resp, err
}
