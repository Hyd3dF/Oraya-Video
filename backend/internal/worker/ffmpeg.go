package worker

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

// HLSRendition describes one quality variant produced by FFmpeg.
type HLSRendition struct {
	Quality string
	Width   int
	Height  int
	Bitrate int // bits/second
	MaxRate int
	BufSize int
}

// MVPRendition720p is the only quality produced by the MVP.
var MVPRendition720p = HLSRendition{
	Quality: "720p",
	Width:   1280,
	Height:  720,
	Bitrate: 2000_000,
	MaxRate: 2500_000,
	BufSize: 4000_000,
}

type FFmpeg struct {
	Bin     string
	ProbeBin string
}

func NewFFmpeg(bin, probeBin string) *FFmpeg {
	return &FFmpeg{Bin: bin, ProbeBin: probeBin}
}

// TranscodeHLS converts an input file to a single-quality HLS package
// (master.m3u8 + playlist.m3u8 + segment_*.ts) under outDir.
//
// Mirrors the ffmpeg invocation documented in oroya.md §Processing Flow.
func (f *FFmpeg) TranscodeHLS(ctx context.Context, input, outDir string, r HLSRendition) error {
	playlist := filepath.Join(outDir, "playlist.m3u8")
	master := "master.m3u8"
	segments := filepath.Join(outDir, "segment_%03d.ts")

	args := []string{
		"-y",
		"-i", input,
		"-vf", fmt.Sprintf("scale=-2:%d,format=yuv420p", r.Height),
		"-c:v", "libx264",
		"-b:v", fmt.Sprintf("%d", r.Bitrate),
		"-maxrate", fmt.Sprintf("%d", r.MaxRate),
		"-bufsize", fmt.Sprintf("%d", r.BufSize),
		"-c:a", "aac",
		"-b:a", "128k",
		"-ar", "48000",
		"-f", "hls",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", segments,
		"-master_pl_name", master,
		playlist,
	}

	cmd := exec.CommandContext(ctx, f.Bin, args...)
	cmd.Dir = outDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w: %s", err, string(out))
	}
	return nil
}

// ProbeDuration returns video duration in seconds via ffprobe.
func (f *FFmpeg) ProbeDuration(ctx context.Context, input string) (float64, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		input,
	}
	cmd := exec.CommandContext(ctx, f.ProbeBin, args...)
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}
	var d float64
	_, err = fmt.Sscanf(string(out), "%f", &d)
	return d, err
}
