import { IKeyPair } from '../../pages/ContactEdit/contactEdit.types';

export interface IKeyManagerProps {
  visibleTypes: EKeyInputType[];

  setAddedKey: (key: IKeyPair) => void;
  addedKey: IKeyPair;
  keyListTitle?: string;
  formik?: any;

  expectedType: EKeyType;
}

export enum EKeyInputType {
  IMPORT = 'Import key',
  ENTER = 'Enter key',
  GENERATE = 'Generate key'
}

export enum EKeyGenerateType {
  RSA_2048 = 'RSA-2048',
  RSA_4096 = 'RSA-4096'
}

export enum EKeyType {
  PUBLIC,
  PRIVATE
}
