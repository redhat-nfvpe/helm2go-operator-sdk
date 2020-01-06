package new

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func doGoScaffold() error {
	// for testing purposes
	if mock {
		return nil
	}
	// calls the operator-sdk with the correct command line arguments
	cmd := getScaffoldCommand()
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error While Building Scaffold")
		return err
	}

	// add the api
	cmd = getAPICommand()
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error While Adding API")
		return err
	}

	// add the controller
	cmd = getControllerCommand()
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error While Adding Controller")
		return err
	}
	// returns the output scaffold directory
	return nil
}

// getScaffoldCommand returns the correct command to build the operator-sdk scaffold
func getScaffoldCommand() *exec.Cmd {
	//scaffoldCommand := exec.Command("operator-sdk", "new", operatorName, "--dep-manager", "dep")
	scaffoldCommand := exec.Command("operator-sdk", "new", operatorName)
	scaffoldCommand.Dir = filepath.Dir(outputDir)
	scaffoldCommand.Stdout = os.Stdout
	scaffoldCommand.Stderr = os.Stderr
	return scaffoldCommand
}

// getAPICommand returns the correct command to add the api of the specified kind and version
func getAPICommand() *exec.Cmd {
	apiCommand := exec.Command("operator-sdk", "add", "api", "--kind", kind, "--api-version", apiVersion)
	apiCommand.Dir = outputDir
	apiCommand.Stdout = os.Stdout
	apiCommand.Stderr = os.Stderr
	return apiCommand
}

// getControllerCommand returns the correct command to add the controller of the specified kind and version
func getControllerCommand() *exec.Cmd {
	controllerCommand := exec.Command("operator-sdk", "add", "controller", "--kind", kind, "--api-version", apiVersion)
	controllerCommand.Dir = outputDir
	controllerCommand.Stdout = os.Stdout
	controllerCommand.Stderr = os.Stderr
	return controllerCommand
}
