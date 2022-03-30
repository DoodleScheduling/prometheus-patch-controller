/*
Copyright 2022 Doodle.

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

package main

import (
	"flag"
	"os"
	"strings"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	metricsinfradoodlecomv1beta1 "github.com/doodlescheduling/k8sprom-patch-controller/api/v1beta1"
	"github.com/doodlescheduling/k8sprom-patch-controller/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(metricsinfradoodlecomv1beta1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

var (
	metricsAddr             = ":9556"
	probesAddr              = ":9557"
	webhookPort             = 9443
	enableLeaderElection    = false
	leaderElectionNamespace string
	namespaces              = ""
	concurrent              = 2
	fieldManager            = "k8sprom-patch-controller"
)

func main() {
	flag.StringVar(&metricsAddr, "metrics-addr", metricsAddr, "The address of the metric endpoint binds to.")
	flag.StringVar(&probesAddr, "probe-addr", probesAddr, "The address of the probe endpoints bind to.")
	flag.StringVar(&fieldManager, "field-manager", fieldManager, "The name of the field maanger used for server side apply https://kubernetes.io/docs/reference/using-api/server-side-apply/.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", enableLeaderElection,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&leaderElectionNamespace, "leader-election-namespace", leaderElectionNamespace,
		"Specify a different leader election namespace. It will use the one where the controller is deployed by default.")
	flag.StringVar(&namespaces, "namespaces", namespaces,
		"The controller listens by default for all namespaces. This may be limited to a comma delimted list of dedicated namespaces.")
	flag.IntVar(&concurrent, "concurrent", concurrent,
		"The number of concurrent reconcile workers. By default this is 2.")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Import flags into viper and bind them to env vars
	// flags are converted to upper-case, - is replaced with _
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		setupLog.Error(err, "Failed parsing command line arguments")
		os.Exit(1)
	}

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	opts := ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      viper.GetString("metrics-addr"),
		Port:                    viper.GetInt("webhoook-port"),
		HealthProbeBindAddress:  viper.GetString("probe-addr"),
		LeaderElection:          viper.GetBool("enable-leader-election"),
		LeaderElectionNamespace: viper.GetString("leader-election-namespace"),
		LeaderElectionID:        "k8sprom-patch-controller.infra.doodle.com",
	}

	ns := strings.Split(viper.GetString("namespaces"), ",")
	if len(ns) > 0 && ns[0] != "" {
		opts.NewCache = cache.MultiNamespacedCacheBuilder(ns)
		setupLog.Info("watching dedicated namespaces", "namespaces", ns)
	} else {
		setupLog.Info("watching all namespaces")
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), opts)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	dynClient, err := dynamic.NewForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to create dynamic client")
		os.Exit(1)
	}

	if err = (&controllers.PrometheusPatchRuleReconciler{
		Client:       mgr.GetClient(),
		DynClient:    dynClient,
		FieldManager: fieldManager,
		Log:          ctrl.Log.WithName("controllers").WithName("PrometheusPatchRule"),
		Scheme:       mgr.GetScheme(),
		Recorder:     mgr.GetEventRecorderFor("PrometheusPatchRule"),
	}).SetupWithManager(mgr, controllers.PrometheusPatchRuleReconcilerOptions{MaxConcurrentReconciles: viper.GetInt("concurrent")}); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "PrometheusPatchRule")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
