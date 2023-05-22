package api

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/api/grpc/messages"
	"github.com/awakari/client-sdk-go/api/grpc/permits"
	"github.com/awakari/client-sdk-go/api/grpc/subscriptions"
	"github.com/awakari/client-sdk-go/api/grpc/writer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientBuilder interface {

	// ServerPublicKey sets the CA certificate to authenticate the Awakari service.
	// Should be used together with ClientKeyPair and ApiUri.
	ServerPublicKey(caCrt []byte) ClientBuilder

	// ClientKeyPair sets the client certificate key pair to allow Awakari service to authenticate the client.
	// Should be used together with ServerPublicKey and ApiUri.
	ClientKeyPair(clientCrt, clientKey []byte) ClientBuilder

	// ApiUri sets the Awakari public API URI. Should be used together with ServerPublicKey and ClientKeyPair.
	// Useful when a client needs every available public API method.
	// Enables additionally the API methods to read the usage limits and permits.
	ApiUri(apiUri string) ClientBuilder

	// ReadUri sets the Awakari messages reading API URI. Overrides any value set by ApiUri.
	// Useful when the specific message reading API is needed by the client.
	ReadUri(readUri string) ClientBuilder

	// SubscriptionsUri sets the Awakari subscriptions API URI. Overrides any value set by ApiUri.
	// Useful when the specific subscriptions management API is needed by the client.
	SubscriptionsUri(subsUri string) ClientBuilder

	// WriteUri sets the Awakari messages publishing API URI. Overrides any value set by ApiUri.
	// Useful when the specific message publishing API is needed by the client.
	WriteUri(writeUri string) ClientBuilder

	// Build instantiates the Client instance and returns it.
	Build() (c Client, err error)
}

type builder struct {
	caCrt     []byte
	clientCrt []byte
	clientKey []byte
	apiUri    string
	readUri   string
	subsUri   string
	writeUri  string
}

func NewClientBuilder() ClientBuilder {
	return &builder{}
}

func (b *builder) ServerPublicKey(caCrt []byte) ClientBuilder {
	b.caCrt = caCrt
	return b
}

func (b *builder) ClientKeyPair(clientCrt, clientKey []byte) ClientBuilder {
	b.clientCrt = clientCrt
	b.clientKey = clientKey
	return b
}

func (b *builder) ApiUri(apiUri string) ClientBuilder {
	b.apiUri = apiUri
	return b
}

func (b *builder) ReadUri(readUri string) ClientBuilder {
	b.readUri = readUri
	return b
}

func (b *builder) SubscriptionsUri(subsUri string) ClientBuilder {
	b.subsUri = subsUri
	return b
}

func (b *builder) WriteUri(writeUri string) ClientBuilder {
	b.writeUri = writeUri
	return b
}

func (b *builder) Build() (c Client, err error) {
	//
	tlsConfig := &tls.Config{}
	if b.caCrt != nil {
		caCrtPool := x509.NewCertPool()
		caCrtPool.AppendCertsFromPEM(b.caCrt)
		tlsConfig.RootCAs = caCrtPool
	}
	if b.clientCrt != nil && b.clientKey != nil {
		var clientKeyPair tls.Certificate
		clientKeyPair, err = tls.X509KeyPair(b.clientCrt, b.clientKey)
		if err == nil {
			tlsConfig.Certificates = []tls.Certificate{
				clientKeyPair,
			}
		}
	}
	var creds credentials.TransportCredentials
	if tlsConfig.RootCAs != nil || tlsConfig.Certificates != nil {
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}
	optsDial := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}
	//
	var connLimits *grpc.ClientConn
	var svcLimits limits.Service
	if b.apiUri != "" {
		connLimits, err = grpc.Dial(b.apiUri, optsDial...)
		clientLimits := limits.NewServiceClient(connLimits)
		svcLimits = limits.NewService(clientLimits)
	}
	//
	var connMsgs *grpc.ClientConn
	if b.readUri != "" {
		connMsgs, err = grpc.Dial(b.readUri, optsDial...)
	} else if b.apiUri != "" {
		connMsgs, err = grpc.Dial(b.apiUri, optsDial...)
	}
	var svcMsgs messages.Service
	if connMsgs != nil {
		clientMsgs := messages.NewServiceClient(connMsgs)
		svcMsgs = messages.NewService(clientMsgs)
	}
	//
	var connPermits *grpc.ClientConn
	var svcPermits permits.Service
	if b.apiUri != "" {
		connPermits, err = grpc.Dial(b.apiUri, optsDial...)
		clientPermits := permits.NewServiceClient(connPermits)
		svcPermits = permits.NewService(clientPermits)
	}
	//
	var connSubs *grpc.ClientConn
	if b.subsUri != "" {
		connSubs, err = grpc.Dial(b.subsUri, optsDial...)
	} else if b.apiUri != "" {
		connSubs, err = grpc.Dial(b.apiUri, optsDial...)
	}
	var svcSubs subscriptions.Service
	if connSubs != nil {
		clientSubs := subscriptions.NewServiceClient(connSubs)
		svcSubs = subscriptions.NewService(clientSubs)
	}
	//
	var connWriter *grpc.ClientConn
	if b.writeUri != "" {
		connWriter, err = grpc.Dial(b.writeUri, optsDial...)
	} else if b.apiUri != "" {
		connWriter, err = grpc.Dial(b.apiUri, optsDial...)
	}
	var svcWriter writer.Service
	if connWriter != nil {
		clientWriter := writer.NewServiceClient(connWriter)
		svcWriter = writer.NewService(clientWriter)
	}
	//
	c = client{
		connLimits:  connLimits,
		connMsgs:    connMsgs,
		connPermits: connPermits,
		connSubs:    connSubs,
		connWriter:  connWriter,
		svcLimits:   svcLimits,
		svcMsgs:     svcMsgs,
		svcPermits:  svcPermits,
		svcSubs:     svcSubs,
		svcWriter:   svcWriter,
	}
	return
}
