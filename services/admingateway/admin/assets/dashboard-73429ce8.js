import{r as g,c as a,a as s,F as h,b as f,u as w,A as y,d as n,m as _,e as E,o as c,t as i,g as k,E as C,f as x}from"./index-e9cb22d4.js";const T={class:"card-body db-layout"},A={class:"db-contanier"},P={class:"db-apps"},B={class:"db-app"},I={class:"db-app-header bd-bottom"},K={class:"db-app-name"},L={class:"db-app-body"},R={class:"row db-row"},V=["onClick"],$={class:"row db-row"},j={class:"col-md-6 db-right"},q={class:"row db-row"},D={class:"col-md-6 db-right"},F={class:"row db-row"},N={class:"col-md-6 db-right"},S={class:"row db-row"},W={class:"col-md-6 db-right"},O={class:"row db-row"},U={class:"col-md-6 db-right"},X={class:"db-app-footer bd-top"},Y=["onClick"],H={__name:"dashboard",setup(z){const b=k();let l=g({apps:[]});function u(e){console.log(e)}function m(){y.getList().then(({code:e,data:o})=>{let t=C.USER_TOKEN_EXPIRE;if(n.isEqual(e,t.code))return b.proxy.$toast({icon:"error",text:t.msg}),_.goLoginPage();let{items:r}=o,v=r.map(d=>(d.created_time=n.formatTime(d.created_time),d.ended_time=n.formatTime(d.ended_time),d.cur_user_count=n.numberWithCommas(d.cur_user_count),d.max_user_count=n.numberWithCommas(d.max_user_count),d.kind=n.isEqual(d.app_type,E.PRIVATE)?"私有云":"公有云",d));l.apps=v})}m();function p(e){x.setCurrent(e),_.goBasePage(e)}return(e,o)=>(c(),a("div",T,[s("div",A,[o[6]||(o[6]=s("div",{class:"db-header nav-underline-border"},[s("strong",{class:"db-header-title"},"我的应用")],-1)),s("div",P,[(c(!0),a(h,null,f(w(l).apps,t=>(c(),a("div",B,[s("div",I,[s("div",K,i(t.name),1)]),s("div",L,[s("div",R,[o[0]||(o[0]=s("div",{class:"col-md-6 db-left"},"App-Key",-1)),s("div",{class:"col-md-6 db-right jgicon jgicon-copy",onClick:r=>u(t.app_key)},i(t.app_key),9,V)]),s("div",$,[o[1]||(o[1]=s("div",{class:"col-md-6 db-left"},"创建时间",-1)),s("div",j,i(t.created_time),1)]),s("div",q,[o[2]||(o[2]=s("div",{class:"col-md-6 db-left"},"已注册用户",-1)),s("div",D,i(t.cur_user_count),1)]),s("div",F,[o[3]||(o[3]=s("div",{class:"col-md-6 db-left"},"授权用户总数",-1)),s("div",N,i(t.max_user_count),1)]),s("div",S,[o[4]||(o[4]=s("div",{class:"col-md-6 db-left"},"部署方式",-1)),s("div",W,i(t.kind),1)]),s("div",O,[o[5]||(o[5]=s("div",{class:"col-md-6 db-left"},"到期时间",-1)),s("div",U,i(t.ended_time),1)])]),s("div",X,[s("a",{class:"btn",onClick:r=>p(t)},"查看明细",8,Y)])]))),256))])])]))}};export{H as default};
