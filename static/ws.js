window.onload = function () {

    var bots = document.getElementsByClassName("bot")
    bots[0].setAttribute("id","bot-active")
    for(i in bots){
        bots[i].onclick = function(){
            document.getElementById("bot-active").removeAttribute("id")
            this.setAttribute("id","bot-active")
        }
    }

    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function appendChat(item,sender) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(toMsg(item,sender));
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    function toMsg(text,style){
        var item = document.createElement("div");
        item.innerHTML = "<div class=\""+style+"\">"+text+"</div>";
        return item
    }
    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        appendChat(msg.value,"user")
        var request = {
            "message":msg.value,
            "bot":parseInt(document.getElementById("bot-active").getAttribute("botID"))
        }
        conn.send(JSON.stringify(request));
        msg.value = "";
        return false;
    };
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            appendChat("<b>Connection closed.</b>","user");
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                appendChat(messages[i],"bot");
            }
        };
    } else {
        appendChat("<b>Your browser does not support WebSockets.</b>","bot");
    }
};