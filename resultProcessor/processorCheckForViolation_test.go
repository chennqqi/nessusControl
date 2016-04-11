package nessusProcessor

import (
	"testing"
)

var (
	nessusResultRow = &Nessus6ResultRow{
		PluginID:       71783,
		CVE:            "CVE-2013-5211",
		CVSS:           5,
		Risk:           "Medium",
		Host:           "localhost",
		Protocol:       "udp",
		Port:           123,
		Name:           "NTP monlist Command Enabled",
		Synopsis:       "The remote network time service could be used for network reconnaissance or abused in a distributed denial of service attack.",
		Description:    "The version of ntpd on the remote host has the 'monlist' command enabled. This command returns a list of recent hosts that have connected to the service. As such, it can be used for network reconnaissance or, along with a spoofed source IP, a distributed denial of service attack.",
		Solution:       "If using NTP from the Network Time Protocol Project, either upgrade to\nNTP 4.2.7-p26 or later, or add 'disable monitor' to the 'ntp.conf'\nconfiguration file and restart the service. Otherwise, contact the\nvendor.\n\nOtherwise, limit access to the affected service to trusted hosts.",
		SeeAlso:        "https://isc.sans.edu/diary/NTP+reflection+attack/17300\nhttp://bugs.ntp.org/show_bug.cgi?id=1532\nhttp://kb.juniper.net/InfoCenter/index?page=content&id=JSA10613",
		PluginOutput:   "",
		OrganizationID: 1,
		RegionID:       2,
	}

	anyMatchCriteria = &MatchCriteria{
		PluginID:                     71783,
		ExternallyAccessible:         true,
		Port:                         []int{514, 123},
		DescriptionRegexp:            []string{"(The|This|That) version of ntpd", "This better not match"},
		OrganizationID:               []int{5, 3, 1},
		RegionID:                     []int{6, 4, 2},
		IgnoreViolationsWithCriteria: false,
		CountIf: "any",
	}

	allMatchCriteria = &MatchCriteria{
		PluginID:                     71783,
		ExternallyAccessible:         true,
		Port:                         []int{514, 123},
		DescriptionRegexp:            []string{"(The|This|That) version of ntpd", "This better not match"},
		OrganizationID:               []int{5, 3, 1},
		RegionID:                     []int{6, 4, 2},
		IgnoreViolationsWithCriteria: false,
		CountIf: "all",
	}

	ignoredAnyMatchCriteria = &MatchCriteria{
		PluginID:                     71783,
		ExternallyAccessible:         true,
		Port:                         []int{514, 123},
		DescriptionRegexp:            []string{"(The|This|That) version of ntpd", "This better not match"},
		OrganizationID:               []int{5, 3, 1},
		RegionID:                     []int{6, 4, 2},
		IgnoreViolationsWithCriteria: true,
		CountIf: "any",
	}
)

func TestPositiveCheckForViolation(t *testing.T) {
	violation := anyMatchCriteria.CheckForViolation(nessusResultRow)
	if !violation {
		t.FailNow()
	}
}
