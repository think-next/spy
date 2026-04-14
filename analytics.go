package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// ============================================
// 数据模型
// ============================================

type Visitor struct {
	Fingerprint string    `json:"fingerprint"`
	FirstSeen   time.Time `json:"firstSeen"`
	LastSeen    time.Time `json:"lastSeen"`
	VisitCount  int       `json:"visitCount"`
	Sessions    []Session `json:"sessions"`
}

type Session struct {
	ID        string       `json:"id"`
	StartTime time.Time    `json:"startTime"`
	EndTime   time.Time    `json:"endTime"`
	Duration  int          `json:"duration"` // seconds
	Pages     []PageView   `json:"pages"`
	Device    DeviceInfo   `json:"device"`
	Geo       GeoInfo      `json:"geo"`
	Referrer  string       `json:"referrer"`
}

type PageView struct {
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Timestamp time.Time `json:"timestamp"`
	Duration  int       `json:"duration"` // seconds on page
}

type DeviceInfo struct {
	Browser    string `json:"browser"`
	OS         string `json:"os"`
	DeviceType string `json:"deviceType"`
	Screen     string `json:"screen"`
	Language   string `json:"language"`
}

type GeoInfo struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	City    string `json:"city"`
}

type CollectPayload struct {
	Fingerprint string     `json:"fp"`
	SessionID   string     `json:"sid"`
	URL         string     `json:"url"`
	Title       string     `json:"title"`
	Referrer    string     `json:"ref"`
	Browser     string     `json:"browser"`
	OS          string     `json:"os"`
	DeviceType  string     `json:"device"`
	Screen      string     `json:"screen"`
	Language    string     `json:"lang"`
	Duration    int        `json:"dur"`
	Type        string     `json:"type"` // "pageview" | "heartbeat" | "leave"
}

type Profile struct {
	Visitor       Visitor    `json:"visitor"`
	Profile       Profiling  `json:"profile"`
}

type Profiling struct {
	// 基础画像
	Label       string `json:"label"`       // 用户标签
	Device      string `json:"device"`      // 主要设备
	Browser     string `json:"browser"`     // 主要浏览器
	OS          string `json:"os"`          // 主要系统
	Screen      string `json:"screen"`      // 屏幕分辨率
	Language    string `json:"language"`    // 语言偏好

	// 行为画像
	VisitCount  int      `json:"visitCount"`  // 访问次数
	SessionCount int     `json:"sessionCount"` // 会话数
	AvgDuration int      `json:"avgDuration"` // 平均停留(秒)
	AvgPages    int      `json:"avgPages"`    // 平均浏览页数
	ActiveHours []int    `json:"activeHours"` // 活跃时段(0-23)
	TopPages    []string `json:"topPages"`    // 常访页面

	// 地域画像
	IP      string `json:"ip"`
	Country string `json:"country"`
	City    string `json:"city"`

	// 频次分类
	Frequency string `json:"frequency"` // new/returning/active/loyal
	Engagement string `json:"engagement"` // low/medium/high/power
}

// ============================================
// 数据存储
// ============================================

type AnalyticsStore struct {
	mu       sync.RWMutex
	filePath string
	visitors map[string]*Visitor
}

var store *AnalyticsStore

func NewAnalyticsStore(path string) *AnalyticsStore {
	s := &AnalyticsStore{filePath: path, visitors: make(map[string]*Visitor)}
	s.load()
	return s
}

func (s *AnalyticsStore) load() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return
	}
	var visitors map[string]*Visitor
	if err := json.Unmarshal(data, &visitors); err != nil {
		return
	}
	s.visitors = visitors
}

func (s *AnalyticsStore) save() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := json.MarshalIndent(s.visitors, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(s.filePath, data, 0644)
}

func (s *AnalyticsStore) RecordEvent(payload CollectPayload, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	defer s.save()

	fp := payload.Fingerprint
	if fp == "" {
		fp = "unknown"
	}

	now := time.Now()
	ip := extractIP(r)

	v, ok := s.visitors[fp]
	if !ok {
		v = &Visitor{
			Fingerprint: fp,
			FirstSeen:   now,
			LastSeen:    now,
			VisitCount:  0,
			Sessions:    []Session{},
		}
		s.visitors[fp] = v
	}

	// 查找或创建 session
	var session *Session
	for i := range v.Sessions {
		if v.Sessions[i].ID == payload.SessionID {
			session = &v.Sessions[i]
			break
		}
	}

	if session == nil && payload.Type == "pageview" {
		newSession := Session{
			ID:        payload.SessionID,
			StartTime: now,
			EndTime:   now,
			Pages:     []PageView{},
			Device: DeviceInfo{
				Browser:    payload.Browser,
				OS:         payload.OS,
				DeviceType: payload.DeviceType,
				Screen:     payload.Screen,
				Language:   payload.Language,
			},
			Geo: GeoInfo{IP: ip},
		}
		v.Sessions = append(v.Sessions, newSession)
		session = &v.Sessions[len(v.Sessions)-1]
		v.VisitCount++
	}

	if session == nil {
		return
	}

	if payload.Type == "pageview" {
		session.Pages = append(session.Pages, PageView{
			URL:       payload.URL,
			Title:     payload.Title,
			Timestamp: now,
			Duration:  payload.Duration,
		})
		if payload.Referrer != "" {
			session.Referrer = payload.Referrer
		}
	}

	session.EndTime = now
	v.LastSeen = now

	// Calculate session duration
	for i := range session.Pages {
		if session.Pages[i].Duration > 0 {
			session.Duration += session.Pages[i].Duration
		}
	}
}

func (s *AnalyticsStore) GetVisitors() []Visitor {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]Visitor, 0, len(s.visitors))
	for _, v := range s.visitors {
		list = append(list, *v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].LastSeen.After(list[j].LastSeen)
	})
	return list
}

func (s *AnalyticsStore) GetProfile(fp string) *Profile {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.visitors[fp]
	if !ok {
		return nil
	}

	p := &Profile{Visitor: *v}

	// Compute profiling
	totalDuration := 0
	totalPages := 0
	hourMap := make(map[int]int)
	pageCountMap := make(map[string]int)
	deviceMap := make(map[string]int)
	browserMap := make(map[string]int)
	osMap := make(map[string]int)
	screenMap := make(map[string]int)
	langMap := make(map[string]int)
	lastIP := ""

	for _, sess := range v.Sessions {
		totalDuration += sess.Duration
		totalPages += len(sess.Pages)
		for _, pv := range sess.Pages {
			hourMap[pv.Timestamp.Hour()]++
			pageCountMap[pv.URL]++
		}
		if sess.Device.Browser != "" {
			browserMap[sess.Device.Browser]++
		}
		if sess.Device.OS != "" {
			osMap[sess.Device.OS]++
		}
		if sess.Device.DeviceType != "" {
			deviceMap[sess.Device.DeviceType]++
		}
		if sess.Device.Screen != "" {
			screenMap[sess.Device.Screen]++
		}
		if sess.Device.Language != "" {
			langMap[sess.Device.Language]++
		}
		if sess.Geo.IP != "" {
			lastIP = sess.Geo.IP
		}
	}

	sessCount := len(v.Sessions)
	if sessCount > 0 {
		p.Profile.AvgDuration = totalDuration / sessCount
		p.Profile.AvgPages = totalPages / sessCount
	}

	p.Profile.VisitCount = v.VisitCount
	p.Profile.SessionCount = sessCount
	p.Profile.Browser = topOf(browserMap)
	p.Profile.OS = topOf(osMap)
	p.Profile.Device = topOf(deviceMap)
	p.Profile.Screen = topOf(screenMap)
	p.Profile.Language = topOf(langMap)
	p.Profile.IP = lastIP

	// Active hours
	var hours []int
	for h, c := range hourMap {
		hours = append(hours, h)
		_ = c
	}
	sort.Ints(hours)
	p.Profile.ActiveHours = hours

	// Top pages
	type kv struct {
		k string
		v int
	}
	var pvs []kv
	for k, c := range pageCountMap {
		pvs = append(pvs, kv{k, c})
	}
	sort.Slice(pvs, func(i, j int) bool { return pvs[i].v > pvs[j].v })
	topN := 5
	if len(pvs) < topN {
		topN = len(pvs)
	}
	p.Profile.TopPages = make([]string, topN)
	for i := 0; i < topN; i++ {
		p.Profile.TopPages[i] = pvs[i].k
	}

	// Frequency classification
	switch {
	case v.VisitCount == 1:
		p.Profile.Frequency = "new"
	case v.VisitCount <= 3:
		p.Profile.Frequency = "returning"
	case v.VisitCount <= 10:
		p.Profile.Frequency = "active"
	default:
		p.Profile.Frequency = "loyal"
	}

	// Engagement classification
	avgDur := p.Profile.AvgDuration
	switch {
	case avgDur == 0 || avgDur < 10:
		p.Profile.Engagement = "low"
	case avgDur < 60:
		p.Profile.Engagement = "medium"
	case avgDur < 180:
		p.Profile.Engagement = "high"
	default:
		p.Profile.Engagement = "power"
	}

	// Generate label
	p.Profile.Label = generateLabel(p.Profile)

	return p
}

func (s *AnalyticsStore) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalVisitors := len(s.visitors)
	totalSessions := 0
	totalPageViews := 0
	browsers := make(map[string]int)
	devices := make(map[string]int)
	os := make(map[string]int)
	freq := map[string]int{"new": 0, "returning": 0, "active": 0, "loyal": 0}
	recent7 := 0
	recent30 := 0
	now := time.Now()

	for _, v := range s.visitors {
		sessCount := len(v.Sessions)
		totalSessions += sessCount
		for _, sess := range v.Sessions {
			totalPageViews += len(sess.Pages)
			if sess.Device.Browser != "" {
				browsers[sess.Device.Browser]++
			}
			if sess.Device.DeviceType != "" {
				devices[sess.Device.DeviceType]++
			}
			if sess.Device.OS != "" {
				os[sess.Device.OS]++
			}
		}
		switch {
		case v.VisitCount == 1:
			freq["new"]++
		case v.VisitCount <= 3:
			freq["returning"]++
		case v.VisitCount <= 10:
			freq["active"]++
		default:
			freq["loyal"]++
		}
		if v.LastSeen.After(now.Add(-7 * 24 * time.Hour)) {
			recent7++
		}
		if v.LastSeen.After(now.Add(-30 * 24 * time.Hour)) {
			recent30++
		}
	}

	return map[string]interface{}{
		"totalVisitors":  totalVisitors,
		"totalSessions":  totalSessions,
		"totalPageViews": totalPageViews,
		"browsers":       browsers,
		"devices":        devices,
		"os":             os,
		"frequency":      freq,
		"recent7":        recent7,
		"recent30":       recent30,
	}
}

func topOf(m map[string]int) string {
	max := 0
	top := ""
	for k, v := range m {
		if v > max {
			max = v
			top = k
		}
	}
	return top
}

func generateLabel(p Profiling) string {
	var tags []string

	switch p.Frequency {
	case "new":
		tags = append(tags, "🆕 新访客")
	case "loyal":
		tags = append(tags, "💎 忠实用户")
	case "active":
		tags = append(tags, "🔥 活跃用户")
	case "returning":
		tags = append(tags, "🔄 回访用户")
	}

	switch p.Engagement {
	case "power":
		tags = append(tags, "⚡ 深度用户")
	case "high":
		tags = append(tags, "📖 高参与")
	case "medium":
		tags = append(tags, "👀 浏览型")
	case "low":
		tags = append(tags, "⏱️ 路过型")
	}

	if p.Device == "Mobile" {
		tags = append(tags, "📱 移动端")
	} else if p.Device == "Tablet" {
		tags = append(tags, "📱 平板")
	} else {
		tags = append(tags, "💻 桌面端")
	}

	return strings.Join(tags, " · ")
}

func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
