package secret

import (
	"fmt"
	"github.com/dfernandezm/myiac/internal/cluster"
	"github.com/dfernandezm/myiac/testutil"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateSecret(t *testing.T) {
	// setup
	cmdLine := testutil.FakeKubernetesRunner("test-output")
	kubernetesRunner := cluster.NewKubernetesRunner(cmdLine)
	secretManager := NewKubernetesSecretManager("default", kubernetesRunner)
	fmt.Printf(secretManager.namespace)
	// given
	filePath := "/tmp/filepath"
	secretName := "test-secret-name"
	os.Create(filePath)

	// when
	secretManager.CreateFileSecret(secretName, filePath)

	// then
	//TODO: should validate snake case
	expectedCreateSecretCmdLine :=
		fmt.Sprintf("kubectl create secret generic %s " +
			"--from-file=%s.json=%s -n default", secretName, secretName, filePath)
	actualCreateSecretCmdLine := cmdLine.CmdLines[0]

	assert.Equal(t, expectedCreateSecretCmdLine, actualCreateSecretCmdLine)
}