package tokenizer

import (
	"io/ioutil"
	"regexp"

	"github.com/sjansen/carpenter/internal/data"
	"github.com/ua-parser/uap-go/uaparser"
)

type Tokenizer struct {
	*regexp.Regexp

	uaparser *uaparser.Parser
}

func (t *Tokenizer) EnableUserAgentParsing() error {
	r, err := data.Assets.Open("regexes.yaml")
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	uap, err := uaparser.NewFromBytes(bytes)
	if err != nil {
		return err
	}

	t.uaparser = uap
	return nil
}

func (t *Tokenizer) Tokenize(line string) map[string]string {
	values := t.FindStringSubmatch(line)
	if values == nil {
		return nil
	}

	names := t.SubexpNames()[1:]
	result := make(map[string]string, len(names))
	for i, key := range names {
		result[key] = values[i+1]
	}

	uagent, ok := result["user_agent"]
	if ok && t.uaparser != nil {
		client := t.uaparser.Parse(uagent)
		result["client_device_family"] = client.Device.Family
		result["client_os_family"] = client.Os.Family
		result["client_os_major"] = client.Os.Major
		result["client_os_minor"] = client.Os.Minor
		result["client_os_patch"] = client.Os.Patch
		result["client_ua_family"] = client.UserAgent.Family
		result["client_ua_major"] = client.UserAgent.Major
		result["client_ua_minor"] = client.UserAgent.Minor
		result["client_ua_patch"] = client.UserAgent.Patch
	}

	return result
}

var ALB = &Tokenizer{regexp.MustCompile(`^` +
	`(?P<type>[^ ]*) ` +
	`(?P<timestamp>[^ ]*) ` +
	`(?P<lb>[^ ]*) ` +
	`(?P<client_ip>[^ ]*):(?P<client_port>[0-9]*) ` +
	`(?P<target_ip>[^ ]*)[:-](?P<target_port>[0-9]*) ` +
	`(?P<request_processing_time>[-.0-9]*) ` +
	`(?P<target_processing_time>[-.0-9]*) ` +
	`(?P<response_processing_time>[-.0-9]*) ` +
	`(?P<lb_status_code>[-0-9]*) ` +
	`(?P<target_status_code>-|[-0-9]*) ` +
	`(?P<received_bytes>[-0-9]*) ` +
	`(?P<sent_bytes>[-0-9]*) ` +
	`"(?P<request_verb>[^ ]*) (?P<request_url>[^ ]*) (?P<request_proto>- |[^ ]*)" ` +
	`"(?P<user_agent>[^"]*)" ` +
	`(?P<tls_cipher>[A-Z0-9-]+) ` +
	`(?P<tls_protocol>[A-Za-z0-9.-]*) ` +
	`(?P<target_group_arn>[^ ]*) ` +
	`"(?P<trace_id>[^"]*)" ` +
	`"(?P<domain_name>[^"]*)" ` +
	`"(?P<chosen_cert_arn>[^"]*)" ` +
	`(?P<matched_rule_priority>[-.0-9]*) ` +
	`(?P<request_creation_time>[^ ]*) ` +
	`"(?P<actions_executed>[^"]*)" ` +
	`"(?P<redirect_url>[^"]*)"` +
	`(?: "(?P<error_reason>[^ ]*)")?` +
	`(?: "(?P<target_list>[^"]*)")?` +
	`(?: "(?P<target_status_code_list>[^"]*)")?` +
	`(?:.*?)$`, // debug with (?P<>...)
), nil}
