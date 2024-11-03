package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    if len(os.Args) == 1 {
        fmt.Println("Please give one or more integers.")
        return
    }

    var min, max int

    arguments := os.Args
    temp, err := strconv.Atoi(arguments[1])
    if err != nil {
        fmt.Println("Error encountered, exiting:")
        fmt.Println(err)
        return
    } else {
        min = temp
        max = temp
    }

    for i := 2; i < len(arguments); i++ {
        n, err := strconv.Atoi(arguments[i])
        if err != nil {
            fmt.Println("Error encountered, exiting:")
            fmt.Println(err)
            return
        }

        if n < min {
            min = n
        }
        if n > max {
            max = n
        }
    }

    fmt.Println("Min:", min)
    fmt.Println("Max:", max)
}