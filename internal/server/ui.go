package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Conduit</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center;flex-wrap:wrap}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.intg{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.intg:hover{border-color:var(--leather)}
.intg-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.intg-name{font-size:.85rem;font-weight:700}
.intg-flow{font-size:.65rem;color:var(--cd);margin-top:.2rem;display:flex;align-items:center;gap:.3rem}
.flow-node{padding:.15rem .4rem;background:var(--bg);border:1px solid var(--bg3)}
.flow-arrow{color:var(--rust);font-size:.8rem}
.intg-meta{font-size:.55rem;color:var(--cm);margin-top:.35rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.intg-actions{display:flex;gap:.3rem;flex-shrink:0}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.badge.active{border-color:var(--green);color:var(--green)}.badge.paused{border-color:var(--gold);color:var(--gold)}.badge.error{border-color:var(--red);color:var(--red)}.badge.disabled{border-color:var(--cm);color:var(--cm)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.btn-sync{border-color:var(--green);color:var(--green)}.btn-sync:hover{background:var(--green);color:#fff}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(3,1fr)}.row2{grid-template-columns:1fr}.toolbar{flex-direction:column}.search{min-width:100%}}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> CONDUIT</h1><button class="btn btn-p" onclick="openForm()">+ New Integration</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search integrations..." oninput="render()">
<select class="filter-sel" id="status-filter" onchange="render()"><option value="">All Status</option><option value="active">Active</option><option value="paused">Paused</option><option value="error">Error</option><option value="disabled">Disabled</option></select>
</div>
<div id="list"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;

async function load(){var r=await fetch(A+'/integrations').then(function(r){return r.json()});items=r.integrations||[];renderStats();render();}

function renderStats(){
var total=items.length;
var active=items.filter(function(i){return i.status==='active'}).length;
var totalSyncs=items.reduce(function(s,i){return s+(i.sync_count||0)},0);
document.getElementById('stats').innerHTML=[
{l:'Integrations',v:total},{l:'Active',v:active,c:active>0?'var(--green)':''},{l:'Total Syncs',v:totalSyncs}
].map(function(x){return '<div class="st"><div class="st-v" style="'+(x.c?'color:'+x.c:'')+'">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('status-filter').value;
var f=items;
if(sf)f=f.filter(function(i){return i.status===sf});
if(q)f=f.filter(function(i){return(i.name||'').toLowerCase().includes(q)||(i.source_type||'').toLowerCase().includes(q)||(i.target_type||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No integrations configured. Create one to connect your data sources.</div>';return;}
var h='';f.forEach(function(i){
h+='<div class="intg"><div class="intg-top"><div style="flex:1">';
h+='<div class="intg-name">'+esc(i.name)+'</div>';
h+='<div class="intg-flow">';
h+='<span class="flow-node">'+esc(i.source_type||'source')+'</span>';
h+='<span class="flow-arrow">&#8594;</span>';
h+='<span class="flow-node">'+esc(i.target_type||'target')+'</span>';
h+='</div>';
h+='</div><div class="intg-actions">';
h+='<button class="btn btn-sm btn-sync" onclick="sync(''+i.id+'')">&#8635; Sync</button>';
h+='<button class="btn btn-sm" onclick="openEdit(''+i.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+i.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
h+='<div class="intg-meta">';
h+='<span class="badge '+i.status+'">'+i.status+'</span>';
h+='<span>'+i.sync_count+' syncs</span>';
if(i.last_sync_at)h+='<span>Last: '+ft(i.last_sync_at)+'</span>';
h+='<span>Created: '+ft(i.created_at)+'</span>';
h+='</div></div>';
});
document.getElementById('list').innerHTML=h;
}

async function sync(id){await fetch(A+'/integrations/'+id+'/sync',{method:'POST'}).catch(function(){});load();}
async function del(id){if(!confirm('Delete this integration?'))return;await fetch(A+'/integrations/'+id,{method:'DELETE'});load();}

function formHTML(intg){
var i=intg||{name:'',source_type:'',target_type:'',config:'{}',status:'active'};
var isEdit=!!intg;
var h='<h2>'+(isEdit?'EDIT INTEGRATION':'NEW INTEGRATION')+'</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'" placeholder="e.g. Postgres to S3 Backup"></div>';
h+='<div class="row2"><div class="fr"><label>Source Type</label><input id="f-source" value="'+esc(i.source_type)+'" placeholder="e.g. postgres, api, s3"></div>';
h+='<div class="fr"><label>Target Type</label><input id="f-target" value="'+esc(i.target_type)+'" placeholder="e.g. s3, webhook, file"></div></div>';
if(isEdit){h+='<div class="fr"><label>Status</label><select id="f-status">';
['active','paused','error','disabled'].forEach(function(s){h+='<option value="'+s+'"'+(i.status===s?' selected':'')+'>'+s.charAt(0).toUpperCase()+s.slice(1)+'</option>';});
h+='</select></div>';}
h+='<div class="fr"><label>Config (JSON)</label><textarea id="f-config" rows="4" style="font-family:var(--mono)">'+esc(i.config||'{}')+'</textarea></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Create')+'</button></div>';
return h;
}

function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');document.getElementById('f-name').focus();}
function openEdit(id){var item=null;for(var j=0;j<items.length;j++){if(items[j].id===id){item=items[j];break;}}if(!item)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(item);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}

async function submit(){
var name=document.getElementById('f-name').value.trim();
if(!name){alert('Name is required');return;}
var body={name:name,source_type:document.getElementById('f-source').value.trim(),target_type:document.getElementById('f-target').value.trim(),config:document.getElementById('f-config').value.trim()};
if(editId){var sel=document.getElementById('f-status');if(sel)body.status=sel.value;
await fetch(A+'/integrations/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/integrations',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();
}

function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
