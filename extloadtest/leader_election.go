package extloadtest

import (
	"context"
	"errors"
	"flag"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"time"
)

var (
	IAmTheLeader bool
)

func createClientset() (*kubernetes.Clientset, string) {
	config, err := rest.InClusterConfig()
	if err == nil {
		log.Info().Msgf("Extension is running inside a cluster, config found")
	} else if errors.Is(err, rest.ErrNotInCluster) {
		log.Info().Msgf("Extension is not running inside a cluster, try local .kube config")
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	}

	if err != nil {
		log.Fatal().Err(err).Msgf("Could not find kubernetes config")
	}

	config.UserAgent = "steadybit-extension-kubernetes"
	config.Timeout = time.Second * 10
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not create kubernetes client")
	}

	info, err := clientset.ServerVersion()
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not fetch server version.")
	}

	log.Info().Msgf("Cluster connected! Kubernetes Server Version %+v", info)

	return clientset, config.APIPath
}

func InitLeaderElection() {
	go func() {
		id := uuid.New().String()
		leaseLockName := "loadtest-leader-election"
		leaseLockNamespace := "steadybit-agent"

		client, _ := createClientset()

		// use a Go context so we can tell the leaderelection code when we
		// want to step down
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// listen for interrupts or the Linux SIGTERM signal and cancel
		// our context, which the leader election code will observe and
		// step down
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-ch
			klog.Info("Received termination, signaling shutdown")
			cancel()
		}()

		// we use the Lease lock type since edits to Leases are less common
		// and fewer objects in the cluster watch "all Leases".
		lock := &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Name:      leaseLockName,
				Namespace: leaseLockNamespace,
			},
			Client: client.CoordinationV1(),
			LockConfig: resourcelock.ResourceLockConfig{
				Identity: id,
			},
		}

		// start the leader election code loop
		leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
			Lock: lock,
			// IMPORTANT: you MUST ensure that any code you have that
			// is protected by the lease must terminate **before**
			// you call cancel. Otherwise, you could have a background
			// loop still running and another process could
			// get elected before your background loop finished, violating
			// the stated goal of the lease.
			ReleaseOnCancel: true,
			LeaseDuration:   60 * time.Second,
			RenewDeadline:   15 * time.Second,
			RetryPeriod:     5 * time.Second,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {
					// we're notified when we start - this is where you would
					// usually put your code
					log.Info().Msgf("i am the leader: %s", id)
					IAmTheLeader = true
				},
				OnStoppedLeading: func() {
					IAmTheLeader = false
					// we can do cleanup here
					log.Info().Msgf("i am not the leader anymore: %s", id)
				},
				OnNewLeader: func(identity string) {
					// we're notified when new leader elected
					if identity == id {
						// I just got the lock
						return
					}
					log.Debug().Msgf("new leader elected: %s", identity)
				},
			},
		})
	}()

}
