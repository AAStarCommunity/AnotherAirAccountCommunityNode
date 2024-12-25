export interface AuthData {
    message: string[];
    signatures: string[];
    pubkeys: string[][];
}

// 如果需要单独使用某些类型，也可以分别定义
export type Message = string[];
export type Signatures = string[];
export type Pubkeys = string[][];
