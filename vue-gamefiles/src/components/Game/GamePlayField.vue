<template>
    <canvas ref="playFieldCanvas" :width="width" :height="height"></canvas>
</template>

<script lang="ts">
import type Player from '@/models/player';
import { defineComponent, ref } from 'vue';
import type { PropType } from 'vue';

export default defineComponent({
    methods: {
        drawAnimationFrame() {
            if (!this.mapPlayers || !this.playFieldCanvas) {
                return;
            }
            const playFieldContext = this.playFieldCanvas.getContext('2d');

            if (playFieldContext) {
                playFieldContext.fillStyle = 'black';

                playFieldContext.clearRect(
                    0,
                    0,
                    this.playFieldCanvas.width,
                    this.playFieldCanvas.height
                );
                playFieldContext.fillRect(
                    0,
                    0,
                    this.playFieldCanvas.width,
                    this.playFieldCanvas.height
                );

                this.mapPlayers.forEach((player: Player) => {
                    this.renderPlayer(playFieldContext, player);
                });
            }
        },
        renderPlayer(context: CanvasRenderingContext2D | null, player: Player) {
            if (!context) {
                return;
            }

            // Draw player positions
            player.positions.forEach(({ x, y }) => {
                context.beginPath();
                context.arc(x, y, this.playerRadius, 0, 2 * Math.PI);
                context.fillStyle = player.color;
                context.fill();
                context.lineWidth = 1;
                context.strokeStyle = '#003300';
                context.stroke();
                context.closePath();
            });
        },
    },
    props: {
        width: {
            type: Number,
            required: true,
        },
        height: {
            type: Number,
            required: true,
        },
        playerRadius: {
            type: Number,
            required: true,
        },
        mapPlayers: {
            type: Object as PropType<Map<string, Player> | null>,
            required: true,
        },
    },
    setup() {
        const playFieldCanvas = ref<HTMLCanvasElement | null>(null);

        return {
            playFieldCanvas,
        };
    },
});
</script>

<style></style>
