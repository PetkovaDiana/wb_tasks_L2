package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

type args struct {
	files   []string
	flags   []string
	pattern string
}

func TestGrepWithoutFlags(t *testing.T) {
	cases := []args{
		{
			files:   []string{"test"},
			pattern: "lines",
		},
		{
			files:   []string{"test2"},
			pattern: "lines",
		},
		{
			files:   []string{"test", "test2"},
			pattern: "lines",
		},
	}

	for _, testCase := range cases {
		command := []string{"run", "task.go"}
		command = append(command, testCase.pattern)
		command = append(command, testCase.files...)

		myOut, err := exec.Command("go", command...).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		command = command[0:0]
		command = append(command, testCase.pattern)
		command = append(command, testCase.files...)

		realOut, err := exec.Command("grep", command...).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		for i := range realOut {
			if myOut[i] != realOut[i] {
				fmt.Println("MyOut:\n", string(myOut))
				fmt.Println("RealOut:\n", string(realOut))

				t.Errorf("%s != %s\nShould: %s\nGot: %s\n", string(realOut[i]), string(myOut[i]), string(realOut[i]), string(myOut[i]))
				os.Exit(1)
			}
		}
	}
}

func TestGrepWithRegex(t *testing.T) {
	t.Run("^lines$ with test2 file", func(t *testing.T) {
		myOut, err := exec.Command("go", "run", "task.go", `^lines$`, "test2").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		realOut, err := exec.Command("grep", `^lines$`, "test2").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		for i := range realOut {
			if myOut[i] != realOut[i] {
				fmt.Println("MyOut:\n", string(myOut))
				fmt.Println("RealOut:\n", string(realOut))

				t.Errorf("%s != %s\nShould: %s\nGot: %s\n", string(realOut[i]), string(myOut[i]), string(realOut[i]), string(myOut[i]))
				os.Exit(1)
			}
		}
	})

	t.Run("lines starts with - in file test", func(t *testing.T) {
		myOut, err := exec.Command("go", "run", "task.go", `^-`, "test").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		realOut, err := exec.Command("grep", `^-`, "test").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		for i := range realOut {
			if myOut[i] != realOut[i] {
				fmt.Println("MyOut:\n", string(myOut))
				fmt.Println("RealOut:\n", string(realOut))

				t.Errorf("%s != %s\nShould: %s\nGot: %s\n", string(realOut[i]), string(myOut[i]), string(realOut[i]), string(myOut[i]))
				os.Exit(1)
			}
		}
	})

	t.Run("with flag -F and ^ sign in the test file", func(t *testing.T) {
		myOut, err := exec.Command("go", "run", "task.go", "-F", `^`, "test").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		realOut, err := exec.Command("grep", "-F", `^`, "test").CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		for i := range realOut {
			if myOut[i] != realOut[i] {
				fmt.Println("MyOut:\n", string(myOut))
				fmt.Println("RealOut:\n", string(realOut))

				t.Errorf("%s != %s\nShould: %s\nGot: %s\n", string(realOut[i]), string(myOut[i]), string(realOut[i]), string(myOut[i]))
				os.Exit(1)
			}
		}
	})
}

func TestGrepWithFlags(t *testing.T) {
	cases := []args{
		{
			files:   []string{"test"},
			pattern: "lines",
			flags:   []string{"-i", "-n"},
		},
		{
			files:   []string{"test"},
			pattern: "lines",
			flags:   []string{"-v", "-c"},
		},
		{
			files:   []string{"test2"},
			pattern: "lines",
			flags:   []string{"-v", "-n", "-A", "10"},
		},
		{
			files:   []string{"test2"},
			pattern: "lines",
			flags:   []string{"-v", "-n", "-A", "1"},
		},
		{
			files:   []string{"test", "test2"},
			pattern: "lines",
			flags:   []string{"-C", "10", "-v"},
		},
	}

	for _, testCase := range cases {
		command := []string{"run", "task.go"}
		command = append(command, testCase.pattern)
		command = append(command, testCase.files...)

		myOut, err := exec.Command("go", command...).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		command = command[0:0]
		command = append(command, testCase.pattern)
		command = append(command, testCase.files...)

		realOut, err := exec.Command("grep", command...).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Starting test failed: %v\n", err)
			os.Exit(1)
		}

		for i := range realOut {
			if myOut[i] != realOut[i] {
				fmt.Println("MyOut:\n", string(myOut))
				fmt.Println("RealOut:\n", string(realOut))

				t.Errorf("%s != %s\nShould: %s\nGot: %s\n", string(realOut[i]), string(myOut[i]), string(realOut[i]), string(myOut[i]))
				os.Exit(1)
			}
		}
	}
}
