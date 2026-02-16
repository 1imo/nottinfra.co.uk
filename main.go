package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	texttemplate "text/template"
	"time"

	"github.com/yuin/goldmark"
)

const baseURL = "https://nottinfra.co.uk"

type Article struct {
	Title          string
	Slug           string
	Description    string
	Keywords       string
	Body           string
	Created        string
	Updated        string
	DisplayUpdated string

	createdTime time.Time
	updatedTime time.Time
}

type IndexData struct {
	News    []Article
	BaseURL string
}

type ListingData struct {
	Items   []Article
	BaseURL string
}

type ArticleData struct {
	Article
	BaseURL    string
	DetailPath string // e.g. "articles", "updates", "signals"
}

func main() {
	articles, err := loadArticles("content/articles")
	if err != nil {
		log.Fatalf("loading articles: %v", err)
	}

	newsItems, err := loadArticles("content/updates")
	if err != nil {
		log.Fatalf("loading updates: %v", err)
	}

	signals, err := loadArticles("content/signals")
	if err != nil {
		log.Fatalf("loading signals: %v", err)
	}

	tmplRobots, err := template.ParseFiles("templates/robots.txt.tmpl")
	if err != nil {
		log.Fatalf("parsing robots template: %v", err)
	}
	tmplSitemap, err := texttemplate.ParseFiles("templates/sitemap.xml.tmpl")
	if err != nil {
		log.Fatalf("parsing sitemap template: %v", err)
	}

	articleMap := make(map[string]Article)
	for _, a := range articles {
		articleMap[a.Slug] = a
	}
	updateMap := make(map[string]Article)
	for _, u := range newsItems {
		updateMap[u.Slug] = u
	}
	signalMap := make(map[string]Article)
	for _, s := range signals {
		signalMap[s.Slug] = s
	}

	htmlFuncs := template.FuncMap{"raw": func(s string) template.HTML { return template.HTML(s) }}
	var tmplIndex, tmplArticles, tmplSignals, tmplArticle *template.Template
	tmplIndex, err = template.New("index.html.tmpl").Funcs(htmlFuncs).ParseFiles("templates/index.html.tmpl")
	if err != nil {
		log.Fatalf("parsing index template: %v", err)
	}
	tmplArticles, err = template.New("articles.html.tmpl").Funcs(htmlFuncs).ParseFiles("templates/articles.html.tmpl")
	if err != nil {
		log.Fatalf("parsing articles template: %v", err)
	}
	tmplSignals, err = template.New("signals.html.tmpl").Funcs(htmlFuncs).ParseFiles("templates/signals.html.tmpl")
	if err != nil {
		log.Fatalf("parsing signals template: %v", err)
	}
	tmplArticle, err = template.New("article.html.tmpl").Funcs(htmlFuncs).ParseFiles("templates/article.html.tmpl")
	if err != nil {
		log.Fatalf("parsing article template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplIndex.Execute(w, IndexData{News: newsItems, BaseURL: baseURL}); err != nil {
			log.Printf("executing index template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if err := tmplRobots.Execute(w, struct{ BaseURL string }{BaseURL: baseURL}); err != nil {
			log.Printf("executing robots template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	http.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		data := struct {
			BaseURL       string
			Articles      []Article
			Updates       []Article
			Signals       []Article
		}{BaseURL: baseURL, Articles: articles, Updates: newsItems, Signals: signals}
		if err := tmplSitemap.Execute(&buf, data); err != nil {
			log.Printf("executing sitemap template: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.Write(buf.Bytes())
	})

	http.HandleFunc("/updates", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", 302)
	})

	http.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplArticles.Execute(w, ListingData{Items: articles, BaseURL: baseURL}); err != nil {
			log.Printf("executing articles template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	http.HandleFunc("/signals", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplSignals.Execute(w, ListingData{Items: signals, BaseURL: baseURL}); err != nil {
			log.Printf("executing signals template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	http.HandleFunc("/updates/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/updates/")
		if slug == "" {
			http.Redirect(w, r, "/", 302)
			return
		}
		up, ok := updateMap[slug]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplArticle.Execute(w, ArticleData{Article: up, BaseURL: baseURL, DetailPath: "updates"}); err != nil {
			log.Printf("executing update detail template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	http.HandleFunc("/signals/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/signals/")
		if slug == "" {
			http.Redirect(w, r, "/signals", 302)
			return
		}
		sig, ok := signalMap[slug]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplArticle.Execute(w, ArticleData{Article: sig, BaseURL: baseURL, DetailPath: "signals"}); err != nil {
			log.Printf("executing signal detail template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	http.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/articles/")
		if slug == "" {
			http.Redirect(w, r, "/articles", 302)
			return
		}
		article, ok := articleMap[slug]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmplArticle.Execute(w, ArticleData{Article: article, BaseURL: baseURL, DetailPath: "articles"}); err != nil {
			log.Printf("executing article template: %v", err)
			http.Error(w, "Internal Server Error", 500)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Serving at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loadArticles(dir string) ([]Article, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var articles []Article
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || e.Name() == "README.md" {
			continue
		}
		path := filepath.Join(dir, e.Name())
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", path, err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		a, err := parseArticle(data)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}

		// Slug is derived from filename (e.g. my-post.md -> my-post)
		a.Slug = strings.TrimSuffix(e.Name(), ".md")

		// Fallback dates from file mod time if not set in frontmatter
		if a.Created == "" {
			a.Created = info.ModTime().Format("2006-01-02")
		}
		if a.Updated == "" {
			a.Updated = info.ModTime().Format("2006-01-02")
		}

		// Parse times for sorting; ignore parse errors
		if t, err := time.Parse("2006-01-02", a.Created); err == nil {
			a.createdTime = t
		}
		if t, err := time.Parse("2006-01-02", a.Updated); err == nil {
			a.updatedTime = t
			a.DisplayUpdated = formatHumanDate(t)
		}

		articles = append(articles, a)
	}

	// Sort newest updated first
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].updatedTime.After(articles[j].updatedTime)
	})

	return articles, nil
}

func parseArticle(data []byte) (Article, error) {
	parts := bytes.SplitN(data, []byte("\n---\n"), 2)
	if len(parts) != 2 {
		return Article{}, fmt.Errorf("missing frontmatter")
	}

	frontmatter := string(parts[0])
	bodyMarkdown := string(bytes.TrimSpace(parts[1]))

	// Simple YAML-like parsing for our fields
	var a Article
	lines := strings.Split(frontmatter, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "---" {
			continue
		}
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			val = strings.Trim(val, `"`)
			switch key {
			case "title":
				a.Title = val
			case "description":
				a.Description = val
			case "keywords":
				a.Keywords = val
			case "created":
				a.Created = val
			case "updated":
				a.Updated = val
			}
		}
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(bodyMarkdown), &buf); err != nil {
		return Article{}, fmt.Errorf("rendering markdown: %w", err)
	}
	a.Body = buf.String()

	return a, nil
}

func formatHumanDate(t time.Time) string {
	day := t.Day()
	var suffix string
	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	default:
		suffix = "th"
	}
	return fmt.Sprintf("%d%s %s %d", day, suffix, t.Format("Jan"), t.Year())
}
