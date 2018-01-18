var token_mark = "__token__"
var user_mark = "__user__"
var _phone = false
var wsUri = "ws://localhost:80/webSocket"
//var wsUri = "ws://47.92.121.219:80/webSocket"

CheckPhone()

// setCookie(key, value, day)
function setCookie(c_name, value, expiredays) {
    var exdate = new Date() 
	exdate.setDate(exdate.getDate() + expiredays) 
	document.cookie = c_name + "=" + escape(value) + ((expiredays == null) ? "": ";expires=" + exdate.toGMTString())
}

function getCookie(c_name) {
    if (document.cookie.length > 0) {
        c_start = document.cookie.indexOf(c_name + "=") 
		if (c_start != -1) {
            c_start = c_start + c_name.length + 1 
			c_end = document.cookie.indexOf(";", c_start) 
			if (c_end == -1) {
				c_end = document.cookie.length
			}
            return unescape(document.cookie.substring(c_start, c_end))
        }
    }
    return ""
}

// 判断是否为手机端
function CheckPhone() {
    if ((navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i))) {
       _phone = true
       console.log("手机端 浏览")
    }
}

function getRandom(min, max)
{
	return Math.floor(Math.random() * (max-1) + min)
}

// arg = "a=6&b=9"    fun(result)
function ajaxCall(arg,url,fun) {
    var xhr = ajaxFunction();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4) {
            if (xhr.status == 200 || xhr.status == 304) {
				if (fun != null) {				
	                fun(xhr.responseText)
				}
            }
        }
    }
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhr.send(arg);
}

function ajaxFunction() {
    var xmlHttp;
    try { // Firefox, Opera 8.0+, Safari
        xmlHttp = new XMLHttpRequest();
    } catch(e) {
        try { // Internet Explorer
            xmlHttp = new ActiveXObject("Msxml2.XMLHTTP");
        } catch(e) {
            try {
                xmlHttp = new ActiveXObject("Microsoft.XMLHTTP");
            } catch(e) {}
        }
    }

    return xmlHttp;
}