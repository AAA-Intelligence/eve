window.onload = function () {

    var bots = document.getElementsByClassName("bot")
   /* var messages = {};
    // test load messages
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", "/getmessages", true); // false for synchronous request
    xmlHttp.send();
    xmlHttp.onload = function () {
        messages = JSON.parse(xmlHttp.responseText)
        if (bots.length > 0)
            changeActiveBot(bots[0])
    }

    for (i in bots) {
        bots[i].onclick = function () {
            changeActiveBot(this)
        }
    }*/

    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");
    var form = document.getElementById("form")

    msg.onkeypress = function (key) {
        if (key.keyCode === 13 && !key.shiftKey) {
            sendMsg()
            return false;
        }
    }

    function appendChat(item, sender) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(toMsg(item, sender));
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    function toMsg(text, style) {
        var item = document.createElement("div");
        item.innerHTML = "<div class=\"msg " + style + "\">" + text + "</div>";
        return item
    }
    form.onsubmit = function () {
        sendMsg()
        return false
    };

    function sendMsg() {
        if (!conn) {
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
        /*messages[botID].push({
            "Content":message
        })*/
        appendChat(message, "user")
        var request = {
            "message": message,
            "bot": botID
        }
        conn.send(JSON.stringify(request));
        msg.value = "";
        return;
    }

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            appendChat("<b>Connection closed.</b>", "user");
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                //TODO store message localy
                appendChat(messages[i], "bot");
            }
        };
    } else {
        appendChat("<b>Your browser does not support WebSockets.</b>", "bot");
    }

    function escapeHtml(html) {
        var text = document.createTextNode(html);
        var div = document.createElement('div');
        div.appendChild(text);
        return div.innerHTML;
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
        for(var index in messages[botID]){
            var msg = messages[botID][index]
            appendChat(msg["Content"],msg["Sender"]===1?"user":"bot")
        }
    }
};