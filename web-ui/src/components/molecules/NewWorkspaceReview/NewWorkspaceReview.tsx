import { Box } from '@material-ui/core';
import { FC, useContext } from 'react';
import { useHistory } from 'react-router-dom';
import NewWorkspaceContext from '../../../context/NewWorkspaceContext';
import {
  ENWAccessControl,
  INWAccessControlContacts,
  INWAccessControlPassword
} from '../../../context/newWorkspaceContext.types';
import SessionContext from '../../../context/SessionContext';
import WorkspacesService from '../../../services/workspaces/workspacesService';
import { INewWorkspaceRequest } from '../../../services/workspaces/workspacesService.types';
import StepButton from '../../atoms/StepButton/StepButton';
import NewWorkspaceInfo from '../NewWorkspaceInfo/NewWorkspaceInfo';
import useSnackbar from '../Snackbar/useSnackbar.hook';
import { INewWorkspaceReviewProps } from './newWorkspaceReview.types';

const NewWorkspaceReview: FC<INewWorkspaceReviewProps> = () => {
  const {
    handleNext,
    handleBack,
    step,
    setSectionTitle,
    setWorkspaceMnemonic,
    workspaceName,
    workspaceType,
    accessControlType,
    accessControl,
    permissions // TODO: ignored for now
  } = useContext(NewWorkspaceContext);

  const { userIdentity } = useContext(SessionContext);

  const constructNewWorkspaceRequest = (): INewWorkspaceRequest => {
    let workspaceRequest: INewWorkspaceRequest = {
      workspaceName,
      workspaceType,
      workspaceAccessControlType: accessControlType,
      workspaceAccessControl: {},
      baseWorkspaceOwnerKeyID: userIdentity.publicKeyID,
      workspaceAdditionalOwnerPublicKeys: []
    };

    if (accessControlType == ENWAccessControl.PASSWORD) {
      workspaceRequest.workspaceAccessControl = {
        password: (accessControl as INWAccessControlPassword).password
      };
    } else {
      workspaceRequest.workspaceAccessControl = {
        contacts: (accessControl as INWAccessControlContacts).contacts
      };
    }

    return workspaceRequest;
  };

  const { openSnackbar } = useSnackbar();
  const history = useHistory();

  const handleConfirm = () => {
    const createNewWorkspace = async () => {
      return await WorkspacesService.createWorkspace(
        constructNewWorkspaceRequest()
      );
    };

    createNewWorkspace()
      .then((response) => {
        setWorkspaceMnemonic(response.mnemonic);

        setSectionTitle('Workspace created!');
        handleNext();
      })
      .catch((err) => {
        history.push('/workspaces');
        openSnackbar('Unable to create new workspace', 'error');
      });
  };

  return (
    <Box display={'flex'} flexDirection={'column'} width={'50%'}>
      <NewWorkspaceInfo
        accessControl={accessControl}
        accessControlType={accessControlType}
        permissions={permissions}
        workspaceName={workspaceName}
        workspaceType={workspaceType}
      />
      <Box display={'flex'} alignItems={'center'}>
        <Box mr={2}>
          <StepButton
            text={'Back'}
            variant={'outlined'}
            disabled={step < 1}
            shouldSubmit={false}
            onClick={() => {
              handleBack();
            }}
          />
        </Box>
        <StepButton
          text={'Confirm'}
          variant={'contained'}
          shouldSubmit={false}
          onClick={() => {
            handleConfirm();
          }}
        />
      </Box>
    </Box>
  );
};

export default NewWorkspaceReview;
