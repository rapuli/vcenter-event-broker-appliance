package fake

import (
	"context"
	"log"
	"os"

	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/color"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/events"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/metrics"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/processor"
	"github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/provider"

	// "github.com/vmware-samples/vcenter-event-broker-appliance/vmware-event-router/internal/stream"
	"github.com/vmware/govmomi/vim25/types"
)

const source = "https://fake.vcenter01.testing.io/sdk"

// verify that VCenter implements the streamer interface
var _ provider.Provider = (*VCenter)(nil)

// VCenter implements the streamer interface
type VCenter struct {
	eventCh <-chan []types.BaseEvent // channel which simulates events
	*log.Logger
}

// NewFakeVCenter returns a fake vcenter event stream provider streaming events
// received from the specified generator channel
func NewFakeVCenter(generator <-chan []types.BaseEvent) *VCenter {
	return &VCenter{
		eventCh: generator,
		Logger:  log.New(os.Stdout, color.Magenta("[Fake vCenter] "), log.LstdFlags),
	}
}

// PushMetrics is a no-op
func (f *VCenter) PushMetrics(context.Context, metrics.Receiver) {}

// Stream streams events generated by the Generator specified in the VCenter
// server
func (f *VCenter) Stream(ctx context.Context, p processor.Processor) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case baseEvent := <-f.eventCh:
			for idx := range baseEvent {
				// process slice in reverse order to maintain Event.Key ordering
				event := baseEvent[len(baseEvent)-1-idx]

				ce, err := events.NewCloudEvent(event, source)
				if err != nil {
					log.Printf("skipping event %v because it could not be converted to CloudEvent format: %v", event, err)
					continue
				}

				err = p.Process(*ce)
				if err != nil {
					f.Logger.Printf("could not process event %v: %v", ce, err)
				}
			}
		}
	}
}

// Shutdown is a no-op
func (f *VCenter) Shutdown(context.Context) error {
	return nil
}
