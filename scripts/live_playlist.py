import os
import time
import shutil

WINDOW_SIZE = 10  # Number of segments to keep in the playlist
MANIFEST_PATH = "manifests/input.m3u8"
SEGMENT_DURATION = 10.0  # Duration of each segment in seconds

def create_live_playlist():
    sequence_number = 0
    total_segments = 59  # Total number of segments available

    while True:
        try:
            with open(MANIFEST_PATH, 'w') as f:
                # Write header
                f.write('#EXTM3U\n')
                f.write('#EXT-X-VERSION:3\n')
                f.write('#EXT-X-TARGETDURATION:12\n')
                f.write(f'#EXT-X-MEDIA-SEQUENCE:{sequence_number}\n')

                # Calculate segment range for current window
                start_segment = sequence_number + 1
                end_segment = min(start_segment + WINDOW_SIZE, total_segments + 1)

                # Write segments
                for i in range(start_segment, end_segment):
                    segment_num = i
                    if segment_num <= total_segments:
                        f.write(f'#EXTINF:{SEGMENT_DURATION},\n')
                        f.write(f'/media/1080p/segment_{segment_num:03d}.ts\n')

                # Check if we need to loop back
                if end_segment >= total_segments:
                    sequence_number = 0
                else:
                    sequence_number += 1

            print(f"Updated playlist with sequence {sequence_number}")
            time.sleep(10)  # Wait for 10 seconds before next update

        except Exception as e:
            print(f"Error updating playlist: {e}")
            time.sleep(1)  # Wait briefly before retrying

if __name__ == "__main__":
    # Ensure manifests directory exists
    os.makedirs("manifests", exist_ok=True)
    
    try:
        create_live_playlist()
    except KeyboardInterrupt:
        print("\nStopping live playlist generation...") 