package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"sort"
	"strings"
)

func renderAnalytics(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	stats := store.GetStats()
	statsJSON, _ := json.Marshal(stats)
	visitors := store.GetVisitors()
	visitorsJSON, _ := json.Marshal(visitors)

	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Spy Analytics - 用户画像分析</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
%s
</style>
</head>
<body>
<div class="container">
<div class="header">
  <h1>🔍 Spy Analytics</h1>
  <p>用户画像分析平台</p>
  <nav class="nav">
    <a href="/">🏠 Home</a>
    <a href="/analytics" class="active">📊 Analytics</a>
    <a href="/game">🎮 Game</a>
    <a href="/hot">🔥 热榜</a>
  </nav>
</div>

<div class="stats-grid">
  <div class="stat-card blue"><div class="stat-label">总访客</div><div class="stat-value" id="sv">0</div></div>
  <div class="stat-card green"><div class="stat-label">总会话</div><div class="stat-value" id="ss">0</div></div>
  <div class="stat-card pink"><div class="stat-label">总浏览量</div><div class="stat-value" id="sp">0</div></div>
  <div class="stat-card blue"><div class="stat-label">7日活跃</div><div class="stat-value" id="s7">0</div></div>
  <div class="stat-card green"><div class="stat-label">30日活跃</div><div class="stat-value" id="s30">0</div></div>
</div>

<div class="section">
  <div class="section-title">📊 用户频次分布</div>
  <div class="freq-grid" id="freqGrid"></div>
</div>

<div class="chart-row">
  <div class="chart-card"><h3>🌐 浏览器分布</h3><div class="bar-chart" id="browserChart"></div></div>
  <div class="chart-card"><h3>📱 设备类型</h3><div class="bar-chart" id="deviceChart"></div></div>
</div>
<div class="chart-row">
  <div class="chart-card"><h3>💻 操作系统</h3><div class="bar-chart" id="osChart"></div></div>
  <div class="chart-card"><h3>📥 嵌入代码</h3><div class="embed-box"><button class="copy-btn" id="copyBtn">📋</button>&lt;script async src="/track.js"&gt;&lt;/script&gt;</div><p style="font-size:.78rem;color:var(--dim);margin-top:.5rem">将此代码嵌入任意网站即可开始追踪</p></div>
</div>

<div class="section">
  <div class="section-title">👥 访客画像 <span style="font-size:.8rem;color:var(--dim);font-weight:400" id="vc"></span></div>
  <div class="visitor-list" id="vl"></div>
</div>
</div>

<div class="modal-overlay" id="pm">
  <div class="modal">
    <button class="modal-close" id="mc">&times;</button>
    <h2 id="mt"></h2>
    <p class="subtitle" id="ms"></p>
    <div class="profile-grid" id="mcon"></div>
  </div>
</div>

<div class="footer"><span>Spy Analytics</span> · Go · Zero Dependencies</div>

<script>
var S=%s;
var V=%s;

var FL={new:'新访客',returning:'回访用户',active:'活跃用户',loyal:'忠实用户'};
var FI={new:'🆕',returning:'🔄',active:'🔥',loyal:'💎'};
var EL={low:'路过型',medium:'浏览型',high:'高参与',power:'深度用户'};
var EI={low:'⏱️',medium:'👀',high:'📖',power:'⚡'};
var FC={new:'badge-new',returning:'badge-returning',active:'badge-active',loyal:'badge-loyal'};

function fmtN(n){return n>=1e3?(n/1e3).toFixed(1)+'k':String(n)}
function fmtD(s){if(s<60)return s+'s';if(s<3600)return Math.floor(s/60)+'m '+(s%%60)+'s';return Math.floor(s/3600)+'h '+Math.floor((s%%3600)/60)+'m'}
function tAgo(t){var d=Date.now()-new Date(t).getTime();if(d<6e4)return '刚刚';if(d<36e5)return Math.floor(d/6e4)+'分钟前';if(d<864e5)return Math.floor(d/36e5)+'小时前';return Math.floor(d/864e5)+'天前'}
function cV(n){if(n===1)return 'new';if(n<=3)return 'returning';if(n<=10)return 'active';return 'loyal'}
function cE(d){if(d<10)return 'low';if(d<60)return 'medium';if(d<180)return 'high';return 'power'}
function esc(s){var d=document.createElement('div');d.textContent=s;return d.innerHTML}

document.getElementById('sv').textContent=fmtN(S.totalVisitors||0);
document.getElementById('ss').textContent=fmtN(S.totalSessions||0);
document.getElementById('sp').textContent=fmtN(S.totalPageViews||0);
document.getElementById('s7').textContent=fmtN(S.recent7||0);
document.getElementById('s30').textContent=fmtN(S.recent30||0);

var fq=Object.entries(S.frequency||{});
document.getElementById('freqGrid').innerHTML=fq.map(function(e){
var i=['new','returning','active','loyal'].indexOf(e[0])+1;
return '<div class="freq-item c'+i+'"><div class="freq-num">'+e[1]+'</div><div class="freq-label">'+FI[e[0]]+' '+FL[e[0]]+'</div></div>';
}).join('');

function rBar(id,data,color){
var mx=Math.max.apply(null,Object.values(data).concat([1]));
var sr=Object.entries(data).sort(function(a,b){return b[1]-a[1]}).slice(0,8);
if(!sr.length){document.getElementById(id).innerHTML='<div class="empty">暂无数据</div>';return;}
document.getElementById(id).innerHTML=sr.map(function(e){
var pct=(e[1]/mx*100).toFixed(1);
return '<div class="bar-row"><span class="bar-label" title="'+esc(e[0])+'">'+esc(e[0])+'</span><div class="bar-track"><div class="bar-fill '+color+'" style="width:'+pct+'%%"><span class="bar-value">'+e[1]+'</span></div></div></div>';
}).join('');
}
rBar('browserChart',S.browsers||{},'blue');
rBar('deviceChart',S.devices||{},'green');
rBar('osChart',S.os||{},'pink');

document.getElementById('vc').textContent='('+V.length+'人)';
var vl=document.getElementById('vl');
if(!V.length){vl.innerHTML='<div class="empty">暂无访客数据，嵌入追踪代码后开始收集</div>';}
else{vl.innerHTML=V.slice(0,50).map(function(v){
var f=cV(v.visitCount),ls=v.sessions[v.sessions.length-1]||{};
var dv=ls.device?ls.device.deviceType:'Unknown';
var di=dv==='Mobile'||dv==='Tablet'?'📱':'💻';
var ad=v.sessions.length?Math.round(v.sessions.reduce(function(s,x){return s+x.duration},0)/v.sessions.length):0;
var en=cE(ad);
return '<div class="visitor-row" data-fp="'+esc(v.fingerprint)+'"><div class="visitor-avatar">'+di+'</div><div class="visitor-info"><h4>'+esc(v.fingerprint.slice(0,16))+'...</h4><p>'+FI[f]+' '+EI[en]+' '+di+' · '+v.visitCount+'次 · '+v.sessions.length+'会话</p></div><span class="visitor-badge '+FC[f]+'">'+FL[f]+'</span><div class="visitor-time">'+tAgo(v.lastSeen)+'</div></div>';
}).join('');}

vl.addEventListener('click',function(e){var row=e.target.closest('.visitor-row');if(row)showProfile(row.dataset.fp);});
document.getElementById('mc').addEventListener('click',closeModal);
document.getElementById('pm').addEventListener('click',function(e){if(e.target===e.currentTarget)closeModal();});
document.getElementById('copyBtn').addEventListener('click',function(){
navigator.clipboard.writeText('<script async src="/track.js"><\/script>').then(function(){
var b=document.getElementById('copyBtn');b.textContent='✓';setTimeout(function(){b.textContent='📋'},1000);
});});

function showProfile(fp){
var v=V.find(function(x){return x.fingerprint===fp});
if(!v)return;
var ls=v.sessions[v.sessions.length-1]||{};
var dev=ls.device||{};
var f=cV(v.visitCount);
var ad=v.sessions.length?Math.round(v.sessions.reduce(function(s,x){return s+x.duration},0)/v.sessions.length):0;
var ap=v.sessions.length?(v.sessions.reduce(function(s,x){return s+x.pages.length},0)/v.sessions.length).toFixed(1):0;
var en=cE(ad);

var hm={};
v.sessions.forEach(function(s){s.pages.forEach(function(p){var h=new Date(p.timestamp).getHours();hm[h]=(hm[h]||0)+1})});
var mH=Math.max.apply(null,Object.values(hm).concat([1]));
var hh='';
for(var h=0;h<24;h++){var c=hm[h]||0;var op=c>0?(0.2+c/mH*0.8).toFixed(2):0.05;
var bg=c>0?'oklch(0.55 0.18 250 / '+op+')':'oklch(0.15 0 0)';
hh+='<div class="hour-cell '+(c>0?'active':'')+'" style="background:'+bg+'" title="'+h+':00 - '+c+'次">'+h+'</div>';}

var pm={};
v.sessions.forEach(function(s){s.pages.forEach(function(p){pm[p.url]=(pm[p.url]||0)+1})});
var tp=Object.entries(pm).sort(function(a,b){return b[1]-a[1]}).slice(0,5);
var tph=tp.length?tp.map(function(e){return '<div style="display:flex;justify-content:space-between;padding:.3rem 0;border-bottom:1px solid var(--border);font-size:.85rem"><span style="color:var(--accent);word-break:break-all;max-width:80%%">'+esc(e[0])+'</span><span style="color:var(--muted);font-family:var(--font-mono)">'+e[1]+'</span></div>'}).join(''):'<span style="color:var(--dim)">暂无数据</span>';

var vh=v.sessions.slice(-5).reverse().map(function(s){return '<div style="padding:.3rem 0;border-bottom:1px solid var(--border);font-size:.82rem"><span style="color:var(--muted);font-family:var(--font-mono)">'+new Date(s.startTime).toLocaleString('zh-CN')+'</span> <span style="color:var(--dim)">· '+s.pages.length+'页 · '+fmtD(s.duration)+'</span></div>'}).join('');

document.getElementById('mt').textContent=FI[f]+' '+FL[f]+' · '+EI[en]+' '+EL[en];
document.getElementById('ms').textContent=fp+' · 首次访问: '+new Date(v.firstSeen).toLocaleString('zh-CN');
document.getElementById('mcon').innerHTML=
'<div class="profile-item"><div class="pi-label">指纹</div><div class="pi-value" style="font-family:var(--font-mono);font-size:.75rem;word-break:break-all">'+esc(fp)+'</div></div>'+
'<div class="profile-item"><div class="pi-label">访问次数</div><div class="pi-value">'+v.visitCount+' 次</div></div>'+
'<div class="profile-item"><div class="pi-label">会话数</div><div class="pi-value">'+v.sessions.length+' 个</div></div>'+
'<div class="profile-item"><div class="pi-label">平均停留</div><div class="pi-value">'+fmtD(ad)+'</div></div>'+
'<div class="profile-item"><div class="pi-label">平均浏览</div><div class="pi-value">'+ap+' 页</div></div>'+
'<div class="profile-item"><div class="pi-label">设备</div><div class="pi-value">'+esc(dev.deviceType||'—')+' · '+esc(dev.browser||'—')+'</div></div>'+
'<div class="profile-item"><div class="pi-label">系统</div><div class="pi-value">'+esc(dev.os||'—')+'</div></div>'+
'<div class="profile-item"><div class="pi-label">分辨率</div><div class="pi-value">'+esc(dev.screen||'—')+'</div></div>'+
'<div class="profile-item"><div class="pi-label">语言</div><div class="pi-value">'+esc(dev.language||'—')+'</div></div>'+
'<div class="profile-item"><div class="pi-label">IP</div><div class="pi-value">'+esc(ls.geo?ls.geo.ip:'')+'</div></div>'+
'<div class="profile-item full"><div class="pi-label">活跃时段分布</div><div class="hour-grid">'+hh+'</div></div>'+
'<div class="profile-item full"><div class="pi-label">常访页面 TOP5</div>'+tph+'</div>'+
'<div class="profile-item full"><div class="pi-label">访问历史</div>'+vh+'</div>';
document.getElementById('pm').classList.add('show');
}

function closeModal(){document.getElementById('pm').classList.remove('show')}

setInterval(function(){location.reload()},60000);
</script>
</body>
</html>`, analyticsCSS, string(statsJSON), string(visitorsJSON))
}

func renderBarChart(title string, data map[string]int, color string) string {
	if len(data) == 0 {
		return ""
	}
	var entries []struct {
		Key   string
		Value int
	}
	for k, v := range data {
		entries = append(entries, struct {
			Key   string
			Value int
		}{k, v})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Value > entries[j].Value })
	if len(entries) > 8 {
		entries = entries[:8]
	}
	max := entries[0].Value
	if max == 0 {
		max = 1
	}
	var h strings.Builder
	for _, e := range entries {
		pct := float64(e.Value) / float64(max) * 100
		h.WriteString(fmt.Sprintf(`<div class="bar-row"><span class="bar-label" title="%s">%s</span><div class="bar-track"><div class="bar-fill %s" style="width:%.1f%%"><span class="bar-value">%d</span></div></div></div>`,
			html.EscapeString(e.Key), html.EscapeString(e.Key), color, pct, e.Value))
	}
	return h.String()
}

const analyticsCSS = `:root{
  --bg:oklch(0.13 0.005 260);
  --surface:oklch(0.17 0.005 260);
  --card:oklch(0.20 0.006 260);
  --border:oklch(0.27 0.006 260);
  --fg:oklch(0.92 0 0);
  --muted:oklch(0.60 0 0);
  --dim:oklch(0.42 0 0);
  --accent:oklch(0.75 0.15 250);
  --accent2:oklch(0.70 0.18 160);
  --accent3:oklch(0.75 0.15 310);
  --radius:0.75rem;
  --font-sans:'Inter',system-ui,-apple-system,sans-serif;
  --font-mono:'JetBrains Mono',ui-monospace,monospace;
}
*{margin:0;padding:0;box-sizing:border-box}
html{scroll-behavior:smooth}
body{font-family:var(--font-sans);background:var(--bg);color:var(--fg);line-height:1.6;-webkit-font-smoothing:antialiased}
a{color:var(--accent);text-decoration:none}
.container{max-width:1200px;margin:0 auto;padding:0 1.5rem}
.header{text-align:center;padding:3rem 1rem 2rem}
.header h1{font-size:2.2rem;font-weight:700;letter-spacing:-0.03em;background:linear-gradient(135deg,var(--fg) 40%,var(--accent) 100%);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text}
.header p{color:var(--muted);font-size:1rem;margin-top:.4rem}
.nav{display:flex;gap:.5rem;justify-content:center;margin-top:1rem;flex-wrap:wrap}
.nav a{padding:.4rem 1rem;border-radius:999px;border:1px solid var(--border);color:var(--muted);font-size:.85rem;transition:all .2s}
.nav a:hover,.nav a.active{background:var(--card);color:var(--fg);border-color:oklch(0.40 0.08 250)}
.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:.75rem;margin:1.5rem 0}
.stat-card{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:1.25rem 1.5rem;position:relative;overflow:hidden}
.stat-card::after{content:'';position:absolute;top:0;right:0;width:60px;height:60px;border-radius:0 var(--radius) 0 60px;opacity:.08}
.stat-card.blue::after{background:var(--accent)}
.stat-card.green::after{background:var(--accent2)}
.stat-card.pink::after{background:var(--accent3)}
.stat-label{font-size:.8rem;color:var(--muted);margin-bottom:.3rem}
.stat-value{font-size:1.8rem;font-weight:700;font-family:var(--font-mono)}
.section{margin:2rem 0}
.section-title{font-size:1.15rem;font-weight:600;margin-bottom:1rem;display:flex;align-items:center;gap:.5rem}
.chart-row{display:grid;grid-template-columns:1fr 1fr;gap:.75rem}
@media(max-width:640px){.chart-row{grid-template-columns:1fr}}
.chart-card{background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:1.25rem 1.5rem}
.chart-card h3{font-size:.95rem;font-weight:600;margin-bottom:1rem}
.bar-chart{display:flex;flex-direction:column;gap:.5rem}
.bar-row{display:flex;align-items:center;gap:.75rem}
.bar-label{width:100px;font-size:.8rem;color:var(--muted);text-align:right;flex-shrink:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.bar-track{flex:1;height:24px;background:oklch(0.15 0 0);border-radius:4px;overflow:hidden;position:relative}
.bar-fill{height:100%;border-radius:4px;transition:width .6s ease;min-width:2px}
.bar-fill.blue{background:linear-gradient(90deg,oklch(0.55 0.18 250),oklch(0.65 0.15 250))}
.bar-fill.green{background:linear-gradient(90deg,oklch(0.55 0.18 160),oklch(0.65 0.15 160))}
.bar-fill.pink{background:linear-gradient(90deg,oklch(0.55 0.18 310),oklch(0.65 0.15 310))}
.bar-value{position:absolute;right:8px;top:50%;transform:translateY(-50%);font-size:.75rem;font-family:var(--font-mono);color:var(--fg);font-weight:600}
.freq-grid{display:grid;grid-template-columns:repeat(4,1fr);gap:.6rem;text-align:center}
.freq-item{background:oklch(0.15 0 0);border-radius:.5rem;padding:1rem .5rem}
.freq-item .freq-num{font-size:1.5rem;font-weight:700;font-family:var(--font-mono)}
.freq-item .freq-label{font-size:.75rem;color:var(--muted);margin-top:.2rem}
.freq-item.c1 .freq-num{color:oklch(0.75 0.18 250)}
.freq-item.c2 .freq-num{color:oklch(0.75 0.18 160)}
.freq-item.c3 .freq-num{color:oklch(0.75 0.18 310)}
.freq-item.c4 .freq-num{color:oklch(0.75 0.15 70)}
.visitor-list{display:flex;flex-direction:column;gap:.5rem}
.visitor-row{display:grid;grid-template-columns:auto 1fr auto auto;gap:1rem;align-items:center;background:var(--card);border:1px solid var(--border);border-radius:var(--radius);padding:.85rem 1.25rem;cursor:pointer;transition:all .2s}
.visitor-row:hover{border-color:oklch(0.40 0.08 250);transform:translateY(-1px)}
.visitor-avatar{width:36px;height:36px;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:1rem;background:oklch(0.25 0.01 260);border:1px solid var(--border)}
.visitor-info h4{font-size:.9rem;font-weight:600}
.visitor-info p{font-size:.78rem;color:var(--muted)}
.visitor-badge{font-size:.72rem;padding:.2rem .6rem;border-radius:999px;font-family:var(--font-mono)}
.badge-new{background:oklch(0.25 0.05 250);color:oklch(0.75 0.15 250)}
.badge-returning{background:oklch(0.25 0.05 160);color:oklch(0.70 0.18 160)}
.badge-active{background:oklch(0.25 0.05 310);color:oklch(0.75 0.15 310)}
.badge-loyal{background:oklch(0.25 0.05 70);color:oklch(0.75 0.12 70)}
.visitor-time{font-size:.75rem;color:var(--dim);font-family:var(--font-mono);text-align:right}
.modal-overlay{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center;padding:1rem}
.modal-overlay.show{display:flex}
.modal{background:var(--surface);border:1px solid var(--border);border-radius:1rem;width:100%;max-width:640px;max-height:85vh;overflow-y:auto;padding:2rem;position:relative}
.modal-close{position:absolute;top:1rem;right:1rem;background:none;border:none;color:var(--muted);font-size:1.5rem;cursor:pointer;line-height:1}
.modal-close:hover{color:var(--fg)}
.modal h2{font-size:1.3rem;font-weight:700;margin-bottom:.3rem}
.modal .subtitle{color:var(--muted);font-size:.85rem;margin-bottom:1.5rem}
.profile-grid{display:grid;grid-template-columns:1fr 1fr;gap:.75rem}
.profile-item{background:var(--card);border-radius:.5rem;padding:.75rem 1rem}
.profile-item .pi-label{font-size:.72rem;color:var(--muted);margin-bottom:.15rem;text-transform:uppercase;letter-spacing:.05em}
.profile-item .pi-value{font-size:.9rem;font-weight:600}
.profile-item.full{grid-column:1/-1}
.hour-grid{display:grid;grid-template-columns:repeat(12,1fr);gap:3px}
.hour-cell{aspect-ratio:1;border-radius:3px;display:flex;align-items:center;justify-content:center;font-size:.6rem;font-family:var(--font-mono);color:var(--muted);transition:all .3s}
.hour-cell.active{color:var(--fg);font-weight:600}
.embed-box{background:oklch(0.10 0 0);border:1px solid var(--border);border-radius:.5rem;padding:1rem 1.25rem;font-family:var(--font-mono);font-size:.8rem;color:var(--accent);position:relative;overflow-x:auto;white-space:pre-wrap;word-break:break-all}
.embed-box .copy-btn{position:absolute;top:.5rem;right:.5rem;background:var(--card);border:1px solid var(--border);border-radius:.375rem;padding:.3rem .6rem;color:var(--muted);cursor:pointer;font-size:.75rem}
.embed-box .copy-btn:hover{color:var(--fg);border-color:var(--accent)}
.footer{text-align:center;padding:2rem 0;color:var(--dim);font-size:.8rem;border-top:1px solid var(--border);margin-top:2rem}
.empty{text-align:center;color:var(--dim);padding:3rem;font-size:.9rem}
@media(max-width:640px){
  .container{padding:0 1rem}
  .stats-grid{grid-template-columns:repeat(2,1fr)}
  .freq-grid{grid-template-columns:repeat(2,1fr)}
  .visitor-row{grid-template-columns:auto 1fr auto;gap:.75rem}
  .visitor-time{display:none}
  .profile-grid{grid-template-columns:1fr}
  .header h1{font-size:1.6rem}
}`
