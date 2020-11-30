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
	"strings"
	"time"

	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	cloudv1 "github.ibm.com/josiah-sams/hack-vpc-operator/api/v1"
	vsi "github.ibm.com/josiah-sams/hack-vpc-operator/controllers/vsi"
)

// VSIReconciler reconciles a VSI object
type VSIReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	vs     *vsi.Service
}

const (
	port             = 443
	targetPort       = 443
	serviceFinalizer = "vsi.cloud.ibm.com"
)

// +kubebuilder:rbac:groups=cloud.ibm.com,resources=vsis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cloud.ibm.com,resources=vsis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cloud.ibm.com,resources=vsis/finalizers,verbs=update;
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/finalizers,verbs=get;list;watch;create;update;patch;delete;
// +kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch;create;update;patch;delete

// Reconcile ..
func (r *VSIReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logt := r.Log.WithValues("vsi", req.NamespacedName)

	// Fetch the Service instance
	instance := &cloudv1.VSI{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	r.Log.Info(instance.Spec.APIKey)

	// Delete if necessary
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// Instance is not being deleted, add the finalizer if not present
		if !containsServiceFinalizer(instance) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, serviceFinalizer)
			if err := r.Update(ctx, instance); err != nil {
				logt.Error(err, "Error adding finalizer", "service", instance.ObjectMeta.Name)
				// TODO(johnstarich): Shouldn't this update the status with the failure message?
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsServiceFinalizer(instance) {

			if r.vs != nil {
				err = vsi.DeleteVSI(r.vs)
				if err != nil {
					logt.Error(err, "Error with DeleteVSI")
				}
			}

			if err := r.deleteService(req, *instance, "vservice"); err != nil {
				logt.Error(err, "Error deleting resource", "service", instance.ObjectMeta.Name)
				// TODO(johnstarich): Shouldn't this return the error so it will be logged?
				return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 10}, nil
			}

			// remove our finalizer from the list and update it.
			instance.ObjectMeta.Finalizers = deleteServiceFinalizer(instance)
			err = r.Update(ctx, instance)
			if err != nil {
				logt.Error(err, "Error removing finalizers")
			}
			return ctrl.Result{}, err
		}
	}

	if r.vs == nil {
		r.vs, err = vsi.GetOrCreate(instance.Spec.APIKey, instance.Spec.SSHKey)
	} else {
		err = vsi.GetOrCreateWithObj(r.vs, instance.Spec.APIKey, instance.Spec.SSHKey)
	}

	if err != nil {
		logt.Error(err, "VSI GetOrCreate failed")
		return ctrl.Result{}, err
	}

	// populate the IP address from the Cloud
	ipaddress := r.vs.GetAddress()

	if ipaddress != instance.Status.IPAddress {
		instance.Status.IPAddress = ipaddress

		if err := r.deleteService(req, *instance, "vservice"); err != nil {
			return ctrl.Result{}, err
		} else if service, err := r.createService(req, *instance, port, targetPort, "vservice"); err != nil {
			return ctrl.Result{}, err
		} else {
			if err := r.createEndpoint(req, service, ipaddress, targetPort); err != nil {
				return ctrl.Result{}, err
			}
			if err := r.createRoute(req, service); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Use the Cloud VSI status
	instanceState := r.vs.GetStatus()

	instance.Status.State = getState(instanceState)
	instance.Status.Message = getState(instanceState)

	_ = r.Status().Update(ctx, instance)

	return ctrl.Result{Requeue: true, RequeueAfter: 20 * time.Second}, nil
}

// containsServiceFinalizer checks if the instance contains service finalizer
func containsServiceFinalizer(instance *cloudv1.VSI) bool {
	for _, finalizer := range instance.ObjectMeta.Finalizers {
		if strings.Contains(finalizer, serviceFinalizer) {
			return true
		}
	}
	return false
}

// deleteServiceFinalizer delete service finalizer
func deleteServiceFinalizer(instance *cloudv1.VSI) []string {
	var result []string
	for _, finalizer := range instance.ObjectMeta.Finalizers {
		if finalizer == serviceFinalizer {
			continue
		}
		result = append(result, finalizer)
	}
	return result
}

func getState(serviceInstanceState string) string {
	if serviceInstanceState == "running" || serviceInstanceState == "active" || serviceInstanceState == "provisioned" {
		return "Online"
	}
	return serviceInstanceState
}

func (r *VSIReconciler) createRoute(req ctrl.Request, service *corev1.Service) error {
	ctx := context.Background()
	logt := r.Log.WithValues("vsi", req.NamespacedName)

	routeRef := service.GetObjectMeta().GetName()
	logt.Info("check routes", "name", routeRef)

	route := &routev1.Route{}

	if err := r.Get(ctx, client.ObjectKey{Namespace: service.Namespace, Name: routeRef}, route); err != nil {
		if !errors.IsNotFound(err) {
			// Error reading the object - requeue the request.
			logt.Info("object cannot be read", "error", err)
			return err
		}
		logt.Info("creating route", "name", routeRef)
		// else : proceed to create the route;
	} else {
		if route.DeletionTimestamp != nil {
			// route is being deleted... nothing to do.
			return nil
		}

		// return;
		return nil
	}

	route = &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      routeRef,
			Namespace: service.Namespace,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: routeRef,
			},
			TLS: &routev1.TLSConfig{
				Termination: routev1.TLSTerminationPassthrough,
			},
		},
	}

	if err := controllerutil.SetControllerReference(service, route, r.Scheme); err != nil {
		logt.Error(err, "unable to set owner reference on new route")
		return err
	}

	if err := r.Create(ctx, route); err != nil {
		logt.Error(err, "failed to update route (retrying)")
		return err
	}

	logt.Info("route created", "name", routeRef)

	return nil
}

func (r *VSIReconciler) createEndpoint(req ctrl.Request, service *corev1.Service, ipaddress string, port int32) error {
	ctx := context.Background()
	logt := r.Log.WithValues("vsi", req.NamespacedName)

	endpointRef := service.GetObjectMeta().GetName()
	logt.Info("check endpoint", "name", endpointRef)

	endpoints := &corev1.Endpoints{}

	if err := r.Get(ctx, client.ObjectKey{Namespace: service.Namespace, Name: endpointRef}, endpoints); err != nil {
		if !errors.IsNotFound(err) {
			// Error reading the object - requeue the request.
			logt.Info("object cannot be read", "error", err)
			return err
		}
		logt.Info("creating endpoint", "name", endpointRef)
		// else : proceed to create the endpoints;
	} else {
		if endpoints.DeletionTimestamp != nil {
			// endpoints is being deleted... nothing to do.
			return nil
		}

		// return;
		return nil
	}

	endpoints = &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      endpointRef,
			Namespace: service.Namespace,
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP: ipaddress,
					},
				},
				Ports: []corev1.EndpointPort{
					{
						Port: port,
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(service, endpoints, r.Scheme); err != nil {
		logt.Error(err, "unable to set owner reference on new endpoints")
		return err
	}

	if err := r.Create(ctx, endpoints); err != nil {
		logt.Error(err, "failed to update endpoints (retrying)")
		return err
	}
	logt.Info("endpoints created", "name", endpointRef)

	return nil
}

func (r *VSIReconciler) deleteService(req ctrl.Request, instance cloudv1.VSI, key string) error {
	ctx := context.Background()
	logt := r.Log.WithValues("vsi", req.NamespacedName)

	serviceRef := instance.GetObjectMeta().GetName() + "-service-" + key
	logt.Info("check service", "name", serviceRef)
	service := &corev1.Service{}

	if err := r.Get(ctx, client.ObjectKey{Namespace: instance.Namespace, Name: serviceRef}, service); err != nil {
		if errors.IsNotFound(err) {
			// The resource is not found; return
			return nil
		} else {
			// error retrieving service details
			return err
		}
	} else { // else : proceed to delete the service;
		if service.DeletionTimestamp != nil {
			// service is being deleted... nothing to do.
			return nil
		}
	}
	if err := r.Delete(ctx, service); err != nil {
		logt.Error(err, "failed to update service (retrying)")
		return err
	}
	logt.Info("service deleted", "name", serviceRef)
	// return
	return nil
}

func (r *VSIReconciler) createService(req ctrl.Request, instance cloudv1.VSI, port int32, targetport int, key string) (*corev1.Service, error) {
	ctx := context.Background()
	logt := r.Log.WithValues("vsi", req.NamespacedName)

	serviceRef := instance.GetObjectMeta().GetName() + "-service-" + key
	logt.Info("check service", "name", serviceRef)
	service := &corev1.Service{}

	if err := r.Get(ctx, client.ObjectKey{Namespace: instance.Namespace, Name: serviceRef}, service); err != nil {
		if !errors.IsNotFound(err) {
			// Error reading the object - requeue the request.
			logt.Info("object cannot be read", "error", err)
			return nil, err
		}
		logt.Info("creating service", "name", serviceRef)
		// else : proceed to create the service;
	} else {
		if service.DeletionTimestamp != nil {
			// service is being deleted... nothing to do.
			return service, nil
		}

		// return the service;
		return service, nil
	}

	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceRef,
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       port,
					TargetPort: intstr.FromInt(targetport),
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(&instance, service, r.Scheme); err != nil {
		logt.Error(err, "unable to set owner reference on new service")
		return nil, err
	}

	if err := r.Create(ctx, service); err != nil {
		logt.Error(err, "failed to update service (retrying)")
		return nil, err
	}
	logt.Info("service created", "name", serviceRef)

	return service, nil
}

// SetupWithManager ..
func (r *VSIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1.VSI{}).
		Complete(r)
}
