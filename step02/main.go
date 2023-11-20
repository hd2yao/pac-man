package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "os/exec"
)

var maze []string

func loadMaze(file string) error {
    f, err := os.Open("step01/" + file)
    if err != nil {
        return err
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        maze = append(maze, line)
    }

    return nil
}

func printScreen() {
    for _, line := range maze {
        fmt.Println(line)
    }
}

func initialise() {
    cbTerm := exec.Command("stty", "cbreak", "-echo")
    cbTerm.Stdin = os.Stdin

    err := cbTerm.Run()
    if err != nil {
        log.Fatalln("unable to activate cbreak mode:", err)
    }
}

func cleanup() {
    cookedTerm := exec.Command("stty", "-cbreak", "echo")
    cookedTerm.Stdin = os.Stdin

    err := cookedTerm.Run()
    if err != nil {
        log.Fatalln("unable to restore cooked mode:", err)
    }
}

func main() {
    // initialize game
    initialise()
    defer cleanup()

    // load resources
    err := loadMaze("maze01.txt")
    if err != nil {
        log.Println("failed to load maze:", err)
        return
    }

    // game loop
    for {
        // update screen
        printScreen()

        // process input

        // process movement

        // process collisions

        // check game over

        // Temp: break infinite loop
        break

        // repeat
    }
}
