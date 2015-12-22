package eye

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"
)

//启动pprof
func StartPprof() {
	pprofsock := os.Getenv("GOPPROFSOCK")
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("[Eye] pprof panic", err)
			}
		}()
		if len(pprofsock) <= 0 {
			panic(`[Eye] pprof env "GOPPROF" not conf`)
		}
		_, err := os.Stat(filepath.Dir(pprofsock))
		if err != nil {
			panic("[Eye] pprof unixsock path not exist")
		}
		profServeMux := http.NewServeMux()
		profServeMux.HandleFunc("/debug/pprof/", pprof.Index)
		profServeMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		profServeMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		profServeMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		profServeMux.HandleFunc("/gc", PrintGCSummary)

		exec.Command("/bin/sh", "-c", "rm "+pprofsock).Run()
		unix, err := net.Listen("unix", pprofsock)
		exec.Command("/bin/sh", "-c", "chmod a+w "+pprofsock).Run()
		if err != nil {
			log.Println("[Eey]Listen error:", err)
		}
		log.Println("[Eey][pprof]Listen as unix:", pprofsock)
		fcgi.Serve(unix, profServeMux)
	}()
}

var startTime = time.Now()
var pid int

func init() {
	pid = os.Getpid()
}

// print gc information to io.Writer
func PrintGCSummary(w http.ResponseWriter, r *http.Request) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	gcstats := &debug.GCStats{PauseQuantiles: make([]time.Duration, 100)}
	debug.ReadGCStats(gcstats)

	printGC(memStats, gcstats, w)
}

func printGC(memStats *runtime.MemStats, gcstats *debug.GCStats, w http.ResponseWriter) {

	if gcstats.NumGC > 0 {
		lastPause := gcstats.Pause[0]
		elapsed := time.Now().Sub(startTime)
		overhead := float64(gcstats.PauseTotal) / float64(elapsed) * 100
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(w, "NumGC:%d Pause:%s Pause(Avg):%s Overhead:%3.2f%% Alloc:%s Sys:%s Alloc(Rate):%s/s Histogram:%s %s %s \n",
			gcstats.NumGC,
			toS(lastPause),
			toS(avg(gcstats.Pause)),
			overhead,
			toH(memStats.Alloc),
			toH(memStats.Sys),
			toH(uint64(allocatedRate)),
			toS(gcstats.PauseQuantiles[94]),
			toS(gcstats.PauseQuantiles[98]),
			toS(gcstats.PauseQuantiles[99]))
	} else {
		// while GC has disabled
		elapsed := time.Now().Sub(startTime)
		allocatedRate := float64(memStats.TotalAlloc) / elapsed.Seconds()

		fmt.Fprintf(w, "Alloc:%s Sys:%s Alloc(Rate):%s/s\n",
			toH(memStats.Alloc),
			toH(memStats.Sys),
			toH(uint64(allocatedRate)))
	}
}

func avg(items []time.Duration) time.Duration {
	var sum time.Duration
	for _, item := range items {
		sum += item
	}
	return time.Duration(int64(sum) / int64(len(items)))
}

// format bytes number friendly
func toH(bytes uint64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", float64(bytes)/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", float64(bytes)/1024/1024/1024)
	}
}

// short string format
func toS(d time.Duration) string {

	u := uint64(d)
	if u < uint64(time.Second) {
		switch {
		case u == 0:
			return "0"
		case u < uint64(time.Microsecond):
			return fmt.Sprintf("%.2fns", float64(u))
		case u < uint64(time.Millisecond):
			return fmt.Sprintf("%.2fus", float64(u)/1000)
		default:
			return fmt.Sprintf("%.2fms", float64(u)/1000/1000)
		}
	} else {
		switch {
		case u < uint64(time.Minute):
			return fmt.Sprintf("%.2fs", float64(u)/1000/1000/1000)
		case u < uint64(time.Hour):
			return fmt.Sprintf("%.2fm", float64(u)/1000/1000/1000/60)
		default:
			return fmt.Sprintf("%.2fh", float64(u)/1000/1000/1000/60/60)
		}
	}

}
