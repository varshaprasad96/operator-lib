// Copyright 2020 The Operator-SDK Authors
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

package conditions

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "github.com/operator-framework/api/pkg/operators/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Get/Create OperatorCondition", func() {
	It("creates the OperatorCondition", func() {
		oc := v1.OperatorCondition{}
		oc.SetName("some-oc")
		oc.SetNamespace("default")
		Expect(k8sClient.Create(context.TODO(), &oc)).To(Succeed())
	})
	It("gets the OperatorCondition", func() {
		oc := v1.OperatorCondition{}
		nn := types.NamespacedName{Name: "some-oc", Namespace: "default"}
		Expect(k8sClient.Get(context.TODO(), nn, &oc)).To(Succeed())
		Expect(oc.GetName()).To(Equal(nn.Name))
	})
})
