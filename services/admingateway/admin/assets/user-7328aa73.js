import{r as g,O as m,c as v,a as s,h as j,w as p,u as o,q as u,g as _,C as h,o as y,t as d,d as f}from"./index-0818f5f9.js";const b={class:"mb-4 jg-as-box"},k=u('<div class="row jg-cb-row jg-as-header"><div class="jg-bk-form"><div class="row jg-asr-row"><div class="col-sm-4 jg-asr-col"><span class="jg-ars-memo">本月峰值 DAU（个）</span><div class="jg-ars-num">1,029</div></div><div class="col-sm-4 jg-asr-col"><span class="jg-ars-memo">昨日新注册用户数（个）</span><div class="jg-ars-num">2,000</div><div class="jg-ars-percent"> 较前一日<span class="jgicon jgicon-ac-up jg-ars-direction">10%</span></div></div><div class="col-sm-4 jg-asr-col"><span class="jg-ars-memo">截至昨日累计用户数（个）</span><div class="jg-ars-num">3,043</div></div></div></div></div>',1),x={class:"jg-as-tools"},w={class:"jg-as-tool"},V=s("div",{class:"jg-as-button"},"7 天",-1),C=s("div",{class:"jg-as-button jg-as-button-active"},"14 天",-1),D=s("div",{class:"jg-as-button"},"30 天",-1),S={class:"jg-as-date jgicon jgicon-date"},A=["onClick"],B={class:"row jg-as-body"},N={__name:"user",setup(E){let i=_(),e=g({range:{start:new Date,end:new Date}}),{asuserchat:r}=i.refs;m(()=>{let t=i.proxy.$echat.init(r);const a=["#5470C6","#EE6666"];let n={legend:{data:["日活","注册用户数"]},tooltip:{trigger:"none",axisPointer:{type:"cross"}},xAxis:{type:"category",boundaryGap:!1,data:["2015-1","2015-2","2015-3","2015-4","2015-5","2015-6","2015-7","2015-8","2015-9","2015-10","2015-11","2015-12"]},yAxis:{type:"value"},series:[{name:"日活",type:"line",smooth:!0,data:[26,59,90,264,287,707,1756,1822,487,188,60,23],lineStyle:{color:a[1]}},{name:"注册用户数",type:"line",smooth:!0,lineStyle:{color:a[0]},data:[39,59,111,187,483,692,2316,466,554,184,103,70]}]};t.setOption(n)});function l(t){return f.formatTime(new Date(t).getTime(),"yyyy-MM-dd")}return(t,a)=>{const n=h("VDatePicker");return y(),v("div",b,[k,s("div",x,[s("div",w,[V,C,D,s("div",S,[j(n,{modelValue:o(e).range,"onUpdate:modelValue":a[0]||(a[0]=c=>o(e).range=c),modelModifiers:{range:!0},class:"jg-as-date-picker"},{default:p(({togglePopover:c})=>[s("div",{class:"jg-as-date-content",onClick:c},d(l(o(e).range.start))+" 至 "+d(l(o(e).range.end)),9,A)]),_:1},8,["modelValue"])])])]),s("div",B,[s("div",{class:"jg-bk-form",ref_key:"asuserchat",ref:r},null,512)])])}}};export{N as default};
