package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.Println("Starting server on :80...")

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/notes", handleNotesPage)
	http.HandleFunc("/calendar", handleCalendarPage)

	// OA system — serve Vue SPA
	oaDist := filepath.Join("oa", "dist")
	http.HandleFunc("/oa", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/oa/", http.StatusMovedPermanently)
	})
	http.Handle("/oa/", http.StripPrefix("/oa/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		relPath := strings.TrimPrefix(r.URL.Path, "/")
		if relPath == "" {
			relPath = "index.html"
		}
		p := filepath.Join(oaDist, relPath)
		if _, err := os.Stat(p); err != nil {
			http.ServeFile(w, r, filepath.Join(oaDist, "index.html"))
			return
		}
		if fi, _ := os.Stat(p); fi.IsDir() {
			http.ServeFile(w, r, filepath.Join(p, "index.html"))
			return
		}
		http.ServeFile(w, r, p)
	})))

	log.Println("OA system mounted at /oa/")

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")
	fmt.Fprint(w, homePage)
}

func handleNotesPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "notes.html")
}

func handleCalendarPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "calendar.html")
}

const homePage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Spy</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
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
  --radius:0.75rem;
  --font-sans:'Inter',system-ui,-apple-system,sans-serif;
}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:var(--font-sans);background:var(--bg);color:var(--fg);line-height:1.6;-webkit-font-smoothing:antialiased}
a{color:var(--accent);text-decoration:none;transition:color .2s}
a:hover{color:oklch(0.85 0.18 250)}
.container{max-width:1100px;margin:0 auto;padding:0 1.5rem}
header{text-align:center;padding:5rem 1rem 3.5rem;position:relative;overflow:hidden}
header::before{content:'';position:absolute;top:-40%;left:50%;transform:translateX(-50%);width:600px;height:600px;background:radial-gradient(circle,oklch(0.60 0.20 250/0.15) 0%,transparent 70%);pointer-events:none}
header h1{font-size:3rem;font-weight:700;letter-spacing:-0.03em;background:linear-gradient(135deg,var(--fg) 40%,var(--accent) 100%);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text;margin-bottom:.75rem}
header p{color:var(--muted);font-size:1.15rem}
.badge{display:inline-flex;align-items:center;gap:.4rem;margin-top:1rem;padding:.35rem .9rem;border-radius:999px;background:oklch(0.25 0.01 260);border:1px solid var(--border);font-size:.82rem;color:var(--muted);cursor:pointer}
.badge:hover{color:var(--fg);border-color:oklch(0.40 0.08 250)}
footer{text-align:center;padding:2.5rem 0;color:var(--dim);font-size:.82rem;border-top:1px solid var(--border)}
@media(max-width:640px){
  header{padding:3.5rem 1rem 2.5rem}
  header h1{font-size:2rem}
  .container{padding:0 1rem}
}
</style>
</head>
<body>
<div class="container">
<header>
  <h1>🐌 Spy</h1>
  <p style="margin-top:1.2rem;font-style:italic;color:oklch(0.65 0.02 250);font-size:1.1rem;letter-spacing:0.04em">天地不仁，以万物为刍狗</p>
  <div style="margin-top:1.5rem;display:flex;gap:.6rem;justify-content:center;flex-wrap:wrap">
    <a href="/notes" class="badge">📝 笔记</a>
    <a href="/calendar" class="badge">📅 日历</a>
    <a href="/oa/" class="badge">🏢 OA系统</a>
  </div>
</header>
<footer><span>Go</span></footer>
</div>
</body>
</html>`
