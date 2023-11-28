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

package controller

import (
	"context"
	"fmt"

	datasetopsv1 "dataset-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DatasetReconciler reconciles a Dataset object
type DatasetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// /////////////////////////////////
func createImporterPod(ctx context.Context, namespacedName types.NamespacedName, dataset *datasetopsv1.Dataset) error {
	clientset, err := newKubeClient()
	if err != nil {
		return err
	}
	httpURL := dataset.Spec.Source.HTTP.URL
	pvcName := namespacedName.Name
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "importer-" + namespacedName.Name,
			Namespace: namespacedName.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "importer",
					Image: "nginx",
					Command: []string{
						"/bin/sh",
						"-c",
						"curl " + httpURL + " -Lo /mnt/a.out",
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "my-volume",
							MountPath: "/mnt",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "my-volume",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
						},
					},
				},
			},
		},
	}
	createdPod, err := clientset.CoreV1().Pods(namespacedName.Namespace).Create(ctx, pod, v1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Pod %s created\n", createdPod.Name)
	return nil
}

// ///////////////////////////
func createPVC(ctx context.Context, namespacedName types.NamespacedName, pvcSpec *corev1.PersistentVolumeClaimSpec) error {
	clientset, err := newKubeClient()
	if err != nil {
		return err
	}

	// Create a PersistentVolumeClaim object
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
		},
		Spec: *pvcSpec,
	}
	createdPVC, err := clientset.CoreV1().PersistentVolumeClaims(namespacedName.Namespace).Create(ctx, pvc, v1.CreateOptions{})
	fmt.Printf("DEBUG, createdPVC=%s", createdPVC)
	return err
}

// / FIXME .. should be new once..
// should get from /var/run/secret serviceAccount token & ca.crt
func newKubeClient() (*kubernetes.Clientset, error) {
	kubeC := "/root/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeC)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

//+kubebuilder:rbac:groups=dataset-ops.my.domain,resources=datasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dataset-ops.my.domain,resources=datasets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dataset-ops.my.domain,resources=datasets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Dataset object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *DatasetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	dataset := &datasetopsv1.Dataset{}
	if err := r.Get(ctx, req.NamespacedName, dataset); err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if dataset.Status.Phase == "Completed" {
		return ctrl.Result{}, nil
	}
	dataset.Status.Phase = "Processing"
	if err := r.Status().Update(ctx, dataset); err != nil {
		return ctrl.Result{}, err
	}
	// TODO(user): your logic here
	err := createPVC(ctx, req.NamespacedName, dataset.Spec.PVC)
	if err != nil {
		panic(err.Error())
	}
	err = createImporterPod(ctx, req.NamespacedName, dataset)
	if err != nil {
		panic(err.Error())
	}
	// FIXME, it's just an example, should wait Pod Completed
	dataset.Status.Phase = "Completed"
	if err := r.Status().Update(ctx, dataset); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatasetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datasetopsv1.Dataset{}).
		Complete(r)
}
