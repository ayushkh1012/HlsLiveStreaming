package adinsertion

import (
    "fmt"
    "livehls/pkg/manifest"
    "livehls/pkg/scte35"
    "livehls/utils"
    "path/filepath"
)


func InsertAds(manifestPath string, adID string, adDuration int, logger *utils.Logger) error {
    m3u8, err := manifest.LoadManifest(manifestPath)
    if err != nil {
        return err
    }

    // Parsing SCTE-35 markers and inserting ads
    for i, segment := range m3u8.Playlist.Segments {
        if segment == nil {
            continue
        }

        if scteTag, ok := segment.Custom["SCTE"]; ok {
            cue, err := scte35.ParseSCTE35(scteTag.String())
            if err != nil {
                logger.Printf("Failed to parse SCTE-35: %v", err)
                continue
            }

            if scte35.HasSpliceInsert(cue) {
                adSegments := generateAdSegments(adID, adDuration)
                m3u8.InsertSegments(uint64(i), adSegments)
            }
        }
    }
    return m3u8.Save(manifestPath)
}

func generateAdSegments(adID string, duration int) []manifest.Segment {
    segments := make([]manifest.Segment, 0)
    remainingDuration := float64(duration)
    segmentIndex := 0
    
    for remainingDuration > 0 {
        segDuration := 10.0
        if remainingDuration < 10.0 {
            segDuration = remainingDuration
        }
        
        segments = append(segments, manifest.Segment{
            URI: filepath.Join("/ads", adID, fmt.Sprintf("1080p/segment_%03d.ts", segmentIndex)),
            Duration: segDuration,
        })
        
        remainingDuration -= segDuration
        segmentIndex++
    }
    return segments
}
