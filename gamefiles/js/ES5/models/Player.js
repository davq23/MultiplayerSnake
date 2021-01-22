var playerRadius = 20;
var playerParts = 4;

var playerUp = 1;
var playerDown = 2;
var playerLeft = 3;
var playerRight = 4;

function Player(canvas, id, name, x, y, color) {
    this.canvas = canvas
    this.ctx = canvas.getContext('2d');

    this.id = id;
    this.name = name;
    this.color = color;
    this.positions = [];
    this.score = 0;

    for (let i = 0; i < playerParts; i++) {
        this.positions.push(new Position(x + (playerRadius * 2 * i), y));
    }

    var selfPlayer = this;

    this.update = function(head) {
        if (head instanceof Position) {
            for (let i = selfPlayer.positions.length - 1; i > 0; i--) {
                selfPlayer.positions[i] = selfPlayer.positions[i - 1];
            }

            selfPlayer.positions[0] = head;
        }
    }

    this.render = function() {
        for (let i = 0; i < selfPlayer.positions.length; i++) {
            selfPlayer.ctx.beginPath();
            selfPlayer.ctx.arc(selfPlayer.positions[i].x, selfPlayer.positions[i].y, playerRadius, 0, 2 * Math.PI);
            selfPlayer.ctx.fillStyle = selfPlayer.color;
            selfPlayer.ctx.fill();
            selfPlayer.ctx.lineWidth = 1;
            selfPlayer.ctx.strokeStyle = '#003300';
            selfPlayer.ctx.stroke();
            selfPlayer.ctx.closePath();
        }

        selfPlayer.ctx.beginPath();
        selfPlayer.ctx.fillStyle = 'white';
        selfPlayer.ctx.font = '50px Lucida Sans';
        selfPlayer.ctx.fillText(selfPlayer.name, selfPlayer.positions[0].x, selfPlayer.positions[0].y, 200);
        selfPlayer.ctx.closePath();
    }

    this.toJSON = function() {
        return {
            player_id: selfPlayer.id,
            color: selfPlayer.color,
            positions: selfPlayer.positions,
            direction: selfPlayer.direction,
        }
    }

    this.setScore = function(score, token) {
        if (score && selfPlayer.score != score) {
            selfPlayer.score = score;
            const sc = document.getElementById(selfPlayer.id + '-score');
            
            if (token) {
                if (localStorage) {
                    localStorage.setItem('game-token', token);
                } else {
                    document.cookie = "game-token=" + token+ "; path=/; secure; samesite=strict";
                }
            }

            if (sc) sc.innerText = selfPlayer.name + " : " + selfPlayer.score;
        }
    }
}