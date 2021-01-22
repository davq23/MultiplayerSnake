var hubList = document.getElementById('hub-list');
var hubRefresh = document.getElementById('hub-refresh');
var createHub = document.getElementById('create-hub');

document.onreadystatechange = function() {
    if (document.readyState == 'complete') {
        function saveToken(response) {
            var token = JSON.parse(response);

            try {
                if (localStorage) {
                    localStorage.setItem('game-token', token.token);
                } else {
                    document.cookie = "game-token=" + token+ "; path=/; secure; samesite=strict";
                }
            } catch(err) {
                alert(err);
            }
        }

        function getEnterHubFunction(hub) {
            return function(event) {
                var body = {
                    'username': document.getElementById('username').value,
                    'hubname': hub.name
                }

                var xhr = new XMLHttpRequest();

                xhr.onreadystatechange = function() {
                    if (this.readyState == 4 && this.status == 200) {
                        saveToken(this.responseText);
                        window.location.href = '/play';
                    } 
                    else if (this.readyState == 4 && this.status != 200) {
                        document.getElementById('create-hub-err').innerText = this.responseText;
                    }
                }

                xhr.withCredentials = true;

                var authToken = '';

                if (localStorage) {
                    authToken = localStorage.getItem('game-token');
                } else {
                    var cookies = document.cookie.split(';=');

                    var gameTokenPos = cookies.findIndex('game-token');

                    if (gameTokenPos >= 0 || gameTokenPos < cookies.length - 1) {
                        authToken = cookies[gameTokenPos + 1];
                    } else {
                        throw 'Unable to send user state';
                    }
                }

                xhr.setRequestHeader("Authorization", authToken);
                xhr.open('POST', '/hubs/join');
                xhr.send(JSON.stringify(body));
            }
        }

        createHub.onclick = async function (event) {
            var body = {
                'username': document.getElementById('username').value,
                'hubname': document.getElementById('hubname').value
            }
        
            var xhr = new XMLHttpRequest();

            xhr.onreadystatechange = function() {
                if (this.readyState == 4 && this.status == 200) {
                    saveToken(this.responseText);
                    window.location.href = '/play';
                } 
                else if (this.readyState == 4 && this.status != 200) {
                    document.getElementById('create-hub-err').innerText = this.responseText;
                }
            }

            xhr.withCredentials = true;

            var authToken = '';

            if (localStorage) {
                authToken = localStorage.getItem('game-token');
            } else {
                var cookies = document.cookie.split(';=');

                var gameTokenPos = cookies.findIndex('game-token');

                if (gameTokenPos >= 0 || gameTokenPos < cookies.length - 1) {
                    authToken = cookies[gameTokenPos + 1];
                } else {
                    throw 'Unable to send user state';
                }
            }

            xhr.setRequestHeader("Authorization", authToken);
            xhr.open('POST', '/hubs/create');
            xhr.send(JSON.stringify(body));
        
        }

        async function getHubs(event) {
            var message = document.createElement('h3');
            message.innerText = '<<'
            message.classList.add('spinning') 

            hubList.appendChild(message);

            var response = await fetch('/hubs');
        
            if (response.status === 200) {
                var fragment = document.createDocumentFragment();

                var hubs = await response.json();
        
                hubList.innerHTML = '';

                if (hubs && hubs.length > 0) {
                    hubs.forEach(hub => {
                        var li = document.createElement('li');
                        li.innerText = hub.name + " : " +  hub.player_num + (hub.player_num !== 1 ? ' players' : ' player');
                        li.classList.add('hub');
    
                        li.onclick = getEnterHubFunction(hub);
            
                        fragment.appendChild(li);
                    });
                } else {
                    message.innerText = 'No hubs available';
                    message.classList.remove('spinning');
                    fragment.appendChild(message)
                }
        
                
                hubList.appendChild(fragment);
            }
        
        
        }
        
        hubRefresh.onclick = getHubs;

        getHubs(null);
    }
}

