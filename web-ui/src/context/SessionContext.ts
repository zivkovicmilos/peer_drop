import { createContext } from 'react';
import { ESearchContext, IUserIdentity } from './sessionContext.types';

export interface ISessionContext<User = IUserIdentity | null> {
  userIdentity: User;
  searchContext: ESearchContext;
  setUserIdentity: (user: IUserIdentity | null) => void;
  setSearchContext: (searchContext: ESearchContext | null) => void;
}

const SessionContext = createContext<ISessionContext>({
  userIdentity: null,
  searchContext: null,
  setUserIdentity: (user: IUserIdentity | null) => {},
  setSearchContext: (searchContext: ESearchContext) => {}
});

export default SessionContext;

export const SessionContextConsumer = SessionContext.Consumer;
