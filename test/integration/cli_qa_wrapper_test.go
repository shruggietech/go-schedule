//go:build windows

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIQAWrapperUsesSelectedPolicyAndRetainsTimeoutOutput(t *testing.T) {
	requireTwin(t, powershellTwin())
	goTool, err := exec.LookPath("go")
	if err != nil {
		t.Skip("Go is required to build the owned wrapper fixture")
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "inputs")
	if err := os.MkdirAll(filepath.Join(bin, "windows"), 0700); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(dir, "main.go")
	exe := filepath.Join(bin, "windows", "gosched-qa.exe")
	program := `package main
import("encoding/json";"fmt";"os";"time")
func main(){
 f,_:=os.OpenFile(os.Getenv("QA_WRAPPER_LOG"),os.O_APPEND|os.O_WRONLY|os.O_CREATE,0600);json.NewEncoder(f).Encode(os.Args[1:]);f.Close()
 if os.Args[1]=="timing-policy"{json.NewEncoder(os.Stdout).Encode(map[string]int{"outer_timeout_ms":150});return}
 fmt.Println("retained command stdout");fmt.Fprintln(os.Stderr,"retained command stderr")
 if err:=os.WriteFile(os.Getenv("QA_WRAPPER_READY"),[]byte("ready"),0600);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
 time.Sleep(10*time.Second)
}
`
	if err := os.WriteFile(source, []byte(program), 0600); err != nil {
		t.Fatal(err)
	}
	if out, err := scriptCommand(goTool, "build", "-o", exe, source).CombinedOutput(); err != nil {
		t.Fatalf("build fixture %v\n%s", err, out)
	}
	wrapper, err := filepath.Abs(filepath.Join("..", "..", "scripts", "cli-qa.ps1"))
	if err != nil {
		t.Fatal(err)
	}
	// Synchronize preparation in a private copy. The production policy and its
	// 150 ms timeout remain unchanged; the child must emit both channels first.
	body, err := os.ReadFile(wrapper)
	if err != nil {
		t.Fatal(err)
	}
	waitBoundary := "$qaTimedOut = -not $qaProcess.WaitForExit($TimeoutMilliseconds)"
	if strings.Count(string(body), waitBoundary) != 1 {
		t.Fatal("wrapper wait boundary changed; review fixture synchronization")
	}
	readiness := `if ($Evaluation) {
        $fixturePreparation = [Diagnostics.Stopwatch]::StartNew()
        while (-not [IO.File]::Exists($env:QA_WRAPPER_READY) -and $fixturePreparation.ElapsedMilliseconds -lt 10000) {
            if ($qaProcess.HasExited) { throw 'Fixture child exited before output readiness.' }
            [Threading.Thread]::Sleep(10)
        }
        if (-not [IO.File]::Exists($env:QA_WRAPPER_READY)) {
            $qaProcess.Kill($true)
            $qaProcess.WaitForExit()
            throw 'Fixture child did not reach output readiness within its preparation bound.'
        }
    }
    ` + waitBoundary
	wrapper = filepath.Join(dir, "cli-qa-fixture.ps1")
	if err := os.WriteFile(wrapper, []byte(strings.Replace(string(body), waitBoundary, readiness, 1)), 0600); err != nil {
		t.Fatal(err)
	}
	sourceRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	budget := filepath.Join(dir, "budget.json")
	if err := os.WriteFile(budget, []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, action := range []string{"capture", "resume"} {
		t.Run(action, func(t *testing.T) {
			log := filepath.Join(dir, action+".jsonl")
			outDir := filepath.Join(dir, action)
			args := []string{"-NoProfile", "-NonInteractive", "-File", wrapper, "-Action", action, "-SourceDirectory", sourceRoot, "-BinaryDirectory", bin, "-OutputDirectory", outDir, "-WSLDistribution", "owned distro", "-WSLUser", "owned user"}
			if action == "capture" {
				args = append(args, "-Budget", budget)
			}
			cmd := scriptCommand("pwsh", args...)
			cmd.Env = append(os.Environ(), "QA_WRAPPER_LOG="+log, "QA_WRAPPER_READY="+filepath.Join(dir, action+"-ready"))
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("timed out command reported success")
			}
			for _, want := range []string{"retained command stdout", "retained command stderr", "original BinaryDirectory", "resume", "allowance (150 ms)"} {
				if !strings.Contains(string(out), want) {
					t.Fatalf("missing %q in timeout result:\n%s", want, out)
				}
			}
			data, err := os.ReadFile(log)
			if err != nil {
				t.Fatal(err)
			}
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			if len(lines) != 2 {
				t.Fatalf("calls: %s", data)
			}
			if !strings.Contains(lines[0], "timing-policy") {
				t.Fatalf("policy not read first: %s", data)
			}
			if action == "capture" && !strings.Contains(lines[0], "--budget") {
				t.Fatal("explicit budget omitted from policy")
			}
			if action == "resume" && !strings.Contains(lines[0], "--resume-out") {
				t.Fatal("original resume budget not selected")
			}
			t.Logf("%s original inputs and terminated output retained", action)
		})
	}
}
