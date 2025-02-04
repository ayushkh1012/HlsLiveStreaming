package adinsertion

import (
    "fmt"
    "livehls/pkg/manifest"  // Correct module name
    "livehls/pkg/scte35"    // Correct module name
    "livehls/utils"         // Correct module name
    "path/filepath"
)


func InsertAds(manifestPath string, adID string, adDuration int, logger *utils.Logger) error {
    // Load the existing manifest
    m3u8, err := manifest.LoadManifest(manifestPath)
    if err != nil {
        return err
    }

    // Parse SCTE-35 markers and insert ads
    for i, segment := range m3u8.Playlist.Segments {
        if segment == nil {
            continue
        }

        // Get SCTE tag value using String() method
        if scteTag, ok := segment.Custom["SCTE"]; ok {
            cue, err := scte35.ParseSCTE35(scteTag.String())
            if err != nil {
                logger.Printf("Failed to parse SCTE-35: %v", err) // Using Printf instead of Error
                continue
            }

            // Check for splice insert command
            if scte35.HasSpliceInsert(cue) {
                // Insert ad segments
                adSegments := generateAdSegments(adID, adDuration)
                m3u8.InsertSegments(uint64(i), adSegments)
            }
        }
    }

    // Save the updated manifest
    return m3u8.Save(manifestPath)
}

func generateAdSegments(adID string, duration int) []manifest.Segment {
    segments := make([]manifest.Segment, 0)
    remainingDuration := float64(duration)
    segmentIndex := 0
    
    for remainingDuration > 0 {
        segDuration := 10.0 // Standard segment duration
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
