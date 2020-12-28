package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type option struct {
	Option string
}

type vm struct {
	Name string
}

var vmDir string
var dir []os.FileInfo

func main() {
	checks()
	checkVirtualMachineDirectory()
	selection := selectOption()

	if selection == "Start a virtual machine" {
		startVirtualMachine()
	} else if selection == "Stop a virtual machine" {
		fmt.Println("TODO")
	} else if selection == "List all running virtual machines" {

	}
}

func checks() {
	// Checking if VMware is installed
	if _, err := os.Stat("/Applications/VMware Fusion.app"); os.IsNotExist(err) {
		log.Fatal(err)
	}

	// Checking if vmrun is available
	cmd := exec.Command("vmrun")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func checkVirtualMachineDirectory() {
	if value, ok := os.LookupEnv("VIRTUAL_MACHINES_DIR"); ok {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		vmDir = homeDir + value
	
		if _, err := os.Stat(vmDir); os.IsNotExist(err) {
			log.Fatal(err)
		}
	
		dir, err = ioutil.ReadDir(vmDir)
		if err != nil {
			log.Fatal(err)
		}
		
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		vmDir = homeDir + "/Virtual Machines.localized"
	
		if _, err := os.Stat(vmDir); os.IsNotExist(err) {
			log.Fatal(err)
		}
	
		dir, err = ioutil.ReadDir(vmDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func selectOption() string {
	options := []option{
		{Option: "Start a virtual machine"},
		{Option: "Stop a virtual machine"},
		{Option: "List all running virtual machines"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F4BE{{ .Option | cyan }}",
		Inactive: "  {{ .Option | cyan }}",
		Selected: "\U0001F4BE {{ .Option | white | cyan}}",
	}

	prompt := promptui.Select{
		Label:     ">>",
		Items:     options,
		Templates: templates,
		Size:      3,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Promt failed %v\n", err)
	}

	return options[i].Option
}

func startVirtualMachine() {
	if len(dir) < 1 {
		fmt.Println("No existing VMs")
		return
	}

	for i, f := range dir {
		if filepath.Ext(f.Name()) == ".vmwarevm" {
			fmt.Println(strconv.Itoa(i-1) + " " + f.Name())
		}
	}

	var input int
	fmt.Scan(&input)

	for i, f := range dir {
		if i-1 == input {
			dir2, err := ioutil.ReadDir(vmDir + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}

			for _, d := range dir2 {
				if filepath.Ext(d.Name()) == ".vmx" {
					fullPath := vmDir + "/" + f.Name()
					vmxImage := d.Name()

					vmrunPath, _ := exec.LookPath("vmrun")
					cmdRun := &exec.Cmd{
						Path:   vmrunPath,
						Args:   []string{vmrunPath, "-T", "fusion", "start", vmxImage, "nogui"},
						Dir:    fullPath,
						Stdout: os.Stdout,
						Stderr: os.Stderr,
					}

					if err := cmdRun.Run(); err != nil {
						fmt.Println("Error:", err)
					}
				}
			}
		}
	}
}

func listRunningVMs() {
	vmrunPath, _ := exec.LookPath("vmrun")
	cmdRun := &exec.Cmd{
		Path:   vmrunPath,
		Args:   []string{vmrunPath, "list"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := cmdRun.Run(); err != nil {
		fmt.Println("Error:", err)
	}

	stdout, err := cmdRun.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("output: ", stdout)
}
