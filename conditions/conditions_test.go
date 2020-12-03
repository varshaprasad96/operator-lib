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
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiv1 "github.com/operator-framework/api/pkg/operators/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kubeclock "k8s.io/apimachinery/pkg/util/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	conditionFoo = "conditionFoo"
)

var _ = Describe("Condition", func() {
	var operatorCond *apiv1.OperatorCondition
	var ns = "default"
	ctx := context.TODO()
	var clock kubeclock.Clock = &kubeclock.RealClock{}
	var transitionTime metav1.Time = metav1.Time{Time: clock.Now()}

	Describe("Get", func() {
		It("should fetch the right condition", func() {
			operatorCond = &apiv1.OperatorCondition{
				ObjectMeta: metav1.ObjectMeta{Name: "operator-condition-1", Namespace: ns},
				Spec:       apiv1.OperatorConditionSpec{},
				Status: apiv1.OperatorConditionStatus{
					Conditions: []metav1.Condition{
						{
							Type:               conditionFoo,
							Status:             metav1.ConditionTrue,
							Reason:             "foo",
							Message:            "The operator is in foo condition",
							LastTransitionTime: transitionTime,
						},
					},
				},
			}
			operatorCond.SetGroupVersionKind(apiv1.GroupVersion.WithKind("OperatorCondition"))

			By("creating the operator condition")
			err := os.Setenv(operatorCondEnvVar, "operator-condition-1")
			Expect(err).NotTo(HaveOccurred())
			readNamespace = func() (string, error) {
				return ns, nil
			}

			k := types.NamespacedName{
				Namespace: ns,
				Name:      "operator-condition-1",
			}

			client, err := client.New(cfg, client.Options{Scheme: sch})

			err = client.Create(ctx, operatorCond)
			Expect(err).NotTo(HaveOccurred())

			op := &apiv1.OperatorCondition{}
			Eventually(func() error {
				return client.Get(ctx, k, op)
			}, cfg.Timeout, time.Millisecond*10).Should(Succeed())
			err = client.Get(ctx, k, op)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
