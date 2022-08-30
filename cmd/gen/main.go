package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/insomniacslk/stabilityai"
	"github.com/spf13/pflag"
)

var (
	flagAPIKey  = pflag.StringP("api-key", "a", "", "Stability AI API key")
	flagEngine  = pflag.StringP("engine", "e", stabilityai.DefaultEngine, "Stability AI engine")
	flagAPIHost = pflag.StringP("api-host", "p", stabilityai.DefaultAPIHost, "Stability AI API host")
	flagWidth   = pflag.Uint64P("width", "W", 512, "Image width")
	flagHeight  = pflag.Uint64P("height", "H", 512, "Image height")
)

func main() {
	pflag.Parse()
	if len(pflag.Args()) == 0 {
		log.Fatal("No prompt specified")
	}
	prompt := strings.Join(pflag.Args(), " ")
	c := stabilityai.NewClient(
		stabilityai.WithAPIKey(*flagAPIKey),
		stabilityai.WithAPIHost(*flagAPIHost),
		stabilityai.WithEngine(*flagEngine),
	)
	if err := c.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	answers, err := c.GenerateImage(prompt, *flagWidth, *flagHeight)
	if err != nil {
		log.Fatalf("Failed to generate: %v", err)
	}
	log.Printf("Received %d answers", len(answers))
	for idx, ans := range answers {
		fmt.Printf("%d) received %d artifacts\n", idx+1, len(ans.Artifacts))
		for _, artifact := range ans.Artifacts {
			fmt.Printf("ID       : %d\n", artifact.GetId())
			fmt.Printf("Type     : %s\n", artifact.GetType())
			fmt.Printf("MIME     : %s\n", artifact.GetMime())
			fmt.Printf("Magic    : %s\n", artifact.GetMagic())
			fmt.Printf("Text     : %s\n", artifact.GetText())
			fmt.Printf("Tokens   : %s\n", artifact.GetTokens())
			data := artifact.GetBinary()
			if data != nil {
				fd, err := ioutil.TempFile("", "stabilityai.*.png")
				if err != nil {
					fmt.Printf("Failed to create temp file: %v\n", err)
				} else {
					if _, err := fd.Write(data); err != nil {
						fmt.Printf("Failed to write to temp file: %v\n", err)
					}
					fmt.Printf("Written to file %s\n", fd.Name())
				}
			}
			if classifier := artifact.GetClassifier(); classifier != nil {
				for _, category := range classifier.GetCategories() {
					fmt.Printf("Category : %s\n", category.Name)
					fmt.Printf("Action   : %s\n", category.Action)
					for _, concept := range category.Concepts {
						fmt.Printf("Concept  : %s, threshold: %f\n", concept.Concept, *concept.Threshold)
					}
				}
			}
			fmt.Println()
		}
	}
}
