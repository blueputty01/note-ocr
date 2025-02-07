package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

var permCmd = &cobra.Command{
	Use:     "perm",
	Aliases: []string{"addition"},
	Short:   "Latex from screenshot",
	Long:    "Infinitely loop parsing latex from screenshot",
	Run: func(cmd *cobra.Command, args []string) {
		Perm()
	},
}

func init() {
	rootCmd.AddCommand(permCmd)
}

func Perm() {
	fmt.Printf("tex-screenshot has started\n")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	testingToken := os.Getenv("TESTING_TOKEN")

	for {
		fmt.Println("Press Enter to initiate a new scan")
		fmt.Scanln() // wait for Enter Key
		err = clipboard.Init()
		if err != nil {
			panic(err)
		}
		imageData := clipboard.Read(clipboard.FmtImage)
		if imageData == nil {
			fmt.Printf("Exiting; no image data on clipboard\n")
			return
		}
		// Replace with your API endpoint
		apiURL := "https://server.simpletex.net/api/latex_ocr"
		// apiURL := "http://localhost:80/post"

		imageName := "image.jpg" // Change this to your desired image name
		// Create a buffer to write the multipart form data
		var b bytes.Buffer
		writer := multipart.NewWriter(&b)

		// Create a part with a custom MIME type
		partHeader := textproto.MIMEHeader{}
		disposition := fmt.Sprintf("form-data; name=\"file\"; filename=\"%s\"", imageName)
		partHeader.Add("Content-Disposition", disposition)
		partHeader.Add("Content-Type", "image/png")
		part, err := writer.CreatePart(partHeader)
		// Create a form file field with the image data
		if err != nil {
			fmt.Println("Error creating form file:", err)
			return
		}

		// Write the image data to the form file field
		_, err = part.Write(imageData)
		if err != nil {
			fmt.Println("Error writing image data:", err)
			return
		}

		// Close the writer to finalize the multipart form
		err = writer.Close()
		if err != nil {
			fmt.Println("Error closing writer:", err)
			return
		}
		// Print the raw multipart form data
		// fmt.Println("Raw multipart form data:")
		// fmt.Println(b.String())
		// Create a new POST request
		req, err := http.NewRequest("POST", apiURL, &b)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the content type to multipart/form-data
		req.Header.Set("Content-Type", writer.FormDataContentType())

		req.Header.Add("token", testingToken)
		// Send the request
		fmt.Printf("sending request...\n")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		// Print the response status and body
		var jsonResponse Response
		if err := json.Unmarshal(body, &jsonResponse); err != nil {
			fmt.Println("Error parsing JSON response:", err)
			return
		}

		if !jsonResponse.Status {
			fmt.Println("Error")
		}

		latex := jsonResponse.Res.Latex

		fmt.Println(latex)

		clipboard.Write(clipboard.FmtText, []byte(latex))

		fmt.Println("Response copied to clipboard", err)
	}
}
