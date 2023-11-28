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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DataVolumeSourceHTTP can be either an http or https endpoint, with an optional basic auth user name and password, and an optional configmap containing additional CAs
type DataVolumeSourceHTTP struct {
	// URL is the URL of the http(s) endpoint
	URL string `json:"url"`
	// SecretRef A Secret reference, the secret should contain accessKeyId (user name) base64 encoded, and secretKey (password) also base64 encoded
	// +optional
	SecretRef string `json:"secretRef,omitempty"`
	// CertConfigMap is a configmap reference, containing a Certificate Authority(CA) public key, and a base64 encoded pem certificate
	// +optional
	CertConfigMap string `json:"certConfigMap,omitempty"`
	// ExtraHeaders is a list of strings containing extra headers to include with HTTP transfer requests
	// +optional
	ExtraHeaders []string `json:"extraHeaders,omitempty"`
	// SecretExtraHeaders is a list of Secret references, each containing an extra HTTP header that may include sensitive information
	// +optional
	SecretExtraHeaders []string `json:"secretExtraHeaders,omitempty"`
}

// DataVolumeBlankImage provides the parameters to create a new raw blank image for the PVC
type DataVolumeBlankImage struct{}

// DataVolumeSource represents the source for our Data Volume, this can be HTTP, Imageio, S3, GCS, Registry or an existing PVC
type DataVolumeSource struct {
	HTTP *DataVolumeSourceHTTP `json:"http,omitempty"`
	/*	NFS      *DataVolumeSourceNFS      `json:"nfs,omitempty"`
		S3       *DataVolumeSourceS3       `json:"s3,omitempty"`
		GCS      *DataVolumeSourceGCS      `json:"gcs,omitempty"`
		Registry *DataVolumeSourceRegistry `json:"registry,omitempty"`
		PVC      *DataVolumeSourcePVC      `json:"pvc,omitempty"`
		Upload   *DataVolumeSourceUpload   `json:"upload,omitempty"`
	*/
	Blank *DataVolumeBlankImage `json:"blank,omitempty"`
}

// DatasetSpec defines the desired state of Dataset
type DatasetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Dataset. Edit dataset_types.go to remove/update
	Foo    string            `json:"foo,omitempty"`
	Source *DataVolumeSource `json:"source,omitempty"`
	//SourceRef is an indirect reference to the source of data for the requested DataVolume
	//PVC is the PVC specification
	PVC *corev1.PersistentVolumeClaimSpec `json:"pvc,omitempty"`
}

// DatasetStatus defines the observed state of Dataset
type DatasetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase    string `json:"phase,omitempty"`
	Progress string `json:"progress,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="The phase the data volume is in"
// +kubebuilder:printcolumn:name="Progress",type="string",JSONPath=".status.progress",description="Transfer progress in percentage if known, N/A otherwise"
// +kubebuilder:printcolumn:name="Restarts",type="integer",JSONPath=".status.restartCount",description="The number of times the transfer has been restarted."
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// Dataset is the Schema for the datasets API
type Dataset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatasetSpec   `json:"spec,omitempty"`
	Status DatasetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DatasetList contains a list of Dataset
type DatasetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dataset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dataset{}, &DatasetList{})
}
