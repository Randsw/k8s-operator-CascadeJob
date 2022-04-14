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

package controllers

import (
	"context"
	"encoding/json"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	//"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cascadev1alpha1 "github.com/Randsw/CascadeManualOperator/api/v1alpha1"
)

// CascadeManualOperatorReconciler reconciles a CascadeManualOperator object
type CascadeManualOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cascade.cascade.net,resources=cascademanualoperators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cascade.cascade.net,resources=cascademanualoperators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cascade.cascade.net,resources=cascademanualoperators/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CascadeManualOperator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CascadeManualOperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("CascadeManualOperator", req.NamespacedName)

	logger.Info("Reconciling CascadeManualOperator", "request name", req.Name, "request namespace", req.Namespace)

	instance := &cascadev1alpha1.CascadeManualOperator{}
	instanceType := "CascadeManualOperator"

	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("Resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get CascadeManualOperator.")
		return ctrl.Result{}, err
	}
	// Not prepared for deletion
	if !instance.GetDeletionTimestamp().IsZero() {
		logger.Info("Prepared for deletition")
		return ctrl.Result{}, nil
	}

	jobReplicas := int32(1)

	// Check if the Job already exists, if not create a new one
	found := &batchv1.Job{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-job", Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new job
		job := r.getJob(instance, &jobReplicas, req.Name, instanceType)
		logger.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
		err = r.Create(ctx, job)
		if err != nil {
			logger.Error(err, "Failed to create new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get Job")
		return ctrl.Result{}, err
	}

	foundMap := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-cm", Namespace: instance.Namespace}, foundMap)
	if err != nil && errors.IsNotFound(err) {
		cm := r.getCm(instance)
		logger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
		err = r.Create(ctx, cm)
		if err != nil {
			logger.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		logger.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	instance.Status = &found.Status

	//return ctrl.Result{}, r.Client.Status().Update(context.TODO(), instance)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CascadeManualOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cascadev1alpha1.CascadeManualOperator{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}

func (r *CascadeManualOperatorReconciler) getJob(instance *cascadev1alpha1.CascadeManualOperator, replicas *int32, reqName, instanceType string) *batchv1.Job {
	var jobAffinity = corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{{
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{{
						Key:      instanceType,
						Operator: "In",
						Values:   []string{reqName},
					}},
				},
				TopologyKey: "kubernetes.io/hostname",
			}},
		},
	}

	var podSpec = instance.Spec.Template
	podSpec.Spec.Affinity = &jobAffinity

	if podSpec.Spec.RestartPolicy == "Always" {
		podSpec.Spec.RestartPolicy = "OnFailure"
	}
	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-job",
			Namespace: instance.Namespace,
			Labels:    instance.Labels,
		},
		Spec: batchv1.JobSpec{
			Parallelism:             replicas,
			Completions:             replicas,
			Template:                podSpec,
			TTLSecondsAfterFinished: instance.Spec.TTLSecondsAfterFinished,
			BackoffLimit:            instance.Spec.BackoffLimit,
			ActiveDeadlineSeconds:   instance.Spec.ActiveDeadlineSeconds,
		},
	}
	ctrl.SetControllerReference(instance, job, r.Scheme)
	return job
}

func (r *CascadeManualOperatorReconciler) getCm(instance *cascadev1alpha1.CascadeManualOperator) *corev1.ConfigMap {
	data, _ := json.Marshal(instance.Spec.ScenarioConfig)
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-cm",
			Namespace: instance.Namespace,
			Labels:    instance.Labels,
		},
		Data: map[string]string{
			"Configuration": string(data),
		},
	}

	ctrl.SetControllerReference(instance, cm, r.Scheme)
	return cm
}
