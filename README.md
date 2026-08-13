# featureflow-go-sdk
Go Client for featureflow

Installation
```shell
go get github.com/featureflow/featureflow-go-sdk@latest
```

Usage
```go
package main

import "github.com/featureflow/featureflow-go-sdk/featureflow"

var client, _ = featureflow.Client("sdk-srv-env-<api_key>", featureflow.Config{})
```

Evaluate using

```go
package main

import "github.com/featureflow/featureflow-go-sdk/featureflow"

//Get user somewhere in your code
func main(){ 
    client, _ := featureflow.Client("sdk-srv-env-<api_key>", featureflow.Config{})
    user, _ := featureflow.NewUserBuilder("userId").
                               WithAttributes("roles", []featureflow.Attribute{"admin", "user"}).
                               WithAttribute("age", 20).
                               Build()
                 
    if client.Evaluate("my-feature", user).Is("on"){ // same as .IsOn(), also use .IsOff() == .Is("off")
        //feature variant is turend on for this user
    }  
}

```


Register features using
```go
package main
import "github.com/featureflow/featureflow-go-sdk/featureflow"

func main(){
    config := featureflow.Config{
        WithFeatures: []featureflow.FeatureRegistration{
            featureflow.WithFeature("feature-1", "off").Build(),
            featureflow.WithFeature("feature-2", "off").
                        AddVariant("key1","Key 1 Title").
                        AddVariant("key1","key 2 Title").
                        Build(),
        },
    }
}
//Note if you don't add 2 variants, it will set the default variants to "on" and "off"
```

### Naming your application

Optionally tag this workload with an application name so the Featureflow dashboard can
attribute SDK usage and flag evaluations to it (Admin → SDKs, and the "Evaluated by"
panel on each feature's statistics tab):

```go
client, _ := featureflow.Client("sdk-srv-env-<api_key>", featureflow.Config{
    Application: "checkout-api",
})
```

The name is a slug — lowercase letters, numbers, `.`, `_` and `-`, at most 64
characters. An invalid value is dropped with a warning and no tag is sent. The
`FEATUREFLOW_APPLICATION` environment variable is used when the option is not set in
code.


