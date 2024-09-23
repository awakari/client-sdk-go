package api

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/awakari/client-sdk-go/api/grpc/limits"
	"github.com/awakari/client-sdk-go/api/grpc/permits"
	"github.com/awakari/client-sdk-go/api/grpc/reader"
	"github.com/awakari/client-sdk-go/api/grpc/resolver"
	"github.com/awakari/client-sdk-go/api/grpc/subscriptions"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ClientBuilder interface {

	// CertAuthority sets the CA to authenticate the Awakari service.
	// Should be used together with ClientKeyPair and ApiUri.
	CertAuthority(caCrt []byte) ClientBuilder

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

	// SubscriptionsUri sets the Awakari subscriptions-proxy API URI. Overrides any value set by ApiUri.
	// Useful when the specific subscriptions management API is needed by the client.
	SubscriptionsUri(subsUri string) ClientBuilder

	// WriterUri sets the Awakari messages publishing API URI. Overrides any value set by ApiUri.
	// Useful when the specific message publishing API is needed by the client.
	WriterUri(writerUri string) ClientBuilder

	Connections(countMax int, idleTimeout time.Duration, maxLifeDuration ...time.Duration) ClientBuilder

	// Build instantiates the Client instance and returns it.
	Build() (c Client, err error)
}

type builder struct {
	caCrt               []byte
	clientCrt           []byte
	clientKey           []byte
	apiUri              string
	readerUri           string
	subsUri             string
	writerUri           string
	connMax             int
	connIdleTimeout     time.Duration
	connMaxLifeDuration time.Duration
}

func NewClientBuilder() ClientBuilder {
	return &builder{
		connMax: 1,
	}
}

func (b *builder) CertAuthority(caCrt []byte) ClientBuilder {
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

func (b *builder) Connections(countMax int, idleTimeout time.Duration, maxLifeDuration ...time.Duration) ClientBuilder {
	b.connMax = countMax
	b.connIdleTimeout = idleTimeout
	if len(maxLifeDuration) > 0 {
		b.connMaxLifeDuration = maxLifeDuration[0]
	}
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
	var connPoolApi *grpcpool.Pool
	if b.apiUri != "" {
		connPoolApi, err = grpcpool.New(func() (*grpc.ClientConn, error) {
			return grpc.NewClient(b.apiUri, optsDial...)
		}, 1, b.connMax, b.connIdleTimeout, b.connMaxLifeDuration)
	}
	var connPoolReader *grpcpool.Pool
	if b.readerUri != "" {
		connPoolReader, err = grpcpool.New(func() (*grpc.ClientConn, error) {
			return grpc.NewClient(b.readerUri, optsDial...)
		}, 1, b.connMax, b.connIdleTimeout, b.connMaxLifeDuration)
	}
	var connPoolSubs *grpcpool.Pool
	if b.subsUri != "" {
		connPoolSubs, err = grpcpool.New(func() (*grpc.ClientConn, error) {
			return grpc.NewClient(b.subsUri, optsDial...)
		}, 1, b.connMax, b.connIdleTimeout, b.connMaxLifeDuration)
	}
	var connPoolWriter *grpcpool.Pool
	if b.writerUri != "" {
		connPoolWriter, err = grpcpool.New(func() (*grpc.ClientConn, error) {
			return grpc.NewClient(b.writerUri, optsDial...)
		}, 1, b.connMax, b.connIdleTimeout, b.connMaxLifeDuration)
	}
	//
	var svcLimits limits.Service
	if b.apiUri != "" {
		clientLimits := limits.NewClientConnPool(connPoolApi)
		svcLimits = limits.NewService(clientLimits)
	}
	//
	var clientReader reader.ServiceClient
	if b.readerUri != "" {
		clientReader = reader.NewClientConnPool(connPoolReader)
	} else if b.apiUri != "" {
		clientReader = reader.NewClientConnPool(connPoolApi)
	}
	var svcReader reader.Service
	if clientReader != nil {
		svcReader = reader.NewService(clientReader)
	}
	//
	var svcPermits permits.Service
	if b.apiUri != "" {
		clientPermits := permits.NewClientConnPool(connPoolApi)
		svcPermits = permits.NewService(clientPermits)
	}
	//
	var clientSubs subscriptions.ServiceClient
	if b.subsUri != "" {
		clientSubs = subscriptions.NewClientConnPool(connPoolSubs)
	} else if b.apiUri != "" {
		clientSubs = subscriptions.NewClientConnPool(connPoolApi)
	}
	var svcSubs subscriptions.Service
	if clientSubs != nil {
		svcSubs = subscriptions.NewService(clientSubs)
	}
	//
	var clientWriter resolver.ServiceClient
	if b.writerUri != "" {
		clientWriter = resolver.NewClientConnPool(connPoolWriter)
	} else if b.apiUri != "" {
		clientWriter = resolver.NewClientConnPool(connPoolApi)
	}
	var svcWriter resolver.Service
	if clientWriter != nil {
		svcWriter = resolver.NewService(clientWriter)
	}
	//
	c = client{
		connPoolApi:    connPoolApi,
		connPoolReader: connPoolReader,
		connPoolSubs:   connPoolSubs,
		connPoolWriter: connPoolWriter,
		svcLimits:      svcLimits,
		svcReader:      svcReader,
		svcPermits:     svcPermits,
		svcSubs:        svcSubs,
		svcWriter:      svcWriter,
	}
	return
}
