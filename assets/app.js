var app = {
  WEBSOCKET_URL: "ws://127.0.0.1:3000/listen",
};

(function() {
  app.init = function () {
      console.log("initializing javascript app");
      app.connect();

      $('#render-btn').click(function(e) {
        app.sendWorldDesc();
        e.preventDefault();
      });
  };

  app.connect = function () {
    var con = new WebSocket(app.WEBSOCKET_URL)
    console.log("connecting");

    con.onopen = function () {
      console.log("onopen");
    };

    con.onerror = function (error) {
      console.log("onerror " + error);
    };

    con.onmessage = function (e) {
      console.log("onmessage " + e.data);
      var message = e.data;
      // TODO check if IMG...
      if(message.startsWith("IMG")) {
        var imagePath = message.replace("IMG","")
        app.displayImage(imagePath) 
      } else {
        alert("UNSUPPORTED MESSAGE: " + message);
      }
    };

    app.wsConnection = con;
  };

  app.sendWorldDesc = function () {
      var wd = app.getWorldDesc();
      app.send("CFG"+wd);
  };
  
  app.send = function (msg) {
    console.log("will send: '" + msg + "'");
    app.wsConnection.send(msg);
  };

  app.getWorldDesc = function () {
    var worldDesc = $("#editor").val();
    return worldDesc;
  };

  app.displayImage = function (imageName) {
    console.log('displaying image: ' + imageName);
    var path = "/renders/" + imageName;
    $("#scene img").attr("src", path);
  };
})();

$(document).ready(function() {
  app.init();
});
