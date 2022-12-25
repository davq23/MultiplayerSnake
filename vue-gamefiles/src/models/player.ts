import type Direction from './direction';
import type Position from './position';

export default interface Player {
    player_id: string;
    name: string;
    color: string;
    positions: Position[];
    direction: Direction;
    score: number | null | undefined;
}
