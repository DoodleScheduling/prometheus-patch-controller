/*
Copyright 2022 Doodle.

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

package v1beta1

import (
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ActiveCondition             = "Active"
	FailedReason                = "Failed"
	InactiveReason              = "Inactive"
	PendingReason               = "Pending"
	ActiveReason                = "Active"
	InvalidPrometheusURLReason  = "InvalidPrometheusURL"
	PrometheusQueryFailedReason = "PrometheusQueryFailed"
	PatchAppliedCondition       = "PatchApplied"
	PatchApplyFailedReason      = "Failed"
	PatchAppliedReason          = "Applied"
	NoPatchFoundReason          = "NoPatchFound"
)

// PrometheusPatchRuleSpec defines the desired state of PrometheusPatchRule
type PrometheusPatchRuleSpec struct {
	// Prometheus holds information about where to find prometheus
	// +required
	Prometheus PrometheusSpec `json:"prometheus"`

	// Interval is the duration in which the expression gets evaluated
	// +required
	Interval metav1.Duration `json:"interval,omitempty"`

	// Expression is the prometheus .query
	// +required
	Expr string `json:"expr,omitempty"`

	// For is a durstion for how long the rule should be in pending before apply patches.
	// +required
	For metav1.Duration `json:"for,omitempty"`

	// .JSON6902Patches define to what target are applied what patches
	// +required
	JSON6902Patches []JSON6902Patch `json:"json6902Patches,omitempty"`

	// Suspend may suspend reconciliation of the resource.
	// +optional
	Suspend bool `json:"suspend,omitempty"`
}

// PrometheusSpec contains specs for accessing prometheus
type PrometheusSpec struct {
	Address string `json:"address"`
}

// JSON6902Patch is a target selector and a list of JSON6902 patches
type JSON6902Patch struct {
	// Patch contains JSON6902 patches with
	// an array of operation objects.
	// +required
	Patch []JSONPatch `json:"patch,omitempty"`

	// Target points to the resources that the patch document should be applied to.
	// +optional
	Target Selector `json:"target,omitempty"`
}

// JSONPatch is a JSON 6902 conform patch
type JSONPatch struct {
	OP    string     `json:"op"`
	Path  string     `json:"path"`
	Value extv1.JSON `json:"value"`
}

// Selector specifies a set of resources. Any resource that matches intersection of all conditions is included in this
// set.
type Selector struct {
	// Group is the API group to select resources from.
	// Together with Version and Kind it is capable of unambiguously identifying and/or selecting resources.
	// https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/api-group.md
	// +optional
	Group string `json:"group,omitempty"`

	// Version of the API Group to select resources from.
	// Together with Group and Kind it is capable of unambiguously identifying and/or selecting resources.
	// https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/api-group.md
	// +optional
	Version string `json:"version,omitempty"`

	// Kind of the API Group to select resources from.
	// Together with Group and Version it is capable of unambiguously
	// identifying and/or selecting resources.
	// https://github.com/kubernetes/community/blob/master/contributors/design-proposals/api-machinery/api-group.md
	// +optional
	Kind string `json:"kind,omitempty"`

	// Namespace to select resources from.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name to match resources with.
	// +optional
	Name string `json:"name,omitempty"`

	// LabelSelector is a string that follows the label selection expression
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#api
	// It matches with the resource labels.
	// +optional
	LabelSelector string `json:"labelSelector,omitempty"`
}

// PrometheusPatchRuleStatus defines the observed state of PrometheusPatchRule
type PrometheusPatchRuleStatus struct {
	// Conditions holds the conditions for the PrometheusPatchRule.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ConditionalResource is a resource with conditions
type conditionalResource interface {
	GetStatusConditions() *[]metav1.Condition
	GetGeneration() int64
}

// setResourceCondition sets the given condition with the given status,
// reason and message on a resource.
func setResourceCondition(resource conditionalResource, condition string, status metav1.ConditionStatus, reason, message string) {
	conditions := resource.GetStatusConditions()

	newCondition := metav1.Condition{
		Type:               condition,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: resource.GetGeneration(),
	}

	apimeta.SetStatusCondition(conditions, newCondition)
}

// PrometheusPatchRuleNotActive
func PrometheusPatchRuleNotActive(rule PrometheusPatchRule, reason, message string) PrometheusPatchRule {
	setResourceCondition(&rule, ActiveCondition, metav1.ConditionFalse, reason, message)
	return rule
}

// PrometheusPatchRuleActive
func PrometheusPatchRuleActive(rule PrometheusPatchRule, reason, message string) PrometheusPatchRule {
	setResourceCondition(&rule, ActiveCondition, metav1.ConditionTrue, reason, message)
	return rule
}

// PrometheusPatchRuleNoPatchApplied
func PrometheusPatchRuleNoPatchApplied(rule PrometheusPatchRule, reason, message string) PrometheusPatchRule {
	setResourceCondition(&rule, PatchAppliedCondition, metav1.ConditionFalse, reason, message)
	return rule
}

// PrometheusPatchRulePatchApplied
func PrometheusPatchRulePatchApplied(rule PrometheusPatchRule, reason string) PrometheusPatchRule {
	setResourceCondition(&rule, PatchAppliedCondition, metav1.ConditionTrue, reason, "")
	return rule
}

// GetStatusConditions returns a pointer to the Status.Conditions slice
func (in *PrometheusPatchRule) GetStatusConditions() *[]metav1.Condition {
	return &in.Status.Conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Active",type="string",JSONPath=".status.conditions[?(@.type==\"Active\")].status",description=""
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type==\"Active\")].reason",description=""
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description=""

// PrometheusPatchRule is the Schema for the patchrules API
type PrometheusPatchRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PrometheusPatchRuleSpec   `json:"spec,omitempty"`
	Status PrometheusPatchRuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PrometheusPatchRuleList contains a list of PrometheusPatchRule
type PrometheusPatchRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PrometheusPatchRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PrometheusPatchRule{}, &PrometheusPatchRuleList{})
}
