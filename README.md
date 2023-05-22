# Contents

1. [Overview](#1-overview)<br/>
2. [Security](#2-security)<br/>
   2.1. [Server Public Key](#21-server-public-key)<br/>
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

## 2.1. Server Public Key

Used to authenticate the Awakari service by the client. A client should fetch it by a public link TODO.

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
       ServerPublicKey(caCrt).
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

See the [int_test.go](int_test.go) for the code example.

## 3.1. Limits

## 3.2. Permits

## 3.3. Subscriptions

## 3.4. Messages

### 3.4.1. Publishing

### 3.4.2. Receiving

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
CLIENT_CERT_PATH=test0.client0.crt \
CLIENT_PRIVATE_KEY_PATH=test0.client0.key \
SERVER_PUBLIC_KEY_PATH=ca.crt \
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
