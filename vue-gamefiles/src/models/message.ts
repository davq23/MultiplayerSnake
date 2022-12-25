import type Player from './player';

export enum MessageType {
    // MessageTracking message to move player
    MessageTracking = 0,
    // MessageMove message to move player
    MessageMove = 1,
    // MessageRegister message to register player
    MessageRegister = 2,
    // MessageUnregister to unregister player
    MessageUnregister = 3,
    // MessageGetPlayers to get all other players
    MessageGetPlayers = 4,
    // MessageRefresh regenerates token
    MessageRefresh = 5,
}

export default interface Message {
    type: MessageType;
    player_info: Player | null | undefined;
    received_at: string | undefined;
    sent_at: string | null | undefined;
    players: Object | null | undefined;
    new_token: string | null | undefined;
}
