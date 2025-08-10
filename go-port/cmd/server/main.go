package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/webgl-water-go/go-port/internal/app"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	// Command line flags
	var (
		port       = flag.Int("port", getEnvInt("PORT", 8080), "HTTP server port")
		assetsPath = flag.String("assets", getEnvString("ASSETS_PATH", "./assets"), "Path to assets directory")
		staticPath = flag.String("static", getEnvString("STATIC_PATH", "./web/static"), "Path to static files directory")
		help       = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		printUsage()
		return
	}

	// Validate paths
	if err := validatePaths(*assetsPath, *staticPath); err != nil {
		log.Fatalf("Path validation failed: %v", err)
	}

	log.Printf("Starting WebGL Water Server")
	log.Printf("Port: %d", *port)
	log.Printf("Assets path: %s", *assetsPath)
	log.Printf("Static path: %s", *staticPath)

	// Create and start server
	server := app.NewServer(*assetsPath, *staticPath, *port)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func printUsage() {
	log.Printf("WebGL Water Tutorial - Go Port")
	log.Printf("")
	log.Printf("Usage: %s [options]", os.Args[0])
	log.Printf("")
	log.Printf("Options:")
	log.Printf("  -port int        HTTP server port (default 8080)")
	log.Printf("  -assets string   Path to assets directory (default ./assets)")
	log.Printf("  -static string   Path to static files directory (default ./web/static)")
	log.Printf("  -help           Show this help message")
	log.Printf("")
	log.Printf("Environment variables:")
	log.Printf("  PORT            HTTP server port")
	log.Printf("  ASSETS_PATH     Path to assets directory")
	log.Printf("  STATIC_PATH     Path to static files directory")
	log.Printf("")
	log.Printf("Example:")
	log.Printf("  %s -port 3000 -assets /path/to/assets -static /path/to/static", os.Args[0])
	log.Printf("")
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Simple integer parsing
		if result, err := strconv.Atoi(value); err == nil {
			return result
		}
	}
	return defaultValue
}

func validatePaths(assetsPath, staticPath string) error {
	// Check if assets path exists or can be created
	if err := ensureDirectory(assetsPath); err != nil {
		return err
	}

	// Check if static path exists or can be created
	if err := ensureDirectory(staticPath); err != nil {
		return err
	}

	return nil
}

func ensureDirectory(path string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Check if directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Printf("Directory %s does not exist, creating it", absPath)
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return err
		}
	}

	return nil
}
