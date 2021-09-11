export interface IIdentityEditProps {
  type?: EIdentityEditType;
}

export enum EIdentityEditType {
  NEW,
  EDIT
}

export interface IIdentityEditParams {
  identityId: string | null;
}