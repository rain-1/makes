package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"path/filepath"
)

type buildrule struct {
	// user input
	output string
	inputs []string
	command string
	
	// calculated attributes
	built bool
	lock *sync.Mutex
}

var parents *Set
var sources *Set

var jobs map[string]buildrule

func Map(vs []string, f func(string) string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}

func parseLine(line string) buildrule {
	// the format is:
	//
	// <output>\t[<input-1>\t...]--<command>
	
	parts := strings.SplitN(line, "\t--\t", 2)
	filenames := strings.Split(parts[0], "\t")
	
	return buildrule{
		output: filepath.Clean(filenames[0]),
		inputs: Map(filenames[1:], filepath.Clean),
		command: parts[1],
		built: false,
		lock: &sync.Mutex{},
	}
}

func parseStdin() buildrule {
	// TODO: ensure input is not ""
	
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		rule := parseLine(scanner.Text())
		
		parents.Add(rule.output)
		for _, input := range rule.inputs {
			parents.Remove(input)
			sources.Add(input)
		}
		sources.Remove(rule.output)
		
		jobs[rule.output] = rule
	}
	
	parents.CleanUp()
	sources.CleanUp()
	
	jobs[""] = buildrule{
		output: "",
		inputs: parents.SetToSlice(),
		command: "echo \"Build successful!\"",
		built: false,
		lock: &sync.Mutex{},
	}
	
	return jobs[""]
}

var cancel chan bool
var cancelled bool

func performBuild(rule buildrule) {
	var err error
	
	var wg sync.WaitGroup
	
	var done = make(chan bool)
	defer close(done)
	
	// open a lock on building this item
	rule.lock.Lock()
	defer rule.lock.Unlock()
	
	if(!rule.built && !cancelled) {
		// Build all the children in a pool of parallel goroutines
		wg.Add(len(rule.inputs))
		go func() {
			wg.Wait()
			done <- true
		}()
		for _, job := range rule.inputs {
			j, present := jobs[job]
			go func(job string, j buildrule, present bool) {
				if(present) {
					performBuild(j)
				} else {
					if(!sources.Contains(job)) {
						fmt.Println("WAS NOT PRESENT", j, job)
					}
				}
				wg.Done()
			}(job, j, present)
		}
		
		// wait for it to finish
		select {
		case <-done:
			break
		case <-cancel:
			fmt.Println("CANCELLING")
			return
		}
		
		// Determine if we should run the command
		var output_stat, input_stat os.FileInfo
		var build = false
		output_stat, err = os.Stat(rule.output)
		if(err != nil) {
			build = true
		} else {
			for _, input := range rule.inputs {
				input_stat, err = os.Stat(input)
				
				// TODO check error
				if(err != nil) {
					fmt.Println("COULD NOT DO", input, "MISSING")
					cancelled = true
					for { cancel <- true }
				}
				
				if(input_stat.ModTime().After(output_stat.ModTime())) {
					build = true
					break
				}
			}
		}
		
		// Run the command
		if(build) {
			cmd := exec.Command("sh", "-c", rule.command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			
			// version that prints the inputs, for debugging
			//fmt.Println(rule.inputs, "\n", "#", rule.command)
			fmt.Println("#", rule.command)
			err := cmd.Run()
			
			if(err != nil) {
				fmt.Println("COULD NOT DO", rule.command, err)
				cancelled = true
				for { cancel <- true }
			}
		}
	}
}

func main() {
	jobs = make(map[string]buildrule)
	parents = NewSet()
	sources = NewSet()
	
	cancel = make(chan bool)
	cancelled = false
	
	top := parseStdin()

	//fmt.Println("sources", sources.SetToSlice())

	performBuild(top)
	
	if(cancelled) {
		os.Exit(1)
	}
}
