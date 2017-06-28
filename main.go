package main

import (
	featureflow "./featureflow"
	"strconv"
	"fmt"
)

func main(){
	client, _ := featureflow.Client("srv-env-9b5fff890c724d119a334a64ed4d2eb2", featureflow.Config{})

	ctx, _ := featureflow.NewContextBuilder("user1").Build()

	key := "something"

	fmt.Printf("%s evaluates to %s and testing IsOn() returns %s",
		key,
		client.Evaluate(key, ctx),
		strconv.FormatBool(client.Evaluate(key, ctx).IsOff()),
	)

}
