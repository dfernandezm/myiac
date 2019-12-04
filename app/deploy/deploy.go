package deploy

import (
	"fmt"
	"strings"
	"github.com/dfernandezm/myiac/app/cluster"
	"github.com/dfernandezm/myiac/app/commandline"
	"github.com/dfernandezm/myiac/app/util"
)

type Deployment struct {
	AppName          string
	ChartPath        string
	DryRun           bool
	HelmSetParams    map[string]string // key value pairs
	HelmReleaseName  string
	HelmValuesParams []string // yaml filenames to pass as --values
}

func DeployApp(appName string, environment string) {
	//TODO: generify
	if appName == "all" {
		deployApps(environment)
	}

	if appName == "moneycolserver" {
		deployMoneyColServer()
	}

	if appName == "moneycolfrontend" {
		deployMoneyColFrontend()
	}

	if appName == "elasticsearch" {
		deployElasticsearch()
	}

	if appName == "traefik" {
		deployTraefik(environment)
	}
}

func deployApps(environment string) {
	deployElasticsearch()
	deployMoneyColServer()
	deployMoneyColFrontend()
	deployTraefik(environment)
}

func deployMoneyColFrontend() {
	//TODO: read releases first, get parameters for chart paths
	releaseName := "esteemed-peacock"
	moneycolPath := "/development/repos/moneycol/"
	deployPath := util.GetHomeDir() + moneycolPath + "frontend/deploy"
	appName := "moneycolfrontend"
	chartPath := fmt.Sprintf("%s/%s/chart", deployPath, appName)
	moneyColFrontendDeploy := Deployment{AppName: appName, ChartPath: chartPath, 
								DryRun: false, HelmReleaseName: releaseName}
	deployApp(&moneyColFrontendDeploy)
}

func deployElasticsearch() {
	releaseName := ""
	moneycolPath := "/development/repos/moneycol/"
	deployPath := util.GetHomeDir() + moneycolPath + "server/deploy"
	appName := "elasticsearch"
	chartPath := fmt.Sprintf("%s/%s/chart", deployPath, appName)
	elasticsearchDeploy := Deployment{AppName: appName, ChartPath: chartPath, 
										DryRun: false, HelmReleaseName: releaseName}
	deployApp(&elasticsearchDeploy)
}

func deployMoneyColServer() {
	releaseName := "ponderous-lion"
	moneycolPath := "/development/repos/moneycol/"
	deployPath := util.GetHomeDir() + moneycolPath + "server/deploy"
	appName := "moneycolserver"
	chartPath := fmt.Sprintf("%s/%s/chart", deployPath, appName)
	moneyColServerDeploy := Deployment{AppName: appName, ChartPath: chartPath, 
										DryRun: false, HelmReleaseName: releaseName}
	deployApp(&moneyColServerDeploy)
}

func deployTraefik(environment string) {
	releaseName := "opining-frog"
	moneycolPath := "/development/repos/moneycol/"
	deployPath := util.GetHomeDir() + moneycolPath + "server/deploy"
	appName := "traefik"
	chartPath := fmt.Sprintf("%s/%s/chart", deployPath, appName)

	//TODO: Set paramaters, separate this
	helmSetParams := make(map[string]string)
	internalIps := cluster.GetInternalIpsForNodes()

	// very flaky --set for ips like this: --set externalIps={ip1\,ip2\,ip3}
	internalIpsForHelmSet := "{" + strings.Join(internalIps, "\\,") + "}"
	helmSetParams["externalIps"] = internalIpsForHelmSet
	deployment := Deployment{AppName: appName, ChartPath: chartPath,
		DryRun:          false,
		HelmReleaseName: releaseName,
		HelmSetParams:   helmSetParams}

	deployApp(&deployment)
	
	if (environment == "dev") {
		// once deployed, repoint dev DNS to any public IP of nodes
		changeDevDns(deployPath)
	}
}

func changeDevDns(deployPath string) {
	publicIps := cluster.GetAllPublicIps()
	aPublicIp := publicIps[0] // any public ip works for this
	devDnsTfFile := deployPath + "/terraform/dns"
	cluster.ApplyDnsIpChange(devDnsTfFile, aPublicIp)
}

func deployApp(deployment *Deployment) {

	var argsTpl = ""
	if deployment.HelmReleaseName == "" {
		argsTpl = "install %s"
	} else {
		argsTpl = "upgrade " + deployment.HelmReleaseName + " %s"
	}

	argsStr := fmt.Sprintf(argsTpl, deployment.ChartPath)

	if len(deployment.HelmValuesParams) > 0 {
		valuesParams := ""
		for _, filePath := range deployment.HelmValuesParams {
			valuesParams += valuesParams + "--values " + filePath + " "
		}
		valuesParams = strings.TrimSpace(valuesParams)
		argsStr = fmt.Sprintf("%s %s", argsStr, valuesParams)
	}

	if len(deployment.HelmSetParams) > 0 {
		setParams := ""
		for k, v := range deployment.HelmSetParams {
			setParams += setParams + "--set " + k + "=" + v + " "
		}
		setParams = strings.TrimSpace(setParams)
		argsStr = fmt.Sprintf("%s %s", argsStr, setParams)
	}

	if deployment.DryRun {
		argsStr += " --debug --dry-run"
	}

	argsArray := strings.Fields(argsStr)
	cmd := commandline.New("helm", argsArray)
	cmd.Run()
	fmt.Printf("Finished deploying app: %s\n\n", deployment.AppName)
}