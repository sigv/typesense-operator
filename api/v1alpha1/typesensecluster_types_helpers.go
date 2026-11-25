package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *TypesenseClusterSpec) GetResources() corev1.ResourceRequirements {
	if s.Resources != nil {
		return *s.Resources
	}

	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("1000m"),
			corev1.ResourceMemory: resource.MustParse("512Mi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("256Mi"),
		},
	}
}

func (s *TypesenseClusterSpec) GetAdditionalServerConfiguration() []corev1.EnvFromSource {
	if s.AdditionalServerConfiguration != nil {
		return []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: *s.AdditionalServerConfiguration,
				},
			},
		}
	}

	return []corev1.EnvFromSource{}
}

func (s *TypesenseClusterSpec) GetCorsDomains() string {
	if s.CorsDomains == nil {
		return ""
	}
	return *s.CorsDomains
}

func (s *TypesenseClusterSpec) GetAffinity(labels map[string]string) *corev1.Affinity {
	if s.Affinity == nil {
		return nil
	}

	affinity := s.Affinity.DeepCopy()

	if affinity.PodAffinity != nil {
		terms := affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution
		for i := range terms {
			setPodAffinityTermLabels(&terms[i], labels)
		}

		weightedTerms := affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution
		for i := range weightedTerms {
			setPodAffinityTermLabels(&weightedTerms[i].PodAffinityTerm, labels)
		}
	}

	if affinity.PodAntiAffinity != nil {
		terms := affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution
		for i := range terms {
			setPodAffinityTermLabels(&terms[i], labels)
		}

		weightedTerms := affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution
		for i := range weightedTerms {
			setPodAffinityTermLabels(&weightedTerms[i].PodAffinityTerm, labels)
		}
	}

	return affinity
}

func (s *TypesenseClusterSpec) GetTopologySpreadConstraints(labels map[string]string) []corev1.TopologySpreadConstraint {
	tscs := make([]corev1.TopologySpreadConstraint, len(s.TopologySpreadConstraints))

	for i, tsc := range s.TopologySpreadConstraints {
		if tsc.LabelSelector == nil {
			tsc.LabelSelector = &metav1.LabelSelector{
				MatchLabels: cloneLabels(labels),
			}
		}
		tscs[i] = tsc
	}
	return tscs
}

func setPodAffinityTermLabels(term *corev1.PodAffinityTerm, labels map[string]string) {
	if term.LabelSelector == nil {
		term.LabelSelector = &metav1.LabelSelector{
			MatchLabels: cloneLabels(labels),
		}
	}
}

func cloneLabels(labels map[string]string) map[string]string {
	result := make(map[string]string, len(labels))
	for key, value := range labels {
		result[key] = value
	}
	return result
}
