/*
Copyright 2023.

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

	"k8s.io/apimachinery/pkg/runtime"
	apimmgmtv1 "no.malvik/apimops/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ApimServiceReconciler reconciles a ApimService object
type ApimServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apimmgmt.no.malvik,resources=apimservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apimmgmt.no.malvik,resources=apimservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apimmgmt.no.malvik,resources=apimservices/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ApimServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var apimmgmtv1 apimmgmtv1.ApimService
	if err := r.Get(ctx, req.NamespacedName, &apimmgmtv1); err != nil {
		l.Error(err, "unable to fetch Foo")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	apiUrl := apimmgmtv1.Spec.ApiUrl

	l.Info("Dies ist ein Test")
	l.Info(apiUrl)
	//resp, err := http.Get(myapim.Spec.ApiUrl)
	/*	resp, err := http.Get("https://conferenceapi.azurewebsites.net")
		if err != nil {
			log.Error(err, "Fehler")
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err, "Fehler 2")
		}

		sb := string(body)
		log.Info(sb)
	*/
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApimServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apimmgmtv1.ApimService{}).
		Complete(r)
}
