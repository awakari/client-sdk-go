# Contents

1. [Overview](#1-overview)<br/>
2. [Usage](#2-usage)<br/>
3. [Contributing](#6-contributing)<br/>
   3.1. [Versioning](#31-versioning)<br/>
   3.2. [Issue Reporting](#32-issue-reporting)<br/>
   3.3. [Testing](#33-testing)<br/>
   &nbsp;&nbsp;&nbsp;3.3.1. [Functional](#331-functional)<br/>
   &nbsp;&nbsp;&nbsp;3.3.2. [Performance](#332-performance)<br/>
   3.4. [Releasing](#34-releasing)<br/>

# 1. Overview

Reference Awakari SDK for a Golang client

# 2. Usage

```go
package main

import (
   clientSdk "github.com/awakari/client-sdk-go"
   "os"
)

func main() {
	
	// load TLS certificates
    var caCrt []byte
	var err error
	caCrt, err = os.ReadFile("ca.crt")
	var clientCrt []byte
	if err == nil {
		clientCrt, err = os.ReadFile("client.crt")
	}
	var clientKey []byte
	if err == nil {
		clientKey, err = os.ReadFile("client.key")
    }
	
	var client clientSdk.Client
	client, err = clientSdk.
		NewBuilder().
		ServerPublicKey(caCrt).
		ClientKeyPair(clientCrt, clientKey).
		ApiUri("awakari.com:443").
		Build()
	if err == nil {
		defer client.Close() 
		// TODO: use the client here to manage subscriptions, publish and receive messages, etc... 
	}
}
```

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
