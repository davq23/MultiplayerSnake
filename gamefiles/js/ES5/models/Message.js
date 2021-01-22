var MessageTracking = 0;
var MessageMove = 1;
var MessageRegister = 2;
var MessageUnregister = 3;
var MessageGetPlayers = 4;
var MessageRefresh = 5;

function Message(type, player) {
    if (typeof(type) == 'number' && player instanceof Player) {
        this.type = type;
        this.player = player;
    }

    var selfMessage = this;

    this.toJSON = function() {
        return {
            type: selfMessage.type, 
            player_info: selfMessage.player
        };
    }
}