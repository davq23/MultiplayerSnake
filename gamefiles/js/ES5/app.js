var canvas = document.getElementById('gameboard');
var control = document.getElementById('gamecontrol');
var scores = document.getElementById('gamescores');
var players = {}

document.onreadystatechange = function (event) {
    if (this.readyState === 'complete') {
        try {
            var ws = null;

            var webSocketRoute = 'wss://rocky-hamlet-16573.herokuapp.com/game?user_info=';

            if (localStorage) {
                webSocketRoute += localStorage.getItem('game-token');

            } else {
                var cookies = document.cookie.split(';=');

                var gameTokenPos = cookies.findIndex('game-token');

                if (gameTokenPos >= 0 || gameTokenPos < cookies.length - 1) {
                    webSocketRoute += cookies[gameTokenPos + 1];
                } else {
                    throw 'Unable to send user state';
                }
            }

            ws = new WebSocket(webSocketRoute);

            var playerControl = new PlayerControl(control);
            var playerScores = new ScoreTable(scores);

            var ctx = canvas.getContext('2d');

            var playerUser = new Player(canvas, '', '', 0, 0, 'black');

            function main() {
                ctx.fillStyle = 'black';

                ctx.clearRect(0, 0, canvas.width, canvas.height);
                ctx.fillRect(0, 0, canvas.width, canvas.height);

                for (var id in players) {
                    if (players.hasOwnProperty(id)) {
                        players[id].render();
                    }
                }
            }

            ws.onopen = function (event) {

            };

            ws.onclose = function () {
                alert('Connection closed');
                window.location.href = '/';
            }

            function moveEvent(direction) {
                return function () {
                    var msg = new Message(MessageMove, playerUser);

                    msg.player.direction = direction
                    ws.send(JSON.stringify(msg))
                }
            }

            playerControl.setEvents(moveEvent(playerUp), moveEvent(playerDown), moveEvent(playerLeft), moveEvent(playerRight));

            ws.onmessage = function (event) {
                var msg = JSON.parse(event.data);

                switch (msg.type) {
                    case MessageGetPlayers:

                        playerUser = msg.player_info

                        for (var p in msg.players) {
                            if (msg.players.hasOwnProperty(p)) {
                                players[p] = new Player(canvas, p, msg.players[p].name, msg.players[p].positions[0].x, msg.players[p].positions[0].y, msg.players[p].color)
                            }
                        }

                        playerScores.render(players);

                        break;
                    case MessageMove:
                        var player = players[msg.player_info.player_id];

                        if (player) {
                            if (msg.player_info.positions.length === player.positions.length) {
                                player.update(new Position(msg.player_info.positions[0].x, msg.player_info.positions[0].y));

                            } else {
                                player.positions = [];

                                for (var i = 0; i < msg.player_info.positions.length; i++) {
                                    player.positions.push(msg.player_info.positions[i]);
                                }
                            }

                            player.setScore(msg.player_info.score, msg.new_token);
                        }


                        break;
                    case MessageRefresh:
                        if (msg.new_token) {
                            if (localStorage) {
                                localStorage.setItem('game-token', msg.new_token);
                            } else {
                                document.cookie = "game-token=" + msg.new_token + "; path=/; secure; samesite=strict";
                            }
                        }
                        break;
                    case MessageRegister:
                        playerUser = new Player(canvas, msg.player_info.player_id, msg.player_info.name, msg.player_info.positions[0].x, msg.player_info.positions[0].y, msg.player_info.color);
                        playerUser.setScore(msg.player_info.score);
                        players[msg.player_info.player_id] = playerUser;

                        var resp = new Message(MessageGetPlayers, playerUser)
                        ws.send(JSON.stringify(resp))

                        break;
                    case MessageTracking:

                        var response = new Message(MessageTracking, playerUser);
                        ws.send(JSON.stringify(response))

                        break;
                    case MessageUnregister:
                        delete players[msg.player_info.player_id];

                        if (playerUser.id === msg.player_info.player_id) {
                            playerControl.anchor.innerHTML = '<h3 style="color:red;">EATEN</h3>';
                        }

                        playerScores.render(players);

                        break;
                }

                //requestAnimationFrame(main);
            }

            playerControl.render();

            setTimeout(requestAnimationFrame(main), 60);

        } catch (err) {
            alert(err);
            console.log(err);
        }

    }
}