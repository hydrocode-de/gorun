export type InputLiteral = string | number | boolean | Date | null
export type InputArray = string[] | number[] | boolean[] | Date[]
export type InputStruct = {[key: string]: any}

export type InputParams = InputLiteral | InputArray | InputStruct