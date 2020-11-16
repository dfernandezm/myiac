package cluster

import (
	"fmt"
	"github.com/dfernandezm/myiac/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFileSecret(t *testing.T) {
	// setup
	cmdLine := testutil.FakeKubernetesRunner("secret default/test-secret created")
	kubernetesRunner := NewKubernetesRunner(cmdLine)

	// given
	filePath := "/tmp/filepath"
	secretName := "test-Secret-Name"

	// when
	kubernetesRunner.CreateFileSecret(secretName, "default", filePath)

	// then
	expectedCreateSecretCmdLine := fmt.Sprintf("kubectl create secret generic %s " +
			"--from-file=%s.json=%s -n default", secretName, secretName, filePath)
	actualCreateSecretCmdLine := cmdLine.CmdLines[0]

	assert.Equal(t, expectedCreateSecretCmdLine, actualCreateSecretCmdLine)
}