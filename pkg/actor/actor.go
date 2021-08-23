package actor

import (
	"context"
	api "github.com/myeung18/algo-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Actor interface {
	Handles([]string) bool
	Act(context.Context, *api.AlgoCoding) error
	GetActionType() api.ActionType
}

func NewOperatorActions(scheme *runtime.Scheme, cl client.Client, config *rest.Config) []Actor {

	return []Actor{}
}
