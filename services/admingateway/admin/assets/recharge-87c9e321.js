import{r,c as m,a,p as v,u as l,q as t,d as o,o as b}from"./index-e9cb22d4.js";const p={class:"mb-4"},g={class:"nav nav-underline-border"},j={class:"nav-item"},f={class:"nav-item"},u={class:"tab-content"},E={__name:"recharge",setup(w){let i={ONLINE:"online",PUBLIC:"public"},c=r({activeTab:i.ONLINE});function e(d){o.extend(c,{activeTab:d})}return(d,s)=>(b(),m("div",p,[a("ul",g,[a("li",j,[a("a",{class:v(["nav-link jgicon jgicon-money1",{active:l(o).isEqual(l(c).activeTab,l(i).ONLINE)}]),onClick:s[0]||(s[0]=n=>e(l(i).ONLINE))},"在线充值",2)]),a("li",f,[a("a",{class:v(["nav-link jgicon jgicon-money2",{active:l(o).isEqual(l(c).activeTab,l(i).PUBLIC)}]),onClick:s[1]||(s[1]=n=>e(l(i).PUBLIC))},"对公汇款",2)])]),a("div",u,[a("div",{class:v(["tab-pane p-3",{active:l(o).isEqual(l(c).activeTab,l(i).ONLINE)}])},s[2]||(s[2]=[t('<div class="row jg-ab-row"><label class="col-sm-1 col-form-label">充值账号</label><div class="col-sm-4"><div class="form-control-plaintext">10929298490</div></div><div class="col-sm-4"></div></div><div class="row jg-ab-row"><label class="col-sm-1 col-form-label">账户余额</label><div class="col-sm-4"><div class="form-control-plaintext">30293.00 元</div></div><div class="col-sm-4"></div></div><div class="row jg-ab-row"><label class="col-sm-1 col-form-label">充值金额</label><div class="col-sm-2"><div class="form-control-plaintext jg-ra-input"><input class="form-control" type="number"><span>元</span></div></div><div class="col-sm-4"></div></div><div class="row jg-ab-row"><label class="col-sm-1 col-form-label">支付类型</label><div class="col-sm-6"><div class="form-control-plaintext jg-ra-paies"><div class="jg-ra-pay jgicon jgicon-wechatpay">微信支付</div><div class="jg-ra-pay jgicon jgicon-alipay jg-ra-pay-selected">支付宝</div></div></div><div class="col-sm-4"></div></div><div class="row jg-ab-row"><label class="col-sm-1 col-form-label"></label><div class="col-sm-6"><div class="jg-ra-qr"></div></div><div class="col-sm-4"></div></div>',5)]),2),a("div",{class:v(["tab-pane p-3",{active:l(o).isEqual(l(c).activeTab,l(i).PUBLIC)}])},s[3]||(s[3]=[t('<div class="row jg-ab-row"><label class="col-sm-1 col-form-label">收款名称</label><div class="col-sm-4"><div class="form-control-plaintext">北京追风未来科技有限公司</div></div><div class="col-sm-4"></div></div><div class="row jg-ab-row"><label class="col-sm-1 col-form-label">收款账号</label><div class="col-sm-4"><div class="form-control-plaintext">110293840383</div></div><div class="col-sm-4"></div></div><div class="row jg-ab-row"><label class="col-sm-1 col-form-label">收款银行</label><div class="col-sm-4"><div class="form-control-plaintext">北京工商银行人民大学支行</div></div><div class="col-sm-4"></div></div>',3)]),2)])]))}};export{E as default};
