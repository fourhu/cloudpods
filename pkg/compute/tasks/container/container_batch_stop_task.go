// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package container

import (
	"context"
	"fmt"

	"yunion.io/x/jsonutils"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/compute/models"
)

func init() {
	taskman.RegisterTask(ContainerBatchStopTask{})
}

type ContainerBatchStopTask struct {
	taskman.STask
}

func (t *ContainerBatchStopTask) OnInit(ctx context.Context, objs []db.IStandaloneModel, data jsonutils.JSONObject) {
	t.SetStage("OnContainersStopComplete", nil)

	params := make([]api.ContainerStopInput, 0)
	t.GetParams().Unmarshal(&params, "params")

	for i := range objs {
		ctr := objs[i].(*models.SContainer)
		if err := ctr.StartStopTask(ctx, t.GetUserCred(), &params[i], t.GetTaskId()); err != nil {
			t.SetStageFailed(ctx, jsonutils.NewString(fmt.Sprintf("stop container %s: %s", ctr.GetName(), err.Error())))
			return
		}
	}
}

func (t *ContainerBatchStopTask) OnContainersStopComplete(ctx context.Context, objs []db.IStandaloneModel, data jsonutils.JSONObject) {
	t.SetStageComplete(ctx, nil)
}
