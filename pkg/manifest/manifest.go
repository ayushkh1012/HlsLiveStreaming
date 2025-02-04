package manifest

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/grafov/m3u8"
)

// Segment represents a media segment in the HLS playlist.
type Segment struct {
    URI      string
    Duration float64
    SeqID    uint64
    SCTE     string
}

// Manifest wraps an m3u8.MediaPlaylist for additional functionalities.
type Manifest struct {
    Playlist *m3u8.MediaPlaylist
}

// LoadManifest loads an HLS media playlist from the specified path.
func LoadManifest(path string) (*Manifest, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    p, _, err := m3u8.DecodeFrom(f, true)
    if err != nil {
        return nil, err
    }

    mediaPlaylist, ok := p.(*m3u8.MediaPlaylist)
    if !ok {
        return nil, err
    }

    return &Manifest{Playlist: mediaPlaylist}, nil
}

// Save writes the current state of the media playlist to the specified path.
func (m *Manifest) Save(path string) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    _, err = f.Write([]byte(m.Playlist.String()))
    return err
}

// InsertSegments inserts new segments into the playlist at the specified sequence ID.
func (m *Manifest) InsertSegments(seqID uint64, segments []Segment) {
    for _, seg := range segments {
        m.Playlist.AppendSegment(&m3u8.MediaSegment{
			URI:            seg.URI,
			Duration:       seg.Duration,
			Title:          "",
			SeqId:          seg.SeqID,
			ProgramDateTime: time.Now(),
		})		
        seqID++
    }
}

// RemoveOldSegments maintains the playlist length within the specified window size.
func (m *Manifest) RemoveOldSegments(windowSize int) {
    if len(m.Playlist.Segments) > windowSize {
        m.Playlist.Segments = m.Playlist.Segments[len(m.Playlist.Segments)-windowSize:]
        m.Playlist.SeqNo += uint64(len(m.Playlist.Segments) - windowSize)
    }
}

// UpdateManifest updates the playlist by inserting new segments and removing old ones.
func (m *Manifest) UpdateManifest(newSegments []Segment, windowSize int) {
    for _, seg := range newSegments {
        m.Playlist.AppendSegment(&m3u8.MediaSegment{
            URI:            seg.URI,
            Duration:       seg.Duration,
            Title:          "",
            SeqId:          seg.SeqID,
            ProgramDateTime: time.Now(),
        })
    }
    m.RemoveOldSegments(windowSize)
}

type ManifestHandler struct {
    mu            sync.Mutex
    windowSize    int
    sequence      int
    totalSegments int
    manifestPath  string
}

func NewManifestHandler(manifestPath string) *ManifestHandler {
    return &ManifestHandler{
        windowSize:    10,
        sequence:      0,
        totalSegments: 59,
        manifestPath:  manifestPath,
    }
}

func (h *ManifestHandler) UpdateManifest() error {
    h.mu.Lock()
    defer h.mu.Unlock()

    content := "#EXTM3U\n"
    content += "#EXT-X-VERSION:3\n"
    content += "#EXT-X-TARGETDURATION:12\n"
    content += fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", h.sequence)

    startSegment := h.sequence + 1
    endSegment := startSegment + h.windowSize

    for i := startSegment; i < endSegment; i++ {
        segNum := ((i - 1) % h.totalSegments) + 1

        // Add ad markers every 4 segments
        if (segNum-1)%4 == 0 && segNum > 1 {
            content += "#EXT-X-DISCONTINUITY\n"
            content += "#EXT-X-CUE-OUT:30.0\n"
            
            // Add ad segments
            for adSeg := 0; adSeg < 3; adSeg++ {
                content += "#EXTINF:10.0,\n"
                content += fmt.Sprintf("/ads/adv1/1080p/segment_%03d.ts\n", adSeg)
            }
            
            content += "#EXT-X-DISCONTINUITY\n"
            content += "#EXT-X-CUE-IN\n"
        }

        content += "#EXTINF:10.0,\n"
        content += fmt.Sprintf("/media/1080p/segment_%03d.ts\n", segNum)
    }

    // Write to file
    if err := os.MkdirAll(filepath.Dir(h.manifestPath), 0755); err != nil {
        return fmt.Errorf("failed to create manifest directory: %v", err)
    }

    if err := os.WriteFile(h.manifestPath, []byte(content), 0644); err != nil {
        return fmt.Errorf("failed to write manifest: %v", err)
    }

    h.sequence = (h.sequence + 1) % h.totalSegments
    return nil
}

func (h *ManifestHandler) Start() {
    ticker := time.NewTicker(10 * time.Second)
    go func() {
        for range ticker.C {
            if err := h.UpdateManifest(); err != nil {
                fmt.Printf("Error updating manifest: %v\n", err)
            }
        }
    }()
}
