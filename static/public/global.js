function loadDoc(token) {
  var url = "http://localhost:9000/github"
  var username = ""
  $.get(url.concat("/user?user_id="+token), function(data) {
    $("#avatar_url").append('<img src=' + parseUser(data)["avatar_url"] + ' width=200px ></img>');
    $("#name").append(parseUser(data)["name"]);
    $("#blog").append('<a href=' + 'http://' + parseUser(data)["blog"] + '>' +parseUser(data)["blog"]+'</a>');
    // $("#json").append(data);
    username = parseUser(data)["login"];
    $.get(url.concat("/user/gists?user_id="+token+"&username="+username), function(data) {
      // $("#json").append(parseGist(data));
      var gists = parseGists(data)
      for (var i = 0; i < gists.length; i++) {
        loadGist(gists[i]);
      }
    });
  });
}


function parseToken(url) {
  var query = url.split("?");
  if (query.length > 1) {
    var temp = query[1].split("=");
    return temp[1];
  } else {
    return ""
  }
}

function storeToken(token) {
  return localStorage.setItem('user_id', token);
}

function getTokenStorage() {
  return localStorage.getItem("user_id");
}

function getTokenURL() {
  return parseToken(document.URL)
}

function tokenExists() {
  if (getTokenStorage() == null || getTokenStorage() == "undefined") {
    return false;
  } else {
    return true;
  }
}

function tokenExistsInURL() {
  if (parseToken(document.URL) == "") {
    return false;
  } else {
    return true;
  }
}

function parseUser(json) {
  var obj = JSON.parse(json);
  return obj;
}

function loadGist(url) {
  $.get(url, function(data) {
    $("#gist").append('<div class="col md-12">'+data+'</div>');
  });
}

function parseGists(json) {
  var obj = JSON.parse(json);
  var gistURLs = []
  for (var i = 0; i < obj.length; i++) {
    gistURLs.push(parseGist(obj[i]))
  }
  return gistURLs
}

function parseGist(gist) {
  var files = gist["files"];
  var gistName = Object.keys(files)[0];
  var url = files[gistName]["raw_url"];
  return url;
}

function logout() {
  $("#logout").hide();
  $("#login").show();
  localStorage.clear();
}

function turnOnLogin() {
  $("#login").show();
}

function turnOnLogout() {
  $("#logout").show();
}
