package provisioning

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var scenarios = []scenario{
	{
		name:                    "scenario 1 Ubuntu Docker 1.13",
		osName:                  "ubuntu",
		containerRuntime:        "docker",
		containerRuntimeVersion: "1.13",
	},

	{
		name:                    "scenario 2 Ubuntu Docker 17.03",
		osName:                  "ubuntu",
		containerRuntime:        "docker",
		containerRuntimeVersion: "17.03",
	},

	{
		name:                    "scenario 3 Ubuntu CRI-O 1.9",
		osName:                  "ubuntu",
		containerRuntime:        "cri-o",
		containerRuntimeVersion: "1.9",
	},

	{
		name:                    "scenario 4 CentOS Docker 1.12",
		osName:                  "centos",
		containerRuntime:        "docker",
		containerRuntimeVersion: "1.12",
	},

	{
		name:                    "scenario 5 CentOS Docker 1.13",
		osName:                  "centos",
		containerRuntime:        "docker",
		containerRuntimeVersion: "1.13",
	},

	{
		name:                    "scenario 6 CoreOS Docker 1.13",
		osName:                  "coreos",
		containerRuntime:        "docker",
		containerRuntimeVersion: "1.13",
	},

	{
		name:                    "scenario 7 CoreOS Docker 17.03",
		osName:                  "coreos",
		containerRuntime:        "docker",
		containerRuntimeVersion: "17.03",
	},
}

type scenario struct {
	// name holds short description of the test scenario, it is also used to create machines and nodes names
	// so please don't put "strange" characters there
	name                    string
	osName                  string
	containerRuntime        string
	containerRuntimeVersion string
}

type scenarioSelector struct {
	osName                  []string
	containerRuntime        []string
	containerRuntimeVersion []string
}

func doesSenarioSelectorMatch(selector *scenarioSelector, testCase scenario) bool {
	for _, selectorOSName := range selector.osName {
		if testCase.osName == selectorOSName {
			return true
		}
	}

	for _, selectorContainerRuntime := range selector.containerRuntime {
		if testCase.containerRuntime == selectorContainerRuntime {
			return true
		}
	}

	for _, selectorContainerRuntimeVersion := range selector.containerRuntimeVersion {
		if testCase.containerRuntime == selectorContainerRuntimeVersion {
			return true
		}
	}

	return false
}

func runScenarios(st *testing.T, excludeSelector *scenarioSelector, testParams []string, manifestPath string, cloudProvider string) {
	for _, testCase := range scenarios {
		if excludeSelector != nil && doesSenarioSelectorMatch(excludeSelector, testCase) {
			continue
		}

		st.Run(testCase.name, func(it *testing.T) {
			testScenario(it, testCase, cloudProvider, testParams, manifestPath)
		})
	}

}

func testScenario(t *testing.T, testCase scenario, cloudProvider string, testParams []string, manifestPath string) {
	t.Parallel()

	kubernetesCompliantName := fmt.Sprintf("%s-%s", testCase.name, cloudProvider)
	kubernetesCompliantName = strings.Replace(kubernetesCompliantName, " ", "-", -1)
	kubernetesCompliantName = strings.ToLower(kubernetesCompliantName)

	scenarioParams := append([]string(nil), testParams...)
	scenarioParams = append(scenarioParams, fmt.Sprintf("<< MACHINE_NAME >>=%s", kubernetesCompliantName))
	scenarioParams = append(scenarioParams, fmt.Sprintf("<< NODE_NAME >>=%s", kubernetesCompliantName))
	scenarioParams = append(scenarioParams, fmt.Sprintf("<< OS_NAME >>=%s", testCase.osName))
	scenarioParams = append(scenarioParams, fmt.Sprintf("<< CONTAINER_RUNTIME >>=%s", testCase.containerRuntime))
	scenarioParams = append(scenarioParams, fmt.Sprintf("<< CONTAINER_RUNTIME_VERSION >>=%s", testCase.containerRuntimeVersion))

	gopath := os.Getenv("GOPATH")
	projectDir := filepath.Join(gopath, "src/github.com/kubermatic/machine-controller")

	kubeConfig := filepath.Join(projectDir, ".kubeconfig")

	err := verify(kubeConfig, manifestPath, scenarioParams, 60*time.Hour)
	if err != nil {
		t.Errorf("verify failed due to error=%v", err)
	}
}