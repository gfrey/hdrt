var app = {
  WEBSOCKET_URL: "ws://127.0.0.1:3000/listen",
};

(function() {
  app.init = function () {
      console.log("initializing javascript app");
      app.connect();
  };

  app.connect = function () {
    var con = new WebSocket(app.WEBSOCKET_URL)
    console.log("connecting");

    con.onopen = function () {
      console.log("onopen");
      con.send("Ping"); 
    };

    con.onerror = function (error) {
      console.log("onerror " + error);
    };

    con.onmessage = function (e) {
      console.log("onmessage " + e.data);
    };
  };
})();

app.init();
