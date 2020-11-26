/*


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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VSISpec defines the desired state of VSI
type VSISpec struct {
	APIKey string `json:"apikey"`
	SSHKey string `json:"sshkey"`
	VPC    string `json:"vpc"`
	Subnet string `json:"subnet"`
}

// VSIStatus defines the observed state of VSI
type VSIStatus struct {
	IPAddress string `json:"ipaddress"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// VSI is the Schema for the vsis API
type VSI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSISpec   `json:"spec,omitempty"`
	Status VSIStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VSIList contains a list of VSI
type VSIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VSI{}, &VSIList{})
}
