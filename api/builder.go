package api

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/api/grpc/permits"
	"github.com/awakari/client-sdk-go/api/grpc/reader"
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

	// ReaderUri sets the Awakari messages reading API URI. Overrides any value set by ApiUri.
	// Useful when the specific message reading API is needed by the client.
	ReaderUri(readerUri string) ClientBuilder

	// SubscriptionsUri sets the Awakari subscriptions API URI. Overrides any value set by ApiUri.
	// Useful when the specific subscriptions management API is needed by the client.
	SubscriptionsUri(subsUri string) ClientBuilder

	// WriterUri sets the Awakari messages publishing API URI. Overrides any value set by ApiUri.
	// Useful when the specific message publishing API is needed by the client.
	WriterUri(writerUri string) ClientBuilder

	// Build instantiates the Client instance and returns it.
	Build() (c Client, err error)
}

type builder struct {
	caCrt     []byte
	clientCrt []byte
	clientKey []byte
	apiUri    string
	readerUri string
	subsUri   string
	writerUri string
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

func (b *builder) ReaderUri(readerUri string) ClientBuilder {
	b.readerUri = readerUri
	return b
}

func (b *builder) SubscriptionsUri(subsUri string) ClientBuilder {
	b.subsUri = subsUri
	return b
}

func (b *builder) WriterUri(writerUri string) ClientBuilder {
	b.writerUri = writerUri
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
	var connReader *grpc.ClientConn
	if b.readerUri != "" {
		connReader, err = grpc.Dial(b.readerUri, optsDial...)
	} else if b.apiUri != "" {
		connReader, err = grpc.Dial(b.apiUri, optsDial...)
	}
	var svcReader reader.Service
	if connReader != nil {
		clientReader := reader.NewServiceClient(connReader)
		svcReader = reader.NewService(clientReader)
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
	if b.writerUri != "" {
		connWriter, err = grpc.Dial(b.writerUri, optsDial...)
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
		connReader:  connReader,
		connPermits: connPermits,
		connSubs:    connSubs,
		connWriter:  connWriter,
		svcLimits:   svcLimits,
		svcReader:   svcReader,
		svcPermits:  svcPermits,
		svcSubs:     svcSubs,
		svcWriter:   svcWriter,
	}
	return
}
