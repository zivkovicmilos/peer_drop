export interface IKeyManagerProps {
  visibleTypes: EKeyInputType[];

  setAddedKeys: (keys: IKeysWrapper) => void;
  addedKeys: IKeysWrapper;
}

export interface IKeysWrapper {
  keys: string[];
}

export enum EKeyInputType {
  IMPORT = 'Import key',
  ENTER = 'Enter key'
}
