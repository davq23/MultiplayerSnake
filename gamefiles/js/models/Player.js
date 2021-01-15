import Position from "./Position.js";

const playerRadius = 20;
const playerParts = 4;

export const playerUp = 1;
export const playerDown = 2;
export const playerLeft = 3;
export const playerRight = 4;

export class Player {
    constructor(canvas, id, name, x, y, color) {
        if (canvas instanceof HTMLCanvasElement && typeof(id) === 'string' && typeof (x) === 'number' && typeof (y) === 'number') {
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
        } else {
            throw new Error('Invalid arguments to Player constructor');
        }
    }

    render() {
        for (let i = 0; i < this.positions.length; i++) {
            this.ctx.beginPath();
            this.ctx.arc(this.positions[i].x, this.positions[i].y, playerRadius, 0, 2 * Math.PI);
            this.ctx.fillStyle = this.color;
            this.ctx.fill();
            this.ctx.lineWidth = 1;
            this.ctx.strokeStyle = '#003300';
            this.ctx.stroke();
            this.ctx.closePath();
        }

        this.ctx.beginPath();
        this.ctx.fillStyle = 'white';
        this.ctx.font = '100px Arial';
        this.ctx.fillText(this.name, this.positions[0].x, this.positions[0].y, 100);
        this.ctx.closePath();
    }

    toJSON() {
        return {
            player_id: this.id,
            color: this.color,
            positions: this.positions,
            direction: this.direction,
        }
    }

    setScore(score) {

        if (score && this.score != score) {
            this.score = score;
            const sc = document.getElementById(this.id + '-score');

            if (sc) sc.innerText = `${this.name} : ${this.score}`;
        }
    }

}