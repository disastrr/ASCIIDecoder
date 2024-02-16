// just go run main.go to start the server
// then go to localhost:8080 in your browser to see the interface
// i also made it pretty (˵¯͒〰¯͒˵)

package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Data structure to hold the latest decoded string
type ViewData struct {
	DecodedString string
	StatusCode    int
}

var viewData ViewData

func main() {
	// Define HTTP endpoints
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/decoder", handleDecoder)
	// Serve static files from the "static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	// Start the server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	} else {
		fmt.Println("The server is running on port 8080")
	}
}

// Handler for the main page (GET /)
func handleMainPage(w http.ResponseWriter, r *http.Request) {
	// Display the main page with the input form and the latest decoded string
	tmpl, err := template.ParseFiles("template.html")

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the HTML template
	err = tmpl.Execute(w, viewData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Handler for decoding the string (POST /decoder)
func handleDecoder(w http.ResponseWriter, r *http.Request) {
	// Retrieve the encoded string from the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	encodedString := r.Form.Get("encodedString")

	// Call the generateArt function to decode the string
	decodedString, err := generateArt(encodedString)
	if err != nil {
		// Handle decoding errors
		viewData.StatusCode = http.StatusBadRequest
		viewData.DecodedString = fmt.Sprintf("Decoding Error: %v", err)
	} else {
		// Update the view data with the latest decoded string
		viewData.DecodedString = decodedString
		viewData.StatusCode = http.StatusAccepted
	}

	// Redirect to the main page after decoding
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func generateArt(input string) (string, error) {
	// looking for "[digit string]"
	re := regexp.MustCompile(`\[(\d+)\s(.*?)\]`)
	matches := re.FindAllStringSubmatch(input, -1)

	// ERROR HANDLING: Return an error if no matches are found
	if len(matches) == 0 {
		return "", errors.New("no matches found in the input")
	}

	for _, match := range matches {
		count, err := strconv.Atoi(match[1])

		// ERROR HANDLING
		if err != nil {
			return "", fmt.Errorf("converting count to integer: %v", err)
		}

		if len(match[2]) == 0 {
			return "", errors.New("the second argument cannot be an empty string")
		}

		if !strings.Contains(match[0], " ") {
			return "", errors.New("the arguments must be separated by a space")
		}

		// Make sure square brackets are not empty and follow the specified pattern
		if match[0] == "[]" || !re.MatchString(match[0]) {
			return "", errors.New("invalid square bracket pattern")
		}

		// converts and replaces
		replacement := strings.Repeat(match[2], count)
		input = strings.Replace(input, match[0], replacement, 1)
	}

	return input, nil
}
