export interface IContactEditProps {
  type?: EContactEditType;
}

export enum EContactEditType {
  NEW,
  EDIT
}

export interface IContactEditParams {
  contactId: string | null;
}

export interface IKeyPair {
  keyID: string;
  publicKey: string;
  privateKey: string;
}