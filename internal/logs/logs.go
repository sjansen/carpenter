package logs

import (
	"regexp"
)

var albRE = regexp.MustCompile(`^` +
	`([^ ]*) ` + // type
	`([^ ]*) ` + // timestamp
	`([^ ]*) ` + // elb
	`([^ ]*):([0-9]*) ` + // client:port
	`([^ ]*)[:-]([0-9]*) ` + // target:port
	`([-.0-9]*) ` + // request_processing_time
	`([-.0-9]*) ` + // target_processing_time
	`([-.0-9]*) ` + // response_processing_time
	`([-0-9]*) ` + // elb_status_code
	`(-|[-0-9]*) ` + // target_status_code
	`([-0-9]*) ` + // received_bytes
	`([-0-9]*) ` + // sent_bytes
	`"([^ ]*) ([^ ]*) (- |[^ ]*)" ` + // "request"
	`"([^"]*)" ` + // "user_agent"
	`([A-Z0-9-]+) ` + // ssl_cipher
	`([A-Za-z0-9.-]*) ` + // ssl_protocol
	`([^ ]*) ` + // target_group_arn
	`"([^"]*)" ` + // "trace_id"
	`"([^"]*)" ` + // "domain_name"
	`"([^"]*)" ` + // "chosen_cert_arn"
	`([-.0-9]*) ` + // matched_rule_priority
	`([^ ]*) ` + // request_creation_time
	`"([^"]*)" ` + // "actions_executed"
	`"([^"]*)"` + // "redirect_url"
	`(?: "([^ ]*)")?` + // "error_reason"
	`(?: "([^"]*)")?` + // "target:port_list"
	`(?: "([^"]*)")?` + // "target_status_code_list"
	`(.*?)$`,
)

var zeroALB = &ALB{}

type ALB struct {
	Proto                  string
	Timestamp              string
	ELB                    string
	ClientIP               string
	ClientPort             string
	TargetIP               string
	TargetPort             string
	RequestProcessingTime  string
	TargetProcessingTime   string
	ResponseProcessingTime string
	ELBStatusCode          string
	TargetStatusCode       string
	ReceivedBytes          string
	SentBytes              string
	RequestVerb            string
	RequestURL             string
	RequestProto           string
	UserAgent              string
	SSLCipher              string
	SSLProtocol            string
	TargetGroupARN         string
	TraceID                string
	DomainName             string
	ChosenCertARN          string
	MatchedRulePriority    string
	RequestCreationTime    string
	ActionsExecuted        string
	RedirectURL            string
	LambdaErrorReason      string
	TargetPortList         string
	TargetStatusCodeList   string
	Extra                  string
}

func (l *ALB) Parse(line string) {
	m := albRE.FindStringSubmatch(line)
	if m == nil {
		*l = *zeroALB
		return
	}
	l.Proto = m[1]
	l.Timestamp = m[2]
	l.ELB = m[3]
	l.ClientIP = m[4]
	l.ClientPort = m[5]
	l.TargetIP = m[6]
	l.TargetPort = m[7]
	l.RequestProcessingTime = m[8]
	l.TargetProcessingTime = m[9]
	l.ResponseProcessingTime = m[10]
	l.ELBStatusCode = m[11]
	l.TargetStatusCode = m[12]
	l.ReceivedBytes = m[13]
	l.SentBytes = m[14]
	l.RequestVerb = m[15]
	l.RequestURL = m[16]
	l.RequestProto = m[17]
	l.UserAgent = m[18]
	l.SSLCipher = m[19]
	l.SSLProtocol = m[20]
	l.TargetGroupARN = m[21]
	l.TraceID = m[22]
	l.DomainName = m[23]
	l.ChosenCertARN = m[24]
	l.MatchedRulePriority = m[25]
	l.RequestCreationTime = m[26]
	l.ActionsExecuted = m[27]
	l.RedirectURL = m[28]
	l.LambdaErrorReason = m[29]
	l.TargetPortList = m[30]
	l.TargetStatusCodeList = m[31]
	l.Extra = m[32]
}
