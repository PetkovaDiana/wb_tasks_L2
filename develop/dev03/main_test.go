package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// args - for testing input
type args struct {
	files []string
	flags []string
}

func TestWithoutFlags(t *testing.T) {
	cases := []args{
		{
			files: []string{"test1"},
		},
		{
			files: []string{"test2"},
		},
		{
			files: []string{"test3"},
		},
		{
			files: []string{"test4"},
		},
		{
			files: []string{"test1", "test2"},
		},
		{
			files: []string{"test1", "test3"},
		},
		{
			files: []string{"test2", "test3"},
		},
		{
			files: []string{"test1", "test2", "test3", "test4", "test5"},
		},
	}

	for _, testCase := range cases {
		command := []string{"run", "task.go"}
		command = append(command, testCase.files...)

		funcOut, err := exec.Command("go", command...).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(2)
		}

		rOut, err := exec.Command("sort", testCase.files...).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(2)
		}
		realOut := strings.ReplaceAll(string(rOut), "\r", "")

		for i, val := range realOut {
			if string(val) != string(funcOut[i]) {
				t.Errorf("%s != %s\nShould: %s\nGot: %s", string(val), string(funcOut[i]), string(val), string(funcOut[i]))
				os.Exit(1)
			}
		}
	}
}

func TestWithFlags(t *testing.T) {
	cases := []args{
		{
			files: []string{"test1"},
			flags: []string{"-k", "5", "-n", "-r"},
		},
		{
			files: []string{"test2"},
			flags: []string{"-k", "3", "-n"},
		},
		{
			files: []string{"test3"},
			flags: []string{"-k", "2", "-r"},
		},
		{
			files: []string{"test4"},
			flags: []string{"-u", "-r"},
		},
		{
			files: []string{"test5"},
			flags: []string{"-M", "-k", "4"},
		},
		{
			files: []string{"test1", "test2"},
			flags: []string{"-k", "9", "-r"},
		},
		{
			files: []string{"test1", "test3"},
			flags: []string{"-u", "-r"},
		},
		{
			files: []string{"test2", "test3"},
			flags: []string{"-k", "2"},
		},
		{
			files: []string{"test1", "test2", "test3"},
			flags: []string{"-n", "-r"},
		},
	}

	for _, testCase := range cases {
		command := []string{"run", "task.go"}
		command = append(command, testCase.flags...)
		command = append(command, testCase.files...)

		funcOut, err := exec.Command("go", command...).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v", err)
			os.Exit(1)
		}

		command = append(testCase.flags, testCase.files...)

		rOut, err := exec.Command("sort", command...).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v", err)
			os.Exit(1)
		}
		realOut := strings.ReplaceAll(string(rOut), "\r", "")

		for i, val := range realOut {
			if string(val) != string(funcOut[i]) {
				t.Errorf("%s != %s\nShould: %s\nGot: %s\nFlags: %s", string(val), string(funcOut[i]), string(val), string(funcOut[i]), testCase.flags)
				os.Exit(1)
			}
		}
	}
}
