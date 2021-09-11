import MomentUtils from '@date-io/moment';
import { CssBaseline, withStyles } from '@material-ui/core';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
import React, { useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import './App.css';
import { SnackbarProvider } from './components/molecules/Snackbar/snackbar.context';
import SessionContext, { ISessionContext } from './context/SessionContext';
import { ESearchContext, IUserIdentity } from './context/sessionContext.types';
import AppRouter from './router/AppRouter';
import { globalStyles } from './shared/assets/styles/global.styles';
import ThemeProvider from './theme/ThemeProvider';

function App() {
  const [userIdentity, setUserIdentity] = useState<IUserIdentity | null>({
    name: 'Milos',
    keyID: '4AEE18F83AFDEB23',
    picture: 'https://static.dw.com/image/58133780_6.jpg'
  });
  const [searchContext, setSearchContext] = useState<ESearchContext>(
    ESearchContext.WORKSPACES
  );

  // TODO load user identity from disk

  const sessionContextValue: ISessionContext = {
    userIdentity,
    searchContext,
    setUserIdentity,
    setSearchContext
  };

  return (
    <SessionContext.Provider value={sessionContextValue}>
      <MuiPickersUtilsProvider utils={MomentUtils}>
        <SnackbarProvider>
          <ThemeProvider>
            <CssBaseline />
            <Router>
              <AppRouter />
            </Router>
          </ThemeProvider>
        </SnackbarProvider>
      </MuiPickersUtilsProvider>
    </SessionContext.Provider>
  );
}

export default withStyles(globalStyles)(App);
