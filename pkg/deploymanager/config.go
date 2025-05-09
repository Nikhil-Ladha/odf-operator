package deploymanager

import (
	"context"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	opv1 "github.com/operator-framework/api/pkg/operators/v1"
	opv1a1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"go.uber.org/zap/zapcore"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	odfv1a1 "github.com/red-hat-storage/odf-operator/api/v1alpha1"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(k8sscheme.AddToScheme(scheme))
	utilruntime.Must(odfv1a1.AddToScheme(scheme))
	utilruntime.Must(opv1.AddToScheme(scheme))
	utilruntime.Must(opv1a1.AddToScheme(scheme))
	utilruntime.Must(extv1.AddToScheme(scheme))
}

type DeployManager struct {
	Client client.Client
	Ctx    context.Context
	Log    logr.Logger
}

// NewDeployManager creates a DeployManager struct with default configuration
func NewDeployManager() (*DeployManager, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		return nil, fmt.Errorf("no KUBECONFIG environment variable set")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	log.SetLogger(zap.New(zap.UseDevMode(true), zap.StacktraceLevel(zapcore.FatalLevel)))

	return &DeployManager{
		Client: client,
		Log:    log.Log.WithName("DeployManager"),
		Ctx:    context.TODO(),
	}, nil
}
