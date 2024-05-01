/*
Copyright 2024.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rmkv1alpha1 "markruler.com/api/v1alpha1"
)

// MachineReconciler reconciles a Machine object
type MachineReconciler struct {
	client.Client
	log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=rmk.markruler.com,resources=machines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rmk.markruler.com,resources=machines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=rmk.markruler.com,resources=machines/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Machine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *MachineReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	
	ctx := context.Background()
	log := r.log.WithValues("machine", req.NamespacedName)

	log.Info("informer => Work Queue => Controller!")

	var machine rmkv1alpha1.Machine

	if err := r.Get(ctx, req.NamespacedName, &machine); err != nil {
		log.Info("error GET Machine", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if machine.Spec.Role == "garbage" {
		if err := r.Delete(ctx, &machine); err != nil {
			log.Info("error DELETE Machine", "deleteName", req.NamespacedName)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		log.Info(">>>> DELETE machine", "deleteName", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	if !machine.Status.Ready {
		log.Info("Machine is not ready")
	}

	if machine.Spec.Role == "" {
		machine.Spec.Role = "garbage"
	}

	if machine.Spec.Role == "worker" {
		machine.Status.Ready = true
	} else {
		machine.Status.Ready = false
	}

	if err := r.Status().Update(ctx, &machine); err != nil {
		log.Info("error UPDATE status", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else {
		log.Info("Update machine", "updateName", req.NamespacedName)
	}


	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rmkv1alpha1.Machine{}).
		Complete(r)
}
