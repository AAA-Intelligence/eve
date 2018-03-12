window.onload = function () {

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
        if (!msg.value || msg.value.trim().length<1) {
            return;
        }
        var message = escapeHtml(msg.value).trim()
        appendChat(message, "user")
        conn.send(message);
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
};