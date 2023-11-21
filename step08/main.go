package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "math/rand"
    "os"
    "os/exec"
    "time"

    "github.com/danicat/simpleansi"
)

var (
    configFile = flag.String("config-file", "config.json", "path to custom configuration file")
    mazeFile   = flag.String("maze-flag", "maze01.txt", "path to custom maze file")
)

// Config holds the emoji configuration
type Config struct {
    Player   string `json:"player"`
    Ghost    string `json:"ghost"`
    Wall     string `json:"wall"`
    Dot      string `json:"dot"`
    Pill     string `json:"pill"`
    Death    string `json:"death"`
    Space    string `json:"space"`
    UseEmoji bool   `json:"use_emoji"`
}

var cfg Config

func loadConfig(file string) error {
    f, err := os.Open(file)
    if err != nil {
        return err
    }
    defer f.Close()

    decoder := json.NewDecoder(f)
    err = decoder.Decode(&cfg)
    if err != nil {
        return err
    }

    return nil
}

// define sprite struct to tracking 2D coordinates(row and column) information
type sprite struct {
    row      int
    col      int
    startRow int
    startCol int
}

var player sprite
var ghosts []*sprite
var maze []string
var score int
var numDots int
var lives = 1

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

    // traverse each character of the maze
    for row, line := range maze {
        for col, char := range line {
            switch char {
            case 'P':
                player = sprite{row, col, row, col}
            case 'G':
                ghosts = append(ghosts, &sprite{row, col, row, col})
            case '.':
                numDots++
            }
        }
    }

    return nil
}

func printScreen() {
    simpleansi.ClearScreen()
    for _, line := range maze {
        for _, char := range line {
            switch char {
            case '#':
                fmt.Print(simpleansi.WithBlueBackground(cfg.Wall))
            case '.':
                fmt.Print(cfg.Dot)
            case 'X':
                fmt.Print(cfg.Pill)
            default:
                fmt.Print(cfg.Space)
            }
        }
        fmt.Println()
    }

    moveCursor(player.row, player.col)
    fmt.Print(cfg.Player)

    for _, ghost := range ghosts {
        moveCursor(ghost.row, ghost.col)
        fmt.Print(cfg.Ghost)
    }

    // 将光标移出迷宫绘图区域
    moveCursor(len(maze)+1, 0)
    fmt.Println("Score:", score, "\tLives:", lives)
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

func readInput() (string, error) {
    buffer := make([]byte, 100)

    cnt, err := os.Stdin.Read(buffer)
    if err != nil {
        return "", err
    }

    if cnt == 1 && buffer[0] == 0x1b {
        return "ESC", nil
    } else if cnt >= 3 { // 方向键的转义序列有 3 个字节长
        // 以 ESC+[ 开头，然后是 A~D 之间的字母
        if buffer[0] == 0x1b && buffer[1] == '[' {
            switch buffer[2] {
            case 'A':
                return "UP", nil
            case 'B':
                return "DOWN", nil
            case 'C':
                return "RIGHT", nil
            case 'D':
                return "LEFT", nil
            }
        }
    }

    return "", nil
}

func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
    newRow, newCol = oldRow, oldCol

    switch dir {
    case "UP":
        newRow = newRow - 1
        if newRow < 0 {
            // 再次回到最下面一行
            newRow = len(maze) - 1
        }
    case "DOWN":
        newRow = newRow + 1
        if newRow == len(maze) {
            newRow = 0
        }
    case "RIGHT":
        newCol = newCol + 1
        if newCol == len(maze[0]) {
            newCol = 0
        }
    case "LEFT":
        newCol = newCol - 1
        if newCol < 0 {
            newCol = len(maze[0]) - 1
        }
    }

    // 先尝试移动，如果新的位置碰巧遇到墙（#），则移动呗取消
    if maze[newRow][newCol] == '#' {
        newRow = oldRow
        newCol = oldCol
    }

    return
}

func movePlayer(dir string) {
    player.row, player.col = makeMove(player.row, player.col, dir)

    // Remove dot from maze
    removeDot := func(row, col int) {
        maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
    }

    switch maze[player.row][player.col] {
    case '.':
        numDots--
        score++
        removeDot(player.row, player.col)
    case 'X':
        score += 10
        removeDot(player.row, player.col)
    }
}

func drawDirection() string {
    dir := rand.Intn(4)
    move := map[int]string{
        0: "UP",
        1: "DOWN",
        2: "RIGHT",
        3: "LEFT",
    }
    return move[dir]
}

func moveGhosts() {
    for _, ghost := range ghosts {
        dir := drawDirection()
        ghost.row, ghost.col = makeMove(ghost.row, ghost.col, dir)
    }
}

func moveCursor(row, col int) {
    if cfg.UseEmoji {
        // 将 col 值缩放2倍，确保每个角色都定位在正确的位置，不过会让迷宫看起来更大
        simpleansi.MoveCursor(row, col*2)
    } else {
        simpleansi.MoveCursor(row, col)
    }
}

func main() {
    flag.Parse()

    // initialize game
    initialise()
    defer cleanup()

    // load resources
    err := loadMaze(*mazeFile)
    if err != nil {
        log.Println("failed to load maze:", err)
        return
    }

    err = loadConfig(*configFile)
    if err != nil {
        log.Println("failed to load configuration:", err)
        return
    }

    // process input (async)
    input := make(chan string)
    go func(ch chan<- string) {
        for {
            input, err := readInput()
            if err != nil {
                log.Print("error reading input:", err)
                ch <- "ESC"
            }
            ch <- input
        }
    }(input)

    // game loop
    for {
        // process movement
        select {
        case inp := <-input:
            if inp == "ESC" {
                lives = 0
            }
            movePlayer(inp)
        default:
        }

        moveGhosts()

        // process collisions
        for _, ghost := range ghosts {
            if player == *ghost {
                lives--
            }
        }

        // update screen
        printScreen()

        // check game over
        if numDots == 0 || lives <= 0 {
            if lives == 0 {
                moveCursor(player.row, player.col)
                fmt.Print(cfg.Death)
                moveCursor(len(maze)+2, 0)
            }
            break
        }

        // repeat
        time.Sleep(200 * time.Millisecond)
    }
}
