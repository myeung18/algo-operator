/*
Copyright 2021.

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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	cachev1alpha1 "github.com/myeugn18/algo-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AlgoCodingReconciler reconciles a AlgoCoding object
type AlgoCodingReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.algo.com,resources=algocodings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.algo.com,resources=algocodings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.algo.com,resources=algocodings/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AlgoCoding object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *AlgoCodingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("algocoding", req.NamespacedName)

	log := r.Log.WithValues("", req.NamespacedName)
	log.Info("Reconciling algocoding instance")

	// your logic here
	instance := &cachev1alpha1.AlgoCoding{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	podList := &corev1.PodList{}
	lbs := map[string]string{
		"app":     instance.Name,
		"version": "v0.0.1",
	}
	labelSelector := labels.SelectorFromSet(lbs)
	listOps := &client.ListOptions{Namespace: req.Namespace, LabelSelector: labelSelector}
	if err := r.Client.List(context.TODO(), podList, listOps); err != nil {
		return ctrl.Result{}, err
	}

	//count the pods that are pending or running as available
	var available []corev1.Pod
	for _, pod := range podList.Items {
		if pod.ObjectMeta.DeletionTimestamp != nil {
			continue //already deleted
		}
		if pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodPending {
			available = append(available, pod)
		}
	}

	availPodName := []string{}
	for _, pod := range available {
		availPodName = append(availPodName, pod.ObjectMeta.Name)
	}

	status := cachev1alpha1.AlgoCodingStatus{
		PodNames: availPodName,
	}
	if !reflect.DeepEqual(instance.Status, status) {
		instance.Status = status
		err = r.Client.Status().Update(context.TODO(), instance)
		if err != nil {
			log.Error(err, "failed to update instance status ")
			return ctrl.Result{}, err
		}
	}

	numAvailable := int32(len(available))
	if numAvailable > instance.Spec.Replicas {
		log.Info("Scaling down the number of pods is more than expected. going to remove few of them. ")
		diff := numAvailable - instance.Spec.Replicas
		dpods := available[:diff]
		for _, dp := range dpods {
			err = r.Client.Delete(context.TODO(), &dp)
			if err != nil {
				log.Error(err, " failed to dele", "pod name", dp.Name)
				return ctrl.Result{}, err
			}
			log.Info("scaling from corresponding service", "Pod", numAvailable, "Service", instance.Spec.Replicas)
			strPort := dp.Name[strings.LastIndex(dp.Name, "-")+1:]
			sName := instance.Name + "-service-" + strPort
			s := &corev1.Service{}
			err := r.Client.Get(context.TODO(), types.NamespacedName{
				Name:      sName,
				Namespace: req.Namespace,
			}, s) //get the service and set to pointer s
			err = r.Client.Delete(context.TODO(), s) //and delete it
			if err != nil {
				if errors.IsNotFound(err) {
					return ctrl.Result{}, nil //it is already gone
				}
				return ctrl.Result{}, err
			}
		} //

		return ctrl.Result{Requeue: true}, nil //reconcile done, check in another loop.
	}

	if numAvailable < instance.Spec.Replicas {
		log.Info("Scaling up is needed as the number of running pod is below expected.")
		//create new pod instance which means creating a new CR for the application
		pod := newPodForCR(instance)
		if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
			return reconcile.Result{}, err
		}

		err = r.Client.Create(context.TODO(), pod)
		if err != nil {
			log.Error(err, "failed to create new pod", "pod.name", pod.Name)
			return ctrl.Result{}, err
		}

		svc := newServiceForPod(instance)
		if err := controllerutil.SetControllerReference(instance, svc, r.Scheme); err != nil {
			return reconcile.Result{}, err
		}
		err = r.Client.Create(context.TODO(), svc)
		if err != nil {
			log.Error(err, "failed to create a service for the pod", "svc.Name", svc.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

var nextPort = 0

func newServiceForPod(cr *cachev1alpha1.AlgoCoding) *corev1.Service {
	strPort := strconv.Itoa(nextPort)
	labels := map[string]string{
		"app": cr.Name,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-service-" + strPort,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       8089,
				TargetPort: intstr.FromInt(8080),
				NodePort:   int32(nextPort),
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}
}

func newPodForCR(cr *cachev1alpha1.AlgoCoding) *corev1.Pod {
	if nextPort == 0 {
		nextPort = 32000
	} else {
		nextPort++
	}

	strPort := strconv.Itoa(nextPort)
	labels := map[string]string{
		"app":      cr.Name,
		"version":  "v0.0.1",
		"nodePort": strPort,
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod-" + strPort,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "algocoding-web",
					Image:   cr.Spec.WebImage,
					Command: []string{"catalina.sh", "run"},
				},
				{
					Name:    "algocoding-db",
					Image:   cr.Spec.DBImage,
					Env:     []corev1.EnvVar{{Name: "MYSQL_USER", Value: "admin"}, {Name: "MYSQL_PASSWORD", Value: "admin"}, {Name: "MYSQL_ROOT_PASSWORD", Value: "root+1"}, {Name: "MYSQL_DATABASE", Value: "registry"}},
					Command: []string{"/entrypoint.sh"},
					Args:    []string{"mysqld"},
				},
			},
		},
	}

}

// SetupWithManager sets up the controller with the Manager.
func (r *AlgoCodingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.AlgoCoding{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}