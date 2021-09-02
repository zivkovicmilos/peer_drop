import MomentUtils from '@date-io/moment';
import { CssBaseline, withStyles } from '@material-ui/core';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
import React, { useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import './App.css';
import { SnackbarProvider } from './components/molecules/Snackbar/snackbar.context';
import SessionContext, { ISessionContext } from './context/SessionContext';
import { IUserIdentity } from './context/sessionContext.types';
import AppRouter from './router/AppRouter';
import { globalStyles } from './shared/assets/styles/global.styles';
import ThemeProvider from './theme/ThemeProvider';

function App() {
  const [userIdentity, setUserIdentity] = useState<IUserIdentity | null>(null);

  const sessionContextValue: ISessionContext = {
    userIdentity,
    setUserIdentity
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
