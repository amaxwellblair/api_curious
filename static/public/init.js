$(document).ready(function() {
  if (tokenExists()) {
    turnOnLogout();
    var token = getTokenStorage();
    loadDoc(token);
  } else if (tokenExistsInURL()) {
    turnOnLogout();
    var token = getTokenURL();
    loadDoc(token);
    storeToken(token);
  } else {
    turnOnLogin();
  }
});
