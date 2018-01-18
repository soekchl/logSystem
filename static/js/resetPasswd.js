
function init() {
	addKeyPress("passwd")
	addKeyPress("phone")
	addKeyPress("code")
}

var checked = 0

function Register() {
	if (checkData(true)) {
		return
	}
	let c = document.getElementById("code-input")
	if (c.value.length != 6) {
		noticeErrSpan("code", "请输入正确的验证码")
		return
	}

	let arg = "phone="+document.getElementById("phone-input").value
	arg += "&passwd="+document.getElementById("passwd-input").value
	arg += "&code="+document.getElementById("code-input").value
//	console.log(arg)
	ajaxCall(arg, "/api/resetPasswd", function (text) {
		if (text == "88") {
			noticeErrSpan("phone", "验证码发送失败，请检查手机号码")
		} else if (text == "0") {
				alert("重置密码成功！")
	            location.replace("/") // 跳转到主页
	          } else if (text == "1") {
	              noticeErrSpan("passwd", "密码长度 6-32 ")	
	          } else if (text == "5") {
	              noticeErrSpan("phone", "请输入正确的手机号码")
	          } else if (text == "2") {
	              noticeErrSpan("phone", "手机号码已注册")	
	          } else if (text == "9") {
	              noticeErrSpan("code", "验证码错误")
	          }
	})
}

// 检查帐号是否可以注册
function sendPhoneCode() {
    // 检查
    if ( !checked && !checkData(false)) {
    	// 验证码获取
        let arg = "phone="+document.getElementById("phone-input").value + "&reset=1"
        ajaxCall(arg, "/api/phoneCode", function (text) {
            if (text == "88") {
                noticeErrSpan("phone", "验证码发送失败，请检查手机号码")
            } else if (text == "0") {
				document.getElementById("show-time").innerHTML = "已发送验证码"
                Countdown()
            } else if (text == "1") {
                noticeErrSpan("phone", "手机号码未注册!")
            } else if (text == "5") {
                noticeErrSpan("phone", "请输入正确的手机号码")
            } else {
            	noticeErrSpan("phone", "请等待 "+text+"秒后 在发送验证码")
            }
        })
    }
}

function Countdown() {
	let count = 60
	checked = count
	OneCountdown() 
	for (var i = 1; i <= count && checked != 0; i++) {
		setTimeout("OneCountdown()", i*1000)
	}
}

function OneCountdown() {
	document.getElementById("show-time").innerHTML = checked--
	if (checked < 0) {
		document.getElementById("show-time").innerHTML = "获取验证码"
		checked = 0
	}
}

function checkData(pFlag) {
	let err = false

	if (pFlag) {
		let pw = document.getElementById("passwd-input")
		if (pw.value.length < 1) {
			noticeErrSpan("passwd", "请输入密码")	
			err = true
		} else if ((pw.value.length < 6) || (pw.value.length > 32))  {
			noticeErrSpan("passwd", "密码长度 6-32 ")	
			err = true
		}
	}

	let phone = document.getElementById("phone-input")
	if (phone.value.length < 1) {
		noticeErrSpan("phone", "请输入手机号码")
		err = true
	} else {
		let pattern= /^1[3|4|5|8|7][0-9]\d{4,8}$/ 
        if ( !pattern.test(phone.value) ) {
        	noticeErrSpan("phone", "请输入正确的手机号码")
        	err = true
        }
	}
	return err
}

// span 提示
function noticeErrSpan(id, msg) {
	document.getElementById(id+"-input").className = "item-input item-input-error"
	document.getElementById(id+"-span").style.display="";
	document.getElementById(id+"-span").innerHTML = msg
	setTimeout("document.getElementById('"+id+"-span').style.display='none'", 1500)
}

function addKeyPress(id) {
	let temp = document.getElementById(id+"-input")
	temp.onkeypress = function (evt) {
		temp.className = "item-input"
	}
}