<template>
    <div>
        <div v-if="message !== ''">{{ message }}</div>
        <websocket-wrapper
            @error="showMessage"
            @message="onMessage"
            :local-message="localMessage"
        >
            <game-play-field
                :width="1000"
                :height="2000"
                :map-players="mapPlayers"
                :player-radius="20"
            />
        </websocket-wrapper>
    </div>
</template>

<script lang="ts">
import type Message from '@/models/message';
import { MessageType } from '@/models/message';
import type Player from '@/models/player';
import { defineComponent } from 'vue';

interface GameViewData {
    currentPlayer: Player | null;
    localMessage: Message | null;
    message: string;
    mapPlayers: Map<string, Player> | null;
}

export default defineComponent({
    data(): GameViewData {
        return {
            currentPlayer: null,
            localMessage: null,
            mapPlayers: null,
            message: '',
        };
    },

    methods: {
        handleGetPlayers(message: Message) {
            if (message.players) {
                this.currentPlayer = message.player_info as Player;
                this.mapPlayers = new Map<string, Player>(
                    Object.entries(message.players)
                );
            }
        },
        handleMove(message: Message) {
            if (message.player_info && this.mapPlayers) {
                const player = this.mapPlayers.get(
                    message.player_info.player_id
                );

                if (player) {
                    player.positions = message.player_info.positions;
                }
            }
        },
        handleRegister(message: Message) {
            if (message.player_info && this.mapPlayers) {
                this.mapPlayers.set(
                    message.player_info.player_id,
                    message.player_info
                );

                const getPlayersMessage = { ...message };
                getPlayersMessage.type = MessageType.MessageGetPlayers;
                this.localMessage = getPlayersMessage;
            }
        },
        handleTracking() {
            const messageTracking = {
                type: MessageType.MessageTracking,
                player_info: this.currentPlayer,
            } as Message;

            this.localMessage = messageTracking;
        },
        handleUnregister(message: Message) {
            if (message.player_info) {
                this.mapPlayers?.delete(message.player_info.player_id);
            }
        },
        onMessage(message: Message) {
            switch (message.type) {
                case MessageType.MessageGetPlayers:
                    this.handleGetPlayers(message);
                    break;

                case MessageType.MessageMove:
                    this.handleMove(message);
                    break;
                case MessageType.MessageRefresh:
                    if (message.new_token) {
                        localStorage.setItem('game-token', message.new_token);
                    }
                    break;
                case MessageType.MessageRegister:
                    this.handleRegister(message);
                    break;
                case MessageType.MessageTracking:
                    this.handleTracking();
                    break;
                case MessageType.MessageUnregister:
                    this.handleUnregister(message);
                    break;
            }
        },
        showMessage(message: string) {
            this.message = message;
        },
    },
});
</script>

<style></style>
