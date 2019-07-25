package knative

import (
	"context"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/knative/eventing-sources/pkg/kncloudevents"
	"github.com/pkg/errors"
)

// EmitKnativeEvent emits (sends) the event you specify via the broker you specify.
// eventType: CloudEvents type (see: https://github.com/cloudevents/spec/blob/v0.2/spec.md#type), e.g.: com.example.object.delete
// eventVersion: Version number of the event. Will be appended to the event's type separated by a forward slash (/).
// sourceID: will be appended to "urn:event:from" to form the fully qualified CloudEvents Source ID (see: https://github.com/cloudevents/spec/blob/v0.2/spec.md#source)
func EmitKnativeEvent(brokerURL, eventType, eventVersion, sourceID string, eventData interface{}) error {
	// Based on https://raw.githubusercontent.com/knative/eventing-sources/master/cmd/heartbeats/main.go

	c, err := kncloudevents.NewDefaultClient(brokerURL)
	if err != nil {
		return errors.WithStack(err)
	}

	source := types.ParseURLRef(
		fmt.Sprintf("urn:event:from:%s", sourceID))

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			Type:   eventType + "/" + eventVersion,
			Source: *source,
		}.AsV02(),
		Data: eventData,
	}

	if _, err := c.Send(context.Background(), event); err != nil {
		return errors.Wrap(err, "failed to send cloudevent")
	}

	return nil
}
