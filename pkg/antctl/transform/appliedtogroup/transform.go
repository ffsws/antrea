// Copyright 2020 Antrea Authors
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

package appliedtogroup

import (
	"io"
	"reflect"

	"github.com/vmware-tanzu/antrea/pkg/antctl/transform"
	"github.com/vmware-tanzu/antrea/pkg/antctl/transform/common"
	networkingv1beta1 "github.com/vmware-tanzu/antrea/pkg/apis/networking/v1beta1"
)

type Response struct {
	Name string                  `json:"name" yaml:"name"`
	Pods []common.GroupMemberPod `json:"pods,omitempty" yaml:"pods,omitempty"`
}

func listTransform(l interface{}) (interface{}, error) {
	groups := l.(*networkingv1beta1.AppliedToGroupList)
	result := []Response{}
	for _, group := range groups.Items {
		o, _ := objectTransform(&group)
		result = append(result, o.(Response))
	}
	return result, nil
}

func objectTransform(o interface{}) (interface{}, error) {
	group := o.(*networkingv1beta1.AppliedToGroup)
	var pods []common.GroupMemberPod
	for _, pod := range group.Pods {
		pods = append(pods, common.GroupMemberPodTransform(pod))
	}
	return Response{Name: group.GetName(), Pods: pods}, nil
}

func Transform(reader io.Reader, single bool) (interface{}, error) {
	return transform.GenericFactory(
		reflect.TypeOf(networkingv1beta1.AppliedToGroup{}),
		reflect.TypeOf(networkingv1beta1.AppliedToGroupList{}),
		objectTransform,
		listTransform,
	)(reader, single)
}
