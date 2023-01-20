// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {app} from '../models';
import {primitive} from '../models';
import {map[string]app} from '../models';

export function AddHost(arg1:string):Promise<void>;

export function DropCollection(arg1:string,arg2:string,arg3:string):Promise<boolean>;

export function DropDatabase(arg1:string,arg2:string):Promise<boolean>;

export function DropIndex(arg1:string,arg2:string,arg3:string,arg4:string):Promise<boolean>;

export function FindItems(arg1:string,arg2:string,arg3:string,arg4:string):Promise<app.findResult>;

export function GetIndexes(arg1:string,arg2:string,arg3:string):Promise<Array<primitive.M>>;

export function Hosts():Promise<map[string]app.Host>;

export function InsertItems(arg1:string,arg2:string,arg3:string,arg4:string):Promise<any>;

export function OpenCollection(arg1:string,arg2:string,arg3:string):Promise<primitive.M>;

export function OpenConnection(arg1:string):Promise<Array<string>>;

export function OpenDatabase(arg1:string,arg2:string):Promise<Array<string>>;

export function RemoveHost(arg1:string):Promise<void>;

export function RemoveItemById(arg1:string,arg2:string,arg3:string,arg4:string):Promise<boolean>;

export function RemoveItems(arg1:string,arg2:string,arg3:string,arg4:string,arg5:boolean):Promise<number>;

export function Settings():Promise<app.Settings>;

export function UpdateHost(arg1:string,arg2:string):Promise<void>;

export function UpdateItems(arg1:string,arg2:string,arg3:string,arg4:string):Promise<number>;

export function UpdateSettings(arg1:string):Promise<app.Settings>;
