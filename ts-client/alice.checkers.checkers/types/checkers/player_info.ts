/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "alice.checkers.checkers";

export interface PlayerInfo {
  index: string;
  wonCount: number;
  lostCount: number;
  forfeitedCount: number;
}

function createBasePlayerInfo(): PlayerInfo {
  return { index: "", wonCount: 0, lostCount: 0, forfeitedCount: 0 };
}

export const PlayerInfo = {
  encode(message: PlayerInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.index !== "") {
      writer.uint32(10).string(message.index);
    }
    if (message.wonCount !== 0) {
      writer.uint32(16).uint64(message.wonCount);
    }
    if (message.lostCount !== 0) {
      writer.uint32(24).uint64(message.lostCount);
    }
    if (message.forfeitedCount !== 0) {
      writer.uint32(32).uint64(message.forfeitedCount);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PlayerInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePlayerInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.string();
          break;
        case 2:
          message.wonCount = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.lostCount = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.forfeitedCount = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PlayerInfo {
    return {
      index: isSet(object.index) ? String(object.index) : "",
      wonCount: isSet(object.wonCount) ? Number(object.wonCount) : 0,
      lostCount: isSet(object.lostCount) ? Number(object.lostCount) : 0,
      forfeitedCount: isSet(object.forfeitedCount) ? Number(object.forfeitedCount) : 0,
    };
  },

  toJSON(message: PlayerInfo): unknown {
    const obj: any = {};
    message.index !== undefined && (obj.index = message.index);
    message.wonCount !== undefined && (obj.wonCount = Math.round(message.wonCount));
    message.lostCount !== undefined && (obj.lostCount = Math.round(message.lostCount));
    message.forfeitedCount !== undefined && (obj.forfeitedCount = Math.round(message.forfeitedCount));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<PlayerInfo>, I>>(object: I): PlayerInfo {
    const message = createBasePlayerInfo();
    message.index = object.index ?? "";
    message.wonCount = object.wonCount ?? 0;
    message.lostCount = object.lostCount ?? 0;
    message.forfeitedCount = object.forfeitedCount ?? 0;
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
