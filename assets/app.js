var app = {
  WEBSOCKET_URL: "ws://127.0.0.1:3000/listen",
  STORE_JSON_KEY: "json_content",
};

(function() {
  app.init = function () {
      console.log("initializing javascript app");
      app.connect();

      $('[data-role=render]').click(function(e) {
        app.resetImage();
        app.sendWorldDesc();
        e.preventDefault();
        return false;
      });

      $('[data-role=abort]').click(function(e) {
        app.abortRendering();
        e.preventDefault();
        return false;
      });

      app.initEditor();
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
    app.saveJSON();
    var wd = app.getWorldDesc();
    app.send("CFG"+wd);
  };
  
  app.send = function (msg) {
    console.log("will send: '" + msg + "'");
    app.wsConnection.send(msg);
    console.log("did send: '" + msg + "'");
  };

  app.getWorldDesc = function () {
    var worldDesc = app.editor.get();
    return JSON.stringify(worldDesc);
  };

  app.handleMessage = function (message) {
    if(message.startsWith("IMG")) {
      var imagePath = "/renders/" + message.replace("IMG","")
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
    var path = imageName;
    $("#scene img").attr("src", path);
    $("#scene a").attr("href", path);
  };

  app.resetImage = function () {
    var w = document.querySelector("#scene img").naturalWidth;
    var h = document.querySelector("#scene img").naturalHeight;
    var ratio = h/w;
    w = 600;
    h = 600*ratio;

    app.displayImage("http://fillmurray.com/" + w + "/" + h);
  };

  app.initEditor = function () {
    var container = document.querySelector("#editor");
    var options = {
       mode: 'code',
        modes: ['code', 'form', 'text', 'tree', 'view'], // allowed modes
        error: function (err) {
          alert(err.toString());
        }
    };
    var json = {};

    app.editor = new JSONEditor(container, options, json);
    app.loadJSON();
  };

  app.loadJSON = function() {
  console.log("loading json");
    var jsontext = localStorage.getItem(app.STORE_JSON_KEY);
    if(jsontext == "" || jsontext == null || jsontext == undefined) {
        jsontext = "{}";
    }
    var json = JSON.parse(jsontext);
    app.editor.set(json);
  };

  app.saveJSON = function() {
    console.log("saving json");
    var json = app.editor.get();
    var jsontext = JSON.stringify(json);
        localStorage.setItem(app.STORE_JSON_KEY, jsontext);
  };

})();

$(document).ready(function() {
  app.init();
});
