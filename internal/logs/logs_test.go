package logs

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseALB(t *testing.T) {
	for tc, expected := range map[string]*ALB{
		"http": {
			Proto:                  "http",
			Timestamp:              "2018-07-02T22:23:00.186641Z",
			ELB:                    "app/my-loadbalancer/50dc6c495c0c9188",
			ClientIP:               "192.168.131.39",
			ClientPort:             "2817",
			TargetIP:               "10.0.0.1",
			TargetPort:             "80",
			RequestProcessingTime:  "0.000",
			TargetProcessingTime:   "0.001",
			ResponseProcessingTime: "0.002",
			ELBStatusCode:          "200",
			TargetStatusCode:       "200",
			ReceivedBytes:          "34",
			SentBytes:              "366",
			RequestVerb:            "GET",
			RequestURL:             "http://www.example.com:80/",
			RequestProto:           "HTTP/1.1",
			UserAgent:              "curl/7.46.0",
			SSLCipher:              "-",
			SSLProtocol:            "-",
			TargetGroupARN: "arn:aws:elasticloadbalancing:us-east-2:123456789012:" +
				"targetgroup/my-targets/73e2d6bc24d8a067",
			TraceID:              "Root=1-58337262-36d228ad5d99923122bbe354",
			DomainName:           "-",
			ChosenCertARN:        "-",
			MatchedRulePriority:  "0",
			RequestCreationTime:  "2018-07-02T22:22:48.364000Z",
			ActionsExecuted:      "forward",
			RedirectURL:          "-",
			LambdaErrorReason:    "-",
			TargetPortList:       "",
			TargetStatusCodeList: "",
		},
	} {
		require := require.New(t)

		line, err := ioutil.ReadFile("testdata/" + tc + ".txt")
		require.NoError(err)

		actual := &ALB{}
		actual.Parse(string(bytes.TrimSpace(line)))
		require.Equal(expected, actual)
	}
}
