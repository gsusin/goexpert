package limiter

import (
	"net/http"
	"strconv"
	"time"
)

type LimitedHandler struct {
	access          map[string]int
	periodInSeconds int
	limitOfRequests int
	start           time.Time
	//Add to .env
	blockInSeconds int
	blockStart     time.Time
	isBlocked      bool
	//Add to .env
	tokens      [2]string
	tokenLimits [2]int
}

func NewLimitedHandler() *LimitedHandler {
	configs, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}

	period, _ := strconv.Atoi(configs.PeriodInSeconds)
	limit, _ := strconv.Atoi(configs.LimitOfRequests)
	return &LimitedHandler{
		make(map[string]int),
		period,
		limit,
		time.Now(),
		1,
		time.Time{},
		false,
		[2]string{"HIGH", "LOW"},
		[2]int{20, 15},
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
		lh.access[ip] = 0
	}

	lh.access[ip]++
	switch {
	case !lh.isBlocked && lh.access[ip] > maxRequests:
		lh.blockStart = time.Now()
		lh.isBlocked = true
		return blockedFunc
	case lh.isBlocked && (*lh).blockIsExpired():
		lh.access[ip] = 0
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
