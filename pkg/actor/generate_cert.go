package actor

import (
	"context"
	api "github.com/myeung18/algo-operator/api/v1alpha1"
	"github.com/myeung18/algo-operator/pkg/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const defaultKeySize = 2048

// We use 366 days on certificate lifetimes to at least match X years,
// otherwise leap years risk putting us just under.
const defaultCALifetime = 10 * 366 * 24 * time.Hour  // ten years
const defaultCertLifetime = 5 * 366 * 24 * time.Hour // five years

// Options settable via command-line flags. See below for defaults.
var keySize int
var caCertificateLifetime time.Duration
var certificateLifetime time.Duration
var allowCAKeyReuse bool
var overwriteFiles bool
var generatePKCS8Key bool

// generateCert issues node and root client certificates via Kubernetes cluster CA
type generateCert struct {
	action

	config   *rest.Config
	CertsDir string
	CAKey    string
}

func (g *generateCert) GetActionType() api.ActionType {
	return api.RequestCertAction
}

func (rc *generateCert) Handles(conds []api.AlgoCodingStatus) bool {
	return true
}

func (rc *generateCert) Act(ctx context.Context, cluster *resource.Cluster) error {

	log := rc.log.WithValues("algoCluster", cluster.ObjectKey())
	if !cluster.Spec().TLSEnabled || cluster.Spec().NodeTLSSecret != "" {
		log.V(DEBUGLEVEL).Info("Skipping TLS cert generation", "enabled", cluster.Spec().TLSEnabled, "secret", cluster.Spec().NodeTLSSecret)
	}



	return nil
}

func newGenerateCert(scheme *runtime.Scheme, cl client.Client, config *rest.Config) Actor {
	return &generateCert{
		action: newAction("", scheme, cl),
		config: config,
	}
}
