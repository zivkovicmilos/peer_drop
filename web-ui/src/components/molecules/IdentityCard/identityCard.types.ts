import { IIdentityResponse } from '../../../services/identities/identitiesService.types';

export interface IIdentityCardProps extends IIdentityResponse {
}

export enum EIdentityCardMenuItem {
  EDIT = 'Edit',
  SHARE = 'Share',
  BACKUP = 'Backup',
  SET_IDENTITY = 'Set identity'
}
