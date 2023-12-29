package manager

import (

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"fmt"
	"io"

	"go.uber.org/zap/zapcore"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	terminalv1 "github.com/joshmeranda/marina-operator/api/v1"
	"github.com/joshmeranda/marina-operator/controllers"
	"github.com/urfave/cli/v2"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(terminalv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func start(ctx *cli.Context) error {
	metricsAddr := ctx.String("metrics-bind-address")
	enableLeaderElection := ctx.Bool("leader-elect")
	probeAddr := ctx.String("health-probe-bind-address")
	webhookPort := ctx.Int("webhook-port")

	opts := zap.Options{
		Development: true,
	}

	if ctx.Bool("quiet") {
		opts.Level = zapcore.WarnLevel
	}

	if ctx.Bool("silent") {
		opts.DestWriter = io.Discard
	}

	if ctx.Bool("verbose") {
		opts.Level = zapcore.DebugLevel
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	var config *rest.Config
	var err error

	if kubeconfig := ctx.String("kubeconfig"); kubeconfig != "" {
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			return fmt.Errorf("failed to get config from kubeconfig: %w", err)
		}
	} else if config, err = rest.InClusterConfig(); err != nil {
		return fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: webhookPort,
		}),
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "261831dc.marina.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return err
	}

	if err = (&controllers.TerminalReconciler{
		Client: mgr.GetClient(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Terminal")
		return err
	}

	if err = (&controllers.UserReconciler{
		Client: mgr.GetClient(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "User")
		return err
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		return err
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctx.Context); err != nil {
		// if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		return err
	}

	return nil
}

func App() cli.App {
	return cli.App{
		Name:        "manager",
		Description: "run the marina operator manager",
		Action:      start,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "metrics-bind-address",
				Usage:   "The address the metric endpoint binds to.",
				EnvVars: []string{"METRICS_BIND_ADDRESS"},
				Value:   ":8080",
			},
			&cli.StringFlag{
				Name:    "health-probe-bind-address",
				Usage:   "The address the probe endpoint binds to.",
				EnvVars: []string{"HEALTH_PROBE_BIND_ADDRESS"},
				Value:   ":8081",
			},
			&cli.BoolFlag{
				Name:    "leader-elect",
				Usage:   "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.",
				EnvVars: []string{"LEADER_ELECT"},
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "kubeconfig",
				Usage:   "the path to the kubeconfig file to use for the terminal",
				EnvVars: []string{"KUBECONFIG"},
				Aliases: []string{"f"},
			},
			&cli.IntFlag{
				Name:    "webhook-port",
				Usage:   "the port for the webhook server to listen on",
				Aliases: []string{"p"},
				EnvVars: []string{"WEBHOOK_PORT"},
				Value:   9443,
			},

			&cli.BoolFlag{
				Name:    "quiet",
				Usage:   "suppress all output except for warnings and errors",
				Aliases: []string{"q"},
			},
			&cli.BoolFlag{
				Name:    "silent",
				Usage:   "suppress all output",
				Aliases: []string{"s"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "run verbosely",
				Aliases: []string{"v"},
			},
		},
	}
}
