# Contents

1. [Overview](#1-overview)<br/>
2. [Security](#2-security)<br/>
   2.1. [Certificate Authority](#21-certificate-authority)<br/>
   2.2. [Client Key Pair](#22-client-key-pair)<br/>
   2.3. [User Identity](#23-user-identity)<br/>
3. [Usage](#3-usage)<br/>
   3.1. [Limits](#31-limits)<br/>
   3.2. [Permits](#32-permits)<br/>
   3.3. [Subscriptions](#33-subscriptions)<br/>
   3.4. [Messages](#34-messages)<br/>
4. [Contributing](#4-contributing)<br/>
   4.1. [Versioning](#41-versioning)<br/>
   4.2. [Issue Reporting](#42-issue-reporting)<br/>
   4.3. [Testing](#43-testing)<br/>
   &nbsp;&nbsp;&nbsp;4.3.1. [Functional](#431-functional)<br/>
   &nbsp;&nbsp;&nbsp;4.3.2. [Performance](#432-performance)<br/>
   4.4. [Releasing](#44-releasing)<br/>

# 1. Overview

Reference Awakari SDK for a Golang client.

# 2. Security

To secure the Awakari public API usage, the mutual TLS encryption is used together with additional user identity.

> **Note**:
>
> Not available for self-hosted core system. 
> Skip the [2. Security](#2-security) section entirely when using self-hosted core system.

## 2.1. Certificate Authority

Used to authenticate the Awakari service by the client. A client should fetch it, for example: [demo instance CA](https://awakari.com/certs/ca-demo.awakari.cloud.crt).

```go
package main

import (
   "github.com/awakari/client-sdk-go/api"
   "os"
   ...
)

func main() {
   ...
   caCrt, err := os.ReadFile("ca.crt")
   if err != nil {
	   panic(err)
   }
   ...
   client, err := api.
       NewClientBuilder().
       ...
       CertAuthority(caCrt).
       ...
       Build()
   ...
}
```

## 2.2. Client Key Pair

Used to authenticate the Group Client. Contains a client's private key and a client's certificate. A client should 
request Awakari contacts to obtain it. 

A client's certificate is used by Awakari to extract the ***DN*** value. This value, e.g. 
`CN=my-service-using-awakari.com` is treated by Awakari as Group Client Identity.

```go
package main

import (
   "github.com/awakari/client-sdk-go/api"
   "os"
   ...
)

func main() {
   ...
   clientCrt, err := os.ReadFile("client.crt")
   if err != nil {
      panic(err)
   }
   clientKey, err := os.ReadFile("client.key")
   if err != nil {
      panic(err)
   }
   ...
   client, err := api.
      NewClientBuilder().
      ...
      ClientKeyPair(clientCrt, clientKey).
      ...
      Build()
   ...
}
```

## 2.3. User Identity

> **Note**:
> 
> * Any Group Client can be used by many users.
> * Awakari doesn't verify a user identity and trusts any user id specified by the client.

A client is required to specify a user id in every API call. The user authentication and authorization are the client's 
responsibility. The good example is to integrate a 3-rd party identity provider and use the `sub` field from a JWT token 
as a user id.

# 3. Usage

See the [int_test.go](int_test.go) for the complete test code example.

Before using the API, it's necessary to initialize the client. 
When using a hybrid deployment the initialization should be like follows:

```go
package main

import (
   "github.com/awakari/client-sdk-go/api"
   "os"
   ...
)

func main() {
   ...
   client, err := api.
       NewClientBuilder().
       ReaderUri("core-reader:50051"). // skip this line if reader API is not used
       SubscriptionsUri("core-subscriptionsproxy:50051"). // skip this line if subscriptions API is not used
       WriterUri("core-resolver:50051"). // skip this line if writer API is not used
       Build()
   if err != nil {
       panic(err)
   }
   defer client.Close()
   ...
}
```

The initialization is a bit different for a serverless API usage:

```go
package main

import (
   "github.com/awakari/client-sdk-go/api"
   "os"
   ...
)

func main() {
   ...
   caCrt, err := os.ReadFile("ca.crt")
   if err != nil {
	   panic(err)
   }
   clientCrt, err := os.ReadFile("client.crt")
   if err != nil {
      panic(err)
   }
   clientKey, err := os.ReadFile("client.key")
   if err != nil {
      panic(err)
   }
   client, err := api.
       NewClientBuilder().
       CertAuthority(caCrt).
       ClientKeyPair(clientCrt, clientKey).
       ApiUri("demo.awakari.com:443").
       Build()
   if err != nil {
       panic(err)
   }
   defer client.Close()
   ...
}
```

## 3.1. Limits

> **Note**:
>
> Limits API is not available for self-hosted core system.
> Skip the [3.1. Limits](#31-limits) section entirely when using self-hosted core system.

Usage limit represents the successful API call count limit. The limit is identified per:
* group id
* user id (optional)
* subject

There are the group-level limits where user id is not specified. All users from the group share the group limit in this
case.

Usage subject may be one of:
* Subscriptions
* Publish Events

```go
package main

import (
   "context"
   "fmt"
   "github.com/awakari/client-sdk-go/api"
   "github.com/awakari/client-sdk-go/model/usage"
   "time"
   ...
)

func main() {
   ...
   var client api.Client // TODO initialize client here
   var userId string     // set this to "sub" field value from an authentication token, for example
   ...
   ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
   defer cancel()
   var l usage.Limit
   var err error
   l, err = client.ReadUsageLimit(ctx, userId, usage.SubjectPublishEvents)
   if err == nil {
      if u.UserId == "" {
         fmt.Printf("group usage publish events limit: %d", l.Count)
      } else {
         fmt.Printf("user specific publish events limit: %d", l.Count)
      }
   }
   ...
}
```

## 3.2. Permits

> **Note**:
>
> Permits API is not available for self-hosted core system.
> Skip the [3.2. Permits](#32-permits) section entirely when using self-hosted core system.

Usage permits represents the current usage statistics (counters) by the subject. Similar to usage limit, the counters
represent the group-level usage when the user id is empty.

```go
package main

import (
   "context"
   "fmt"
   "github.com/awakari/client-sdk-go/api"
   "github.com/awakari/client-sdk-go/model/usage"
   "time"
   ...
)

func main() {
   ...
   var client api.Client // TODO initialize client here
   var userId string     // set this to "sub" field value from an authentication token, for example
   ...
   ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
   defer cancel()
   var u usage.Usage
   var err error
   u, err = client.ReadUsage(ctx, userId, usage.SubjectSubscriptions)
   if err == nil {
      if u.UserId == "" {
         fmt.Printf("group subscriptions usage: %d", l.Count)
      } else {
         fmt.Printf("user specific subscriptions usage: %d", l.Count)
      }
   }
   ...
}
```

## 3.3. Subscriptions

```go
package main

import (
   "context"
   "fmt"
   "github.com/awakari/client-sdk-go/api"
   "github.com/awakari/client-sdk-go/model/usage"
   "time"
   ...
)

func main() {
   ...
   var client api.Client // TODO initialize client here
   var userId string     // set this to "sub" field value from an authentication token, for example
   ...
   ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Minute)
   defer cancel()
   
   // Create a subscription
   var subId string
   var err error
   subData := subscription.Data{
      Description: "my subscription",
	  Enabled:     true,
      Condition: condition.NewBuilder().
         MatchAttrKey("tags").
         MatchText("SpaceX").
         BuildTextCondition(),
   }
   subId, err = client.CreateSubscription(ctx, userId, subData)
   
   // Update the subscription mutable fields
   upd := subscription.Data{
      Description: "my disabled subscription",
      Enabled:     false,
   }
   err = client.UpdateSubscription(ctx, userId, subId, upd)
   
   // Delete the subscription
   err = client.DeleteSubscription(ctx, userId, subId)
   if err != nil {
      panic(err)
   }
   
   // Search own subscription ids
   var ids []string
   limit := uint32(10)
   ids, err = client.Search(ctx, userId, limit, "")
   if err != nil {
      panic(err)
   } 
   for _, id := range ids {
      // Read the subscription details
      subData, err = client.Read(ctx, userId, id)
      if err == nil {
         panic(err)
      }
      fmt.Printf("subscription %d details: %+v\n", id, subData)
   }
   
   ...
}
```

## 3.4. Messages

### 3.4.1. Publishing

```go
package main

import (
   "context"
   "fmt"
   "github.com/awakari/client-sdk-go/api"
   "github.com/awakari/client-sdk-go/model/usage"
   "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
   "time"
   ...
)

func main() {
   ...
   var client api.Client // TODO initialize client here
   var userId string     // set this to "sub" field value from an authentication token, for example
   ...
   ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
   defer cancel()
   var ws model.WriteStream[*pb.CloudEvent]
   ws, err = client.OpenMessagesWriter(ctx, userId)
   if err == nil {
      panic(err)
   }
   defer ws.Close()
   msgs := []*pb.CloudEvent{
      {
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
      },
   }
   
   var writtenCount uint32
   var n uint32
   for writtenCount < uint32(len(msgs)) {
      n, err = ws.WriteBatch(msgs)
      if err != nil {
         break
      }
      writtenCount += n
   }
   if err != nil {
      panic(err)
   }
   ...
}
```

### 3.4.2. Receiving

```go
package main

import (
   "context"
   "fmt"
   "github.com/awakari/client-sdk-go/api"
   "github.com/awakari/client-sdk-go/model/usage"
   "time"
...
)

func main() {
   ...
   var client api.Client // TODO initialize client here
   var userId string     // set this to "sub" field value from an authentication token, for example
   batchSize := uint32(16)
   ...
   ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
   defer cancel()
   var r model.ReadStream[*pb.CloudEvent]
   r, err = client.OpenMessagesReader(ctx, userId, subId, batchSize)
   if err != nil {
      panic(err)
   }
   defer r.Close()
   var msgs []*pb.CloudEvent
   for {
      msgs, err = r.Read()
      if err != nil {
         break
      }
      fmt.Printf("subscription %s - received the next messages batch: %+v\n", subId, msgs)
   }
   if err != nil {
      panic(err)
   }
   ...
}
```

# 4. Contributing

## 4.1. Versioning

The library follows the [semantic versioning](http://semver.org/).
The single source of the version info is the git tag:
```shell
git describe --tags --abbrev=0
```

## 4.2. Issue Reporting

TODO

## 4.3. Testing

### 4.3.1. Functional

```shell
API_URI=api.local:443 \
CA_PATH=ca.crt \
CLIENT_CERT_PATH=test0.client0.crt \
CLIENT_PRIVATE_KEY_PATH=test0.client0.key \
make test
```

### 4.3.2. Performance

TODO

## 4.4. Releasing

To release a new version (e.g. `1.2.3`) it's enough to put a git tag:
```shell
git tag -v1.2.3
git push --tags
```

The corresponding CI job is started to build a docker image and push it with the specified tag (+latest).
