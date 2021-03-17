package handlers

import (
	ethereumv1alpha1 "github.com/kotalco/kotal/apis/ethereum/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func nodeFromSpec(spec *ethereumv1alpha1.NodeSpec, name string) *ethereumv1alpha1.Node {
	return &ethereumv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: *spec,
	}
}
