import{o as a,c as d,a as o,t as v,Q as _,p as h,k as p}from"./index-0818f5f9.js";const g={class:"modal-dialog modal-dialog-centered"},u={class:"modal-content"},f={class:"modal-header"},k={class:"jg-title"},j={class:"modal-body"},w={class:"modal-footer"},C={key:0,class:"modal-backdrop show"},y={__name:"dialog",props:["show","title"],emits:["close","save"],setup(i,{emit:l}){const e=i,t=l;function n(){t("save")}function c(){t("hide")}return(r,s)=>(a(),d("div",null,[o("div",{class:h(["modal fade",{"jg-modal show":e.show}])},[o("div",g,[o("div",u,[o("div",f,[o("div",k,v(e.title),1),o("div",{class:"jgicon jgicon-close",onClick:s[0]||(s[0]=m=>c())})]),o("div",j,[_(r.$slots,"default")]),o("div",w,[o("div",{class:"jg-button",onClick:s[1]||(s[1]=m=>n())},"保存")])])])],2),e.show?(a(),d("div",C)):p("",!0)]))}};export{y as _};
