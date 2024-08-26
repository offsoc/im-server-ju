import{r as g,P as d,c as o,a as e,F as v,b as _,u as a,t as i,q as f,d as h,o as n,j as b,G as k}from"./index-0818f5f9.js";const y={class:"mb-4"},j={class:"card-body"},C=e("ul",{class:"nav nav-underline-border"},[e("li",{class:"nav-item"},[e("a",{class:"nav-link active jgicon jgicon-list"},"待开票记录")])],-1),w={class:"tab-content rounded-bottom"},N={class:"tab-pane p-3 active jg-wait-list"},E={class:"table jg-table"},A={class:"jg-td-c"},T=e("th",null,"充值主体",-1),x=e("th",{class:"jg-td-c"},"充值金额",-1),I=e("th",{class:"jg-td-c"},"可开票金额",-1),L=e("th",{class:"jg-td-c"},"充值方式",-1),O=e("th",{class:"jg-td-c"},"充值时间",-1),P={class:"jg-td-c"},V=["checked","onChange"],B={class:"jg-td-c"},R={class:"jg-td-c"},$={class:"jg-td-c"},q={class:"jg-td-c"},D={class:"tab-pane p-3 active jg-tab-line"},F={class:"form-check form-check-inline"},S=["value","onChange"],W={class:"form-check-label"},G={class:"form-check form-check-inline"},M=e("label",{class:"form-check-label jg-invoice-label"},"开票金额合计：",-1),U={class:"jg-invoice-label jg-invoice-amount"},Y=f('<div class="row g-2 jg-row"><div class="col-md"><div class="form-floating"><input class="form-control" type="email" placeholder="发票抬头"><label>发票抬头</label></div></div><div class="col-md"><div class="form-floating"><input class="form-control" type="email" placeholder="纳税人识别号"><label>纳税人识别号</label></div></div></div><div class="row g-2 jg-row"><div class="col-md"><div class="form-floating"><input class="form-control" type="email" placeholder="接收人"><label>接收人</label></div></div><div class="col-md"><div class="form-floating"><input class="form-control" type="email" placeholder="手机号"><label>手机号</label></div></div></div>',2),z={class:"row g-2 jg-row"},H={class:"col-md"},J={class:"form-floating"},K=e("input",{class:"form-control",type:"email",placeholder:"邮箱地址"},null,-1),Q=e("div",{class:"col-md"},[e("div",{class:"form-floating"},[e("input",{class:"form-control",type:"email",placeholder:"开票备注"}),e("label",null,"开票备注")])],-1),ee={__name:"draw",setup(X){let l=g({records:[{id:1,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"},{id:2,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"},{id:3,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"},{id:4,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"},{id:5,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"},{id:6,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"},{id:7,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"},{id:8,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"},{id:9,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"},{id:10,title:"北京未来科技有限公司",checked:!1,recharge:25e3,invoice:25e3,type:"微信",time:"2023-10-10 23:03"}],radios:[{name:"type",value:d.ONLINE,label:"增值税普通发票（电子）"},{name:"type",value:d.PAPER,label:"增值税普通发票（纸质）"}],total:0,invoiceType:d.ONLINE,checkAll:!1});function m(s){l.invoiceType=s}function p(s){l.records.map(t=>(h.isEqual(t.id,s.id)&&(t.checked=!s.checked),t));let c=0;h.each(l.records,t=>{t.checked&&(c+=t.invoice)}),l.total=h.numberWithCommas(c)}function u(){l.checkAll=!l.checkAll;let s=0;l.records.map(c=>(c.checked=l.checkAll,l.checkAll&&(s+=c.invoice),c)),l.total=h.numberWithCommas(s)}return(s,c)=>(n(),o("div",y,[e("div",j,[C,e("div",w,[e("div",N,[e("table",E,[e("thead",null,[e("tr",null,[e("th",A,[e("input",{class:"form-check-input jg-td-select",type:"checkbox",onChange:c[0]||(c[0]=t=>u())},null,32)]),T,x,I,L,O])]),e("tbody",null,[(n(!0),o(v,null,_(a(l).records,t=>(n(),o("tr",null,[e("td",P,[e("input",{class:"form-check-input jg-td-select",checked:t.checked,type:"checkbox",onChange:r=>p(t)},null,40,V)]),e("td",null,i(t.title),1),e("td",B,i(t.recharge),1),e("td",R,i(t.invoice),1),e("td",$,i(t.type),1),e("td",q,i(t.time),1)]))),256))])])]),e("div",D,[(n(!0),o(v,null,_(a(l).radios,t=>(n(),o("div",F,[b(e("input",{class:"form-check-input",type:"radio",name:"radio.name",value:t.value,"onUpdate:modelValue":c[1]||(c[1]=r=>a(l).invoiceType=r),onChange:r=>m(t.value)},null,40,S),[[k,a(l).invoiceType]]),e("label",W,i(t.label),1)]))),256)),e("div",G,[M,e("label",U,i(a(l).total),1)])]),Y,e("div",z,[e("div",H,[e("div",J,[K,e("label",null,i(a(l).invoiceType==a(d).ONLINE?"邮箱地址":"收件地址"),1)])]),Q])])])]))}};export{ee as default};
