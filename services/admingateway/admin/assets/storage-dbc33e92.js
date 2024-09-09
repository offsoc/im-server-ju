import{_ as q}from"./form-f22f1c7b.js";import{B as p,r as R,z as $,w as T,i as z,A as b,d as m,R as D,C as F,o,a as n,t as v,u as l,c as i,b as k,F as g,k as S,p as j,h as M,g as W}from"./index-e9cb22d4.js";const I={class:"mb-4 app-base jg-cb-box"},L={class:"row jg-cb-row jg-cb-header"},O={class:"jg-cb-form jg-file-form"},P={class:"jg-cb-form-item"},U={class:"col-sm-7"},G={class:"warn"},H={class:"jg-cb-form-item"},J={class:"col-sm-4 store-form"},Q={class:"form-check-inline store_inline"},X=["value","checked"],Y={class:"form-check-label"},Z={key:0,class:"form-check-inline store_inline"},ee={class:"col-sm-3 jg-cb-btns"},te={class:"md-6"},ne={class:"nav nav-underline-border",role:"tablist"},ae=["onClick"],se={class:"tab-content rounded-bottom"},ie={__name:"storage",setup(le){let w=z(),h=W(),{currentRoute:{_rawValue:{params:{app_key:_}}}}=w,d=[{type:"aws",name:"AWS",state:p({access_key:"",secret_key:"",endpoint:"",region:"",bucket:""}),fields:[{name:"access_key",label:"Access Key",type:"input_text"},{name:"secret_key",label:"Secret Key",type:"input_text"},{name:"endpoint",label:"Endpoint",type:"input_text"},{name:"region",label:"Region",type:"input_text"},{name:"bucket",label:"Buket Name",type:"input_text"}]},{type:"qiniu",name:"七牛云",state:p({access_key:"",secret_key:"",bucket:""}),fields:[{name:"access_key",label:"Access Key",type:"input_text"},{name:"secret_key",label:"Secret Key",type:"input_text"},{name:"domain",label:"Domain Name",type:"input_text"},{name:"bucket",label:"Bucket Name",type:"input_text"}]},{type:"oss",name:"阿里云",state:p({access_key:"",secret_key:"",endpoint:"",bucket:""}),fields:[{name:"access_key",label:"Access Key",type:"input_text"},{name:"secret_key",label:"Secret Key",type:"input_text"},{name:"endpoint",label:"Endpoint",type:"input_text"},{name:"bucket",label:"Bucket name",type:"input_text"}]}],e=R({checkedValue:{},current:"",channels:{oss:{name:"阿里云"},aws:{name:"AWS"},qiniu:{name:"七牛云"}},fileConfs:[],formTitle:"添加配置"});const c=p(d[0].type);function E(a){c.value=a.type,x()}function C(){return e.fileConfs.find(t=>t.channel==c.value)}function B(a){let t=C();b.setStorageConfig({app_key:_,channel:c.value,conf:a}).then(()=>{let r=c.value,s=0,u=e.channels[r].name,y={channel:c,enable:s,name:u};if(!t)return e.fileConfs.push(y),h.proxy.$toast({icon:"success",text:"添加成功，请选择存储类型后，保存设置"});h.proxy.$toast({icon:"success",text:"修改成功，请选择存储类型后，保存设置"})})}async function x(){const a=await b.getStorageConfig({app_key:_,channel:c.value});m.forEach(d,t=>{t.type===c.value&&a.data&&a.data.conf&&(t.state.value={...t.state.value,...a.data.conf})}),e.formTitle=C()?"修改配置":"添加配置"}x();function A(a){let t=a.target.value;e.checkedValue.value=t}function N(){b.getEnableStorage({app_key:_}).then(a=>{let{code:t,data:r={}}=a;if(m.isEqual(t,D.SUCCESS)){let{file_confs:s}=r,u={channel:""},y=m.map(s,f=>{let K=e.channels[f.channel]||{name:""};return f.enable&&(u=f),{...f,name:K.name}});e.current=u.channel,e.fileConfs=y}})}function V(){let a=e.checkedValue.value;b.setEnableStorage({app_key:_,channel:a}).then(()=>{e.current=a,h.proxy.$toast({icon:"success",text:"保存成功"})})}return N(),(a,t)=>{const r=F("n-flex");return o(),$(r,{vertical:""},{default:T(()=>[n("div",I,[n("div",L,[n("div",O,[n("div",P,[t[0]||(t[0]=n("label",{class:"col-sm-1 col-form-label jg-form-item-label"}," 正在使用存储 ",-1)),n("div",U,[n("span",G,v(l(e).channels[l(e).current]&&l(e).channels[l(e).current].name||"未设置"),1)])]),n("div",H,[t[2]||(t[2]=n("label",{class:"col-sm-1 col-form-label"}," 设置存储类型 ",-1)),n("div",J,[(o(!0),i(g,null,k(l(e).fileConfs,s=>(o(),i("div",Q,[n("input",{class:"form-check-input",type:"radio",name:"setting.name",value:s.channel,checked:s.enable,onChange:A},null,40,X),n("label",Y,v(s.name),1)]))),256)),l(e).fileConfs.length==0?(o(),i("div",Z,t[1]||(t[1]=[n("span",{class:"warn fs-12"},"请优先添加文件存储配置",-1)]))):S("",!0)]),n("div",ee,[l(e).fileConfs.length>0?(o(),i("div",{key:0,class:"jg-button jg-button-bg warn-bg",onClick:V},"保存设置")):S("",!0)])])])])]),n("div",te,[n("ul",ne,[(o(!0),i(g,null,k(l(d),s=>(o(),i("li",{class:"nav-item sw-nav-item",onClick:u=>E(s)},[n("a",{class:j(["nav-link jgicon jgicon-free",{active:l(m).isEqual(c.value,s.type)}])}," 添加 "+v(s.name)+" 配置",3)],8,ae))),256))]),n("div",se,[(o(!0),i(g,null,k(l(d),s=>(o(),i("div",{class:j(["tab-pane p-3",{active:l(m).isEqual(c.value,s.type)}])},[M(q,{fields:s.fields,state:s.state.value,"btn-type":1,title:l(e).formTitle,onSave:B},null,8,["fields","state","title"])],2))),256))])])]),_:1})}}};export{ie as default};
