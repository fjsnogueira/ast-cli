package main

import (
	"log"
	"os"

	"github.com/checkmarx/ast-cli/internal/commands"
	"github.com/checkmarx/ast-cli/internal/params"
	"github.com/checkmarx/ast-cli/internal/wrappers"
	"github.com/checkmarx/ast-cli/internal/wrappers/configuration"
	"github.com/spf13/viper"
)

const (
	successfulExitCode = 0
	failureExitCode    = 1
)

func main() {
	bindKeysToEnvAndDefault()
	configuration.LoadConfiguration()
	scans := viper.GetString(params.ScansPathKey)
	groups := viper.GetString(params.GroupsPathKey)
	logs := viper.GetString(params.LogsPathKey)
	projects := viper.GetString(params.ProjectsPathKey)
	results := viper.GetString(params.ResultsPathKey)
	uploads := viper.GetString(params.UploadsPathKey)

	scansWrapper := wrappers.NewHTTPScansWrapper(scans)
	groupsWrapper := wrappers.NewHTTPGroupsWrapper(groups)
	logsWrapper := wrappers.NewLogsWrapper(logs)
	uploadsWrapper := wrappers.NewUploadsHTTPWrapper(uploads)
	projectsWrapper := wrappers.NewHTTPProjectsWrapper(projects)
	resultsWrapper := wrappers.NewHTTPResultsWrapper(results)
	authWrapper := wrappers.NewAuthHTTPWrapper()
	resultsPredicatesWrapper := wrappers.NewResultsPredicatesHTTPWrapper()

	astCli := commands.NewAstCLI(
		scansWrapper,
		resultsPredicatesWrapper,
		uploadsWrapper,
		projectsWrapper,
		resultsWrapper,
		authWrapper,
		logsWrapper,
		groupsWrapper,
	)
	err := astCli.Execute()
	exitIfError(err)
	os.Exit(successfulExitCode)
}

func exitIfError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(failureExitCode)
	}
}

func bindKeysToEnvAndDefault() {
	for _, b := range params.EnvVarsBinds {
		err := viper.BindEnv(b.Key, b.Env)
		if err != nil {
			exitIfError(err)
		}
		viper.SetDefault(b.Key, b.Default)
	}
}
