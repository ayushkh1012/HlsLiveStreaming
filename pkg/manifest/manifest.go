package manifest

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"

    "github.com/grafov/m3u8"
)

type Segment struct {
    URI      string
    Duration float64
    SeqID    uint64
    SCTE     string
}

type Manifest struct {
    Playlist *m3u8.MediaPlaylist
}

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

func (m *Manifest) Save(path string) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    _, err = f.Write([]byte(m.Playlist.String()))
    return err
}

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
func (m *Manifest) RemoveOldSegments(windowSize int) {
    if len(m.Playlist.Segments) > windowSize {
        m.Playlist.Segments = m.Playlist.Segments[len(m.Playlist.Segments)-windowSize:]
        m.Playlist.SeqNo += uint64(len(m.Playlist.Segments) - windowSize)
    }
}

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
    sequence      uint64
    totalSegments int
    manifestPath  string
    adDuration    int
    initialAd     bool
}

func NewManifestHandler(manifestPath string) *ManifestHandler {
    return &ManifestHandler{
        windowSize:    5,        // Show 5 segments at a time
        sequence:      0,
        totalSegments: 59,       // Total number of segments
        manifestPath:  manifestPath,
        adDuration:   30,       // 30 second ads (3 segments of 10 seconds each)
        initialAd:    true,
    }
}

func (h *ManifestHandler) UpdateManifest() error {
    h.mu.Lock()
    defer h.mu.Unlock()

    content := "#EXTM3U\n"
    content += "#EXT-X-VERSION:3\n"
    content += "#EXT-X-TARGETDURATION:10\n"
    content += fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n\n", h.sequence)

    // Always insert initial ad for sequence 0
    if h.sequence == 0 {
        content += "# Pre-roll Ad Break\n"
        content += "#EXT-X-DISCONTINUITY\n"
        content += fmt.Sprintf("#EXT-X-CUE-OUT:DURATION=%d\n", h.adDuration)
        for i := 0; i < 3; i++ {
            content += "#EXTINF:10.0,\n"
            content += fmt.Sprintf("/ads/adv1/1080p/segment_%03d.ts\n", i)
        }
        content += "#EXT-X-DISCONTINUITY\n"
        content += "#EXT-X-CUE-IN\n\n"
        h.initialAd = false
    }

    startSegment := h.sequence + 1
    endSegment := startSegment + uint64(h.windowSize)

    // Add segments within the window
    for i := startSegment; i < endSegment; i++ {
        segNum := ((i - 1) % uint64(h.totalSegments)) + 1

        // Add regular content segment
        content += "#EXTINF:10.0,\n"
        content += fmt.Sprintf("/media/1080p/segment_%03d.ts\n", segNum)

        // Add mid-roll ad break every 10 segments
        if segNum%10 == 0 {
            content += fmt.Sprintf("\n# Mid-roll Ad Break %d\n", segNum/10)
            content += "#EXT-X-DISCONTINUITY\n"
            content += fmt.Sprintf("#EXT-X-CUE-OUT:DURATION=%d\n", h.adDuration)
            for adSeg := 0; adSeg < 3; adSeg++ {
                content += "#EXTINF:10.0,\n"
                content += fmt.Sprintf("/ads/adv1/1080p/segment_%03d.ts\n", adSeg)
            }
            content += "#EXT-X-DISCONTINUITY\n"
            content += "#EXT-X-CUE-IN\n\n"
        }
    }

    if err := os.MkdirAll(filepath.Dir(h.manifestPath), 0755); err != nil {
        return fmt.Errorf("failed to create manifest directory: %v", err)
    }

    if err := os.WriteFile(h.manifestPath, []byte(content), 0644); err != nil {
        return fmt.Errorf("failed to write manifest: %v", err)
    }

    h.sequence = (h.sequence + 1) % uint64(h.totalSegments)
    return nil
}

func (h *ManifestHandler) Start() {
    if err := h.UpdateManifest(); err != nil {
        fmt.Printf("Error in initial manifest update: %v\n", err)
    }

    // Update manifest every 10 seconds
    ticker := time.NewTicker(10 * time.Second)
    go func() {
        for range ticker.C {
            if err := h.UpdateManifest(); err != nil {
                fmt.Printf("Error updating manifest: %v\n", err)
            }
        }
    }()
}
