import { IIdentityResponse } from '../../../services/identities/identitiesService.types';

export interface IIdentityCardProps extends IIdentityResponse {
  triggerUpdate: boolean;
  setTriggerUpdate: (triggerUpdate: boolean) => void;
}

export enum EIdentityCardMenuItem {
  EDIT = 'Edit',
  SHARE = 'Share',
  BACKUP = 'Backup',
  DELETE = 'Delete',
  SET_AS_PRIMARY = 'Set as primary'
}
