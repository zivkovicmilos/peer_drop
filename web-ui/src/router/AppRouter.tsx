import { Redirect, Route, Switch } from 'react-router-dom';
import AppLayout from '../components/layouts/AppLayout/AppLayout';
import Contacts from '../components/pages/Contacts/Contacts';
import Identities from '../components/pages/Identities/Identities';
import Settings from '../components/pages/Settings/Settings';
import Workspaces from '../components/pages/Workspaces/Workspaces';

const AppRouter = () => (
  <AppLayout>
    <Switch>
      <Route path={'/workspaces'} exact={true} component={Workspaces} />
      <Route path={'/contacts'} exact={true} component={Contacts} />
      <Route path={'/identities'} exact={true} component={Identities} />
      <Route path={'/settings'} exact={true} component={Settings} />
      <Redirect to={'/workspaces'} />
    </Switch>
  </AppLayout>
);

export default AppRouter;
