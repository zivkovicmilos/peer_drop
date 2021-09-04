export interface IKeyManagerProps {
  visibleTypes: EKeyInputType[];
}

export enum EKeyInputType {
  IMPORT = 'Import key',
  ENTER = 'Enter key',
  GENERATE = 'Generate key'
}
