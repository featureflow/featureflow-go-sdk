# featureflow-go-sdk
Go Client for featureflow


Usage
```go
package main

import "github.com/featureflow/featureflow-go-sdk/featureflow"

var client, _ = featureflow.Client("srv-env-<api_key>", featureflow.Config{})
```

Evaluate using

```go
package main

import "github.com/featureflow/featureflow-go-sdk/featureflow"

//Get user somewhere in your code
func main(){ 
    client, _ := featureflow.Client("srv-env-<api_key>", featureflow.Config{})
    user, _ := featureflow.NewUserBuilder("userId").
                               WithAttributes("roles", []string{"admin", "user"}).
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


