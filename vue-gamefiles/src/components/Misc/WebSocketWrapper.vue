<template>
    <slot></slot>
</template>

<script lang="ts">
import type Message from '@/models/message';
import { defineComponent } from 'vue';
import type { PropType } from 'vue';

interface WebsocketWrapperData {
    webSocketClient: WebSocket | null;
}

export default defineComponent({
    data(): WebsocketWrapperData {
        return {
            webSocketClient: null,
        };
    },

    emits: {
        error(message: string): boolean {
            return message.length > 0;
        },
        message(wsText: Message) {
            return (wsText.type as number) >= 0;
        },
    },

    methods: {
        initiateWebSocket(webSocketUrl: string): Promise<WebSocket> {
            return new Promise(() => {
                const websocketClient = new WebSocket(webSocketUrl);

                return websocketClient;
            });
        },
        onIncomingWebSocketMessage(event: MessageEvent) {
            const message = JSON.parse(event.data) as Message;

            this.$emit('message', message);
        },
    },

    mounted() {
        this.initiateWebSocket(this.webSocketUrl)
            .then((webSocketClient: WebSocket) => {
                this.webSocketClient = webSocketClient;

                this.webSocketClient.onmessage =
                    this.onIncomingWebSocketMessage;
            })
            .catch((reason: string) => {
                this.$emit('error', reason);
            });
    },

    props: {
        webSocketUrl: {
            type: String,
            required: true,
        },
        localMessage: {
            type: Object as PropType<Message | null>,
            required: true,
        },
    },

    unmounted() {
        this.webSocketClient?.close();
    },

    watch: {
        localMessage(newValue: Message | null) {
            if (newValue) {
                this.webSocketClient?.send(JSON.stringify(newValue));
            }
        },
    },
});
</script>

<style></style>
