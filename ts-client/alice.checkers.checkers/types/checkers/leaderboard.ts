/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { WinningPlayer } from "./winning_player";

export const protobufPackage = "alice.checkers.checkers";

export interface Leaderboard {
  winners: WinningPlayer[];
}

function createBaseLeaderboard(): Leaderboard {
  return { winners: [] };
}

export const Leaderboard = {
  encode(message: Leaderboard, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.winners) {
      WinningPlayer.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Leaderboard {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLeaderboard();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.winners.push(WinningPlayer.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Leaderboard {
    return { winners: Array.isArray(object?.winners) ? object.winners.map((e: any) => WinningPlayer.fromJSON(e)) : [] };
  },

  toJSON(message: Leaderboard): unknown {
    const obj: any = {};
    if (message.winners) {
      obj.winners = message.winners.map((e) => e ? WinningPlayer.toJSON(e) : undefined);
    } else {
      obj.winners = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Leaderboard>, I>>(object: I): Leaderboard {
    const message = createBaseLeaderboard();
    message.winners = object.winners?.map((e) => WinningPlayer.fromPartial(e)) || [];
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };
