package main

import "os"
import "fmt"

func main() {
    var env = os.Getenv("MEDKIT_ENV")
    if env == "" {
        fmt.Println("You have no environment set.  Please define an environment variable MEDKIT_ENV")
    } else {
        fmt.Println("Your environment is " + env)
    }
}
