/*
Copyright 2022.

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

package v1alpha1

import (
	//batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CascadeModule struct {
	// Configuration parameter for Cascade Module
	// +patchMergeKey=name
	// +patchStrategy=merge
	ModuleName    string            `json:"modulename"`
	Configuration map[string]string `json:"configuration"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Specifies the duration in seconds relative to the startTime that the job may be active
	// before the system tries to terminate it; value must be positive integer
	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty" protobuf:"varint,3,opt,name=activeDeadlineSeconds"`

	// Specifies the number of retries before marking this job failed.
	// Defaults to 6
	// +optional
	BackoffLimit *int32 `json:"backoffLimit,omitempty" protobuf:"varint,7,opt,name=backoffLimit"`

	// ttlSecondsAfterFinished limits the lifetime of a Job that has finished
	// execution (either Complete or Failed). If this field is set,
	// ttlSecondsAfterFinished after the Job finishes, it is eligible to be
	// automatically deleted. When the Job is being deleted, its lifecycle
	// guarantees (e.g. finalizers) will be honored. If this field is unset,
	// the Job won't be automatically deleted. If this field is set to zero,
	// the Job becomes eligible to be deleted immediately after it finishes.
	// This field is alpha-level and is only honored by servers that enable the
	// TTLAfterFinished feature.
	// +optional
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty" protobuf:"varint,8,opt,name=ttlSecondsAfterFinished"`

	// Describes the pod that will be created when executing a job.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Template corev1.PodTemplateSpec `json:"template" protobuf:"bytes,6,opt,name=template"`
}

type CascadeScenario struct {
	// CascadeModules list with configuration parameters
	// +patchMergeKey=name
	// +patchStrategy=merge
	CascadeModules []CascadeModule `json:"cascademodules"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CascadeManualOperatorSpec defines the desired state of CascadeManualOperator
type CascadeManualOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Specifies the duration in seconds relative to the startTime that the job may be active
	// before the system tries to terminate it; value must be positive integer
	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty" protobuf:"varint,3,opt,name=activeDeadlineSeconds"`

	// Specifies the number of retries before marking this job failed.
	// Defaults to 6
	// +optional
	BackoffLimit *int32 `json:"backoffLimit,omitempty" protobuf:"varint,7,opt,name=backoffLimit"`

	// ttlSecondsAfterFinished limits the lifetime of a Job that has finished
	// execution (either Complete or Failed). If this field is set,
	// ttlSecondsAfterFinished after the Job finishes, it is eligible to be
	// automatically deleted. When the Job is being deleted, its lifecycle
	// guarantees (e.g. finalizers) will be honored. If this field is unset,
	// the Job won't be automatically deleted. If this field is set to zero,
	// the Job becomes eligible to be deleted immediately after it finishes.
	// This field is alpha-level and is only honored by servers that enable the
	// TTLAfterFinished feature.
	// +optional
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty" protobuf:"varint,8,opt,name=ttlSecondsAfterFinished"`

	// Describes the pod that will be created when executing a job.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	Template corev1.PodTemplateSpec `json:"template" protobuf:"bytes,6,opt,name=template"`

	ScenarioConfig CascadeScenario `json:"scenarioconfig"`
}

// CascadeManualOperatorStatus defines the observed state of CascadeManualOperator
type CascadeManualOperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Active    int32 `json:"active"`
	Succeeded int32 `json:"succeeded"`
	Failed int32 `json:"failed"`
	Result string `json:"result"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Active Jobs",type="string",JSONPath=".status.active",description="The active status of this Scenario"
//+kubebuilder:printcolumn:name="Succeeded Jobs",type="string",JSONPath=".status.succeeded",description="The succeeded status of this Scenario"
//+kubebuilder:printcolumn:name="Failed Jobs",type="string",JSONPath=".status.failed",description="The failed status of this Scenario"
//+kubebuilder:printcolumn:name="Scenario Result",type="string",JSONPath=".status.result",description="The result of this Scenario"
// CascadeManualOperator is the Schema for the cascademanualoperators API
type CascadeManualOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CascadeManualOperatorSpec   `json:"spec,omitempty"`
	Status CascadeManualOperatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CascadeManualOperatorList contains a list of CascadeManualOperator
type CascadeManualOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CascadeManualOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CascadeManualOperator{}, &CascadeManualOperatorList{})
}
