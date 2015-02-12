var app = {
  WEBSOCKET_URL: "ws://127.0.0.1:3000/listen",
};

(function() {
  app.init = function () {
      console.log("initializing javascript app");
      app.connect();

      $('[data-role=render]').click(function(e) {
        app.resetImage();
        app.sendWorldDesc();
        e.preventDefault();
      });

      $('[data-role=abort]').click(function(e) {
        app.abortRendering();
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
      app.handleMessage(message);
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
    console.log("did send: '" + msg + "'");
  };

  app.getWorldDesc = function () {
    var worldDesc = $("#editor").val();
    return worldDesc;
  };

  app.handleMessage = function (message) {
    if(message.startsWith("IMG")) {
      var imagePath = message.replace("IMG","")
      app.displayImage(imagePath) 
    } else {
      alert("UNSUPPORTED MESSAGE: " + message);
    }
  };

  app.abortRendering = function () {
    app.send("ABORT");
  };

  app.displayImage = function (imageName) {
    console.log('displaying image: ' + imageName);
    var path = "/renders/" + imageName;
    $("#scene img").attr("src", path);
  };

  app.resetImage = function () {
    var w = document.querySelector("#scene img").naturalWidth;
    var h = document.querySelector("#scene img").naturalHeight;
    var ratio = h/w;
    w = 600;
    h = 600*ratio;

    $("#scene img").attr("src", "http://fillmurray.com/" + w + "/" + h);
  };
})();

$(document).ready(function() {
  app.init();
});
