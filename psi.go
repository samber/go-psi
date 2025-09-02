package psi

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func AllPSIStats() (PSIStatsResource, error) {
	cpu, err := PSIStatsForResource(CPU)
	if err != nil {
		return PSIStatsResource{}, err
	}

	memory, err := PSIStatsForResource(Memory)
	if err != nil {
		return PSIStatsResource{}, err
	}

	io, err := PSIStatsForResource(IO)
	if err != nil {
		return PSIStatsResource{}, err
	}

	return PSIStatsResource{
		CPU:    &cpu,
		Memory: &memory,
		IO:     &io,
	}, nil
}

// PSIStatsForResource reads pressure stall information for the specified
// resource from /proc/pressure/<resource>. At time of writing this can be
// either "cpu", "memory" or "io".
func PSIStatsForResource(resource Resource) (PSIStats, error) {
	data, err := ReadFileNoStat(ResourceToPath(resource))
	if err != nil {
		return PSIStats{}, fmt.Errorf("psi_stats: unavailable for %q: %w", resource, err)
	}

	return parsePSIStats(resource, bytes.NewReader(data))
}

// parsePSIStats parses the specified file for pressure stall information.
func parsePSIStats(resource Resource, r io.Reader) (PSIStats, error) {
	psiStats := PSIStats{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()
		prefix := strings.Split(l, " ")[0]
		switch prefix {
		case "some":
			psi := PSILine{}
			_, err := fmt.Sscanf(l, fmt.Sprintf("some %s", lineFormat), &psi.Avg10, &psi.Avg60, &psi.Avg300, &psi.Total)
			if err != nil {
				return PSIStats{}, err
			}
			psiStats.Some = &psi
		case "full":
			psi := PSILine{}
			_, err := fmt.Sscanf(l, fmt.Sprintf("full %s", lineFormat), &psi.Avg10, &psi.Avg60, &psi.Avg300, &psi.Total)
			if err != nil {
				return PSIStats{}, err
			}
			psiStats.Full = &psi
		default:
			// If we encounter a line with an unknown prefix, ignore it and move on
			// Should new measurement types be added in the future we'll simply ignore them instead
			// of erroring on retrieval
			continue
		}
	}

	return psiStats, nil
}

// ReadFileNoStat uses io.ReadAll to read contents of entire file.
// This is similar to os.ReadFile but without the call to os.Stat, because
// many files in /proc and /sys report incorrect file sizes (either 0 or 4096).
// Reads a max file size of 1024kB.  For files larger than this, a scanner
// should be used.
func ReadFileNoStat(filename string) ([]byte, error) {
	const maxBufferSize = 1024 * 1024

	if !strings.HasPrefix(filename, "/proc/") {
		return nil, fmt.Errorf("file %q is not in /proc", filename)
	}

	// bearer:disable go_gosec_filesystem_filereadtaint
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck

	reader := io.LimitReader(f, maxBufferSize)
	return io.ReadAll(reader)
}

func CompareThreshold(threshold int, current int) bool {
	return threshold >= 0 && current > threshold
}
