package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	Memory "github.com/coderick14/ARMed/Memory"
	"io"
	"os"
	"strings"
)

var helpString = `ARMed version 1.0
Author : https://github.com/coderick14

ARMed is a very basic emulator of the 64-bit LEGv8 architecture written in Golang
USAGE : ARMed [--all] SOURCE_FILE

The --all flag will show all register values after an instruction, with updated ones in color.
In absence of this flag, it will show only updated registers.

Found a bug? Feel free to raise an issue on https://github.com/coderick14/ARMed
Contributions welcome :)`

func main() {
	var err error
	helpPtr := flag.Bool("help", false, "Display help")
	allPtr := flag.Bool("all", false, "Display all registers")

	flag.Parse()

	if *helpPtr == true {
		fmt.Println(helpString)
		return
	}

	if len(flag.Args()) == 0 {
		err = errors.New("Error : Missing filename.\n Type ARMed --help for further help")
		fmt.Println(err)
		return
	}

	fileName := flag.Args()[0]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file : ", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString(';')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error while reading file : ", err)
			return
		}
		line = strings.TrimSpace(strings.TrimRight(line, ";"))
		Memory.InstructionMem.Instructions = append(Memory.InstructionMem.Instructions, line)
	}

	Memory.InitRegisters()

	for _, _ = range Memory.InstructionMem.Instructions {
		Memory.SaveRegisters()
		err = Memory.InstructionMem.ValidateAndExecuteInstruction()
		if err != nil {
			fmt.Println(err)
			return
		}
		Memory.ShowRegisters(*allPtr)
	}

}
