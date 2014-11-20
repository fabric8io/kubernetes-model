package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	buildapi "github.com/openshift/origin/pkg/build/api"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
)

type Schema struct {
	Build            buildapi.Build
	BuildConfig      buildapi.BuildConfig
	Deployment       deployapi.Deployment
	DeploymentConfig deployapi.DeploymentConfig
}

func main() {
	schema, err := GenerateSchema(reflect.TypeOf(Schema{}))
	if err != nil {
		fmt.Errorf("An error occurred: %v", err)
		return
	}
	b, _ := json.Marshal(&schema)
	result := string(b)
	result = strings.Replace(result, "\"additionalProperty\":", "\"additionalProperties\":", -1)
	fmt.Println(result)
}
