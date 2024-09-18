package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	pm "github.com/takanoriyanagitani/go-post-many-http"

	fd "github.com/takanoriyanagitani/go-post-many-http/source/raw/fs/dir"

	hs "github.com/takanoriyanagitani/go-post-many-http/sender/net/tcp/http/std"
)

func getenvOrAlt(key, alt string) string {
	val, found := os.LookupEnv(key)
	switch found {
	case found:
		return val
	default:
		return alt
	}
}

func parseIntOrAlt(s string, alt int) int {
	parsed, e := strconv.Atoi(s)
	switch e {
	case nil:
		return parsed
	default:
		return alt
	}
}

const DefaultLogLevel slog.Level = slog.LevelInfo

var logLevelMap map[string]slog.Level = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

var LogLevel slog.Level = func(levelString string) slog.Level {
	level, ok := logLevelMap[levelString]
	switch ok {
	case true:
		return level
	default:
		return DefaultLogLevel
	}
}(getenvOrAlt("ENV_LOG_LEVEL", "info"))

var BodyDir string = getenvOrAlt("ENV_TRUSTED_DIR_NAME", "./requests.d")

var TargetUrl string = getenvOrAlt("ENV_TARGET_URL", "http://localhost")
var TargetTyp string = getenvOrAlt("ENV_TARGET_TYP", "text/plain")

var SaveRequestOnError string = getenvOrAlt("ENV_SAVE_REQ", "false")
var SaveNameOnError string = getenvOrAlt("ENV_SAVE_NAME", "./err.dat")

var SaveReqOnError bool = "true" == SaveRequestOnError

var MaxBodySize int = parseIntOrAlt(
	getenvOrAlt("ENV_MAX_BODY_SIZE", "1048576"),
	1048576,
)

func main() {
	slog.SetLogLoggerLevel(LogLevel)

	slog.Info("log level set.", "log level", LogLevel)

	slog.Info("body dir set.", "ENV_TRUSTED_DIR_NAME", BodyDir)
	slog.Info("max body size set.", "ENV_MAX_BODY_SIZE", MaxBodySize)
	slog.Info("target url set.", "ENV_TARGET_URL", TargetUrl)
	slog.Info("target typ set.", "ENV_TARGET_TYP", TargetTyp)

	slog.Info("save request on error set.", "ENV_SAVE_REQ", SaveRequestOnError)
	slog.Info("save name on error set.", "ENV_SAVE_NAME", SaveNameOnError)

	rawSrcFs := fd.RawSourceManyFsDir{
		TrustedDirName: ".",
		ReadDirFS:      os.DirFS(BodyDir).(fs.ReadDirFS),
		MaxBodySize:    int64(MaxBodySize),
	}
	var rawSrc pm.RawRequestSourceMany = rawSrcFs.ToRawRequestSourceMany(
		io.EOF,
	)
	var simpleSrc pm.SimpleRequestSourceMany = rawSrc.ToSimpleReqSource(
		TargetUrl,
		TargetTyp,
	)
	var ctx context.Context = context.Background()
	var requests <-chan pm.SimpleRequest = simpleSrc.ToChannel(
		ctx,
		func(e error) (end bool) { return io.EOF == e },
		func(e error) { log.Fatalf("error: %v\n", e) },
	)

	var senderStd hs.SenderStdST = hs.SenderStdNewSTdefault(http.DefaultClient)
	var sender pm.SenderST = senderStd.ToSender()
	tot, e := sender.SendManyEx(
		ctx,
		requests,
		func(res pm.TinyResponse, req pm.SimpleRequest) error {
			var code int = res.StatusCode
			switch code {
			case http.StatusOK:
				return nil
			default:
				slog.Info("rejected.", "url", req.Url)
				slog.Info("rejected.", "type", req.ContentType)
				slog.Info("rejected.", "body len", len(req.Body))
				if SaveReqOnError {
					_ = os.WriteFile(
						SaveNameOnError,
						req.Body,
						0644,
					)
				}
				return fmt.Errorf("unexpected status: %v", code)
			}
		},
	)
	if nil != e {
		log.Fatalf("%v\n", e)
	}
	slog.Info("sent.", "count", tot)
}
