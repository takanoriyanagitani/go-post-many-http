package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
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

var BodyDir string = getenvOrAlt("ENV_TRUSTED_DIR_NAME", "./requests.d")

var TargetUrl string = getenvOrAlt("ENV_TARGET_URL", "http://localhost")
var TargetTyp string = getenvOrAlt("ENV_TARGET_TYP", "text/plain")

var MaxBodySize int = parseIntOrAlt(
	getenvOrAlt("ENV_MAX_BODY_SIZE", "1048576"),
	1048576,
)

func main() {
	log.Printf("body dir(ENV_TRUSTED_DIR_NAME): %s\n", BodyDir)
	log.Printf("max body size(ENV_MAX_BODY_SIZE): %v\n", MaxBodySize)
	log.Printf("target url(ENV_TARGET_URL): %s\n", TargetUrl)
	log.Printf("target typ(ENV_TARGET_TYP): %s\n", TargetTyp)

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
				log.Printf("rejected url:      %s\n", req.Url)
				log.Printf("rejected type:     %s\n", req.ContentType)
				log.Printf("rejected body len: %v\n", len(req.Body))
				return fmt.Errorf("unexpected status: %v", code)
			}
		},
	)
	if nil != e {
		log.Fatalf("%v\n", e)
	}
	log.Printf("sent count: %v\n", tot)
}
