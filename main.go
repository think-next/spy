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
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
:root{
  --bg:oklch(0.13 0.005 260);
  --surface:oklch(0.17 0.005 260);
  --card:oklch(0.20 0.006 260);
  --border:oklch(0.27 0.006 260);
  --fg:oklch(0.92 0 0);
  --muted:oklch(0.60 0 0);
  --dim:oklch(0.42 0 0);
  --accent:oklch(0.75 0.15 250);
  --accent2:oklch(0.70 0.18 160);
  --glow:oklch(0.60 0.20 250 / 0.15);
  --radius:0.75rem;
  --font-sans:'Inter',system-ui,-apple-system,sans-serif;
  --font-mono:'JetBrains Mono',ui-monospace,monospace;
}
*{margin:0;padding:0;box-sizing:border-box}
html{scroll-behavior:smooth}
body{font-family:var(--font-sans);background:var(--bg);color:var(--fg);line-height:1.6;-webkit-font-smoothing:antialiased}
a{color:var(--accent);text-decoration:none;transition:color .2s}
a:hover{color:oklch(0.85 0.18 250)}

.container{max-width:1100px;margin:0 auto;padding:0 1.5rem}

/* Hero */
header{text-align:center;padding:5rem 1rem 3.5rem;position:relative;overflow:hidden}
header::before{content:'';position:absolute;top:-40%;left:50%;transform:translateX(-50%);width:600px;height:600px;background:radial-gradient(circle,var(--glow) 0%,transparent 70%);pointer-events:none}
header h1{font-size:3rem;font-weight:700;letter-spacing:-0.03em;background:linear-gradient(135deg,var(--fg) 40%,var(--accent) 100%);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text;margin-bottom:.75rem;animation:fadeDown .6s ease-out}
header p{color:var(--muted);font-size:1.15rem;animation:fadeDown .6s ease-out .1s both}
.badge{display:inline-flex;align-items:center;gap:.4rem;margin-top:1rem;padding:.35rem .9rem;border-radius:999px;background:oklch(0.25 0.01 260);border:1px solid var(--border);font-size:.82rem;color:var(--muted);font-family:var(--font-mono);animation:fadeDown .6s ease-out .2s both}
.badge .dot{width:7px;height:7px;border-radius:50%;background:#3fb950;animation:pulse 2s ease-in-out infinite}

/* Sections */
section{margin-bottom:3rem;animation:fadeUp .6s ease-out .2s both}
h2{font-size:1.25rem;font-weight:600;margin-bottom:1.25rem;display:flex;align-items:center;gap:.5rem}
h2 .icon{font-size:1.1rem}

/* Grid */
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:.75rem}

/* Cards */
.card{display:block;background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:1.15rem 1.35rem;transition:transform .2s ease,border-color .2s ease,box-shadow .2s ease;position:relative;overflow:hidden}
.card::before{content:'';position:absolute;inset:0;background:linear-gradient(135deg,oklch(0.60 0.15 250/0.04) 0%,transparent 50%);opacity:0;transition:opacity .3s}
.card:hover{transform:translateY(-2px);border-color:oklch(0.40 0.08 250);box-shadow:0 8px 30px oklch(0 0 0/0.3),0 0 0 1px oklch(0.40 0.08 250/0.3)}
.card:hover::before{opacity:1}
.card h3{font-size:.95rem;font-weight:600;color:var(--fg);margin-bottom:.3rem;position:relative}
.card .desc{color:var(--muted);font-size:.85rem;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden;margin-bottom:.4rem}
.card .meta{display:flex;align-items:center;gap:.6rem;color:var(--dim);font-size:.78rem;font-family:var(--font-mono);position:relative}
.card .meta .star{color:oklch(0.80 0.16 85)}

/* Footer */
footer{text-align:center;padding:2.5rem 0;color:var(--dim);font-size:.82rem;border-top:1px solid var(--border)}
footer span{color:var(--muted)}

.empty{color:var(--dim);font-style:italic;padding:2rem;text-align:center}

/* Animations */
@keyframes fadeDown{from{opacity:0;transform:translateY(-15px)}to{opacity:1;transform:translateY(0)}}
@keyframes fadeUp{from{opacity:0;transform:translateY(15px)}to{opacity:1;transform:translateY(0)}}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}

/* Responsive */
@media(max-width:640px){
  header{padding:3.5rem 1rem 2.5rem}
  header h1{font-size:2rem}
  .grid{grid-template-columns:1fr}
  .container{padding:0 1rem}
}
</style>
</head>
<body>
<div class="container">
<header>
  <h1>think-next</h1>
  <p>Notes, code, and everything in between</p>
  <div style="display:flex;gap:.6rem;justify-content:center;margin-top:1rem;flex-wrap:wrap">
    <div class="badge"><span class="dot"></span> always building</div>
    <a href="/game" class="badge" style="text-decoration:none;cursor:pointer">🎮 追光者</a>
  </div>
</header>

<section>
  <h2><span class="icon">✦</span> Notion Notes</h2>
  <div class="grid">
`)
	if len(pages) == 0 {
		fmt.Fprint(w, `<div class="empty">No public notes yet</div>`)
	}
	for i, p := range pages {
		fmt.Fprintf(w, `    <a href="%s" class="card" style="animation:fadeUp .5s ease-out %dms both"><h3>%s</h3><div class="meta">%s</div></a>
`, html.EscapeString(p.URL), html.EscapeString(p.Title), p.UpdatedAt, (i+1)*80)
	}
	fmt.Fprint(w, `  </div>
</section>

<section style="animation:fadeUp .6s ease-out .3s both">
  <h2><span class="icon">⌘</span> GitHub Repos</h2>
  <div class="grid">
`)
	if len(repos) == 0 {
		fmt.Fprint(w, `<div class="empty">No public repos yet</div>`)
	}
	for i, r := range repos {
		desc := r.Description
		if desc == "" {
			desc = "No description"
		}
		lang := r.Language
		if lang == "" {
			lang = "—"
		}
		fmt.Fprintf(w, `    <a href="%s" class="card" style="animation:fadeUp .5s ease-out %dms both"><h3>%s</h3><div class="desc">%s</div><div class="meta"><span>%s</span><span class="star">★ %d</span></div></a>
`, html.EscapeString(r.HTMLURL), html.EscapeString(r.Name), html.EscapeString(desc), lang, r.Stars, (i+1)*80)
	}
	fmt.Fprint(w, `  </div>
</section>

<footer>
<nav style="display:flex;justify-content:center;gap:1.5rem;margin-bottom:1rem">
<a href="/" style="color:var(--muted);font-size:.9rem">🏠 Home</a>
<a href="/game" style="color:var(--muted);font-size:.9rem">🎮 Game</a>
</nav>
<span>Go</span> · Notion · GitHub
</footer>
</div>
</body>
</html>`)
}

func renderGame(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "game.html")
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
		http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		renderGame(w, r)
	})
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
