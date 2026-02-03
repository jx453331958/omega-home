package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type MetaResponse struct {
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

// FetchMeta fetches website metadata (favicon and description)
func FetchMeta(c *gin.Context) {
	targetURL := c.Query("url")
	if targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	// Validate URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	// Fetch the page with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create request: " + err.Error()})
		return
	}

	// Add browser-like headers to avoid 403 Forbidden
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch url: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch url: status " + resp.Status})
		return
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse html: " + err.Error()})
		return
	}

	meta := MetaResponse{}

	// Extract favicon
	// Priority: link[rel="icon"], link[rel="shortcut icon"], link[rel="apple-touch-icon"], fallback to /favicon.ico
	faviconURL := ""
	doc.Find("link[rel~='icon'], link[rel='shortcut icon'], link[rel='apple-touch-icon']").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if href, exists := s.Attr("href"); exists && href != "" {
			faviconURL = href
			return false // break
		}
		return true
	})

	// If no favicon found in meta tags, use default /favicon.ico
	if faviconURL == "" {
		faviconURL = "/favicon.ico"
	}

	// Convert relative URL to absolute
	if !strings.HasPrefix(faviconURL, "http://") && !strings.HasPrefix(faviconURL, "https://") {
		baseURL := parsedURL.Scheme + "://" + parsedURL.Host
		if strings.HasPrefix(faviconURL, "/") {
			meta.Icon = baseURL + faviconURL
		} else {
			meta.Icon = baseURL + "/" + faviconURL
		}
	} else {
		meta.Icon = faviconURL
	}

	// Extract description
	// Priority: meta[name="description"], meta[property="og:description"], <title> as fallback
	description := ""

	// Try meta description
	doc.Find("meta[name='description'], meta[property='og:description']").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if content, exists := s.Attr("content"); exists && content != "" {
			description = content
			return false // break
		}
		return true
	})

	// Fallback to title
	if description == "" {
		// Try og:title first
		doc.Find("meta[property='og:title']").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if content, exists := s.Attr("content"); exists && content != "" {
				description = content
				return false
			}
			return true
		})

		// Fallback to <title>
		if description == "" {
			title := doc.Find("title").First().Text()
			description = strings.TrimSpace(title)
		}
	}

	meta.Description = description

	c.JSON(http.StatusOK, meta)
}
