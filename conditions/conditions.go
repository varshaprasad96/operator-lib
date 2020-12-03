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
	"fmt"
	"os"

	apiv1 "github.com/operator-framework/api/pkg/operators/v1"
	"github.com/operator-framework/operator-lib/internal/utils"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var readNamespace = utils.GetOperatorNamespace

const (
	// operatorCondEnvVar is the env variable which
	// contains the name of the Condition CR associated to the operator,
	// set by OLM.
	operatorCondEnvVar = "OPERATOR_CONDITION_NAME"
)

type condition struct {
	namespacedName types.NamespacedName
	condType       conditionType
	client         client.Client
}

var _ Condition = &condition{}

func NewCondition(cl client.Client, condType conditionType) (Condition, error) {
	objKey, err := GetNamespacedName()
	if err != nil {
		return nil, err
	}
	return &condition{
		namespacedName: *objKey,
		condType:       condType,
		client:         cl,
	}, nil
}

func (c *condition) Get(ctx context.Context, condType conditionType) (*metav1.Condition, error) {
	operatorCond := &apiv1.OperatorCondition{}
	err := c.client.Get(ctx, c.namespacedName, operatorCond)
	if err != nil {
		return nil, err
	}
	con := meta.FindStatusCondition(operatorCond.Status.Conditions, string(condType))

	if con == nil {
		return nil, fmt.Errorf("conditionType %s not found", string(condType))
	}
	return con, nil
}

func (c *condition) Set(ctx context.Context, condType conditionType, status metav1.ConditionStatus, option ...Option) error {
	return nil
}

// GetNamespacedName returns the NamespacedName of the CR. It returns an error
// when the name of the CR cannot be found from the environment variable set by
// OLM. Hence, GetNamespacedName() can provide the NamespacedName when the operator
// is running on cluster and is being managed by OLM. If running locally, operator
// writers are encouraged to skip this method or gracefully handle the errors by logging
// a message.
func GetNamespacedName() (*types.NamespacedName, error) {
	conditionName := os.Getenv(operatorCondEnvVar)
	if conditionName == "" {
		return nil, fmt.Errorf("could not determine operator condition name: environment variable %s not set", operatorCondEnvVar)
	}
	operatorNs, err := readNamespace()
	if err != nil {
		return nil, fmt.Errorf("could not determine operator namespace: %v", err)
	}
	return &types.NamespacedName{Name: conditionName, Namespace: operatorNs}, nil
}
