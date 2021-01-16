import PlayerControl from "./gui/control.js";
import Scores from "./gui/scores.js";


import {
    Message,
    MessageMove,
    MessageRegister,
    MessageUnregister,
    MessageGetPlayers,
    MessageTracking
} from "./models/Message.js";

import {
    Player,
    playerDown,
    playerLeft,
    playerRight,
    playerUp
} from "./models/Player.js";
import Position from "./models/Position.js";

const canvas = document.getElementById('gameboard');
const control = document.getElementById('gamecontrol');
const scores = document.getElementById('gamescores');
const players = new Map();

document.onreadystatechange = function (event) {
    if (this.readyState === 'complete') {
        try {
            const ws = new WebSocket(`wss://rocky-hamlet-16573.herokuapp.com/game`);

            const playerControl1 = new PlayerControl(control);
            const playerScores = new Scores(scores);

            const ctx = canvas.getContext('2d');

            let player1 = new Player(canvas, '', '', 0, 0, 'black');

            let frames = 0;

            ws.onopen = function (event) {

            };

            ws.onclose = function () {
                alert('Connection closed');
                window.location.href = '/';
            }

            function moveEvent(direction) {
                return async function () {
                    const msg = new Message(MessageMove, player1);

                    msg.player.direction = direction
                    ws.send(JSON.stringify(msg))
                }
            }

            playerControl1.setEvents(moveEvent(playerUp), moveEvent(playerDown), moveEvent(playerLeft), moveEvent(playerRight));

            ws.onmessage = function (event) {
                const msg = JSON.parse(event.data);

                switch (msg.type) {
                    case MessageRegister:
                        player1 = new Player(canvas, msg.player_info.player_id, msg.player_info.name, msg.player_info.positions[0].X, msg.player_info.positions[0].Y, msg.player_info.color);
                        players.set(msg.player_info.player_id, player1);

                        const resp = new Message(MessageGetPlayers, player1)
                        ws.send(JSON.stringify(resp))

                    case MessageMove:
                        const player = players.get(msg.player_info.player_id);

                        if (player) {
                            if (msg.player_info.positions.length === player.positions.length) {
                                player.update(new Position(msg.player_info.positions[0].X, msg.player_info.positions[0].Y));

                            } else {
                                this.positions.length = 0;

                                this.positions = new Array(msg.player_info.positions.length);

                                for (let i = 0; i < msg.player_info.positions.length; i++) {
                                    this.positions[i] = new Position( msg.player_info.positions[i].X, msg.player_info.positions[i].Y);                                    
                                }
                            }
                            player.setScore(msg.player_info.score);
                        }

                       
                        break;
                    case MessageGetPlayers:
                        const entries = Object.entries(msg.players)
                        player1 = msg.player_info

                        entries.forEach(function (entry) {
                            players.set(entry[0], new Player(canvas, entry[0], entry[1].name, entry[1].positions[0].X, entry[1].positions[0].Y, entry[1].color));
                        })

                        playerScores.render(players);

                        break;

                    case MessageTracking:

                        const response = new Message(MessageTracking, player1);
                        ws.send(JSON.stringify(response))

                        break;
                    case MessageUnregister:
                        players.delete(msg.player_info.player_id);

                        if (player1.id === msg.player_info.player_id) {
                            playerControl1.anchor.innerHTML = `<h3 style="color:red;">EATEN</h3>`; 
                        }

                        playerScores.render(players);

                        break;
                }
            }

            playerControl1.render();

            function main() {
                ctx.fillStyle = 'black';

                ctx.clearRect(0, 0, canvas.width, canvas.height);
                ctx.fillRect(0, 0, canvas.width, canvas.height);

                players.forEach(function (value, key, map) {
                    value.render()
                })

                frames++
                frames %= Number.MAX_SAFE_INTEGER
                requestAnimationFrame(main);
            }

            requestAnimationFrame(main);
        } catch (err) {
            console.log(err);
        }

    }
}