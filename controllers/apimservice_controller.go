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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armapimanagement "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement"
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
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	apiUrl := apimmgmtv1.Spec.ApiUrl
	l.Info(apiUrl)

	// Get Azure credentials
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		l.Error(err, "unable to get Azure credentials")
		return ctrl.Result{}, err
	}

	// Get subscription ID
	subscriptionID := apimmgmtv1.Spec.SubscriptionId

	// Create a new API Management client
	// Create the client
	client, err := armapimanagement.NewAPIClient(subscriptionID, cred, nil)
	if err != nil {
		l.Error(err, "unable to get API Management client")
		return ctrl.Result{}, err
	}
	client.BeginCreateOrUpdate(ctx,
		apimmgmtv1.Spec.ResourceGroup,
		apimmgmtv1.Spec.ServiceName,
		"conference",
		armapimanagement.APICreateOrUpdateParameter{
			Properties: &armapimanagement.APICreateOrUpdateProperties{
				Description: to.Ptr(apimmgmtv1.Spec.Description),
				DisplayName: to.Ptr(apimmgmtv1.Spec.DisplayName),
				Path:        to.Ptr("conf"),
				Protocols: []*armapimanagement.Protocol{
					to.Ptr(armapimanagement.Protocol("https")),
					to.Ptr(armapimanagement.Protocol("http"))},
				Format:    to.Ptr(armapimanagement.ContentFormat("swagger-link-json")),
				Value:     to.Ptr(apiUrl),
				IsCurrent: to.Ptr(true),
				IsOnline:  to.Ptr(true),
			},
		},
		&armapimanagement.APIClientBeginCreateOrUpdateOptions{IfMatch: nil})

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApimServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apimmgmtv1.ApimService{}).
		Complete(r)
}
