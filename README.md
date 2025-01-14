# ttlmap

This package provides a thread-safe implementation of a Time-To-Live (TTL) map in Go. 

## Features

* **Thread-safe:** Uses a `sync.RWMutex` to ensure safe concurrent access.
* **Key-value storage:** Stores key-value pairs with associated TTL durations.
* **Automatic expiration:** Automatically removes entries after their TTL has expired.
* **Error handling:** Returns an error if the key does not exist.

## Usage

```go
package main

import (
        "fmt"
        "time"

        "github.com/your-username/ttlmap" // Replace with actual import path
)

func main() {
        ttlCache := ttlmap.NewTTLCache()

        // Add a key-value pair with a TTL of 5 seconds
        ttlCache.Add(1, "Hello", 5*time.Second)

        // Get the value
        value, err := ttlCache.Get(1)
        if err != nil {
                fmt.Println("Error:", err)
                return
        }
        fmt.Println("Value:", value)

        // Remove the entry
        err = ttlCache.Remove(1)
        if err != nil {
                fmt.Println("Error:", err)
        }
}
```

## License

The ttlmap package is licensed under the MIT License. See the LICENSE file for details.
