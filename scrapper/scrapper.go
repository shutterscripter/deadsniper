package scrapper

import (
	"deadsniper/config"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

const (
	maxBodyCheck   = 8192
	maxBodySoft404 = 50000 
)

func isSoft404(body []byte) bool {
	s := string(body)
	if len(s) > maxBodyCheck {
		s = s[:maxBodyCheck]
	}
	lower := strings.ToLower(s)

	phrases := []string{"404 not found", "error 404", "404.", "that's an error", "page not found"}
	for _, p := range phrases {
		if strings.Contains(lower, p) {
			return true
		}
	}
	idx := strings.Index(lower, "404")
	if idx < 0 {
		return false
	}
	start := idx - 80
	if start < 0 {
		start = 0
	}
	end := idx + 120
	if end > len(lower) {
		end = len(lower)
	}
	near := lower[start:end]
	return strings.Contains(near, "not found")
}

func VisitUrl(url string) (deadLinks, blockedByBot []string, err error) {
	if url == "" {
		return nil, nil, fmt.Errorf("URL is required")
	}

	var mu sync.Mutex
	deadLinks = make([]string, 0)
	blockedByBot = make([]string, 0)

	threads := config.DefaultConfig.Threads
	if threads < 1 {
		threads = 1
	}

	c := colly.NewCollector()
	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(config.DefaultConfig.Timeout) * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: threads,
		RandomDelay: time.Duration(config.DefaultConfig.Delay*1000) * time.Millisecond,
	})

	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	})

	okStatus := map[int]bool{200: true, 201: true, 204: true, 301: true, 302: true, 304: true}
	c.OnResponse(func(r *colly.Response) {
		
		dead := false
		if r.StatusCode == 403 {
			mu.Lock()
			blockedByBot = append(blockedByBot, r.Request.URL.String())
			mu.Unlock()
		} else if r.StatusCode >= 400 || !okStatus[r.StatusCode] {
			dead = true
		} else if r.StatusCode == 200 && strings.Contains(strings.ToLower(r.Headers.Get("Content-Type")), "text/html") && len(r.Body) < maxBodySoft404 {
			dead = isSoft404(r.Body)
		}
		if dead {
			mu.Lock()
			deadLinks = append(deadLinks, r.Request.URL.String())
			mu.Unlock()
		}
	})


	c.OnError(func(r *colly.Response, err error) {
		if r == nil || r.Request == nil {
			return
		}
		
		if r.StatusCode == 403 {
			mu.Lock()
			blockedByBot = append(blockedByBot, r.Request.URL.String())
			mu.Unlock()
		} else if r.StatusCode >= 400 {
			mu.Lock()
			deadLinks = append(deadLinks, r.Request.URL.String())
			mu.Unlock()
		}
	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		if h.Request.URL.String() != url {
			return
		}
		link := h.Request.AbsoluteURL(h.Attr("href"))
		if link != "" {
			c.Visit(link)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Processing", r.URL.String())
		if config.DefaultConfig.Verbose {
			fmt.Println("Processing", r.URL.String())
		}
	})

	c.Visit(url)
	c.Wait()

	return deadLinks, blockedByBot, nil
}