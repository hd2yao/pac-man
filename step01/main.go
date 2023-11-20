package main

import (
    "bufio"
    "fmt"
    "os"
)

var maze []string

func loadMaze(file string) error {
    f, err := os.Open(file)
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

func main() {

}
