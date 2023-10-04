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

package controllers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
	//"strings"
	//"io/ioutil"

	//scp "github.com/bramvdbogaerde/go-scp"
	"github.com/sfreiberg/simplessh"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	fivegv1alpha1 "github.com/Lim-Ho-Keun/infra-operator/api/v1alpha1"
)

// InfraReconciler reconciles a Infra object
type InfraReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=fiveg.kt.com,resources=infras,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fiveg.kt.com,resources=infras/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fiveg.kt.com,resources=infras/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Infra object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile

/*
type access struct {
	login    string
	password string
}
var loginAccess []access
loginAccess = append(loginAccess, access{"root", "swinfra12"})
var client *simplessh.Client

// Try to connect with first password, then tried second else fails gracefully

for _, credentials := range loginAccess {
	if client, err = simplessh.ConnectWithPassword("172.10.1.23", credentials.login, credentials.password); err == nil {
		break
	}
}

if err != nil {
	fmt.Println("888888888888888888888888")
	return ctrl.Result{}, err
}

defer client.Close()
*/

/*
	type access struct {
		login    string
		password string
	}

var loginAccess []access
var sshclient *simplessh.Client
*/
func (r *InfraReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)
	log.Info("")

	var err2 error

	instance := &fivegv1alpha1.Infra{}
	//err := r.Get(context.TODO(), req.NamespacedName, instance)

	type access struct {
		login    string
		password string
	}

	var loginAccess []access
	loginAccess = append(loginAccess, access{"root", "swinfra12"})
	var sshclient *simplessh.Client

	// Try to connect with first password, then tried second else fails gracefully

	for _, credentials := range loginAccess {
		if sshclient, err2 = simplessh.ConnectWithPassword("172.10.1.23", credentials.login, credentials.password); err2 == nil {
			break
		}
	}

	if err2 != nil {
		return ctrl.Result{}, err2
	}
	/*
		client, err2 := scp.NewClientBySSH(sshclient)
		if err2 != nil {
			fmt.Println("Error creating new SSH session from existing connection", err)
		}
	*/
	defer sshclient.Close()

	err := r.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/vs-service.yaml"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/gateway.yaml"); err2 != nil {
			}

			//istio 삭제하기

			if _, err2 := sshclient.Exec("/root/files/istioctl x uninstall --purge -y"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/istio_namespace.yaml"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/sc-wffc-legacy-csi.yaml"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/cr_hostpath.yaml -n hostpath-provisioner"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/operator.yaml -n hostpath-provisioner"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/webhook.yaml -n hostpath-provisioner"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/sc_namespace.yaml"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/cert-manager.yaml"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("kubectl delete -f /root/files/multus-daemonset-thick.yml"); err2 != nil {
			}

			if _, err2 := sshclient.Exec("rm -rf /root/files"); err2 != nil {
			}

			//if _, err2 := sshclient.Exec("rm -rf /var/hpvolumes"); err2 != nil {
			//}

			//if _, err2 := sshclient.Exec("rm -rf /root/files"); err2 != nil {
			//}

			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		} else {
			return reconcile.Result{}, err
		}
	}

	// Now run the commands on the remote machine:
	//exec.Command("bash", "-c", "apt-get install -y sshpass")
	/*
		fmt.Println("files in /")
		files, err := ioutil.ReadDir("/")
		if err != nil {
		}

		for _, file := range files {
			fmt.Println(file.Name(), file.IsDir())
		}

		fmt.Println("files in files")
		files2, err := ioutil.ReadDir("/files")
		if err != nil {
		}
		for _, file2 := range files2 {
			fmt.Println(file2.Name(), file2.IsDir())
		}
		var cmderr error
		cmd := exec.Command("bash", "-c", "ls -al")
		cmd.Stdout = os.Stdout
		if cmderr = cmd.Run(); cmderr != nil {
			fmt.Println(cmderr)
		}

		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	*/
	/*
		cmd = exec.Command("apt-get install -y sshpass")
		cmd.Stdout = os.Stdout
		if cmderr := cmd.Run(); cmderr != nil {
			fmt.Println(cmderr)
		}
	*/
	var cmderr error
	cmd := exec.Command("bash", "-c", "sshpass -p swinfra12 scp -r -o StrictHostKeyChecking=no /files root@172.10.1.23:/root/files")
	cmd.Stdout = os.Stdout
	if cmderr = cmd.Run(); cmderr != nil {
		fmt.Println(cmderr)
	}

	time.Sleep(10 * time.Second)
	//cmd = exec.Command("sshpass")
	/*
		cmd.Stdout = os.Stdout
		if cmderr = cmd.Run(); cmderr != nil {
			fmt.Println(cmderr)
		}
	*/
	/*
		cmd = exec.Command("sshpass -p swinfra12 scp -r -o StrictHostKeyChecking=no /files/vs-service.yaml root@172.10.1.23:/root/files/")
		cmd.Stdout = os.Stdout
		if cmderr = cmd.Run(); cmderr != nil {
			fmt.Println(cmderr)
		}
	*/
	/*
		if _, err2 := sshclient.Exec("kubectl apply -f /webhook.yaml"); err2 != nil {
		}
		if _, err2 := sshclient.Exec("kubectl apply -f /operator.yaml"); err2 != nil {
		}
		if _, err2 := sshclient.Exec("kubectl apply -f /root/asu/smf/workload/dp-service-smf.yaml"); err2 != nil {
		}
		if _, err2 := sshclient.Exec("kubectl apply -f /root/asu/smf/workload/dp-service-smf.yaml"); err2 != nil {
		}
		if _, err2 := sshclient.Exec("kubectl apply -f /root/asu/smf/workload/dp-service-smf.yaml"); err2 != nil {
		}
		if _, err2 := sshclient.Exec("kubectl apply -f /root/asu/smf/workload/dp-service-smf.yaml"); err2 != nil {
		}
	*/
	//if _, err2 := sshclient.Exec("sudo mkdir /var/hpvolumes"); err2 != nil {
	//}

	if _, err2 := sshclient.Exec("kubectl apply -f /root/files/multus-daemonset-thick.yml"); err2 != nil {
	}

	if _, err2 := sshclient.Exec("kubectl apply -f /root/files/cert-manager.yaml"); err2 != nil {
	}

	if _, err2 := sshclient.Exec("kubectl apply -f /root/files/sc_namespace.yaml"); err2 != nil {
	}

	time.Sleep(120 * time.Second)

	if _, err2 := sshclient.Exec(" kubectl apply -f /root/files/webhook.yaml -n hostpath-provisioner"); err2 != nil {
		fmt.Println("webhook")
		fmt.Println(err2)
	}

	time.Sleep(15 * time.Second)

	if _, err2 := sshclient.Exec("kubectl apply -f /root/files/operator.yaml -n hostpath-provisioner"); err2 != nil {
	}

	time.Sleep(20 * time.Second)

	if _, err2 := sshclient.Exec("kubectl apply -f /root/files/cr_hostpath.yaml -n hostpath-provisioner"); err2 != nil {
		fmt.Println(err2)
	}

	if _, err2 := sshclient.Exec("kubectl apply -f /root/files/sc-wffc-legacy-csi.yaml"); err2 != nil {
	}

	if _, err2 := sshclient.Exec("kubectl apply -f /root/files/istio_namespace.yaml"); err2 != nil {
	}

	if _, err2 := sshclient.Exec("/root/files/istioctl install --set profile=default -f /root/files/install_default.yaml -y"); err2 != nil {
		fmt.Println("istioctl")
		fmt.Println(err2)
	}

	if _, err2 := sshclient.Exec("kubectl apply -f /root/files/gateway.yaml"); err2 != nil {
	}

	if _, err2 := sshclient.Exec("kubectl apply -f /root/files/vs-service.yaml"); err2 != nil {
	}

	//if _, err2 := sshclient.Exec("kubectl apply -f /files/"); err2 != nil {
	//}

	/*
		if _, err2 := sshclient.Exec("/files/istioctl install --set profile=default -f /files/install_default.yaml"); err2 != nil {
		}

		if _, err2 := sshclient.Exec("sudo mkdir /root/hpvolume"); err2 != nil {
		}
	*/
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InfraReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fivegv1alpha1.Infra{}).
		Complete(r)
}
