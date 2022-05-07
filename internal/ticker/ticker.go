package ticker

import (
	"fmt"
	"github.com/gonevo/logium/internal/config"
	"log"
	"os"
	"sync"
	"time"
)

// Tick starts debug tickers
func Tick(c *config.Config) {

	if len(c.Sources) == 0 {
		return
	}

	ticker := time.NewTicker(time.Second)
	done := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(len(c.Sources))

	for _, source := range c.Sources {
		go func(source config.Source) {
			defer wg.Done()
			file, err := os.OpenFile(source.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					if _, err := file.WriteString(fmt.Sprintf("%s\n", t)); err != nil {
						log.Println(err)
					}
				}
			}
		}(source)
	}

	wg.Wait()
	log.Println("done")
}
