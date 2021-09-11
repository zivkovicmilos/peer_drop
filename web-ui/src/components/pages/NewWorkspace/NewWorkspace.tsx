import { Box, IconButton } from '@material-ui/core';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import { FC, useState } from 'react';
import NewWorkspaceContext, {
  INewWorkspaceContext
} from '../../../context/NewWorkspaceContext';
import {
  ENewWorkspaceType,
  INWAccessControlContacts,
  INWAccessControlPassword,
  INWPermissions
} from '../../../context/newWorkspaceContext.types';
import Link from '../../atoms/Link/Link';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import NewWorkspaceSteps from '../../molecules/NewWorkspaceSteps/NewWorkspaceSteps';
import { ENewWorkspaceStep, INewWorkspaceProps } from './newWorkspace.types';

const NewWorkspace: FC<INewWorkspaceProps> = () => {
  const steps: ENewWorkspaceStep[] = [
    ENewWorkspaceStep.PARAMS,
    ENewWorkspaceStep.SECURITY,
    ENewWorkspaceStep.PERMISSIONS,
    ENewWorkspaceStep.REVIEW
  ];

  const [step, setStep] = useState<number>(0);
  const [workspaceName, setWorkspaceName] = useState<string>('');
  const [workspaceType, setWorkspaceType] = useState<ENewWorkspaceType>(
    ENewWorkspaceType.SEND_ONLY
  );

  const [accessControl, setAccessControl] = useState<
    INWAccessControlContacts | INWAccessControlPassword
  >({
    contactIDs: []
  });

  const [permissions, setPermissions] = useState<INWPermissions>({
    autocloseWorkspace: {
      active: false
    },
    enforcePeerLimit: {
      active: false
    },
    additionalOwners: {
      active: false
    }
  });

  const newWorkspaceContextValue: INewWorkspaceContext = {
    workspaceName,
    workspaceType,
    setWorkspaceName,
    setWorkspaceType,
    accessControl,
    setAccessControl,
    permissions,
    setPermissions
  };

  return (
    <NewWorkspaceContext.Provider value={newWorkspaceContextValue}>
      <Box display={'flex'} flexDirection={'column'}>
        <Box display={'flex'} alignItems={'center'}>
          <Link to={'/workspaces'}>
            <IconButton>
              <ArrowBackRoundedIcon
                style={{
                  fill: 'black'
                }}
              />
            </IconButton>
          </Link>
          <Box>
            <PageTitle title={'New Workspace'} />
          </Box>
        </Box>
        <Box display={'flex'} width={'100%'} mt={4}>
          <Box>This is the area</Box>
          <Box ml={'auto'}>
            <NewWorkspaceSteps
              currentStep={steps[step]}
              currentStepIndex={step}
              steps={steps}
            />
          </Box>
        </Box>
      </Box>
    </NewWorkspaceContext.Provider>
  );
};

export default NewWorkspace;
