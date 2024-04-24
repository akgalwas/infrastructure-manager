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

package controller

import (
	"context"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardener_apis "github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	"github.com/go-logr/logr"
	"github.com/kyma-project/infrastructure-manager/internal/controller/providers"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// GardenerClusterProvisioningRequestReconciler reconciles a GardenerClusterProvisioningRequest object
type GardenerClusterProvisioningRequestReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Log         logr.Logger
	ShootClient gardener_apis.ShootInterface
	Enabled     bool
}

//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=gardenerclusterprovisioningrequests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=gardenerclusterprovisioningrequests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructuremanager.kyma-project.io,resources=gardenerclusterprovisioningrequests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GardenerClusterProvisioningRequest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *GardenerClusterProvisioningRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var provisioningRequest imv1.GardenerClusterProvisioningRequest

	err := r.Get(ctx, req.NamespacedName, &provisioningRequest)

	if err != nil {
		r.Log.Error(err, "unable to fetch GardenerClusterProvisioningRequest")
		return ctrl.Result{}, err
	} else {
		r.Log.Info("Reconciling GardenerClusterProvisioningRequest", "Name", provisioningRequest.Name, "Namespace", provisioningRequest.Namespace)
	}

	shoot := &v1beta1.Shoot{}
	shoot.Spec = provisioningRequest.Shoot
	shoot.Name = provisioningRequest.Name
	shoot.Spec.Provider.ControlPlaneConfig = &runtime.RawExtension{Raw: providers.GetGCPControlPlane()}
	shoot.Spec.Networking = &v1beta1.Networking{}
	shoot.Spec.Networking.Nodes = providers.PtrTo("10.180.0.0/16")

	shoot.Spec.Provider.InfrastructureConfig = &runtime.RawExtension{Raw: providers.GetRawGCPInfrastructureConfig()}

	if err != nil {
		r.Log.Error(err, "unable to map GardenerClusterProvisioningRequest to Shoot")
		return ctrl.Result{}, err
	}

	createdShoot, provisioningErr := r.ShootClient.Create(ctx, shoot, v1.CreateOptions{})

	if provisioningErr != nil {
		r.Log.Error(provisioningErr, "unable to create Shoot")
		return ctrl.Result{}, provisioningErr
	}
	r.Log.Info("Shoot created", "Name", createdShoot.Name, "Namespace", createdShoot.Namespace)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GardenerClusterProvisioningRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&imv1.GardenerClusterProvisioningRequest{}, builder.WithPredicates(predicate.Or(
			predicate.LabelChangedPredicate{},
			predicate.AnnotationChangedPredicate{},
			predicate.GenerationChangedPredicate{}),
		)).
		Complete(r)
}
