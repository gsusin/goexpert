package limiter

import (
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type LimitedHandler struct {
	AccessStorage
	AccessState
}

type AccessStorage interface {
	SetZero(ip string)
	Increment(ip string)
	Get(ip string) int
}

type AccessState struct {
	periodInSeconds int
	limitOfRequests int
	start           time.Time
	blockInSeconds  int
	blockStart      time.Time
	isBlocked       bool
	tokens          [2]string
	tokenLimits     [2]int
}

func NewLimitedHandler(as AccessStorage) *LimitedHandler {
	_, filename, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(filename), "..")
	configs, err := LoadConfig(rootPath)
	if err != nil {
		panic(err)
	}

	period, _ := strconv.Atoi(configs.PeriodInSeconds)
	limit, _ := strconv.Atoi(configs.LimitOfRequests)
	block, _ := strconv.Atoi(configs.BlockInSeconds)
	tokenLimit1, _ := strconv.Atoi(configs.TokenLimit1)
	tokenLimit2, _ := strconv.Atoi(configs.TokenLimit2)
	return &LimitedHandler{
		as,
		AccessState{
			period,
			limit,
			time.Now(),
			block,
			time.Time{},
			false,
			[2]string{"HIGH", "LOW"},
			[2]int{tokenLimit1, tokenLimit2},
		},
	}
}

func (lh *LimitedHandler) LimitedFunc(w http.ResponseWriter, r *http.Request, f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	//ip, _, _ := strings.Cut(r.RemoteAddr, ":")
	ip := r.RemoteAddr

	maxRequests := 0
	if code := r.Header.Get("API_KEY"); code != "" {
		for i, token := range lh.tokens {
			if code == token {
				maxRequests = lh.tokenLimits[i]
				break
			}
		}
	}
	if maxRequests == 0 {
		maxRequests = lh.limitOfRequests
	}

	if (*lh).periodIsExpired() {
		lh.start = time.Now()
		lh.SetZero(ip)
	}

	lh.Increment(ip)
	switch {
	case !lh.isBlocked && lh.Get(ip) > maxRequests:
		lh.blockStart = time.Now()
		lh.isBlocked = true
		return blockedFunc
	case lh.isBlocked && (*lh).blockIsExpired():
		lh.SetZero(ip)
		lh.isBlocked = false
		return f
	case lh.isBlocked:
		return blockedFunc
	default:
		return f
	}
}

func blockedFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(429)
	w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
}

func (lh *LimitedHandler) periodIsExpired() bool {
	return time.Since(lh.start).Milliseconds() > int64(lh.periodInSeconds*1000)
}

func (lh *LimitedHandler) blockIsExpired() bool {
	return time.Since(lh.blockStart).Milliseconds() > int64(lh.blockInSeconds*1000)
}
