/*
Copyright 2023 zsh.

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

package controller

import (
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "github.com/lengdanran/hello-operator/api/v1"
)

// HelloAppReconciler reconciles a HelloApp object
type HelloAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.zsh.io,resources=helloapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.zsh.io,resources=helloapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.zsh.io,resources=helloapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HelloApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *HelloAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)

	// TODO(user): your logic here
	logger := log.FromContext(ctx)

	// 获取 HelloApp
	app := &appsv1.HelloApp{}
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		if errors.IsNotFound(err) {
			// HelloApp CRD is deleted, delete the associated Pods
			logger.Info("HelloApp CRD is deleted, delete the associated Pods")
			return ctrl.Result{}, nil
		}

		logger.Error(err, "fail to get the Application")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, err
	}
	// 获取当前的pod数量
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(req.Namespace), client.MatchingLabels{"label": app.Spec.Label}); err != nil {
		return ctrl.Result{}, err
	}
	podCount := int32(len(podList.Items))

	if podCount != app.Spec.Replicas {
		logger.Info("Updating pod Replicas.....")
		logger.Info(fmt.Sprintf("Current pod cnt = %v", podCount))
		// Adjust the Pod count based on the desired replicas
		if podCount < app.Spec.Replicas {
			logger.Info(fmt.Sprintf("Less than desired replicas %v", app.Spec.Replicas))
			diff := app.Spec.Replicas - podCount
			for i := int32(0); i < diff; i++ {
				pod := newPodForCR(app)
				if err := controllerutil.SetControllerReference(app, pod, r.Scheme); err != nil {
					return ctrl.Result{}, err
				}
				if err := r.Create(ctx, pod); err != nil {
					return ctrl.Result{}, err
				}
			}
		} else if podCount > app.Spec.Replicas {
			logger.Info(fmt.Sprintf("More than desired replicas %v", app.Spec.Replicas))
			diff := podCount - app.Spec.Replicas

			for i, pod := range podList.Items {
				err := r.Delete(ctx, &pod)
				if err != nil {
					logger.Info("Delete pod failed......")
					return ctrl.Result{}, err
				}
				if int32(i) == diff-1 {
					break
				}
			}
			return ctrl.Result{}, nil
		}
	}
	logger.Info(fmt.Sprintf("Current pod cnt matches with desired cnt %v", app.Spec.Replicas))
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelloAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.HelloApp{}).
		Complete(r)
}

func newPodForCR(cr *appsv1.HelloApp) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", cr.Name),
			Namespace:    cr.Namespace,
			Labels:       cr.Labels,
		},
		Spec: cr.Spec.Template.Spec,
	}
}
