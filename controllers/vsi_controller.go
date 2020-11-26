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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"

	cloudv1 "github.ibm.com/josiah-sams/hack-vpc-operator/api/v1"
)

// VSIReconciler reconciles a VSI object
type VSIReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cloud.ibm.com,resources=vsis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloud.ibm.com,resources=vsis/status,verbs=get;update;patch

func (r *VSIReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("vsi", req.NamespacedName)

	// Fetch the Service instance
	instance := &cloudv1.VSI{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	r.Log.Info(instance.Spec.APIKey)

	instance.Status.IPAddress = "8.8.8.8"

	_ = r.Status().Update(ctx, instance)

	return ctrl.Result{}, nil
}

func (r *VSIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1.VSI{}).
		Complete(r)
}
