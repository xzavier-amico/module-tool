package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xzavier-amico/module-tool/internal/converter"
	"github.com/xzavier-amico/module-tool/internal/format"
)

type convertReq struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type convertResp struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

func Start(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(indexHTML))
	})
	http.HandleFunc("/api/formats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(format.All)
	})
	http.HandleFunc("/api/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", 405)
			return
		}
		var req convertReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(convertResp{Error: err.Error()})
			return
		}
		if !format.Valid(req.From) || !format.Valid(req.To) {
			json.NewEncoder(w).Encode(convertResp{Error: "invalid format"})
			return
		}
		res := converter.Convert(req.From, req.To, req.Content)
		json.NewEncoder(w).Encode(convertResp{Result: res})
	})
	fmt.Println("module-tool web UI listening on :" + port)
	http.ListenAndServe(":"+port, nil)
}

const indexHTML = `<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>module-tool converter</title>
<style>
body{font-family:system-ui;margin:20px}
textarea{width:100%;height:260px;font-family:monospace}
select,button{padding:6px;margin:4px}
.row{display:flex;gap:10px;align-items:center}
</style>
</head>
<body>
<h2>Module Converter</h2>
<div class="row">
<label>From:</label>
<select id="from"></select>
<label>To:</label>
<select id="to"></select>
<button onclick="convert()">Convert</button>
</div>
<h3>Source</h3>
<textarea id="src"></textarea>
<h3>Result</h3>
<textarea id="dst" readonly></textarea>
<script>
const formats=["surge","loon","qx","egern","stash"];
function init(){formats.forEach(f=>{from.add(new Option(f,f));to.add(new Option(f,f));});to.value="loon";}
async function convert(){
 const res=await fetch('/api/convert',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({from:from.value,to:to.value,content:src.value})});
 const j=await res.json();
 if(j.error){alert(j.error);return;}
 dst.value=j.result;
}
init();
</script>
</body>
</html>`
