import { Player } from "./Player.js";


export const MessageTracking = 0;
export const MessageMove = 1;
export const MessageRegister = 2;
export const MessageUnregister = 3;
export const MessageGetPlayers = 4;


export class Message {
    constructor(type, player) {
        if (typeof(type) == 'number' && player instanceof Object) {
            this.type = type;
            this.player = player;
        } else {
            throw new Error("Incorrect args")
        }
    }

    toJSON() {
        return {
            type: this.type, 
            player_info: this.player
        };
    }
}