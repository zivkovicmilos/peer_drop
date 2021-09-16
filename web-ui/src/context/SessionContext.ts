import { createContext } from 'react';
import { IIdentityResponse } from '../services/identities/identitiesService.types';
import { ESearchContext } from './sessionContext.types';

export interface ISessionContext<Identity = IIdentityResponse | null> {
  userIdentity: Identity;
  searchContext: ESearchContext;
  setUserIdentity: (user: IIdentityResponse | null) => void;
  setSearchContext: (searchContext: ESearchContext | null) => void;
}

const SessionContext = createContext<ISessionContext>({
  userIdentity: null,
  searchContext: null,
  setUserIdentity: (user: IIdentityResponse | null) => {},
  setSearchContext: (searchContext: ESearchContext) => {}
});

export default SessionContext;

export const SessionContextConsumer = SessionContext.Consumer;
