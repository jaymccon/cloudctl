package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configures cloud resource providers",
}

func init() {
	RootCmd.AddCommand(ConfigureCmd)
}

func Edit() {
	// TODO: work out filename, populate initial file, parse edited file and check if it changed
	fpath := os.TempDir() + "/thetemporaryfile.txt"
	f, err := os.Create(fpath)
	if err != nil {
		log.Printf("1")
		log.Fatal(err)
	}
	f.Close()

	editor := exec.Command("vim", fpath)
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout
	editor.Stderr = os.Stderr
	err = editor.Start()
	if err != nil {
		log.Printf("2")
		log.Fatal(err)
	}
	err = editor.Wait()
	if err != nil {
		log.Printf("Error while editing. Error: %v\n", err)
	} else {
		log.Printf("Successfully edited.")
	}
}
