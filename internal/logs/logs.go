package logs

import (
	"regexp"
)

type Parser struct {
	*regexp.Regexp
}

func (p *Parser) Parse(line string) map[string]string {
	names := p.SubexpNames()[1:]
	result := make(map[string]string, len(names))

	values := p.FindStringSubmatch(line)
	for i, key := range names {
		result[key] = values[i+1]
	}

	return result
}

var ALB = &Parser{regexp.MustCompile(`^` +
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
	`(?:.*?)$`,
)}
