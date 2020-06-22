package mpawscognitoidp

import (
	"errors"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

const (
	namespace              = "AWS/Cognito"
	metricsTypeSum         = "Sum"
	metricsTypeAverage     = "Average"
	metricsTypeMaximum     = "Maximum"
	metricsTypeMinimum     = "Minimum"
	metricsTypeSampleCount = "Sample Count"

)

type metrics struct {
	CloudWatchName string
	MackerelName   string
	Type           string
}

// CognitoIdpPlugin mackerel plugin for aws cognito
type CognitoIdpPlugin struct {
	Name            string
	Prefix          string
	PoolId          string
	PoolClientId    string

	AccessKeyID     string
	SecretAccessKey string
	Region          string
	CloudWatch      *cloudwatch.CloudWatch
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p CognitoIdpPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "cognito-idp"
	}
	return p.Prefix
}

// prepare creates CloudWatch instance
func (p *CognitoIdpPlugin) prepare() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	config := aws.NewConfig()
	if p.AccessKeyID != "" && p.SecretAccessKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(p.AccessKeyID, p.SecretAccessKey, ""))
	}
	if p.Region != "" {
		config = config.WithRegion(p.Region)
	}

	p.CloudWatch = cloudwatch.New(sess, config)

	return nil
}

// getLastPoint fetches a CloudWatch metric and parse
func (p CognitoIdpPlugin) getLastPoint(metric metrics) (float64, error) {
	now := time.Now()

	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("UserPool"),
			Value: aws.String(p.PoolId),
		},
		{
			Name:  aws.String("UserPoolClient"),
			Value: aws.String(p.PoolClientId),
		},
	}

	response, err := p.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: dimensions,
		StartTime:  aws.Time(now.Add(time.Duration(300) * 3 * time.Second * -1)),
		EndTime:    aws.Time(now),
		MetricName: aws.String(metric.CloudWatchName),
		Period:     aws.Int64(300),
		Statistics: []*string{aws.String(metric.Type)},
		Namespace:  aws.String(namespace),
	})
	if err != nil {
		return 0, err
	}

	datapoints := response.Datapoints
	if len(datapoints) == 0 {
		return 0, errors.New("fetched no datapoints")
	}

	latest := new(time.Time)
	var latestVal float64
	for _, dp := range datapoints {
		if dp.Timestamp.Before(*latest) {
			continue
		}

		latest = dp.Timestamp
		switch metric.Type {
		case metricsTypeSum:
			latestVal = *dp.Sum
		case metricsTypeAverage:
			latestVal = *dp.Average
		case metricsTypeSampleCount:
			latestVal = *dp.SampleCount
		}
	}

	return latestVal, nil
}

// FetchMetrics fetch the metrics
func (p CognitoIdpPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})

	for _, met := range [...]metrics{
		{CloudWatchName: "SignUpSuccesses", MackerelName: "SignUpSuccesses", Type: metricsTypeSum},
		{CloudWatchName: "SignUpSuccesses", MackerelName: "SignUpParcentageOfSuccessful", Type: metricsTypeAverage},
		{CloudWatchName: "SignUpSuccesses", MackerelName: "SignUpSampleCount", Type: metricsTypeSampleCount},
		{CloudWatchName: "SignUpThrottles", MackerelName: "SignUpThrottles", Type: metricsTypeSum},
		{CloudWatchName: "SignInSuccesses", MackerelName: "SignInSuccesses", Type: metricsTypeSum},
		{CloudWatchName: "SignInSuccesses", MackerelName: "SignInParcentageOfSuccessful", Type: metricsTypeAverage},
		{CloudWatchName: "SignInSuccesses", MackerelName: "SignInSampleCount", Type: metricsTypeSampleCount},
		{CloudWatchName: "SignInThrottles", MackerelName: "SignInThrottles", Type: metricsTypeSum},
		{CloudWatchName: "TokenRefreshSuccesses", MackerelName: "TokenRefreshSuccesses", Type: metricsTypeSum},
		{CloudWatchName: "TokenRefreshSuccesses", MackerelName: "TokenRefreshParcentageOfSuccessful", Type: metricsTypeAverage},
		{CloudWatchName: "TokenRefreshSuccesses", MackerelName: "TokenRefreshSampleCount", Type: metricsTypeSampleCount},
		{CloudWatchName: "TokenRefreshThrottles", MackerelName: "TokenRefreshThrottles", Type: metricsTypeSum},
		{CloudWatchName: "FederationSuccesses", MackerelName: "FederationSuccesses", Type: metricsTypeSum},
		{CloudWatchName: "FederationSuccesses", MackerelName: "FederationParcentageOfSuccessful", Type: metricsTypeAverage},
		{CloudWatchName: "FederationSuccesses", MackerelName: "FederationSampleCount", Type: metricsTypeSampleCount},
		{CloudWatchName: "FederationThrottles", MackerelName: "FederationThrottles", Type: metricsTypeSum},
	} {
		v, err := p.getLastPoint(met)
		if err == nil {
			stat[met.MackerelName] = v
		} else {
			log.Printf("%s: %s", met, err)
		}
	}
	return stat, nil
}

// GraphDefinition of CognitoIdpPlugin
func (p CognitoIdpPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.Prefix)
	labelPrefix = strings.Replace(labelPrefix, "-", " ", -1)

	var graphdef = map[string]mp.Graphs{
		"signup": {
			Label: ("signup"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "SignUpSuccesses", Label: "Success Count"},
				{Name: "SignUpParcentageOfSuccessful", Label: "Parcentage Of Successfull"},
				{Name: "SignUpSampleCount", Label: "Request Count"},
				{Name: "SignUpThrottles", Label: "Throttles Count"},
			},
		},
		"signin": {
			Label: ("signin"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "SignInSuccesses", Label: "Success Count"},
				{Name: "SignInParcentageOfSuccessful", Label: "Parcentage Of Successfull"},
				{Name: "SignInSampleCount", Label: "Request Count"},
				{Name: "SignInThrottles", Label: "Throttles Count"},
			},
		},
		"tokenRefresh": {
			Label: ("tokenRefresh"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "TokenRefreshSuccesses", Label: "Success Count"},
				{Name: "TokenRefreshParcentageOfSuccessful", Label: "Parcentage Of Successfull"},
				{Name: "TokenRefreshSampleCount", Label: "Request Count"},
				{Name: "TokenRefreshThrottles", Label: "Throttles Count"},
			},
		},
		"federation": {
			Label: ("federation"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "FederationSuccesses", Label: "Success Count"},
				{Name: "FederationParcentageOfSuccessful", Label: "Parcentage Of Successfull"},
				{Name: "FederationSampleCount", Label: "Request Count"},
				{Name: "FederationThrottles", Label: "Throttles Count"},
			},
		},
	}
	return graphdef
}

// Do the plugin
func Do() {
	optAccessKeyID := flag.String("access-key-id", "", "AWS Access Key ID")
	optSecretAccessKey := flag.String("secret-access-key", "", "AWS Secret Access Key")
	optRegion := flag.String("region", "", "AWS Region")
	optPoolId := flag.String("pool-id", "", "Pool Id")
	optPoolClientId := flag.String("pool-client-id", "", "Pool Client Id")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optPrefix := flag.String("metric-key-prefix", "kinesis-streams", "Metric key prefix")
	flag.Parse()

	var plugin CognitoIdpPlugin

	plugin.AccessKeyID = *optAccessKeyID
	plugin.SecretAccessKey = *optSecretAccessKey
	plugin.Region = *optRegion
	plugin.Prefix = *optPrefix
	plugin.PoolId = *optPoolId
	plugin.PoolClientId = *optPoolClientId

	err := plugin.prepare()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(plugin)
	helper.Tempfile = *optTempfile

	helper.Run()
}
