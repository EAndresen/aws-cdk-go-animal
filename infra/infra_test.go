package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestAwsCdkGoStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewLambdaCronStack(app, "MyStack", nil)

	// THEN
	bytes, err := json.Marshal(app.Synth(nil).GetStackArtifact(stack.ArtifactId()).Template())
	if err != nil {
		t.Error(err)
	}

	template := gjson.ParseBytes(bytes)
	displayName := template.Get("Resources.AnimalTableA30B0E36.Properties.BillingMode").String()
	assert.Equal(t, "PAY_PER_REQUEST", displayName)
}
