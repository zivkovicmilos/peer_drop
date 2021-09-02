import {createContext} from 'react';
import {IUserIdentity} from "./sessionContext.types";

export interface ISessionContext<User = IUserIdentity | null> {
  userIdentity: User;
  setUserIdentity: (user: IUserIdentity | null) => void;
}

const SessionContext = createContext<ISessionContext>({
  userIdentity: null,
  setUserIdentity: (user: IUserIdentity | null) => {
  },
});

export default SessionContext;

export const SessionContextConsumer = SessionContext.Consumer;



