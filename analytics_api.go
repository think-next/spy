package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleCollect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var payload CollectPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid payload"})
		return
	}

	store.RecordEvent(payload, r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleTrackJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write([]byte(trackJSSource))
}

func handleAPIProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fp := r.URL.Path[len("/api/profile/"):]
	if fp == "" {
		http.Error(w, "missing fingerprint", http.StatusBadRequest)
		return
	}

	profile := store.GetProfile(fp)
	if profile == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func handleAPIStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(store.GetStats())
}

func handleAPIVisitors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(store.GetVisitors())
}

const trackJSSource = `var _spy_endpoint=location.origin+"/api/collect";
(function(){
  var sid=Math.random().toString(36).substr(2,12);
  var fp=null;
  var lastPage='';

  function getFP(){
    try{
      var c=navigator.userAgent||''+navigator.language||''+screen.width+'x'+screen.height;
      var h=[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0];
      for(var i=0;i<c.length;i++){h[i%16]=(h[i%16]+c.charCodeAt(i))&0xff}
      var r='';for(var i=0;i<16;i++){r+=h[i].toString(16).padStart(2,'0')}
      return r;
    }catch(e){return 'err'}
  }

  function detectBrowser(){
    var ua=navigator.userAgent;
    if(ua.indexOf('Edg')>-1)return 'Edge';
    if(ua.indexOf('Chrome')>-1)return 'Chrome';
    if(ua.indexOf('Firefox')>-1)return 'Firefox';
    if(ua.indexOf('Safari')>-1)return 'Safari';
    return 'Other';
  }

  function detectOS(){
    var ua=navigator.userAgent;
    if(ua.indexOf('Windows')>-1)return 'Windows';
    if(ua.indexOf('Mac')>-1)return 'macOS';
    if(ua.indexOf('Linux')>-1)return 'Linux';
    if(ua.indexOf('Android')>-1)return 'Android';
    if(/iPhone|iPad|iPod/.test(ua))return 'iOS';
    return 'Other';
  }

  function detectDevice(){
    var ua=navigator.userAgent;
    if(/Mobi|Android.*Mobile|iPhone|iPod/.test(ua))return 'Mobile';
    if(/iPad|Android(?!.*Mobile)/.test(ua))return 'Tablet';
    return 'Desktop';
  }

  function send(type,opts){
    if(!fp)return;
    var data={
      fp:fp,sid:sid,type:type,
      url:location.href,title:document.title||'',
      ref:document.referrer||'',
      browser:detectBrowser(),os:detectOS(),
      device:detectDevice(),
      screen:screen.width+'x'+screen.height,
      lang:navigator.language||''
    };
    if(opts){for(var k in opts)data[k]=opts[k]}
    var r=new XMLHttpRequest();
    r.open('POST',_spy_endpoint,true);
    r.setRequestHeader('Content-Type','application/json');
    r.send(JSON.stringify(data));
  }

  function trackPage(){
    var url=location.href;
    if(url===lastPage)return;
    lastPage=url;
    send('pageview');
  }

  fp=getFP();
  if(!fp)return;

  trackPage();

  var interval=setInterval(function(){send('heartbeat')},30000);
  var timer=setInterval(function(){trackPage()},2000);

  var oldPush=history.pushState;
  var oldReplace=history.replaceState;
  history.pushState=function(){
    oldPush.apply(this,arguments);
    setTimeout(trackPage,100);
  };
  history.replaceState=function(){
    oldReplace.apply(this,arguments);
    setTimeout(trackPage,100);
  };
  window.addEventListener('popstate',function(){setTimeout(trackPage,100)});

  var unload=false;
  window.addEventListener('beforeunload',function(){
    if(unload)return;
    unload=true;
    clearInterval(interval);
    clearInterval(timer);
    send('leave');
  });

  document.addEventListener('visibilitychange',function(){
    if(document.visibilityState==='visible'){trackPage();}
  });
})();`

func init() {
	log.Println("[analytics] module initialized")
}
