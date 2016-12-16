angular.module('deploymentSSH', [])
    .controller('DeploymentSSHController', ['$scope', '$http', '$stateParams', '$location',
        function($scope, $http, $stateParams, $location) {
            var sock = new WebSocket("ws://" + $location.host() + ":" + $location.port() + "/deployments/" + $stateParams.name + "/vms/" + $stateParams.vm_name + "/ssh");
            sock.onerror = function(e) {
                console.log("socket error", e);
            };
            // wait for the socket to open before starting the terminal
            // or there will be ordering issues :/
            sock.onopen = function(e) {
                var term = new Terminal({
                    cols: 120,
                    rows: 30,
                    useStyle: true,
                    screenKeys: true
                });
                term.open(document.getElementById("bash"))
                term.on('title', function(title) {
                    document.title = title;
                });
                // pass data using base64 encoding
                // this is fragile: it will not work with non-ascii text!
                // the Go backend is correctly treating pty IO as opaque
                // byte arrays, while term.js uses javascript strings that
                // are utf16, while the pty is usually utf8.
                // I have some Go code that converts to utf16 before sending but
                // it's ugly and wrong. The right answer is to refactor term.js to use
                // ArrayBuffer with uint8 and convert runes on the fly on the client
                term.on('data', function(data) {
                    sock.send(btoa(data));
                });
                sock.onmessage = function(msg) {
                    term.write(atob(msg.data));
                };
            };
        }
    ])