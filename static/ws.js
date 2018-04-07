window.onload = function () {

    var bots = document.getElementsByClassName("bot")
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");
    var form = document.getElementById("form");
    var sendBtn = document.getElementById("send");
	
    msg.onkeypress = function (key) {
        // send message on press enter
        if (key.keyCode === 13 && !key.shiftKey) {
            sendMsg()
            return false;
        }
    }

    function appendChat(item, sender) {
        log.appendChild(toMsg(item, sender));
        scrollChatToBottom();
    }

    function toMsg(text, style) {
        var item = document.createElement("p");
        var time = new Date()
        var mytime = time.getHours()+":"+time.getMinutes();
        item.innerHTML = "<p class=\"msg " + style + "\">" + text +"<span class=\"timestamp\">"+mytime+"</span>"+"</p>";
        return item
    }
    form.onsubmit = function () {
        sendMsg()
        return false
    };

    function sendMsg() {
        if (!conn || locked) {
            return;
        }
        if (!msg.value || msg.value.trim().length < 1) {
            return;
        }
        var message = escapeHtml(msg.value).trim()
        var activeBot = document.getElementById("bot-active")
        if (!activeBot)
            return
        var botID = parseInt(activeBot.getAttribute("botID"))
        appendChat(message, "user")
        var request = {
            "message": message,
            "bot": botID
        }
        conn.send(JSON.stringify(request));
        msg.value = "";
        lock();
        return;
    }

    var locked = false
    function lock(){
        sendBtn.disabled = true;
        locked = true
    }
    function unlock(){
        sendBtn.disabled = false;
        locked = false
    }

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            appendChat("<b>Connection closed.</b>", "user");
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                appendChat(messages[i], "bot");
            }
            unlock();
        };
    } else {
        appendChat("<b>Your browser does not support WebSockets.</b>", "bot");
    }

    function escapeHtml(html) {
        var text = document.createTextNode(html);
        var p = document.createElement('p');
        p.appendChild(text);
        return p.innerHTML;
    }

    function changeActiveBot(newActive) {
        var active = document.getElementById("bot-active")
        if (active == newActive)
            return
        // remove is there allready is a active one
        active && active.removeAttribute("id")
        newActive.setAttribute("id", "bot-active")
        document.getElementById("log").innerHTML = ""
        var botID = parseInt(newActive.getAttribute("botid"))
        for (var index in messages[botID]) {
            var msg = messages[botID][index]
            appendChat(msg["Content"], msg["Sender"] === 1 ? "user" : "bot")
        }
    }

    function scrollChatToBottom() {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
    scrollChatToBottom();
	
};

	function showPopup() {
		var popup = document.getElementById("popup");
		var contentAr = document.getElementsByClassName("content");
		var content = contentAr[0];
		
		//Show popup
		popup.style["visibility"] = "visible";
		popup.style["display"] = "block";
		
		//Change background (disable interaction)
		content.style["pointer-events"] = "none";
		content.style["-webkit-touch-callout"] = "none";
		content.style["-webkit-user-select"] = "none";
		content.style["-khtml-user-select"] = "none";
		content.style["-moz-user-select"] = "none";
		content.style["-ms-user-select"] = "none";
		content.style["user-select"] = "none";
		//Change background (apply blur)
		content.style["filter"] = "progid:DXImageTransform.Microsoft.Blur(PixelRadius='3')";
		content.style["-webkit-filter"] = "url(#blur-filter)";
		content.style["filter"] = "url(#blur-filter)";
		content.style["-webkit-filter"] = "blur(3px)";
		content.style["filter"] = "blur(3px)";
	}
	
	function hidePopup() {
		var popup = document.getElementById("popup");
		var contentAr = document.getElementsByClassName("content");
		var content = contentAr[0];
		
		//Hide popup
		popup.style["visibility"] = "";
		popup.style["display"] = "";
		
		//Change background (enable interaction)
		content.style["pointer-events"] = "";
		content.style["-webkit-touch-callout"] = "";
		content.style["-webkit-user-select"] = "";
		content.style["-khtml-user-select"] = "";
		content.style["-moz-user-select"] = "";
		content.style["-ms-user-select"] = "";
		content.style["user-select"] = "";
		//Change background (remove blur)
		content.style["filter"] = "";
		content.style["-webkit-filter"] = "";
	}
	
	function fetchJSON(sex, callback) {
		fetch(`./getRandomName?sex=`+sex)
		.then(response => response.json())
		.then(json => callback(null, json.result))
		.catch(error => callback(error, null))
	}
	
	function genName() {
		/*id = Math.floor(Math.random() * 5);
		names = ["Lina", "Laura", "Lisa", "Loreen", "Linda"];*/
		
		//Fetch selected sex
		if (document.getElementById("switch_left").checked) {
			sex = 0;
		} else {
			sex = 1;
		}
		
		fetchJSON(sex, (error, result) => {
			if (error) 
				console.log(error)
			else 
				name = result[Name]
		})
		
		if (sex==0) name = "Laura";
		else name = "Peter";
		setNameOnCreation(name);
	}
	
	function setNameOnCreation(newName) {
		var namefield = document.getElementById("generatedName");
		namefield.innerHTML = newName+"<button onclick='genName()'></button>";		
	}
	
	function onSexChange() {
		genName();
		//+ Change picture!
	}