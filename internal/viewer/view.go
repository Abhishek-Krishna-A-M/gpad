package viewer

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Abhishek-Krishna-A-M/gpad/internal/editor"
)

func View(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	text := string(data)

	if rendered, err := RenderCustom(text); err == nil {
		fmt.Println(rendered)
		return promptEdit(path)
	}

	fmt.Println(text)
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

