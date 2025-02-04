package scte35

import (
    "fmt"

    "github.com/Comcast/scte35-go/pkg/scte35"
)

// ParseSCTE35 decodes a base64-encoded SCTE-35 message and checks for splice insert
func ParseSCTE35(encodedStr string) (*scte35.SpliceInfoSection, error) {
    // Decode the base64 string into a SpliceInfoSection
    spliceInfo, err := scte35.DecodeBase64(encodedStr)
    if err != nil {
        return nil, fmt.Errorf("error decoding SCTE-35 base64 string: %w", err)
    }

    return spliceInfo, nil
}

// HasSpliceInsert checks if the SCTE-35 message contains a splice insert command
func HasSpliceInsert(splice *scte35.SpliceInfoSection) bool {
    if splice.SpliceCommand == nil {
        return false
    }
    
    // Check if it's a SpliceInsert command by type assertion
    _, isSpliceInsert := splice.SpliceCommand.(*scte35.SpliceInsert)
    return isSpliceInsert
}
