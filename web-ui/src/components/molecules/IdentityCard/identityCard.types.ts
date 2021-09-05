import { IIdentity } from '../../pages/Identities/identities.types';

export interface IIdentityCardProps
  extends Pick<
    IIdentity,
    'picture' | 'name' | 'publicKeyID' | 'numWorkspaces' | 'creationDate'
  > {}

export enum EIdentityCardMenuItem {
  EDIT = 'Edit',
  SHARE = 'Share',
  BACKUP = 'Backup',
  SET_IDENTITY = 'Set identity'
}
