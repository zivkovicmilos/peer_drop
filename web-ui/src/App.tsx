import MomentUtils from '@date-io/moment';
import { CssBaseline, withStyles } from '@material-ui/core';
import { MuiPickersUtilsProvider } from '@material-ui/pickers';
import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import './App.css';
import { SnackbarProvider } from './components/molecules/Snackbar/snackbar.context';
import useSnackbar from './components/molecules/Snackbar/useSnackbar.hook';
import SessionContext, { ISessionContext } from './context/SessionContext';
import { ESearchContext } from './context/sessionContext.types';
import AppRouter from './router/AppRouter';
import IdentitiesService from './services/identities/identitiesService';
import { IIdentityResponse } from './services/identities/identitiesService.types';
import { globalStyles } from './shared/assets/styles/global.styles';
import ThemeProvider from './theme/ThemeProvider';

function App() {
  const [userIdentity, setUserIdentity] = useState<IIdentityResponse | null>(
    null
  );
  const [searchContext, setSearchContext] = useState<ESearchContext>(
    ESearchContext.WORKSPACES
  );

  const { openSnackbar } = useSnackbar();

  useEffect(() => {
    const fetchPrimaryIdentity = async () => {
      return await IdentitiesService.getPrimaryIdentity();
    };

    fetchPrimaryIdentity()
      .then((response) => {
        setUserIdentity(response);
      })
      .catch((err) => {
        if (err && err.response && err.response.status == 404) {
          setUserIdentity(null);
        } else {
          openSnackbar('Unable to load primary identity', 'error');
        }
      });
  }, []);

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
