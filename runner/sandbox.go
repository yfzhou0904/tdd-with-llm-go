package runner

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
)

type Sandbox interface {
	ReadExistingCode() (testCode, implCode string, err error)
	WriteAndTest(testCode, implCode string, showDiff bool) (output string, err error)
}

type LocalSandbox struct {
	ImplCodePath string
	TestCodePath string
}

func (s LocalSandbox) ReadExistingCode() (testCode, implCode string, err error) {
	testCodeBytes, err := os.ReadFile(s.TestCodePath)
	if err != nil {
		return "", "", err
	}

	implCodeBytes, err := os.ReadFile(s.ImplCodePath)
	if err != nil {
		return "", "", err
	}

	return string(testCodeBytes), string(implCodeBytes), nil
}

func (s LocalSandbox) WriteAndTest(testCode, implCode string, showDiff bool) (output string, err error) {
	testFile, err := os.Create(s.TestCodePath + ".new")
	if err != nil {
		return "", err
	}
	defer testFile.Close()
	_, err = testFile.WriteString(testCode)
	if err != nil {
		return "", err
	}

	implFile, err := os.Create(s.ImplCodePath + ".new")
	if err != nil {
		return "", err
	}
	_, err = implFile.WriteString(implCode)
	if err != nil {
		return "", err
	}

	var cmd *exec.Cmd

	if showDiff {
		for _, f := range []string{s.ImplCodePath, s.TestCodePath} {
			slog.Info(fmt.Sprintf("diffing %s", f))
			cmd = exec.Command("diff", "-y", f, f+".new")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	}

	for _, f := range []string{s.ImplCodePath, s.TestCodePath} {
		if err = exec.Command("mv", f+".new", f).Run(); err != nil {
			slog.Error("failed to mv", "err", err, "src", f+".new", "dst", f)
		}
	}

	startDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	os.Chdir(path.Dir(s.ImplCodePath))
	defer os.Chdir(startDir)

	outputBytes, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return string(outputBytes), err
	}
	exec.Command("gofmt", "-w", ".").Run()
	exec.Command("goimports", "-w", ".").Run()
	outputBytes, err = exec.Command("go", "test", ".", "-v").CombinedOutput()
	if err != nil {
		return string(outputBytes), err
	}
	return string(outputBytes), nil
}
