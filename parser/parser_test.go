package parser

import (
	"fmt"
	"github.com/vishen/go-jang/tokenizer"
	"testing"
)

func createNode(nt tokenizer.TokenType, value string) tokenizer.Token {
	return tokenizer.Token{Token_type: nt, Value: value}
}

func checkNodeTagAndText(node1, node2 *Node) bool {

	if len(node1.Children) != len(node2.Children) {
		fmt.Printf("[Error] Children length mismatch. %d != %d\n", len(node1.Children), len(node2.Children))
	}

	if node1.Tag != node2.Tag {
		fmt.Printf("[Error] Tag mismatch. %s != %s\n", node1.Tag, node2.Tag)

		return false
	}

	if node1.Text != node2.Text {
		fmt.Printf("[Error] Text mismatch. %s != %s\n", node1.Text, node2.Text)

		return false
	}

	for i, _ := range node1.Children {
		if !checkNodeTagAndText(node1.Children[i], node2.Children[i]) {
			return false
		}
	}

	return true
}

func TestParserSingleNode(t *testing.T) {
	test_case := []tokenizer.Token{
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.Text, "Hello World"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
	}

	wanted := Node{
		Tag:  "html",
		Text: "Hello World",
	}

	actual := Parser(test_case)

	if !checkNodeTagAndText(actual, &wanted) {
		t.Error("Failed")
	}
}

func TestParserSingleChildNode(t *testing.T) {
	test_case := []tokenizer.Token{
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "h1"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.Text, "Hello World"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "h1"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
	}

	wanted := Node{
		Tag: "html",

		Children: []*Node{
			&Node{
				Tag:  "h1",
				Text: "Hello World",
			},
		},
	}

	actual := Parser(test_case)

	if !checkNodeTagAndText(actual, &wanted) {
		t.Error("Failed")
	}
}

func TestParserMultipleChildrenNode(t *testing.T) {
	test_case := []tokenizer.Token{
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "h1"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.Text, "Hello World"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "h1"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "h2"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.Text, "World Hello"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "h2"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
	}

	wanted := Node{
		Tag: "html",

		Children: []*Node{
			&Node{
				Tag:  "h1",
				Text: "Hello World",
			},
			&Node{
				Tag:  "h2",
				Text: "World Hello",
			},
		},
	}

	actual := Parser(test_case)

	if !checkNodeTagAndText(actual, &wanted) {
		t.Error("Failed")
	}
}

func TestParserIgnoreUntilHTMLNode(t *testing.T) {
	test_case := []tokenizer.Token{
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "!doctype"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.Text, "Hello World"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
	}

	wanted := Node{
		Tag:  "html",
		Text: "Hello World",
	}

	actual := Parser(test_case)

	if !checkNodeTagAndText(actual, &wanted) {
		t.Error("Failed")
	}
}

func TestParserNonClosingElement(t *testing.T) {
	test_case := []tokenizer.Token{
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "br"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "p"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "h1"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.Text, "Hello World"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "h1"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "br"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "p"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.Tag, "br"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.CloseTag, ">"),
		createNode(tokenizer.OpenTag, "<"),
		createNode(tokenizer.ForwardSlash, "/"),
		createNode(tokenizer.Tag, "html"),
		createNode(tokenizer.CloseTag, ">"),
	}

	wanted := Node{
		Tag: "html",

		Children: []*Node{
			&Node{
				Tag: "br",
			},
			&Node{
				Tag: "p",
				Children: []*Node{
					&Node{
						Tag:  "h1",
						Text: "Hello World",
					},
					&Node{
						Tag: "br",
					},
				},
			},
			&Node{
				Tag: "br",
			},
		},
	}

	actual := Parser(test_case)

	if !checkNodeTagAndText(actual, &wanted) {
		t.Error("Failed")
	}
}

// func TestPaserTokenizerGoogleFrontpage(t *testing.T) {
// 	html := `<!doctype html><html itemscope="" itemtype="http://schema.org/WebPage" lang="en-AU"><head><meta content="/logos/doodles/2015/emmy-noethers-133rd-birthday-5681045017985024-hp.jpg" itemprop="image"><title>Google</title><script>(function(){window.google={kEI:'RqwPVa38Noif8QXvtIH4Ag',kEXPI:'3700249,3700362,3700366,4011559,4020347,4021338,4024681,4028717,4028932,4029043,4029054,4029141,4029515,4029803,4029815,4030036,4030154,4030415,4030440,4031303,4031321,4031390,4031392,4031430,4031608,4031627,4031643,4031706,8500393,8500572,8501182,10200083,10200682,10200980',authuser:0,kSID:'RqwPVa38Noif8QXvtIH4Ag'};google.kHL='en-AU';})();(function(){google.lc=[];google.li=0;google.getEI=function(a){for(var b;a&&(!a.getAttribute||!(b=a.getAttribute("eid")));)a=a.parentNode;return b||google.kEI};google.getLEI=function(a){for(var b=null;a&&(!a.getAttribute||!(b=a.getAttribute("leid")));)a=a.parentNode;return b};google.https=function(){return"https:"==window.location.protocol};google.ml=function(){};google.time=function(){return(new Date).getTime()};google.log=function(a,b,e,f,l){var d=new Image,h=google.lc,g=google.li,c="",m=google.ls||"";d.onerror=d.onload=d.onabort=function(){delete h[g]};h[g]=d;if(!e&&-1==b.search("&ei=")){var k=google.getEI(f),c="&ei="+k;-1==b.search("&lei=")&&((f=google.getLEI(f))?c+="&lei="+f:k!=google.kEI&&(c+="&lei="+google.kEI))}a=e||"/"+(l||"gen_204")+"?atyp=i&ct="+a+"&cad="+b+c+m+"&zx="+google.time();/^http:/i.test(a)&&google.https()?(google.ml(Error("a"),!1,{src:a,glmm:1}),delete h[g]):(d.src=a,google.li=g+1)};google.y={};google.x=function(a,b){google.y[a.id]=[a,b];return!1};google.load=function(a,b,e){google.x({id:a+n++},function(){google.load(a,b,e)})};var n=0;})();google.kCSI={};var _gjwl=location;function _gjuc(){var a=_gjwl.href.indexOf("#");if(0<=a&&(a=_gjwl.href.substring(a),0<a.indexOf("&q=")||0<=a.indexOf("#q="))&&(a=a.substring(1),-1==a.indexOf("#"))){for(var d=0;d<a.length;){var b=d;"&"==a.charAt(b)&&++b;var c=a.indexOf("&",b);-1==c&&(c=a.length);b=a.substring(b,c);if(0==b.indexOf("fp="))a=a.substring(0,d)+a.substring(c,a.length),c=d;else if("cad=h"==b)return 0;d=c}_gjwl.href="/search?"+a+"&cad=h";return 1}return 0}
// 	function _gjh(){!_gjuc()&&window.google&&google.x&&google.x({id:"GJH"},function(){google.nav&&google.nav.gjh&&google.nav.gjh()})};window._gjh&&_gjh();</script><style>#gbar,#guser{font-size:13px;padding-top:1px !important;}#gbar{height:22px}#guser{padding-bottom:7px !important;text-align:right}.gbh,.gbd{border-top:1px solid #c9d7f1;font-size:1px}.gbh{height:0;position:absolute;top:24px;width:100%}@media all{.gb1{height:22px;margin-right:.5em;vertical-align:top}#gbar{float:left}}a.gb1,a.gb4{text-decoration:underline !important}a.gb1,a.gb4{color:#00c !important}.gbi .gb4{color:#dd8e27 !important}.gbf .gb4{color:#900 !important}</style><style>body,td,a,p,.h{font-family:arial,sans-serif}body{margin:0;overflow-y:scroll}#gog{padding:3px 8px 0}td{line-height:.8em}.gac_m td{line-height:17px}form{margin-bottom:20px}.h{color:#36c}.q{color:#00c}.ts td{padding:0}.ts{border-collapse:collapse}em{font-weight:bold;font-style:normal}.lst{height:25px;width:496px}.gsfi,.lst{font:18px arial,sans-serif}.gsfs{font:17px arial,sans-serif}.ds{display:inline-box;display:inline-block;margin:3px 0 4px;margin-left:4px}input{font-family:inherit}a.gb1,a.gb2,a.gb3,a.gb4{color:#11c !important}body{background:#fff;color:black}a{color:#11c;text-decoration:none}a:hover,a:active{text-decoration:underline}.fl a{color:#36c}a:visited{color:#551a8b}a.gb1,a.gb4{text-decoration:underline}a.gb3:hover{text-decoration:none}#ghead a.gb2:hover{color:#fff !important}.sblc{padding-top:5px}.sblc a{display:block;margin:2px 0;margin-left:13px;font-size:11px}.lsbb{background:#eee;border:solid 1px;border-color:#ccc #999 #999 #ccc;height:30px}.lsbb{display:block}.ftl,#fll a{display:inline-block;margin:0 12px}.lsb{background:url(/images/srpr/nav_logo80.png) 0 -258px repeat-x;border:none;color:#000;cursor:pointer;height:30px;margin:0;outline:0;font:15px arial,sans-serif;vertical-align:top}.lsb:active{background:#ccc}.lst:focus{outline:none}</style><script></script></head><body bgcolor="#fff"><script>(function(){var src='/images/nav_logo176.png';var iesg=false;document.body.onload = function(){window.n && window.n();if (document.images){new Image().src=src;}
// 	if (!iesg){document.f&&document.f.q.focus();document.gbqf&&document.gbqf.q.focus();}
// 	}
// 	})();</script><div id="mngb">    <div id=gbar><nobr><b class=gb1>Search</b> <a class=gb1 href="http://www.google.com.au/imghp?hl=en&tab=wi">Images</a> <a class=gb1 href="http://maps.google.com.au/maps?hl=en&tab=wl">Maps</a> <a class=gb1 href="https://play.google.com/?hl=en&tab=w8">Play</a> <a class=gb1 href="http://www.youtube.com/?gl=AU&tab=w1">YouTube</a> <a class=gb1 href="http://news.google.com.au/nwshp?hl=en&tab=wn">News</a> <a class=gb1 href="https://mail.google.com/mail/?tab=wm">Gmail</a> <a class=gb1 href="https://drive.google.com/?tab=wo">Drive</a> <a class=gb1 style="text-decoration:none" href="http://www.google.com.au/intl/en/options/"><u>More</u> &raquo;</a></nobr></div><div id=guser width=100%><nobr><span id=gbn class=gbi></span><span id=gbf class=gbf></span><span id=gbe></span><a href="http://www.google.com.au/history/optout?hl=en" class=gb4>Web History</a> | <a  href="/preferences?hl=en" class=gb4>Settings</a> | <a target=_top id=gb_70 href="https://accounts.google.com/ServiceLogin?hl=en&continue=http://www.google.com.au/%3Fgfe_rd%3Dcr%26ei%3DRqwPVd-ZK7Hu8weLw4DIAw" class=gb4>Sign in</a></nobr></div><div class=gbh style=left:0></div><div class=gbh style=right:0></div>    </div><center><span id="prt" style="display:block"> <div><style>.pmoabs{background-color:#fff;border:1px solid #E5E5E5;color:#666;font-size:13px;padding-bottom:20px;position:absolute;right:2px;top:3px;z-index:986}#pmolnk{border-radius:2px;-moz-border-radius:2px;-webkit-border-radius:2px}.kd-button-submit{border:1px solid #3079ed;background-color:#4d90fe;background-image:-webkit-gradient(linear,left top,left bottom,from(#4d90fe),to(#4787ed));background-image:-webkit-linear-gradient(top,#4d90fe,#4787ed);background-image:-moz-linear-gradient(top,#4d90fe,#4787ed);background-image:-ms-linear-gradient(top,#4d90fe,#4787ed);background-image:-o-linear-gradient(top,#4d90fe,#4787ed);background-image:linear-gradient(top,#4d90fe,#4787ed);filter:progid:DXImageTransform.Microsoft.gradient(startColorStr='#4d90fe',EndColorStr='#4787ed')}.kd-button-submit:hover{border:1px solid #2f5bb7;background-color:#357ae8;background-image:-webkit-gradient(linear,left top,left bottom,from(#4d90fe),to(#357ae8));background-image:-webkit-linear-gradient(top,#4d90fe,#357ae8);background-image:-moz-linear-gradient(top,#4d90fe,#357ae8);background-image:-ms-linear-gradient(top,#4d90fe,#357ae8);background-image:-o-linear-gradient(top,#4d90fe,#357ae8);background-image:linear-gradient(top,#4d90fe,#357ae8);filter:progid:DXImageTransform.Microsoft.gradient(startColorStr='#4d90fe',EndColorStr='#357ae8')}.kd-button-submit:active{-webkit-box-shadow:inset 0 1px 2px rgba(0,0,0,0.3);-moz-box-shadow:inset 0 1px 2px rgba(0,0,0,0.3);box-shadow:inset 0 1px 2px rgba(0,0,0,0.3)}#pmolnk a{color:#fff;display:inline-block;font-weight:bold;padding:5px 20px;text-decoration:none;white-space:nowrap}.xbtn{color:#999;cursor:pointer;font-size:23px;line-height:5px;padding-top:5px}.padi{padding:0 8px 0 10px}.padt{padding:5px 20px 0 0;color:#444}.pads{text-align:left;max-width:200px}</style> <div class="pmoabs" id="pmocntr2" style="behavior:url(#default#userdata);display:none"> <table border="0"> <tr> <td colspan="2"> <div class="xbtn" onclick="google.promos&&google.promos.toast&& google.promos.toast.cpc()" style="float:right">&times;</div> </td> </tr> <tr> <td class="padi" rowspan="2"> <img src="/images/icons/product/chrome-48.png"> </td> <td class="pads">A faster way to browse the web</td> </tr> <tr> <td class="padt"> <div class="kd-button-submit" id="pmolnk"> <a href="/chrome/index.html?hl=en&amp;brand=CHNG&amp;utm_source=en-hpp&amp;utm_medium=hpp&amp;utm_campaign=en" onclick="google.promos&&google.promos.toast&& google.promos.toast.cl()">Install Google Chrome</a> </div> </td> </tr> </table> </div> <script type="text/javascript">(function(){var a={o:{}};a.o.qa=50;a.o.oa=10;a.o.Y="body";a.o.Oa=!0;a.o.Ra=function(b,c){var d=a.o.Ea();a.o.Ga(d,b,c);a.o.Sa(d);a.o.Oa&&a.o.Pa(d)};a.o.Sa=function(b){(b=a.o.$(b))&&0<b.forms.length&&b.forms[0].submit()};a.o.Ea=function(){var b=document.createElement("iframe");b.height=0;b.width=0;b.style.overflow="hidden";b.style.top=b.style.left="-100px";b.style.position="absolute";document.body.appendChild(b);return b};a.o.$=function(b){return b.contentDocument||b.contentWindow.document};a.o.Ga=function(b,c,d){b=a.o.$(b);b.open();d=["<",a.o.Y,'><form method=POST action="',d,'">'];for(var e in c)c.hasOwnProperty(e)&&d.push('<textarea name="',e,'">',c[e],"</textarea>");d.push("</form></",a.o.Y,">");b.write(d.join(""));b.close()};a.o.ba=function(b,c){c>a.o.oa?google&&google.ml&&google.ml(Error("ogcdr"),!1,{cause:"timeout"}):b.contentWindow?a.o.Qa(b):window.setTimeout(function(){a.o.ba(b,c+1)},a.o.qa)};a.o.Qa=function(b){document.body.removeChild(b)};a.o.Pa=function(b){a.o.Ca(b,"load",function(){a.o.ba(b,0)})};a.o.Ca=function(b,c,d){b.addEventListener?b.addEventListener(c,d,!1):b.attachEvent&&b.attachEvent("on"+c,d)};var m={Va:0,D:1,F:2,K:5};a.k={};a.k.M={ka:"i",J:"d",ma:"l"};a.k.A={N:"0",G:"1"};a.k.O={L:1,J:2,I:3};a.k.v={ea:"a",ia:"g",C:"c",ya:"u",xa:"t",N:"p",pa:"pid",ga:"eid",za:"at"};a.k.la=window.location.protocol+"//www.google.com/_/og/promos/";a.k.ha="g";a.k.Aa="z";a.k.S=function(b,c,d,e){var f=null;switch(c){case m.D:f=window.gbar.up.gpd(b,d,!0);break;case m.K:f=window.gbar.up.gcc(e)}return null==f?0:parseInt(f,10)};a.k.Ka=function(b,c,d){return c==m.D?null!=window.gbar.up.gpd(b,d,!0):!1};a.k.P=function(b,c,d,e,f,h,k,l){var g={};g[a.k.v.N]=b;g[a.k.v.ia]=c;g[a.k.v.ea]=d;g[a.k.v.za]=e;g[a.k.v.ga]=f;g[a.k.v.pa]=1;k&&(g[a.k.v.C]=k);l&&(g[a.k.v.ya]=l);if(h)g[a.k.v.xa]=h;else return google.ml(Error("knu"),!1,{cause:"Token is not found"}),null;return g};a.k.V=function(b,c,d){if(b){var e=c?a.k.ha:a.k.Aa;c&&d&&(e+="?authuser="+d);a.o.Ra(b,a.k.la+e)}};a.k.Fa=function(b,c,d,e,f,h,k){b=a.k.P(c,b,a.k.M.J,a.k.O.J,d,f,null,e);a.k.V(b,h,k)};a.k.Ia=function(b,c,d,e,f,h,k){b=a.k.P(c,b,a.k.M.ka,a.k.O.L,d,f,e,null);a.k.V(b,h,k)};a.k.Na=function(b,c,d,e,f,h,k,l,g,n){switch(c){case m.K:window.gbar.up.dpc(e,f);break;case m.D:window.gbar.up.spd(b,d,1,!0);break;case m.F:g=g||!1,l=l||"",h=h||0,k=k||a.k.A.G,n=n||0,a.k.Fa(e,h,k,f,l,g,n)}};a.k.La=function(b,c,d,e,f){return c==m.D?0<d&&a.k.S(b,c,e,f)>=d:!1};a.k.Ha=function(b,c,d,e,f,h,k,l,g,n){switch(c){case m.K:window.gbar.up.iic(e,f);break;case m.D:c=a.k.S(b,c,d,e)+1;window.gbar.up.spd(b,d,c.toString(),!0);break;case m.F:g=g||!1,l=l||"",h=h||0,k=k||a.k.A.N,n=n||0,a.k.Ia(e,h,k,1,l,g,n)}};a.k.Ma=function(b,c,d,e,f,h){b=a.k.P(c,b,a.k.M.ma,a.k.O.I,d,e,null,null);a.k.V(b,f,h)};var p={Ta:"a",Wa:"l",Ua:"c",fa:"d",I:"h",L:"i",gb:"n",G:"x",cb:"ma",eb:"mc",fb:"mi",Xa:"pa",Ya:"pc",$a:"pi",bb:"pn",ab:"px",Za:"pd",hb:"gpa",jb:"gpi",kb:"gpn",lb:"gpx",ib:"gpd"};a.i={};a.i.s={na:"hplogo",wa:"pmocntr2"};a.i.A={va:"0",G:"1",da:"2"};a.i.p=document.getElementById(a.i.s.wa);a.i.ja=16;a.i.ra=2;a.i.ta=20;google.promos=google.promos||{};google.promos.toast=google.promos.toast||{};a.i.H=function(b){a.i.p&&(a.i.p.style.display=b?"":"none",a.i.p.parentNode&&(a.i.p.parentNode.style.position=b?"relative":""))};a.i.ca=function(b){try{if(a.i.p&&b&&b.es&&b.es.m){var c=window.gbar.rtl(document.body)?"left":"right";a.i.p.style[c]=b.es.m-a.i.ja+a.i.ra+"px";a.i.p.style.top=a.i.ta+"px"}}catch(d){google.ml(d,!1,{cause:a.i.w+"_PT"})}};google.promos.toast.cl=function(){try{a.i.Q==m.F&&a.k.Ma(a.i.T,a.i.B,a.i.A.da,a.i.X,a.i.U,a.i.W),window.gbar.up.sl(a.i.B,a.i.w,p.I,a.i.R(),1)}catch(b){google.ml(b,!1,{cause:a.i.w+"_CL"})}};google.promos.toast.cpc=function(){try{a.i.p&&(a.i.H(!1),a.k.Na(a.i.p,a.i.Q,a.i.s.Z,a.i.T,a.i.Da,a.i.B,a.i.A.G,a.i.X,a.i.U,a.i.W),window.gbar.up.sl(a.i.B,a.i.w,p.fa,a.i.R(),1))}catch(b){google.ml(b,!1,{cause:a.i.w+"_CPC"})}};a.i.aa=function(){try{if(a.i.p){var b=276,c=document.getElementById(a.i.s.na);c&&(b=Math.max(b,c.offsetWidth));var d=parseInt(a.i.p.style.right,10)||0;a.i.p.style.visibility=2*(a.i.p.offsetWidth+d)+b>document.body.clientWidth?"hidden":""}}catch(e){google.ml(e,!1,{cause:a.i.w+"_HOSW"})}};a.i.Ba=function(){var b=["gpd","spd","aeh","sl"];if(!window.gbar||!window.gbar.up)return!1;for(var c=0,d;d=b[c];c++)if(!(d in window.gbar.up))return!1;return!0};a.i.Ja=function(){return a.i.p.currentStyle&&"absolute"!=a.i.p.currentStyle.position};google.promos.toast.init=function(b,c,d,e,f,h,k,l,g,n,q,r){try{a.i.Ba()?a.i.p&&(e==m.F&&!l==!g?(google.ml(Error("tku"),!1,{cause:"zwieback: "+g+", gaia: "+l}),a.i.H(!1)):(a.i.s.C="toast_count_"+c+(q?"_"+q:""),a.i.s.Z="toast_dp_"+c+(r?"_"+r:""),a.i.w=d,a.i.B=b,a.i.Q=e,a.i.T=c,a.i.Da=f,a.i.X=l?l:g,a.i.U=!!l,a.i.W=k,a.k.Ka(a.i.p,e,a.i.s.Z,c)||a.k.La(a.i.p,e,h,a.i.s.C,c)||a.i.Ja()?a.i.H(!1):(a.k.Ha(a.i.p,e,a.i.s.C,c,f,a.i.B,a.i.A.va,a.i.X,a.i.U,a.i.W),n||(window.gbar.up.aeh(window,"resize",a.i.aa),window.lol=
// 	a.i.aa,window.gbar.elr&&a.i.ca(window.gbar.elr()),window.gbar.elc&&window.gbar.elc(a.i.ca),a.i.H(!0)),window.gbar.up.sl(a.i.B,a.i.w,p.L,a.i.R())))):google.ml(Error("apa"),!1,{cause:a.i.w+"_INIT"})}catch(t){google.ml(t,!1,{cause:a.i.w+"_INIT"})}};a.i.R=function(){var b=a.k.S(a.i.p,a.i.Q,a.i.s.C,a.i.T);return"ic="+b};})();</script> <script type="text/javascript">(function(){var sourceWebappPromoID=144002;var sourceWebappGroupID=5;var payloadType=5;var cookieMaxAgeSec=2592000;var dismissalType=5;var impressionCap=25;var gaiaXsrfToken='';var zwbkXsrfToken='';var kansasDismissalEnabled=false;var sessionIndex=0;var invisible=false;window.gbar&&gbar.up&&gbar.up.r&&gbar.up.r(payloadType,function(show){if (show){google.promos.toast.init(sourceWebappPromoID,sourceWebappGroupID,payloadType,dismissalType,cookieMaxAgeSec,impressionCap,sessionIndex,gaiaXsrfToken,zwbkXsrfToken,invisible,'0612');}
// 	});})();</script> </div> </span><br clear="all" id="lgpd"><div id="lga"><a href="/search?site=&amp;ie=UTF-8&amp;q=Emmy+Noether&amp;oi=ddle&amp;ct=emmy-noethers-133rd-birthday-5681045017985024-hp&amp;hl=en&amp;sa=X&amp;ei=RqwPVa38Noif8QXvtIH4Ag&amp;ved=0CAMQNg"><img alt="Emmy Noether's 133rd Birthday" border="0" height="230" src="/logos/doodles/2015/emmy-noethers-133rd-birthday-5681045017985024-hp.jpg" title="Emmy Noether's 133rd Birthday" width="552" id="hplogo" onload="window.lol&&lol()"><br></a><br></div><form action="/search" name="f"><table cellpadding="0" cellspacing="0"><tr valign="top"><td width="25%">&nbsp;</td><td align="center" nowrap=""><input name="ie" value="ISO-8859-1" type="hidden"><input value="en-AU" name="hl" type="hidden"><input name="source" type="hidden" value="hp"><div class="ds" style="height:32px;margin:4px 0"><input style="color:#000;margin:0;padding:5px 8px 0 6px;vertical-align:top" autocomplete="off" class="lst" value="" title="Google Search" maxlength="2048" name="q" size="57"></div><br style="line-height:0"><span class="ds"><span class="lsbb"><input class="lsb" value="Google Search" name="btnG" type="submit"></span></span><span class="ds"><span class="lsbb"><input class="lsb" value="I'm Feeling Lucky" name="btnI" onclick="if(this.form.q.value)this.checked=1; else top.location='/doodles/'" type="submit"></span></span></td><td class="fl sblc" align="left" nowrap="" width="25%"><a href="/advanced_search?hl=en-AU&amp;authuser=0">Advanced search</a><a href="/language_tools?hl=en-AU&amp;authuser=0">Language tools</a></td></tr></table><input id="gbv" name="gbv" type="hidden" value="1"></form><div id="gac_scont"></div><div style="font-size:83%;min-height:3.5em"><br></div><span id="footer"><div style="font-size:10pt"><div style="margin:19px auto;text-align:center" id="fll"><a href="/intl/en/ads/">Advertising&nbsp;Programmes</a><a href="/services/">Business Solutions</a><a href="https://plus.google.com/115477067087672475993" rel="publisher">+Google</a><a href="/intl/en/about.html">About Google</a><a href="http://www.google.com.au/setprefdomain?prefdom=US&amp;sig=0_RmwyyqpPzRcmKTOhGt4963EmCsc%3D" id="fehl">Google.com</a></div></div><p style="color:#767676;font-size:8pt">&copy; 2015 - <a href="/intl/en/policies/privacy/">Privacy</a> - <a href="/intl/en/policies/terms/">Terms</a></p></span></center><div id="xjsd"></div><div id="xjsi" data-jiis="bp"><script>(function(){function c(b){window.setTimeout(function(){var a=document.createElement("script");a.src=b;document.getElementById("xjsd").appendChild(a)},0)}google.dljp=function(b,a){google.xjsu=b;c(a)};google.dlj=c;})();(function(){window.google.xjsrm=[];})();if(google.y)google.y.first=[];if(!google.xjs){window._=window._||{};window._._DumpException=function(e){throw e};if(google.timers&&google.timers.load.t){google.timers.load.t.xjsls=new Date().getTime();}google.dljp('/xjs/_/js/k\x3dxjs.hp.en_US.votPZMqb6rk.O/m\x3dsb_he,d/rt\x3dj/d\x3d1/t\x3dzcms/rs\x3dACT90oGWU4duK4Q5aR8xKQcV-eqzv5iDFw','/xjs/_/js/k\x3dxjs.hp.en_US.votPZMqb6rk.O/m\x3dsb_he,d/rt\x3dj/d\x3d1/t\x3dzcms/rs\x3dACT90oGWU4duK4Q5aR8xKQcV-eqzv5iDFw');google.xjs=1;}google.pmc={"sb_he":{"agen":true,"cgen":true,"client":"heirloom-hp","dh":true,"ds":"","exp":"msedr","fl":true,"host":"google.com.au","jam":0,"jsonp":true,"msgs":{"cibl":"Clear Search","dym":"Did you mean:","lcky":"I\u0026#39;m Feeling Lucky","lml":"Learn more","oskt":"Input tools","psrc":"This search was removed from your \u003Ca href=\"/history\"\u003EWeb History\u003C/a\u003E","psrl":"Remove","sbit":"Search by image","srch":"Google Search"},"ovr":{},"pq":"","refoq":true,"refpd":true,"rfs":[],"scd":10,"sce":5,"stok":"PpkDJ7nIPhdxAY_6iejWckoqe7o"},"d":{}};google.y.first.push(function(){if(google.med){google.med('init');google.initHistory();google.med('history');}});if(google.j&&google.j.en&&google.j.xi){window.setTimeout(google.j.xi,0);}
// 	</script></div></body></html>`

// 	tokens := Tokenizer(html)

// 	root := Parser(tokens)

// 	fmt.Println(root.children[0].parent)
// 	// WalkNode(root)
// }
