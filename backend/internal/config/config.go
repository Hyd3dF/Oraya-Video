package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	CORS      CORSConfig
	Supabase  SupabaseConfig
	Storage   StorageConfig
	Worker    WorkerConfig
	Admin     AdminConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port      string
	Env       string
	LogLevel  string
	PublicURL string // backend'in dışarıdan görünen URL'si — frontend için
}

type CORSConfig struct {
	AllowedOrigins []string
}

type SupabaseConfig struct {
	URL            string
	AnonKey        string
	ServiceRoleKey string
	JWTSecret      string
}

type StorageConfig struct {
	BucketRaw        string
	BucketHLS        string
	BucketThumbnails string
}

type WorkerConfig struct {
	TempDir     string
	Concurrency int
	FFmpegBin   string
	FFprobeBin  string
}

type AdminConfig struct {
	APIToken string
}

type RateLimitConfig struct {
	RPS   int
	Burst int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:      getEnv("SERVER_PORT", "8080"),
			Env:       getEnv("SERVER_ENV", "development"),
			LogLevel:  getEnv("LOG_LEVEL", "info"),
			PublicURL: getEnv("PUBLIC_API_URL", "http://localhost:8080"),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),
		},
		Supabase: SupabaseConfig{
			URL:            mustEnv("SUPABASE_URL"),
			AnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
			ServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
			JWTSecret:      getEnv("SUPABASE_JWT_SECRET", ""),
		},
		Storage: StorageConfig{
			BucketRaw:        getEnv("STORAGE_BUCKET_RAW", "raw-videos"),
			BucketHLS:        getEnv("STORAGE_BUCKET_HLS", "hls-videos"),
			BucketThumbnails: getEnv("STORAGE_BUCKET_THUMBNAILS", "thumbnails"),
		},
		Worker: WorkerConfig{
			TempDir:     getEnv("WORKER_TEMP_DIR", os.TempDir()),
			Concurrency: getEnvInt("WORKER_CONCURRENCY", 2),
			FFmpegBin:   getEnv("FFMPEG_BIN", "ffmpeg"),
			FFprobeBin:  getEnv("FFPROBE_BIN", "ffprobe"),
		},
		Admin: AdminConfig{
			APIToken: getEnv("ADMIN_API_TOKEN", ""),
		},
		RateLimit: RateLimitConfig{
			RPS:   getEnvInt("RATE_LIMIT_RPS", 10),
			Burst: getEnvInt("RATE_LIMIT_BURST", 20),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.Supabase.URL == "" {
		return fmt.Errorf("SUPABASE_URL is required")
	}
	return nil
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func mustEnv(key string) string {
	return getEnv(key, "")
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
