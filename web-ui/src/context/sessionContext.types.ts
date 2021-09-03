export interface IUserIdentity {
  name: string;
  keyID: string;
  picture: string;
}

export enum ESearchContext {
  WORKSPACES = 'workspaces',
  CONTACTS = 'contacts',
  IDENTITIES = 'identities'
}
