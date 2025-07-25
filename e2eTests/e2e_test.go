package e2eTests

import (
	"agl/pkg/agl"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	tassert "github.com/stretchr/testify/assert"
)

func spawnGoRunFromBytes(source []byte, programArgs []string) ([]byte, error) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "gorun")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir) // clean up

	// Create a .go file inside it
	tmpFile := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(tmpFile, source, 0644)
	if err != nil {
		return nil, err
	}

	coreFile := filepath.Join(tmpDir, "aglCore.go")
	err = os.WriteFile(coreFile, []byte(agl.GenCore()), 0644)
	if err != nil {
		return nil, err
	}

	// Run `go run` on the file with additional arguments
	cmdArgs := append([]string{"run", tmpFile, coreFile}, programArgs...)
	cmd := exec.Command("go", cmdArgs...)
	return cmd.Output()
}

func testGenOutput(src string) string {
	fset, f := agl.ParseSrc(src)
	env := agl.NewEnv()
	agl.NewInferrer(env).InferFile("", f, fset, true)
	outSrc := agl.NewGenerator(env, f, fset).Generate()
	out := agl.Must(spawnGoRunFromBytes([]byte(outSrc), nil))
	return string(out)
}

func Test1(t *testing.T) {
	src := `package main
import "fmt"
func main() {
	a := 1
	fmt.Println(a)
}`
	tassert.Equal(t, "1\n", testGenOutput(src))
}
