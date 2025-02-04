package main

import (
    "fmt"
    "livehls/utils"
    "log"
    "net/http"
)

func main() {
    // Load configuration
    config, err := utils.LoadConfig("config/config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    logger := utils.NewLogger()

    // Set up HTTP server
    // Static file handlers
    http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(config.Paths.Media))))
    http.Handle("/ads/", http.StripPrefix("/ads/", http.FileServer(http.Dir(config.Paths.Ads))))
    http.Handle("/manifests/", http.StripPrefix("/manifests/", http.FileServer(http.Dir(config.Paths.Manifests))))

    // Root handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        
        // Simple HTML response
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html>
            <head>
                <title>Live HLS Server</title>
            </head>
            <body>
                <h1>Live HLS Server</h1>
                <p>Available endpoints:</p>
                <ul>
                    <li><a href="/media/">/media/</a> - Media files</li>
                    <li><a href="/ads/">/ads/</a> - Advertisement files</li>
                    <li><a href="/manifests/">/manifests/</a> - HLS manifests</li>
                </ul>
            </body>
        </html>
        `)
    })

    // Start server
    addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
    logger.Printf("Starting server at %s", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        logger.Fatalf("Server failed: %v", err)
    }
}
