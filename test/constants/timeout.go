package constants

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type timeouts struct {
	Zero                      *metav1.Duration
	QuickTaskRun              *metav1.Duration
	SmallDVCreation           *metav1.Duration
	DefaultTaskRun            *metav1.Duration
	TaskRunExtraWaitDelay     *metav1.Duration
	PipelineRunExtraWaitDelay *metav1.Duration
	WaitBeforeExecutingVM     *metav1.Duration
	WaitForVMStart            *metav1.Duration
}

var Timeouts = timeouts{
	Zero:                      &metav1.Duration{0 * time.Second},
	TaskRunExtraWaitDelay:     &metav1.Duration{5 * time.Minute},
	SmallDVCreation:           &metav1.Duration{15 * time.Minute},
	QuickTaskRun:              &metav1.Duration{5 * time.Minute},
	DefaultTaskRun:            &metav1.Duration{10 * time.Minute},
	WaitBeforeExecutingVM:     &metav1.Duration{30 * time.Second},
	WaitForVMStart:            &metav1.Duration{5 * time.Minute},
	PipelineRunExtraWaitDelay: &metav1.Duration{30 * time.Minute},
}
