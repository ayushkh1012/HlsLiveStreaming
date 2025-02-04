package manifest

import (
    "os"
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
