package cli

import (
	"fmt"

	"github.com/bitrise-io/go-pathutil"
	"github.com/bitrise-io/stepman/git"
	"github.com/bitrise-io/stepman/paths"
	"github.com/codegangsta/cli"
)

func update(c *cli.Context) {
	fmt.Println("Update")

	stepsSpecDir := pathutil.UserHomeDir() + paths.STEP_COLLECTION_DIR
	exists, err := pathutil.IsPathExists(stepsSpecDir)
	if err != nil {
		fmt.Println("Failed to update Stepman:", err)
		return
	}
	if exists == false {
		fmt.Println("Stepman is not initialized")
		return
	}

	err = git.DoGitPull(stepsSpecDir)
	if err != nil {
		fmt.Println("Failed tp update Stepman:", err)
		return
	}
	fmt.Println("Stepman updated")
}
