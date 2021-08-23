package actor

import (
	"context"
	"github.com/go-logr/logr"
	api "github.com/myeung18/algo-operator/api/v1alpha1"
	"github.com/myeung18/algo-operator/pkg/resource"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Different logging levels
var DEBUGLEVEL = int(zapcore.DebugLevel)
var WARNLEVEL = int(zapcore.WarnLevel)

type Actor interface {
	Handles([]api.AlgoCodingStatus) bool
	Act(context.Context, *resource.Cluster) error
	GetActionType() api.ActionType
}

func NewOperatorActions(scheme *runtime.Scheme, cl client.Client, config *rest.Config) []Actor {

	return []Actor{
		newGenerateCert(scheme, cl, config),
	}
}

var Log = logf.Log.WithName("action")

func newAction(atype string, scheme *runtime.Scheme, cl client.Client) action {
	return action{
		log:    Log.WithValues("action", atype),
		client: cl,
		scheme: scheme,
	}
}

type action struct {
	log    logr.Logger
	client client.Client
	scheme *runtime.Scheme
}
