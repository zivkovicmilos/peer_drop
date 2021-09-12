import { Box, IconButton } from '@material-ui/core';
import ArrowBackRoundedIcon from '@material-ui/icons/ArrowBackRounded';
import { FC, useState } from 'react';
import NewWorkspaceContext, {
  INewWorkspaceContext
} from '../../../context/NewWorkspaceContext';
import {
  ENewWorkspaceType,
  ENWAccessControl,
  INWAccessControlContacts,
  INWAccessControlPassword,
  INWPermissions
} from '../../../context/newWorkspaceContext.types';
import Link from '../../atoms/Link/Link';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import NewWorkspaceParameters from '../../molecules/NewWorkspaceParameters/NewWorkspaceParameters';
import NewWorkspacePermissions from '../../molecules/NewWorkspacePermissions/NewWorkspacePermissions';
import NewWorkspaceReview from '../../molecules/NewWorkspaceReview/NewWorkspaceReview';
import NewWorkspaceSecurity from '../../molecules/NewWorkspaceSecurity/NewWorkspaceSecurity';
import NewWorkspaceSteps from '../../molecules/NewWorkspaceSteps/NewWorkspaceSteps';
import NewWorkspaceSuccess from '../../molecules/NewWorkspaceSuccess/NewWorkspaceSuccess';
import { ENewWorkspaceStep, INewWorkspaceProps } from './newWorkspace.types';

const NewWorkspace: FC<INewWorkspaceProps> = () => {
  const steps: ENewWorkspaceStep[] = [
    ENewWorkspaceStep.PARAMS,
    ENewWorkspaceStep.SECURITY,
    ENewWorkspaceStep.PERMISSIONS,
    ENewWorkspaceStep.REVIEW
  ];

  const [step, setStep] = useState<number>(0);
  const [sectionTitle, setSectionTitle] = useState<string>('New Workspace');
  const [workspaceMnemonic, setWorkspaceMnemonic] = useState<string>('');
  const [workspaceName, setWorkspaceName] = useState<string>('');
  const [workspaceType, setWorkspaceType] = useState<ENewWorkspaceType>(
    ENewWorkspaceType.SEND_ONLY
  );

  const [accessControl, setAccessControl] = useState<
    INWAccessControlContacts | INWAccessControlPassword
  >({
    contacts: []
  });

  const [accessControlType, setAccessControlType] = useState<ENWAccessControl>(
    ENWAccessControl.SPECIFIC_CONTACTS
  );

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

  const handleBack = () => {
    if (step > 0) {
      let newStep = step - 1;
      setStep(newStep);
    }
  };

  const handleNext = () => {
    if (step < steps.length) {
      let newStep = step + 1;
      setStep(newStep);
    }
  };

  const newWorkspaceContextValue: INewWorkspaceContext = {
    step,
    setStep,
    workspaceName,
    workspaceType,
    setWorkspaceName,
    setWorkspaceType,
    accessControl,
    accessControlType,
    setAccessControl,
    setAccessControlType,
    permissions,
    setPermissions,
    sectionTitle,
    setSectionTitle,
    workspaceMnemonic,
    setWorkspaceMnemonic,
    handleBack,
    handleNext
  };

  const renderSection = () => {
    switch (step) {
      case 0:
        // Parameters
        return <NewWorkspaceParameters />;
      case 1:
        // Security
        return <NewWorkspaceSecurity />;
      case 2:
        // Permissions
        return <NewWorkspacePermissions />;
      case 3:
        // Review
        return <NewWorkspaceReview />;
      case 4:
        // Success
        return <NewWorkspaceSuccess />;
    }
  };

  return (
    <NewWorkspaceContext.Provider value={newWorkspaceContextValue}>
      <Box display={'flex'} flexDirection={'column'}>
        <Box display={'flex'} alignItems={'center'}>
          <Link to={'/workspaces'}>
            <IconButton
              classes={{
                root: 'iconButtonRoot'
              }}
            >
              <ArrowBackRoundedIcon
                style={{
                  fill: 'black'
                }}
              />
            </IconButton>
          </Link>
          <Box ml={2}>
            <PageTitle title={sectionTitle} />
          </Box>
        </Box>
        <Box display={'flex'} width={'100%'} mt={4}>
          {renderSection()}

          <Box
            ml={'auto'}
            display={sectionTitle == 'New Workspace' ? 'flex' : 'none'}
          >
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
