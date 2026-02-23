package scrapper

import (
	"deadsniper/config"
	"fmt"
	"net"
	"net/http"
	neturl "net/url"
	"strconv"
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

func VisitUrl(rawURL string) (deadLinks, blockedByBot []string, err error) {
	if rawURL == "" {
		return nil, nil, fmt.Errorf("URL is required")
	}

	baseURL, err := neturl.Parse(rawURL)
	if err != nil || baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, nil, fmt.Errorf("invalid URL: %s", rawURL)
	}

	var mu sync.Mutex
	deadLinks = make([]string, 0)
	blockedByBot = make([]string, 0)
	deadSeen := make(map[string]struct{})
	blockedSeen := make(map[string]struct{})
	queued := make(map[string]struct{})

	threads := config.DefaultConfig.Threads
	if threads < 1 {
		threads = 1
	}

	maxDepth := config.DefaultConfig.MaxDepth
	if maxDepth < 0 {
		maxDepth = 0
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
		reqURL := r.Request.URL.String()

		dead := false
		if r.StatusCode == 403 {
			mu.Lock()
			if _, ok := blockedSeen[reqURL]; !ok {
				blockedSeen[reqURL] = struct{}{}
				blockedByBot = append(blockedByBot, reqURL)
			}
			mu.Unlock()
		} else if r.StatusCode >= 400 || !okStatus[r.StatusCode] {
			dead = true
		} else if r.StatusCode == 200 && strings.Contains(strings.ToLower(r.Headers.Get("Content-Type")), "text/html") && len(r.Body) < maxBodySoft404 {
			dead = isSoft404(r.Body)
		}
		if dead {
			mu.Lock()
			if _, ok := deadSeen[reqURL]; !ok {
				deadSeen[reqURL] = struct{}{}
				deadLinks = append(deadLinks, reqURL)
			}
			mu.Unlock()
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		if r == nil || r.Request == nil {
			return
		}

		reqURL := r.Request.URL.String()
		if r.StatusCode == 403 {
			mu.Lock()
			if _, ok := blockedSeen[reqURL]; !ok {
				blockedSeen[reqURL] = struct{}{}
				blockedByBot = append(blockedByBot, reqURL)
			}
			mu.Unlock()
		} else if r.StatusCode >= 400 {
			mu.Lock()
			if _, ok := deadSeen[reqURL]; !ok {
				deadSeen[reqURL] = struct{}{}
				deadLinks = append(deadLinks, reqURL)
			}
			mu.Unlock()
		} else {
			mu.Lock()
			if _, ok := deadSeen[reqURL]; !ok {
				deadSeen[reqURL] = struct{}{}
				deadLinks = append(deadLinks, reqURL)
			}
			mu.Unlock()
		}
	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		if h.Request.Ctx.Get("crawl") != "1" {
			return
		}

		currentDepth, _ := strconv.Atoi(h.Request.Ctx.Get("depth"))
		link := h.Request.AbsoluteURL(h.Attr("href"))
		if link == "" {
			return
		}

		parsedLink, perr := neturl.Parse(link)
		if perr != nil || parsedLink.Scheme == "" || parsedLink.Host == "" {
			return
		}
		parsedLink.Fragment = ""
		normalizedLink := parsedLink.String()

		mu.Lock()
		if _, ok := queued[normalizedLink]; ok {
			mu.Unlock()
			return
		}
		queued[normalizedLink] = struct{}{}
		mu.Unlock()

		nextCtx := colly.NewContext()
		crawlNext := false
		if config.DefaultConfig.Recursive {
			sameDomain := strings.EqualFold(parsedLink.Hostname(), baseURL.Hostname())
			if sameDomain && currentDepth < maxDepth {
				crawlNext = true
			}
		}
		if crawlNext {
			nextCtx.Put("crawl", "1")
			nextCtx.Put("depth", strconv.Itoa(currentDepth+1))
		} else {
			nextCtx.Put("crawl", "0")
			nextCtx.Put("depth", strconv.Itoa(currentDepth))
		}

		if err := c.Request("GET", normalizedLink, nil, nextCtx, nil); err != nil && config.DefaultConfig.Verbose {
			fmt.Println("Skipping", normalizedLink, "-", err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		if config.DefaultConfig.Verbose {
			fmt.Println("Processing", r.URL.String())
		}
	})

	rootCtx := colly.NewContext()
	rootCtx.Put("crawl", "1")
	rootCtx.Put("depth", "0")

	mu.Lock()
	queued[rawURL] = struct{}{}
	mu.Unlock()

	if err := c.Request("GET", rawURL, nil, rootCtx, nil); err != nil {
		return nil, nil, err
	}

	c.Wait()

	return deadLinks, blockedByBot, nil
}
