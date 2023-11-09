/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	corev1 "k8s.io/api/core/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	v1beta1 "github.com/doodlescheduling/prometheuspatch-controller/api/v1beta1"
	// +kubebuilder:scaffold:imports
)

type prometheusContainer struct {
	testcontainers.Container
	URI string
}

func setupPrometheusContainer(ctx context.Context) (*prometheusContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "prom/prometheus:v2.34.0",
		ExposedPorts: []string{"9090/tcp"},
		WaitingFor:   wait.ForListeningPort("9090"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "9090")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())
	return &prometheusContainer{Container: container, URI: uri}, nil
}

var _ = Describe("PrometheusPatchRule tests", func() {
	const (
		timeout  = time.Second * 30
		interval = time.Second * 1
	)

	var (
		container *prometheusContainer
		err       error
	)

	container, err = setupPrometheusContainer(context.TODO())
	Expect(err).NotTo(HaveOccurred(), "failed to start prometheus container")
	Describe("fails if it has an invalid prometheus", func() {
		var (
			createdRule *v1beta1.PrometheusPatchRule
			keyRule     types.NamespacedName
		)

		It("creates PrometheusPatchRule successfully", func() {
			keyRule = types.NamespacedName{
				Name:      "rule-" + randStringRunes(5),
				Namespace: "default",
			}
			createdRule = &v1beta1.PrometheusPatchRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      keyRule.Name,
					Namespace: keyRule.Namespace,
				},
				Spec: v1beta1.PrometheusPatchRuleSpec{
					Prometheus: v1beta1.PrometheusSpec{
						Address: ":",
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), createdRule)).Should(Succeed())
		})

		It("fails reconcile because Active condition is False with InvalidPrometheusURL reason", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 1 &&
					got.Status.Conditions[0].Reason == v1beta1.InvalidPrometheusURLReason &&
					got.Status.Conditions[0].Status == "False" &&
					got.Status.Conditions[0].Type == v1beta1.ActiveCondition
			}, timeout, interval).Should(BeTrue())
		})
	})

	Describe("fails if prometheus expression is invalid", func() {
		var (
			createdRule *v1beta1.PrometheusPatchRule
			keyRule     types.NamespacedName
		)

		It("creates PrometheusPatchRule successfully", func() {
			keyRule = types.NamespacedName{
				Name:      "rule-" + randStringRunes(5),
				Namespace: "default",
			}
			createdRule = &v1beta1.PrometheusPatchRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      keyRule.Name,
					Namespace: keyRule.Namespace,
				},
				Spec: v1beta1.PrometheusPatchRuleSpec{
					Expr: "invalid)",
					Prometheus: v1beta1.PrometheusSpec{
						Address: container.URI,
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), createdRule)).Should(Succeed())
		})

		It("fails reconcile because Active condition is False with PrometheusQueryFailed reason", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 1 &&
					got.Status.Conditions[0].Reason == v1beta1.PrometheusQueryFailedReason &&
					got.Status.Conditions[0].Status == "False" &&
					got.Status.Conditions[0].Type == v1beta1.ActiveCondition
			}, timeout, interval).Should(BeTrue())
		})
	})

	Describe("rule is inactive if expression does not return samples", func() {
		var (
			createdRule *v1beta1.PrometheusPatchRule
			keyRule     types.NamespacedName
		)

		It("creates PrometheusPatchRule successfully", func() {
			keyRule = types.NamespacedName{
				Name:      "rule-" + randStringRunes(5),
				Namespace: "default",
			}
			createdRule = &v1beta1.PrometheusPatchRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      keyRule.Name,
					Namespace: keyRule.Namespace,
				},
				Spec: v1beta1.PrometheusPatchRuleSpec{
					Expr: "non_existing_metric > 0",
					Prometheus: v1beta1.PrometheusSpec{
						Address: container.URI,
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), createdRule)).Should(Succeed())
		})

		It("Active is False with reason Pending", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 1 &&
					got.Status.Conditions[0].Reason == v1beta1.InactiveReason &&
					got.Status.Conditions[0].Status == "False" &&
					got.Status.Conditions[0].Type == v1beta1.ActiveCondition
			}, timeout, interval).Should(BeTrue())
		})
	})

	Describe("rule is inactive if expression does not return samples", func() {
		var (
			createdRule *v1beta1.PrometheusPatchRule
			keyRule     types.NamespacedName
		)

		It("creates PrometheusPatchRule successfully", func() {
			keyRule = types.NamespacedName{
				Name:      "rule-" + randStringRunes(5),
				Namespace: "default",
			}
			createdRule = &v1beta1.PrometheusPatchRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      keyRule.Name,
					Namespace: keyRule.Namespace,
				},
				Spec: v1beta1.PrometheusPatchRuleSpec{
					Expr: "non_existing_metric > 0",
					Prometheus: v1beta1.PrometheusSpec{
						Address: container.URI,
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), createdRule)).Should(Succeed())
		})

		It("Active is False with reason Pending", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 1 &&
					got.Status.Conditions[0].Reason == v1beta1.InactiveReason &&
					got.Status.Conditions[0].Status == "False" &&
					got.Status.Conditions[0].Type == v1beta1.ActiveCondition
			}, timeout, interval).Should(BeTrue())
		})
	})

	Describe("rule is active if expression returns samples", func() {
		var (
			createdRule *v1beta1.PrometheusPatchRule
			keyRule     types.NamespacedName
		)

		duration, err := time.ParseDuration("5s")
		Expect(err).NotTo(HaveOccurred(), "failed to parse interval duration")

		It("creates PrometheusPatchRule successfully", func() {
			keyRule = types.NamespacedName{
				Name:      "rule-" + randStringRunes(5),
				Namespace: "default",
			}
			createdRule = &v1beta1.PrometheusPatchRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      keyRule.Name,
					Namespace: keyRule.Namespace,
				},
				Spec: v1beta1.PrometheusPatchRuleSpec{
					Expr: "prometheus_build_info > 0",
					Interval: metav1.Duration{
						Duration: duration,
					},
					Prometheus: v1beta1.PrometheusSpec{
						Address: container.URI,
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), createdRule)).Should(Succeed())
		})

		It("Active condition is True with reason Active", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 2 &&
					got.Status.Conditions[0].Reason == v1beta1.ActiveReason &&
					got.Status.Conditions[0].Status == "True" &&
					got.Status.Conditions[0].Type == v1beta1.ActiveCondition
			}, timeout, interval).Should(BeTrue())
		})

		It("PatchesApplied condition is False since there are no patches defined", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 2 &&
					got.Status.Conditions[1].Reason == v1beta1.NoPatchFoundReason &&
					got.Status.Conditions[1].Status == "False" &&
					got.Status.Conditions[1].Type == v1beta1.PatchAppliedCondition
			}, timeout, interval).Should(BeTrue())
		})
	})

	Describe("patch is applied to single resource selector", func() {
		var (
			createdRule *v1beta1.PrometheusPatchRule
			keyRule     types.NamespacedName
		)

		duration, err := time.ParseDuration("5s")
		Expect(err).NotTo(HaveOccurred(), "failed to parse interval duration")

		It("creates PrometheusPatchRule successfully", func() {
			keyRule = types.NamespacedName{
				Name:      "rule-" + randStringRunes(5),
				Namespace: "default",
			}
			createdRule = &v1beta1.PrometheusPatchRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      keyRule.Name,
					Namespace: keyRule.Namespace,
				},
				Spec: v1beta1.PrometheusPatchRuleSpec{
					Expr: "prometheus_build_info > 0",
					Interval: metav1.Duration{
						Duration: duration,
					},
					JSON6902Patches: []v1beta1.JSON6902Patch{
						v1beta1.JSON6902Patch{
							Target: v1beta1.Selector{
								Version:  "v1",
								Resource: "namespaces",
								Name:     "default",
							},
							Patch: []v1beta1.JSONPatch{
								v1beta1.JSONPatch{
									OP:   "add",
									Path: "/metadata/annotations",
									Value: extv1.JSON{
										Raw: []byte(`{"foo":"bar"}`),
									},
								},
							},
						},
					},
					Prometheus: v1beta1.PrometheusSpec{
						Address: container.URI,
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), createdRule)).Should(Succeed())
		})

		It("PatchesApplied condition is True", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 2 &&
					got.Status.Conditions[1].Reason == v1beta1.PatchAppliedReason &&
					got.Status.Conditions[1].Status == "True" &&
					got.Status.Conditions[1].Type == v1beta1.PatchAppliedCondition
			}, timeout, interval).Should(BeTrue())
		})

		It("actually has resource patched", func() {
			got := &corev1.Namespace{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), types.NamespacedName{
					Name: "default",
				}, got)

				if val, ok := got.Annotations["foo"]; ok {
					return val == "bar"
				}

				return false
			}, timeout, interval).Should(BeTrue())
		})

		It("has as correct field manager set", func() {
			got := &corev1.Namespace{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), types.NamespacedName{
					Name: "default",
				}, got)

				for _, v := range got.ManagedFields {
					if v.Manager == "test-suite" {
						return true
					}
				}

				return false
			}, timeout, interval).Should(BeTrue())
		})
	})

	Describe("multiple patches are applied to multiple resource selector", func() {
		var (
			createdRule *v1beta1.PrometheusPatchRule
			keyRule     types.NamespacedName
		)

		duration, err := time.ParseDuration("5s")
		Expect(err).NotTo(HaveOccurred(), "failed to parse interval duration")

		It("creates PrometheusPatchRule successfully", func() {
			keyRule = types.NamespacedName{
				Name:      "rule-" + randStringRunes(5),
				Namespace: "default",
			}
			createdRule = &v1beta1.PrometheusPatchRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      keyRule.Name,
					Namespace: keyRule.Namespace,
					Labels: map[string]string{
						"selector": "foo",
					},
				},
				Spec: v1beta1.PrometheusPatchRuleSpec{
					Expr: "prometheus_build_info > 0",
					Interval: metav1.Duration{
						Duration: duration,
					},
					JSON6902Patches: []v1beta1.JSON6902Patch{
						v1beta1.JSON6902Patch{
							Target: v1beta1.Selector{
								Version:  "v1",
								Resource: "namespaces",
								Name:     "default",
							},
							Patch: []v1beta1.JSONPatch{
								v1beta1.JSONPatch{
									OP:   "add",
									Path: "/metadata/annotations",
									Value: extv1.JSON{
										Raw: []byte(`{"foo":"bar"}`),
									},
								},
							},
						},
						v1beta1.JSON6902Patch{
							Target: v1beta1.Selector{
								Group:         "metrics.infra.doodle.com",
								Version:       "v1beta1",
								Resource:      "prometheuspatchrules",
								LabelSelector: "selector=foo",
							},
							Patch: []v1beta1.JSONPatch{
								v1beta1.JSONPatch{
									OP:   "add",
									Path: "/metadata/annotations",
									Value: extv1.JSON{
										Raw: []byte(`{"foo":"bar2"}`),
									},
								},
							},
						},
					},
					Prometheus: v1beta1.PrometheusSpec{
						Address: container.URI,
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), createdRule)).Should(Succeed())
		})

		It("PatchesApplied condition is True", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 2 &&
					got.Status.Conditions[1].Reason == v1beta1.PatchAppliedReason &&
					got.Status.Conditions[1].Status == "True" &&
					got.Status.Conditions[1].Type == v1beta1.PatchAppliedCondition
			}, timeout, interval).Should(BeTrue())
		})

		It("actually has namespace resource patched", func() {
			got := &corev1.Namespace{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), types.NamespacedName{
					Name: "default",
				}, got)

				if val, ok := got.Annotations["foo"]; ok {
					return val == "bar"
				}

				return false
			}, timeout, interval).Should(BeTrue())
		})

		It("actually has prometheuspatchrules patched", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				if val, ok := got.Annotations["foo"]; ok {
					return val == "bar2"
				}

				return false
			}, timeout, interval).Should(BeTrue())
		})
	})

	Describe("Condition active is True but condition PatchApplied is False if target finds no resources", func() {
		var (
			createdRule *v1beta1.PrometheusPatchRule
			keyRule     types.NamespacedName
		)

		duration, err := time.ParseDuration("5s")
		Expect(err).NotTo(HaveOccurred(), "failed to parse interval duration")

		It("creates PrometheusPatchRule successfully", func() {
			keyRule = types.NamespacedName{
				Name:      "rule-" + randStringRunes(5),
				Namespace: "default",
			}
			createdRule = &v1beta1.PrometheusPatchRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      keyRule.Name,
					Namespace: keyRule.Namespace,
				},
				Spec: v1beta1.PrometheusPatchRuleSpec{
					Expr: "prometheus_build_info > 0",
					Interval: metav1.Duration{
						Duration: duration,
					},
					JSON6902Patches: []v1beta1.JSON6902Patch{
						v1beta1.JSON6902Patch{
							Target: v1beta1.Selector{
								Group:     "metrics.infra.doodle.com",
								Version:   "v1beta1",
								Resource:  "prometheuspatchrules",
								Name:      "does-not-exist",
								Namespace: keyRule.Namespace,
							},
							Patch: []v1beta1.JSONPatch{
								v1beta1.JSONPatch{
									OP:   "add",
									Path: "/metadata/annotations",
									Value: extv1.JSON{
										Raw: []byte(`{"foo":"bar2"}`),
									},
								},
							},
						},
					},
					Prometheus: v1beta1.PrometheusSpec{
						Address: container.URI,
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), createdRule)).Should(Succeed())
		})

		It("PatchesApplied condition is False", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 2 &&
					got.Status.Conditions[1].Reason == v1beta1.PatchApplyFailedReason &&
					got.Status.Conditions[1].Status == "False" &&
					got.Status.Conditions[1].Type == v1beta1.PatchAppliedCondition
			}, timeout, interval).Should(BeTrue())
		})
	})

	Describe("patch is not applied before transition from pending into active", func() {
		var (
			createdRule *v1beta1.PrometheusPatchRule
			keyRule     types.NamespacedName
		)

		duration, err := time.ParseDuration("8s")
		Expect(err).NotTo(HaveOccurred(), "failed to parse interval duration")

		It("creates PrometheusPatchRule successfully", func() {
			keyRule = types.NamespacedName{
				Name:      "rule-" + randStringRunes(5),
				Namespace: "default",
			}
			createdRule = &v1beta1.PrometheusPatchRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      keyRule.Name,
					Namespace: keyRule.Namespace,
				},
				Spec: v1beta1.PrometheusPatchRuleSpec{
					Expr: "prometheus_build_info > 0",
					Interval: metav1.Duration{
						Duration: duration,
					},
					For: metav1.Duration{
						Duration: duration,
					},
					JSON6902Patches: []v1beta1.JSON6902Patch{
						v1beta1.JSON6902Patch{
							Target: v1beta1.Selector{
								Version:  "v1",
								Resource: "namespaces",
								Name:     "default",
							},
							Patch: []v1beta1.JSONPatch{
								v1beta1.JSONPatch{
									OP:   "add",
									Path: "/metadata/annotations",
									Value: extv1.JSON{
										Raw: []byte(`{"foo":"bar"}`),
									},
								},
							},
						},
					},
					Prometheus: v1beta1.PrometheusSpec{
						Address: container.URI,
					},
				},
			}

			Expect(k8sClient.Create(context.Background(), createdRule)).Should(Succeed())
		})

		It("Active condition is True and Pending", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 1 &&
					got.Status.Conditions[0].Reason == v1beta1.PendingReason &&
					got.Status.Conditions[0].Status == "True" &&
					got.Status.Conditions[0].Type == v1beta1.ActiveCondition
			}, timeout, interval).Should(BeTrue())
		})

		It("eventually gets lifted to active while beeing pending for 5s", func() {
			got := &v1beta1.PrometheusPatchRule{}
			Eventually(func() bool {
				_ = k8sClient.Get(context.Background(), keyRule, got)

				return len(got.Status.Conditions) == 2 &&
					got.Status.Conditions[0].Reason == v1beta1.ActiveReason &&
					got.Status.Conditions[0].Status == "True" &&
					got.Status.Conditions[0].Type == v1beta1.ActiveCondition
			}, timeout, interval).Should(BeTrue())
		})
	})
})
