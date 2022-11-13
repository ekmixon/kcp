/*
Copyright 2021 The KCP Authors.

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

package workspace

import (
	"context"
	"time"

	"github.com/kcp-dev/logicalcluster/v2"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"

	tenancyv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/tenancy/v1alpha1"
	conditionsv1alpha1 "github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/apis/conditions/v1alpha1"
	"github.com/kcp-dev/kcp/pkg/apis/third_party/conditions/util/conditions"
)

type phaseReconciler struct {
	getThisWorkspace func(ctx context.Context, cluster logicalcluster.Name) (*tenancyv1alpha1.ThisWorkspace, error)

	requeueAfter func(workspace *tenancyv1alpha1.ClusterWorkspace, after time.Duration)
}

func (r *phaseReconciler) reconcile(ctx context.Context, workspace *tenancyv1alpha1.ClusterWorkspace) (reconcileStatus, error) {
	logger := klog.FromContext(ctx).WithValues("reconciler", "phase")

	switch workspace.Status.Phase {
	case tenancyv1alpha1.WorkspacePhaseScheduling:
		if workspace.Status.Cluster != "" {
			workspace.Status.Phase = tenancyv1alpha1.WorkspacePhaseInitializing
		}
	case tenancyv1alpha1.WorkspacePhaseInitializing:
		logger = logger.WithValues("cluster", workspace.Status.Cluster)

		this, err := r.getThisWorkspace(ctx, logicalcluster.New(workspace.Status.Cluster))
		if err != nil && !apierrors.IsNotFound(err) {
			return reconcileStatusStopAndRequeue, err
		} else if apierrors.IsNotFound(err) {
			logger.Info("ThisWorkspace disappeared")
			conditions.MarkFalse(workspace, tenancyv1alpha1.WorkspaceInitialized, tenancyv1alpha1.WorkspaceInitializedWorkspaceDisappeared, conditionsv1alpha1.ConditionSeverityError, "ThisWorkspace disappeared")
			return reconcileStatusContinue, nil
		}

		workspace.Status.Initializers = this.Status.Initializers

		if initializers := workspace.Status.Initializers; len(initializers) > 0 {
			after := time.Since(this.CreationTimestamp.Time) / 5
			logger.V(3).Info("ThisWorkspace still has initializers, requeueing", "initializers", initializers, "after", after)
			conditions.MarkFalse(workspace, tenancyv1alpha1.WorkspaceInitialized, tenancyv1alpha1.WorkspaceInitializedInitializerExists, conditionsv1alpha1.ConditionSeverityInfo, "Initializers still exist: %v", workspace.Status.Initializers)
			r.requeueAfter(workspace, after)
			return reconcileStatusContinue, nil
		}

		logger.V(3).Info("ThisWorkspace is ready")
		workspace.Status.Phase = tenancyv1alpha1.WorkspacePhaseReady
		conditions.MarkTrue(workspace, tenancyv1alpha1.WorkspaceInitialized)
	}

	return reconcileStatusContinue, nil
}