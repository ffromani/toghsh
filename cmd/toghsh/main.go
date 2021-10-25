package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"

	"github.com/nektos/act/pkg/model"
)

func main() {
	var listJobs bool
	var jobID string
	flag.BoolVarP(&listJobs, "list", "L", false, "list available jobs and exit")
	flag.StringVarP(&jobID, "job-id", "J", "", "process job")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s path/to/workflow.yaml\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	if !listJobs && jobID == "" {
		fmt.Fprintf(os.Stderr, "missing job id\n")
		os.Exit(2)
	}

	src, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open %q: %v\n", os.Args[1], err)
		os.Exit(4)
	}
	defer src.Close()

	wf, err := model.ReadWorkflow(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read the workflow on %q: %v\n", os.Args[1], err)
		os.Exit(4)
	}

	if listJobs {
		for name := range wf.Jobs {
			fmt.Printf("job: %q\n", name)
		}
		os.Exit(0)
	}

	job, ok := wf.Jobs[jobID]
	if !ok {
		fmt.Fprintf(os.Stderr, "unable to find job %q in %q\n", jobID, os.Args[1])
		os.Exit(8)
	}

	fmt.Printf("### setup environment as per job %q\n", jobID)
	for key, value := range job.Environment() {
		fmt.Printf("export %s=%s\n", key, value)
	}
	fmt.Println()

	for idx, step := range job.Steps {
		fmt.Printf("### order=%d - ID=%q - name=%q\n", idx, step.ID, step.Name)
		if step.Run == "" {
			fmt.Println("### nothing to run\n")
			continue
		}
		fmt.Println(step.Run)
	}
}
