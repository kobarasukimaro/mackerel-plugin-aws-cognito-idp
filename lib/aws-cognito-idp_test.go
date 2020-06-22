package mpawscognitoidp

import (
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

func ExampleCognitoIdpPlugin_GraphDefinition() {
	var cognitoidp CognitoIdpPlugin
	cognitoidp.Prefix = "test-pool"
	helper := mp.NewMackerelPlugin(cognitoidp)
	helper.OutputDefinitions()

	// Output:
	// # mackerel-agent-plugin
	// {"graphs":{"test-pool.federation":{"label":"federation","unit":"integer","metrics":[{"name":"FederationSuccesses","label":"Success Count","stacked":false},{"name":"FederationParcentageOfSuccessful","label":"Parcentage Of Successfull","stacked":false},{"name":"FederationSampleCount","label":"Request Count","stacked":false},{"name":"FederationThrottles","label":"Throttles Count","stacked":false}]},"test-pool.signin":{"label":"signin","unit":"integer","metrics":[{"name":"SignInSuccesses","label":"Success Count","stacked":false},{"name":"SignInParcentageOfSuccessful","label":"Parcentage Of Successfull","stacked":false},{"name":"SignInSampleCount","label":"Request Count","stacked":false},{"name":"SignInThrottles","label":"Throttles Count","stacked":false}]},"test-pool.signup":{"label":"signup","unit":"integer","metrics":[{"name":"SignUpSuccesses","label":"Success Count","stacked":false},{"name":"SignUpParcentageOfSuccessful","label":"Parcentage Of Successfull","stacked":false},{"name":"SignUpSampleCount","label":"Request Count","stacked":false},{"name":"SignUpThrottles","label":"Throttles Count","stacked":false}]},"test-pool.tokenRefresh":{"label":"tokenRefresh","unit":"integer","metrics":[{"name":"TokenRefreshSuccesses","label":"Success Count","stacked":false},{"name":"TokenRefreshParcentageOfSuccessful","label":"Parcentage Of Successfull","stacked":false},{"name":"TokenRefreshSampleCount","label":"Request Count","stacked":false},{"name":"TokenRefreshThrottles","label":"Throttles Count","stacked":false}]}}}
}
