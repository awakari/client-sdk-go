# Contents

1. [Overview](#1-overview)<br/>
2. [Usage](#2-usage)<br/>
   2.1. [Prerequisites](#21-prerequisites)<br/>
3. [Contributing](#6-contributing)<br/>
   3.1. [Versioning](#31-versioning)<br/>
   3.2. [Issue Reporting](#32-issue-reporting)<br/>
   3.3. [Testing](#33-testing)<br/>
   &nbsp;&nbsp;&nbsp;3.3.1. [Functional](#331-functional)<br/>
   &nbsp;&nbsp;&nbsp;3.3.2. [Performance](#332-performance)<br/>
   3.4. [Releasing](#34-releasing)<br/>

# 1. Overview

Reference Awakari SDK for a Golang client.

## 1.1. Security

To secure the Awakari public API usage, the mutual TLS encryption is used together with additional user authentication.

### 1.1.1. Server Public Key

Used to authenticate the Awakari service by the client. A client should fetch it by a public link TODO.

### 1.1.2. Client Key Pair

Used to authenticate the Group Client. Contains a client's private key and a client's certificate. A client should 
request Awakari contacts to obtain it. 

A client's certificate is used by Awakari to extract the ***DN*** value. This value, e.g. 
`CN=my-service-using-awakari.com` is treated by Awakari as Group Client Identity.

### 1.1.3. User Identity

> **Info**:
> 
> Any Group Client can be used by multiple users.

A client should specify the user id in every API call. Awakari doesn't verify a user identity and trusts any user id 
specified by the client. Hence, a client should authenticate the users by themselves, for example integrating a 3-rd 
party identity provider. The good example is to use the `sub` field from a JWT token as a user id.

# 2. Usage

See the [int_test.go](int_test.go) for the code example.

## 2.1. Prerequisites

The Awakari public API requires:
1. Server public key and client key pair for the mutual TLS encryption

# 3. Contributing

## 3.1. Versioning

The library follows the [semantic versioning](http://semver.org/).
The single source of the version info is the git tag:
```shell
git describe --tags --abbrev=0
```

## 3.2. Issue Reporting

TODO

## 3.3. Testing

### 3.3.1. Functional

```shell
API_URI=api.local:443 \
CLIENT_CERT_PATH=test0.client0.crt \
CLIENT_PRIVATE_KEY_PATH=test0.client0.key \
SERVER_PUBLIC_KEY_PATH=ca.crt \
make test
```

### 3.3.2. Performance

TODO

## 3.4. Releasing

To release a new version (e.g. `1.2.3`) it's enough to put a git tag:
```shell
git tag -v1.2.3
git push --tags
```

The corresponding CI job is started to build a docker image and push it with the specified tag (+latest).
