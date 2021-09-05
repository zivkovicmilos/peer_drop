import { IIdentity } from '../../pages/Identities/identities.types';

export interface IIdentityCardProps
  extends Pick<
    IIdentity,
    'picture' | 'name' | 'publicKeyID' | 'numWorkspaces' | 'creationDate'
  > {}{

}