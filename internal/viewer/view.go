package viewer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
)

func View(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    rendered, _ := RenderCustom(string(data))

    // Use 'less' as a pager so users can scroll long notes
    cmd := exec.Command("less", "-R") // -R keeps the ANSI colors
    cmd.Stdin = strings.NewReader(rendered)
    cmd.Stdout = os.Stdout
    err = cmd.Run()
    
    // After they quit 'less', ask if they want to edit
    return promptEdit(path)
}

func promptEdit(path string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEdit this file? (y/N): ")

	resp, _ := reader.ReadString('\n')
	resp = strings.TrimSpace(strings.ToLower(resp))

	if resp == "y" || resp == "yes" {
		return editor.Open(path)
	}

	return nil
}

func ViewRaw(text string) {
	if rendered, err := RenderCustom(text); err == nil {
		fmt.Println(rendered)
		return
	}
	fmt.Println(text)
}

