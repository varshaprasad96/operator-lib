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
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiv1 "github.com/operator-framework/api/pkg/operators/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestSource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "Conditions Suite", []Reporter{printer.NewlineReporter{}, printer.NewProwReporter("Conditions Suite")})
}

var testenv *envtest.Environment
var cfg *rest.Config
var clientset *kubernetes.Clientset
var sch = runtime.NewScheme()
var err error

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	// Add operator apiv1 to scheme
	err = apiv1.AddToScheme(sch)
	Expect(err).NotTo(HaveOccurred())

	testenv = &envtest.Environment{}
	testenv.CRDInstallOptions = envtest.CRDInstallOptions{
		Paths: []string{"olm/"},
	}

	cfg, err = testenv.Start()
	Expect(err).NotTo(HaveOccurred())

	// Remove this, testenv.Start installs CRD
	obj, err := envtest.InstallCRDs(cfg, testenv.CRDInstallOptions)
	Expect(err).NotTo(HaveOccurred())
	Expect(obj).NotTo(BeNil())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	Expect(testenv.Stop()).To(Succeed())
})
