package resource

import (
	api "github.com/myeung18/algo-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

const (
	RELATED_IMAGE_PREFIX         = "RELATED_IMAGE_COCKROACH_"
	NotSupportedVersion          = "not_supported_version"
	CrdbContainerImageAnnotation = "crdb.io/containerimage"
	CrdbVersionAnnotation        = "crdb.io/version"
	CrdbHistoryAnnotation        = "crdb.io/history"
	CrdbRestartAnnotation        = "crdb.io/restart"
	CrdbCertExpirationAnnotation = "crdb.io/certexpiration"
	CrdbRestartTypeAnnotation    = "crdb.io/restarttype"
)


type Cluster struct {

	cr *api.AlgoCoding
	scheme *runtime.Scheme
	initTime metav1.Time
}

func (cluster Cluster) Spec() *api.AlgoCodingSpec {
	return cluster.cr.Spec.DeepCopy()
}

func (cluster Cluster) Status() *api.AlgoCodingStatus {
	return cluster.cr.Status.DeepCopy()
}
func (cluster Cluster) Name() string {
	return cluster.cr.Name
}

func (cluster Cluster) Namespace() string {
	return cluster.cr.Namespace
}

func (cluster Cluster) ObjectKey() types.NamespacedName {
	return types.NamespacedName{Namespace: cluster.Namespace(), Name: cluster.Name()}
}
