package constants

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type timeouts struct {
	Zero                 *metav1.Duration
	QuickTaskRun         *metav1.Duration
	SmallBlankDVCreation *metav1.Duration
	DefaultTaskRun       *metav1.Duration
}

var Timeouts = timeouts{
	Zero:                 &metav1.Duration{0 * time.Second},
	SmallBlankDVCreation: &metav1.Duration{15 * time.Minute},
	QuickTaskRun:         &metav1.Duration{5 * time.Minute},
	DefaultTaskRun:       &metav1.Duration{10 * time.Minute},
}
