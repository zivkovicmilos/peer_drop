import { Redirect, Route, Switch } from 'react-router-dom';
import AppLayout from '../components/layouts/AppLayout/AppLayout';
import ContactEdit from '../components/pages/ContactEdit/ContactEdit';
import { EContactEditType } from '../components/pages/ContactEdit/contactEdit.types';
import Contacts from '../components/pages/Contacts/Contacts';
import Identities from '../components/pages/Identities/Identities';
import Settings from '../components/pages/Settings/Settings';
import Workspaces from '../components/pages/Workspaces/Workspaces';

const AppRouter = () => (
  <AppLayout>
    <Switch>
      <Route path={'/workspaces'} exact={true} component={Workspaces} />

      <Route path={'/contacts'} exact={true} component={Contacts} />
      <Route path={'/contacts/new'} exact={true}>
        <ContactEdit type={EContactEditType.NEW} />
      </Route>
      <Route path={'/contacts/:contactId/edit'} exact={true}>
        <ContactEdit type={EContactEditType.EDIT} />
      </Route>

      <Route path={'/identities'} exact={true} component={Identities} />
      <Route path={'/settings'} exact={true} component={Settings} />
      <Redirect to={'/workspaces'} />
    </Switch>
  </AppLayout>
);

export default AppRouter;
