export default class Position {
    constructor(x, y) {
        if (typeof(x) === 'number' &&  typeof(y) === 'number') {
            this.x = x;
            this.y = y;
        }
    }

    toJSON() {
        return {
            x: this.x,
            y: this.y
        }
    }
}