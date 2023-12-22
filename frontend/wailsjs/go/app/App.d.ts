// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {app} from '../models';
import {map[string]app} from '../models';
import {menu} from '../models';
import {context} from '../models';
import {ui} from '../models';

export function AddHost(arg1:string):Promise<string>;

export function Aggregate(arg1:string,arg2:string,arg3:string,arg4:string,arg5:string):Promise<void>;

export function AskConfirmation(arg1:string):Promise<boolean>;

export function ChooseDirectory(arg1:string):Promise<string>;

export function CountItems(arg1:string,arg2:string,arg3:string,arg4:string):Promise<app.CountItemsResult>;

export function CreateIndex(arg1:string,arg2:string,arg3:string,arg4:string):Promise<string>;

export function DropCollection(arg1:string,arg2:string,arg3:string):Promise<boolean>;

export function DropDatabase(arg1:string,arg2:string):Promise<boolean>;

export function DropIndex(arg1:string,arg2:string,arg3:string,arg4:string):Promise<boolean>;

export function DuplicateCollection(arg1:string,arg2:string,arg3:string,arg4:string,arg5:string,arg6:string):Promise<boolean>;

export function Environment():Promise<app.EnvironmentInfo>;

export function ExecuteShellScript(arg1:string,arg2:string,arg3:string,arg4:string):Promise<app.ExecuteShellScriptResult>;

export function FindItems(arg1:string,arg2:string,arg3:string,arg4:string):Promise<app.FindItemsResult>;

export function GetIndexes(arg1:string,arg2:string,arg3:string):Promise<app.GetIndexesResult>;

export function HostLogs(arg1:string,arg2:string):Promise<app.HostLogsResult>;

export function Hosts():Promise<map[string]app.Host>;

export function InsertItems(arg1:string,arg2:string,arg3:string,arg4:string):Promise<any>;

export function Menu():Promise<menu.Menu>;

export function OpenCollection(arg1:string,arg2:string,arg3:string):Promise<app.OpenCollectionResult>;

export function OpenConnection(arg1:string):Promise<app.OpenConnectionResult>;

export function OpenDatabase(arg1:string,arg2:string):Promise<app.OpenDatabaseResult>;

export function OpenShellScript():Promise<string>;

export function PerformDump(arg1:string):Promise<boolean>;

export function PerformFindExport(arg1:string,arg2:string,arg3:string,arg4:string):Promise<boolean>;

export function PurgeLogDirectory():Promise<void>;

export function RemoveHost(arg1:string):Promise<boolean>;

export function RemoveItemById(arg1:string,arg2:string,arg3:string,arg4:string):Promise<boolean>;

export function RemoveItems(arg1:string,arg2:string,arg3:string,arg4:string,arg5:boolean):Promise<number>;

export function RemoveQuery(arg1:string):Promise<void>;

export function RemoveView(arg1:string):Promise<void>;

export function RenameCollection(arg1:string,arg2:string,arg3:string,arg4:string):Promise<boolean>;

export function ReportSharedStateVariable(arg1:string,arg2:string):Promise<void>;

export function SaveQuery(arg1:string):Promise<string>;

export function SaveShellScript(arg1:string,arg2:string,arg3:string,arg4:string,arg5:boolean):Promise<app.SaveShellScriptResult>;

export function SavedQueries():Promise<map[string]app.SavedQuery>;

export function Settings():Promise<app.Settings>;

export function Startup(arg1:context.Context,arg2:ui.UI):Promise<void>;

export function TruncateCollection(arg1:string,arg2:string,arg3:string):Promise<boolean>;

export function UpdateFoundDocument(arg1:string,arg2:string,arg3:string,arg4:string,arg5:string):Promise<boolean>;

export function UpdateHost(arg1:string,arg2:string):Promise<boolean>;

export function UpdateItems(arg1:string,arg2:string,arg3:string,arg4:string):Promise<number>;

export function UpdateQueries(arg1:string):Promise<boolean>;

export function UpdateSettings(arg1:string):Promise<app.Settings>;

export function UpdateViewStore(arg1:string):Promise<void>;

export function Views():Promise<app.ViewStore>;
