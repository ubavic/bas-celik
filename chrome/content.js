if(document.body.dataset.celik){
    var port = chrome.runtime.connect();
    port.onMessage.addListener(function(msg) {
        msg.source = 'bascelik';
        window.postMessage(msg)
    });
    port.onDisconnect.addListener(function(){
        window.postMessage({sourece: "bascelik", error: 5, message: "Host Disconnected"});
    })
}
