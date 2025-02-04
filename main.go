package main

import (
    "fmt"
    "livehls/pkg/manifest"
    "livehls/utils"
    "log"
    "net/http"
    "path/filepath"
)

func main() {
    // Load configuration
    config, err := utils.LoadConfig("config/config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    logger := utils.NewLogger()

    // Initialize manifest handler with live streaming
    manifestHandler := manifest.NewManifestHandler("manifests/input.m3u8")
    manifestHandler.Start() // Start the live manifest updates

    // Set up HTTP server with live streaming headers
    http.Handle("/media/", addLiveHeaders(http.StripPrefix("/media/", http.FileServer(http.Dir(config.Paths.Media)))))
    http.Handle("/ads/", addLiveHeaders(http.StripPrefix("/ads/", http.FileServer(http.Dir(config.Paths.Ads)))))
    http.Handle("/manifests/", addLiveHeaders(http.StripPrefix("/manifests/", http.FileServer(http.Dir(config.Paths.Manifests)))))

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

func addLiveHeaders(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        w.Header().Set("Expires", "0")
        
        if filepath.Ext(r.URL.Path) == ".m3u8" {
            w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
        } else if filepath.Ext(r.URL.Path) == ".ts" {
            w.Header().Set("Content-Type", "video/mp2t")
        }
        
        handler.ServeHTTP(w, r)
    })
}
