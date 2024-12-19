package main

import (
	"bufio"
	"bytes"
	"ducktail/internal/server"
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"sync"
)

//go:embed website/index.html
var indexHTML string

//go:embed website/static/css/style.css
var styleCSS string

//go:embed website/static/js/script.js
var scriptJS string

func main() {
	logsChannel := make(chan string)
	var wg sync.WaitGroup

	serverInstance := server.API{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		opts := server.ServerOpts{LogsCh: logsChannel, IndexHTML: indexHTML, StyleCSS: styleCSS, ScriptJS: scriptJS}
		serverInstance.Start(opts)
	}()

	reader := bufio.NewReader(os.Stdin)
	var buffer bytes.Buffer

	go func() {
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil && err.Error() != "EOF" {
				log.Println("Error reading stdin:", err)
				break
			}

			// Accumulate the data in the buffer
			buffer.Write(line)

			// Try to unmarshal the accumulated data
			var jsonObj map[string]interface{}
			if err := json.Unmarshal(buffer.Bytes(), &jsonObj); err == nil {
				// Successfully unmarshalled log
				formattedJSON, err := json.MarshalIndent(jsonObj, "", "  ")
				if err != nil {
					log.Println("Error formatting JSON:", err)
					continue
				}

				logsChannel <- string(formattedJSON)

				// Reset the buffer after sending the log
				buffer.Reset()
			}

			if err != nil {
				break
			}
		}

		close(logsChannel)
	}()

	wg.Wait()
}
