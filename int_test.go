package client_sdk_go

import (
	"context"
	"github.com/awakari/client-sdk-go/api"
	"github.com/awakari/client-sdk-go/model"
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
	"github.com/awakari/client-sdk-go/model/usage"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"os"
	"testing"
	"time"
)

const userId = "test-user-0"

var serverPublicKeyPath = os.Getenv("SERVER_PUBLIC_KEY_PATH")
var clientCertPath = os.Getenv("CLIENT_CERT_PATH")
var clientPrivateKeyPath = os.Getenv("CLIENT_PRIVATE_KEY_PATH")
var apiUri = os.Getenv("API_URI")

func TestPublicApiUsage(t *testing.T) {

	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()

	// load TLS certificates 1st
	caCrt, err := os.ReadFile(serverPublicKeyPath)
	if err != nil {
		panic(err)
	}
	clientCrt, err := os.ReadFile(clientCertPath)
	if err != nil {
		panic(err)
	}
	clientKey, err := os.ReadFile(clientPrivateKeyPath)
	if err != nil {
		panic(err)
	}

	// init the client
	var client api.Client
	client, err = api.
		NewClientBuilder().
		ServerPublicKey(caCrt).
		ClientKeyPair(clientCrt, clientKey).
		ApiUri(apiUri).
		Build()
	require.Nil(t, err)
	defer client.Close()

	// Get the initial Subscriptions API Usage
	var usageSubsStart usage.Usage
	usageSubsStart, err = client.ReadUsage(ctx, userId, usage.SubjectSubscriptions)
	assert.Nil(t, err)

	// Create a Subscription
	subData := subscription.Data{
		Metadata: subscription.Metadata{
			Description: "test subscription 0",
			Enabled:     true,
		},
		Condition: condition.NewBuilder().
			MatchAttrKey("tags").
			MatchAttrValuePattern("neutrino").
			MatchAttrValuePartial().
			BuildKiwiTreeCondition(),
	}
	var subId string
	subId, err = client.CreateSubscription(ctx, userId, subData)
	assert.Nil(t, err)

	// Check the Subscriptions API Usage change
	var usageSubs usage.Usage
	usageSubs, err = client.ReadUsage(ctx, userId, usage.SubjectSubscriptions)
	assert.Nil(t, err)
	assert.Equal(t, usageSubsStart.Count+1, usageSubs.Count)
	assert.Equal(t, usageSubsStart.CountTotal+1, usageSubs.CountTotal)

	// Open a Read Stream
	var rs model.ReadStream[*pb.CloudEvent]
	rs, err = client.ReadMessages(ctx, userId, subId)
	assert.Nil(t, err)
	defer rs.Close()

	// Get the initial Publish Messages API Usage
	var usagePubMsgsStart usage.Usage
	usagePubMsgsStart, err = client.ReadUsage(ctx, userId, usage.SubjectPublishMessages)
	assert.Nil(t, err)

	// Write a Message
	var ws model.WriteStream[*pb.CloudEvent]
	ws, err = client.WriteMessages(ctx, userId)
	assert.Nil(t, err)
	defer ws.Close()
	msgSend := &pb.CloudEvent{
		Id:          uuid.NewString(),
		Source:      "http://arxiv.org/abs/2305.06364",
		SpecVersion: "1.0",
		Type:        "com.github.awakari.producer-rss",
		Attributes: map[string]*pb.CloudEventAttributeValue{
			"summary": {
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: "<p>We propose that the dark matter of our universe could be sterile neutrinos which reside within the twin sector of a mirror twin Higgs model. In our scenario, these particles are produced through a version of the Dodelson-Widrow mechanism that takes place entirely within the twin sector, yielding a dark matter candidate that is consistent with X-ray and gamma-ray line constraints. Furthermore, this scenario can naturally avoid the cosmological problems that are typically encountered in mirror twin Higgs models. In particular, if the sterile neutrinos in the Standard Model sector decay out of equilibrium, they can heat the Standard Model bath and reduce the contributions of the twin particles to $N_\\mathrm{eff}$. Such decays also reduce the effective temperature of the dark matter, thereby relaxing constraints from large-scale structure. The sterile neutrinos included in this model are compatible with the seesaw mechanism for generating Standard Model neutrino masses. </p> ",
				},
			},
			"tags": {
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: "neutrino dark matter cosmology higgs standard model dodelson-widrow",
				},
			},
			"title": {
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: "Twin Sterile Neutrino Dark Matter. (arXiv:2305.06364v1 [hep-ph])",
				},
			},
		},
		Data: &pb.CloudEvent_TextData{
			TextData: "",
		},
	}
	var writtenCount uint32
	writtenCount, err = ws.WriteBatch([]*pb.CloudEvent{msgSend})
	assert.Equal(t, uint32(1), writtenCount)
	assert.Nil(t, err)

	// Check the Publish Messages API Usage change
	var usagePubMsgs usage.Usage
	usagePubMsgs, err = client.ReadUsage(ctx, userId, usage.SubjectPublishMessages)
	assert.Nil(t, err)
	assert.Equal(t, usagePubMsgsStart.Count+1, usagePubMsgs.Count)
	assert.Equal(t, usagePubMsgsStart.CountTotal+1, usagePubMsgs.CountTotal)

	// Read the Message by the Subscription
	var msgRead *pb.CloudEvent
	msgRead, err = rs.Read()
	assert.Nil(t, err)
	if err == nil {
		assert.Equal(t, msgSend.Id, msgRead.Id)
	}

	// Delete the Subscription to clean up
	err = client.DeleteSubscription(ctx, userId, subId)
	assert.Nil(t, err)

	// Check the Subscriptions API Usage change
	usageSubs, err = client.ReadUsage(ctx, userId, usage.SubjectSubscriptions)
	assert.Nil(t, err)
	assert.Equal(t, usageSubsStart.Count, usageSubs.Count)
	assert.Equal(t, usageSubsStart.CountTotal, usageSubs.CountTotal)
}

func TestInternalWriter(t *testing.T) {

	if os.Getenv("CI") != "" {
		t.Skip("Skipping test in CI environment")
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-awakari-group-id", "test-group-0")

	var client api.Client
	var err error
	client, err = api.
		NewClientBuilder().
		WriteUri(apiUri).
		Build()
	require.Nil(t, err)
	defer client.Close()

	// Write a Message
	var ws model.WriteStream[*pb.CloudEvent]
	ws, err = client.WriteMessages(ctx, userId)
	assert.Nil(t, err)
	defer ws.Close()
	msgSend := &pb.CloudEvent{
		Id:          uuid.NewString(),
		Source:      "http://arxiv.org/abs/2305.06364",
		SpecVersion: "1.0",
		Type:        "com.github.awakari.producer-rss",
		Attributes: map[string]*pb.CloudEventAttributeValue{
			"summary": {
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: "<p>We propose that the dark matter of our universe could be sterile neutrinos which reside within the twin sector of a mirror twin Higgs model. In our scenario, these particles are produced through a version of the Dodelson-Widrow mechanism that takes place entirely within the twin sector, yielding a dark matter candidate that is consistent with X-ray and gamma-ray line constraints. Furthermore, this scenario can naturally avoid the cosmological problems that are typically encountered in mirror twin Higgs models. In particular, if the sterile neutrinos in the Standard Model sector decay out of equilibrium, they can heat the Standard Model bath and reduce the contributions of the twin particles to $N_\\mathrm{eff}$. Such decays also reduce the effective temperature of the dark matter, thereby relaxing constraints from large-scale structure. The sterile neutrinos included in this model are compatible with the seesaw mechanism for generating Standard Model neutrino masses. </p> ",
				},
			},
			"tags": {
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: "neutrino dark matter cosmology higgs standard model dodelson-widrow",
				},
			},
			"title": {
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: "Twin Sterile Neutrino Dark Matter. (arXiv:2305.06364v1 [hep-ph])",
				},
			},
		},
		Data: &pb.CloudEvent_TextData{
			TextData: "",
		},
	}
	var writtenCount uint32
	writtenCount, err = ws.WriteBatch([]*pb.CloudEvent{msgSend})
	assert.Equal(t, uint32(1), writtenCount)
	assert.Nil(t, err)
}
