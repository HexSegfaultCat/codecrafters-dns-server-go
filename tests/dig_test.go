package dnstest

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"tests/common"
)

func TestQueryWithDig(t *testing.T) {
	port := common.InitializeDnsServer()

	expectedUrl := "codecrafters.io."
	expectedIp := "8.8.8.8"

	command := exec.Command(
		"dig",
		"@127.0.0.1",
		"+answer",
		"-p", strconv.Itoa(int(port)),
		expectedUrl,
	)

	stdout, err := command.Output()
	if err != nil {
		t.Error(err)
	}

	answerParts := []string{}
	outputLines := strings.Split(string(stdout), "\n")
	for _, line := range outputLines {
		if line != "" && !strings.HasPrefix(line, ";") {
			answerParts = strings.Split(line, "\t")
			break
		}
	}

	if answerParts[0] != expectedUrl {
		t.Errorf("Expected URL to be %s, but got %s", answerParts[0], expectedUrl)
	}
	if answerParts[4] != expectedIp {
		t.Errorf("Expected IP to be %s, but got %s", answerParts[4], expectedIp)
	}
}
