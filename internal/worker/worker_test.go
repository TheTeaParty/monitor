package worker

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func Test_doRequests(t *testing.T) {

	urls := []string{
		"google.com",
		"youtube.com",
		"facebook.com",
		"baidu.com",
		"wikipedia.org",
		"yahoo.com",
		"tmall.com",
		"amazon.com",
		"twitter.com",
		"live.com",
		"instagram.com",
		"google.com",
	}

	for i := 0; i < 10; i++ {
		urls = append(urls, urls...)
	}

	startTime := time.Now()
	results := doRequests(urls, runtime.GOMAXPROCS(runtime.NumCPU()))
	assert.Len(t, results, len(urls))
	seconds := time.Since(startTime).Seconds()

	t.Log(fmt.Sprintf("%v urls completed in %v seconds", len(urls), seconds))
}
