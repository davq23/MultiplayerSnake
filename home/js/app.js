const hubList = document.getElementById('hub-list');
const hubRefresh = document.getElementById('hub-refresh');
const createHub = document.getElementById('create-hub');

document.onreadystatechange = function() {
    if (document.readyState == 'complete') {
        async function saveToken(response) {
            const token = await response.json();
            localStorage.setItem('game-token', token.token);
            alert(localStorage.getItem('game-token'));
        }

        createHub.onclick = async function (event) {
            const body = {
                'username': document.getElementById('username').value,
                'hubname': document.getElementById('hubname').value
            }
        
            const response = await fetch('/hubs/create', {
                'method': 'post',
                'credentials': 'include',
                'headers': new Headers({
                    'Authorization': localStorage.getItem('game-token'), 
                }),
                'body': JSON.stringify(body),
            });
        
        
            if (response.status == 200) {
                await saveToken(response);
                window.location.href = '/play';
            } else {
                document.getElementById('create-hub-err').innerText = await response.text();
            }
        
        }

        async function getHubs(event) {
            const response = await fetch('/hubs');
        
            if (response.status === 200) {
                const fragment = document.createDocumentFragment();

                const message = document.createElement('h3');
                message.innerText = 'Loading'
                message.classList.add('spinning') 

                hubList.appendChild(message);


                const hubs = await response.json();
        
                hubList.innerHTML = '';

                if (hubs && hubs.length > 0) {
                    hubs.forEach(hub => {
                        const li = document.createElement('li');
                        li.innerText = `${hub.name} : ${hub.player_num} ${hub.player_num !== 1 ? 'players' : 'player'}`;
                        li.classList.add('hub');
    
                        
                        li.onclick = async function(event) {
                            const body = {
                                'username': document.getElementById('username').value,
                                'hubname': hub.name
                            }
    
                            const resp = await fetch('/hubs/join', {
                                'method': 'post',
                                'credentials': 'include',
                                'headers': new Headers({
                                    'Authorization': localStorage.getItem('game-token'), 
                                }),
                                'body': JSON.stringify(body),
                            });
            
                            if (resp.status === 200) {
                                await saveToken(response);
                                window.location.href = '/play';
                            } else {
                                document.getElementById('join-hub-err').innerText = await response.text();
                            }
                        };
            
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

