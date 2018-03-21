window.onload = function () {

    var bots = document.getElementsByClassName("bot")
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");
    var form = document.getElementById("form");
    var sendBtn = document.getElementById("send");

    msg.onkeypress = function (key) {
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
        var item = document.createElement("div");
        item.innerHTML = "<div class=\"msg " + style + "\">" + text + "</div>";
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