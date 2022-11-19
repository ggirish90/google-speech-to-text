// Command quickstart generates an audio file with the content "Hello, World!".
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	option "google.golang.org/api/option"
)

func main() {
	// Instantiates a client.
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile("leafy-pilot-367607-9f0b82dc8406.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	reqList := &texttospeechpb.ListVoicesRequest{}
	respList, err := client.ListVoices(ctx, reqList)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(respList.GetVoices()))
	for _, val := range respList.GetVoices() {
		go func(val *texttospeechpb.Voice) {
			fmt.Println("voice", val.SsmlGender, val.LanguageCodes[0])
			// Perform the text-to-speech request on the text input with the selected
			// voice parameters and audio file type.
			req := texttospeechpb.SynthesizeSpeechRequest{
				// Set the text input to be synthesized.
				Input: &texttospeechpb.SynthesisInput{
					InputSource: &texttospeechpb.SynthesisInput_Text{Text: `"Currently, most of DevOps engineers use Python as their DevOps programming language. But Python came with some compile-time and service-scaling issues.
		
					For example, simply upgrade to a new version of Python (Python2 to Python3, Python3.7 to Python3.11) might cause your existing script to stop working. When compatibility issue happens, it is not easy to roll back to an old version of Python.
					
					Now in cloud era, Go has become the de facto language for cloud native orchestration and applications. Go comes with all the tools you need to make huge strides in the reliability of your tooling and ability to scale."`},
				},
				// Build the voice request, select the language code ("en-US") and the SSML
				// voice gender ("neutral").
				Voice: &texttospeechpb.VoiceSelectionParams{
					LanguageCode: val.LanguageCodes[0],
					SsmlGender:   val.SsmlGender,
				},
				// Select the type of audio file you want returned.
				AudioConfig: &texttospeechpb.AudioConfig{
					AudioEncoding: texttospeechpb.AudioEncoding_MP3,
				},
			}

			resp, err := client.SynthesizeSpeech(ctx, &req)
			if err != nil {
				log.Fatal(err)
			}

			// The resp's AudioContent is binary.
			filename := val.LanguageCodes[0] + ".mp3"
			err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Audio content written to file: %v\n", filename)
			wg.Done()
		}(val)
	}
	wg.Wait()
}
