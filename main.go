package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"os"
	"net/http"
	"sync"
	"time"
)

var (
	githubUser   = "think-next"
	cacheTTL     = 10 * time.Minute
	maxNotion    = 10
	maxRepos     = 12
	notionToken  string
	githubToken  string
)

type notionPage struct {
	ID        string
	Title     string
	URL       string
	UpdatedAt string
}

type notionResponse struct {
	Results []struct {
		ID        string `json:"id"`
		URL       string `json:"url"`
		UpdatedAt string `json:"last_edited_time"`
		Properties map[string]struct {
			Title []struct {
				PlainText string `json:"plain_text"`
			} `json:"title"`
		} `json:"properties"`
	} `json:"results"`
}

type repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Stars       int    `json:"stargazers_count"`
	Language    string `json:"language"`
}

type cache struct {
	mu      sync.RWMutex
	pages   []notionPage
	repos   []repo
	expired bool
	timer   *time.Timer
}

var c = &cache{}

func refreshCache() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if pages, err := fetchNotion(); err == nil {
			c.mu.Lock()
			c.pages = pages
			c.mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		if repos, err := fetchGitHub(); err == nil {
			c.mu.Lock()
			c.repos = repos
			c.mu.Unlock()
		}
	}()

	wg.Wait()
	c.mu.Lock()
	c.expired = false
	c.timer = time.AfterFunc(cacheTTL, refreshCache)
	c.mu.Unlock()
}

func fetchNotion() ([]notionPage, error) {
	req, err := http.NewRequest("POST", "https://api.notion.com/v1/search", http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+notionToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var nr notionResponse
	if err := json.NewDecoder(resp.Body).Decode(&nr); err != nil {
		return nil, err
	}

	var pages []notionPage
	for _, r := range nr.Results {
		var title string
		for _, prop := range r.Properties {
			if len(prop.Title) > 0 {
				title = prop.Title[0].PlainText
				break
			}
		}
		if title == "" {
			title = "Untitled"
		}
		pages = append(pages, notionPage{
			ID:        r.ID,
			Title:     title,
			URL:       r.URL,
			UpdatedAt: r.UpdatedAt[:10],
		})
	}
	return pages, nil
}

func fetchGitHub() ([]repo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/users/"+githubUser+"/repos?per_page="+fmt.Sprint(maxRepos)+"&sort=updated&type=owner", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func renderPage(w http.ResponseWriter) {
	c.mu.RLock()
	pages := c.pages
	repos := c.repos
	c.mu.RUnlock()

	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>think-next</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#0d1117;color:#c9d1d9;line-height:1.6}
a{color:#58a6ff;text-decoration:none}a:hover{text-decoration:underline}
.container{max-width:960px;margin:0 auto;padding:2rem 1.5rem}
header{text-align:center;padding:3rem 0 2rem;border-bottom:1px solid #21262d;margin-bottom:2rem}
header h1{font-size:2.5rem;color:#f0f6fc;margin-bottom:.5rem}
header p{color:#8b949e;font-size:1.1rem}
section{margin-bottom:2.5rem}
h2{color:#f0f6fc;font-size:1.5rem;margin-bottom:1rem;padding-bottom:.5rem;border-bottom:1px solid #21262d}
.card{background:#161b22;border:1px solid #21262d;border-radius:8px;padding:1rem 1.25rem;margin-bottom:.75rem;transition:border-color .2s}
.card:hover{border-color:#58a6ff}
.card h3{font-size:1.05rem;margin-bottom:.25rem}
.card p{color:#8b949e;font-size:.9rem}
.card .meta{color:#484f58;font-size:.8rem;margin-top:.25rem}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:.75rem}
.lang{display:inline-block;width:10px;height:10px;border-radius:50%;margin-right:6px}
footer{text-align:center;padding:2rem 0;color:#484f58;font-size:.85rem;border-top:1px solid #21262d;margin-top:2rem}
.empty{color:#484f58;font-style:italic;padding:1rem}
@media(max-width:640px){header h1{font-size:1.8rem}.container{padding:1rem}}
</style>
</head>
<body>
<div class="container">
<header>
<h1>🚀 think-next</h1>
<p>个人空间 · 笔记与项目</p>
</header>

<section>
<h2>📝 Notion 笔记</h2>
<div class="grid">
`)
	if len(pages) == 0 {
		fmt.Fprint(w, `<div class="empty">暂无公开笔记</div>`)
	}
	for _, p := range pages {
		fmt.Fprintf(w, `<a href="%s" class="card"><h3>%s</h3><div class="meta">%s</div></a>
`, html.EscapeString(p.URL), html.EscapeString(p.Title), p.UpdatedAt)
	}
	fmt.Fprint(w, `</div></section>

<section>
<h2>📦 GitHub 仓库</h2>
<div class="grid">
`)
	if len(repos) == 0 {
		fmt.Fprint(w, `<div class="empty">暂无公开仓库</div>`)
	}
	for _, r := range repos {
		lang := r.Language
		if lang == "" {
			lang = "—"
		}
		fmt.Fprintf(w, `<a href="%s" class="card"><h3>⭐ %s</h3><p>%s</p><div class="meta">%s · %d ★</div></a>
`, html.EscapeString(r.HTMLURL), html.EscapeString(r.Name), html.EscapeString(r.Description), lang, r.Stars)
	}
	fmt.Fprint(w, `</div></section>

<footer>Built with Go · Powered by Notion &amp; GitHub</footer>
</div>
</body>
</html>`)
}

func main() {
	notionToken = os.Getenv("NOTION_TOKEN")
	githubToken = os.Getenv("GITHUB_TOKEN")
	log.Println("Starting server on :80...")
	refreshCache()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=300")
		renderPage(w)
	})
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
