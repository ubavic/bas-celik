var HOST_NAME = 'com.github.ubavic.bas_celik';
var EXTENSION_NAME = HOST_NAME;

var nativePort;
var ports = []

function openNativePort(portId) {
    var port = chrome.runtime.connectNative(HOST_NAME);

    // receive from host, forward to content
    port.onMessage.addListener(function(message) {
        console.log("Background received from host", message);
        ports.forEach(function(port){
            port.postMessage(message);
        })
    });

    // handle host disconnect
    port.onDisconnect.addListener(function() {
        var error = chrome.runtime.lastError;
        if (error) {
            console.log("Host error", error.message);
        }
        console.log("Host disconnected");
        ports.forEach(function(port){
            port.disconnect();
            ports.splice(ports.indexOf(port), 1);
        })
        nativePort = null
    });

    nativePort = port;
}


chrome.runtime.onConnect.addListener(function(port) {
    port.onDisconnect.addListener(function() {
      ports.splice(ports.indexOf(port), 1);
      console.log(ports);
    });
    ports.push(port);
    console.log(ports);
    if(!nativePort){
        openNativePort();
    }
});