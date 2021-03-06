// +build integration,aws

package integration_test

import (
	"context"
	"os"
	"testing"

	config "github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/config/v1alpha1"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/metrics"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/processor"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
)

const (
	fakeVCenterName = "https://fakevc-01:443/sdk"
)

// implement metrics interface
type fakeReceiver struct {
}

func (f *fakeReceiver) Receive(_ *metrics.EventStats) {
}

var (
	ctx          context.Context
	awsProcessor processor.Processor
	receiver     *fakeReceiver
)

func TestAWS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS EventBridge Suite")
}

var _ = BeforeSuite(func() {
	ctx = context.Background()

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecret := os.Getenv("AWS_SECRET_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	awsBus := os.Getenv("AWS_EVENT_BUS")
	awsRule := os.Getenv("AWS_RULE_ARN")

	Expect(awsAccessKey).ToNot(BeEmpty(), "env var AWS_ACCESS_KEY to authenticate against AWS EventBridge must be set")
	Expect(awsSecret).ToNot(BeEmpty(), "env var AWS_SECRET_KEY to authenticate against AWS EventBridge must be set")
	Expect(awsRegion).ToNot(BeEmpty(), "env var AWS_REGION for AWS EventBridge must be set")
	Expect(awsBus).ToNot(BeEmpty(), "env var AWS_EVENT_BUS for AWS EventBridge must be set")
	Expect(awsRule).ToNot(BeEmpty(), "env var AWS_RULE_ARN for AWS EventBridge must be set")

	cfg := &config.ProcessorConfigEventBridge{
		EventBus: awsBus,
		Region:   awsRegion,
		RuleARN:  awsRule,
		Auth: &config.AuthMethod{
			Type: config.AWSAccessKeyAuth,
			AWSAccessKeyAuth: &config.AWSAccessKeyAuthMethod{
				AccessKey: awsAccessKey,
				SecretKey: awsSecret,
			},
		},
	}

	receiver = &fakeReceiver{}
	p, err := processor.NewEventBridgeProcessor(ctx, cfg, receiver, processor.WithAWSVerbose(true))
	Expect(err).NotTo(HaveOccurred())
	awsProcessor = p
})

var _ = AfterSuite(func() {})
