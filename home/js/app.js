const hubList = document.getElementById('hub-list');
const hubRefresh = document.getElementById('hub-refresh');
const createHub = document.getElementById('create-hub');

document.onreadystatechange = function() {
    if (document.readyState == 'complete') {
        createHub.onclick = async function (event) {
            const body = {
                'username': document.getElementById('username').value,
                'hubname': document.getElementById('hubname').value
            }
        
            const response = await fetch('/hubs/create', {
                'method': 'post',
                'credentials': 'include',
                'body': JSON.stringify(body)
            });
        
        
            if (response.status == 200) {
                window.location.href = '/play';
            } else {
                document.getElementById('create-hub-err').innerText = await response.text();
            }
        
        }

        async function getHubs(event) {
            const response = await fetch('/hubs');
        
            if (response.status === 200) {
                const hubs = await response.json();
        
                console.log(hubs);
        
                const fragment = document.createDocumentFragment();
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
                                method: 'post',
                                credentials: 'include',
                                body: JSON.stringify(body)
                            });
            
                            if (resp.status === 200) {
                                window.location.href = '/play';
                            } else {
                                document.getElementById('join-hub-err').innerText = await response.text();
                            }
                        };
            
                        fragment.appendChild(li);
                    });
                } else {
                    const h3 = document.createElement('h3');
                    h3.innerText = 'No hubs available';
                    fragment.appendChild(h3);
                }
        
                
                hubList.appendChild(fragment);
            }
        
        
        }
        
        hubRefresh.onclick = getHubs;

        getHubs(null);
    }
}

