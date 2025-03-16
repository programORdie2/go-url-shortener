# Go URL Shortener
This is a basic URL shortener application built in Go to learn the basics of Go. It demonstrates the use of Go's `net/http` package for handling HTTP requests and `sqlite` as a lightweight database for storing URLs. The application allows users to shorten URLs and redirect to the original URLs using the shortened ones.

## Features

- Shorten long URLs to short codes
- Redirect to the original URL using the short code
- Simple web interface for URL shortening
- Persistent storage using SQLite

## Prerequisites

- Go installed

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/programordie2/go-url-shortener.git
   cd go-url-shortener
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the application:

   ```bash
   go run main.go
   ```

4. Open your web browser and navigate to `http://localhost:8080` to access the URL shortener.

## Project Structure

- `main.go`: The main server implementation using Go's `net/http` package.
- `static/`: Contains static files for the web interface.

## Usage

1. Enter the URL you want to shorten in the input box on the web page.
2. Click the "Shorten" button.
3. The shortened URL will be displayed, which you can use to redirect to the original URL.

## License

This project is licensed under the MIT License.
