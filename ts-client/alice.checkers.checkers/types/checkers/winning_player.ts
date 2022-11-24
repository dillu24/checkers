/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "alice.checkers.checkers";

export interface WinningPlayer {
  playerAddress: string;
  wonCount: number;
  dateAdded: string;
}

function createBaseWinningPlayer(): WinningPlayer {
  return { playerAddress: "", wonCount: 0, dateAdded: "" };
}

export const WinningPlayer = {
  encode(message: WinningPlayer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.playerAddress !== "") {
      writer.uint32(10).string(message.playerAddress);
    }
    if (message.wonCount !== 0) {
      writer.uint32(16).uint64(message.wonCount);
    }
    if (message.dateAdded !== "") {
      writer.uint32(26).string(message.dateAdded);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WinningPlayer {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWinningPlayer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.playerAddress = reader.string();
          break;
        case 2:
          message.wonCount = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.dateAdded = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): WinningPlayer {
    return {
      playerAddress: isSet(object.playerAddress) ? String(object.playerAddress) : "",
      wonCount: isSet(object.wonCount) ? Number(object.wonCount) : 0,
      dateAdded: isSet(object.dateAdded) ? String(object.dateAdded) : "",
    };
  },

  toJSON(message: WinningPlayer): unknown {
    const obj: any = {};
    message.playerAddress !== undefined && (obj.playerAddress = message.playerAddress);
    message.wonCount !== undefined && (obj.wonCount = Math.round(message.wonCount));
    message.dateAdded !== undefined && (obj.dateAdded = message.dateAdded);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<WinningPlayer>, I>>(object: I): WinningPlayer {
    const message = createBaseWinningPlayer();
    message.playerAddress = object.playerAddress ?? "";
    message.wonCount = object.wonCount ?? 0;
    message.dateAdded = object.dateAdded ?? "";
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
