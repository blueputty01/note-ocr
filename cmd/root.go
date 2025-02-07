package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

var rootCmd = &cobra.Command{
	Use:   "tex-screenshot",
	Short: "tex-screenshot is a cli tool for converting screenshots to LaTex",
	Long:  "tex-screenshot is a cli tool for converting screenshots to LaTex",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// content is a struct which contains a file's name, its type and its data.
type content struct {
	fname string
	ftype string
	fdata []byte
}

func Execute() {
	fmt.Printf("tex-screenshot has started")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	testingToken := os.Getenv("TESTING_TOKEN")

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}
	imageData := clipboard.Read(clipboard.FmtImage)

	// Replace with your API endpoint
	// apiURL := "https://server.simpletex.net/api/latex_ocr"
	apiURL := "http://localhost:80/post"

	imageName := "image.jpg" // Change this to your desired image name
	// Create a buffer to write the multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Create a part with a custom MIME type
	partHeader := make(textproto.MIMEHeader)
	disposition := fmt.Sprintf("form-data; name=\"files\"; filename=\"%s\"", imageName)
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

	writer.WriteField("data", "ahhhhhh hello world")

	// Close the writer to finalize the multipart form
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}
	// Print the raw multipart form data
	fmt.Println("Raw multipart form data:")
	fmt.Println(b.String())
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
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response status and body
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Println("Response Body:")
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}
	jsonPretty, _ := json.MarshalIndent(jsonResponse, "", "  ")
	fmt.Println(string(jsonPretty))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing tex-screenshot '%s'\n", err)
		os.Exit(1)
	}
}
